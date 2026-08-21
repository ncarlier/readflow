package sanitizer

import (
	"net"
	"net/url"
	"strconv"
	"strings"
	"unicode"

	"github.com/go-shiori/dom"
	"github.com/microcosm-cc/bluemonday"
	"github.com/ncarlier/readflow/pkg/logger"
	"github.com/rs/zerolog"
	"github.com/tdewolff/minify/v2"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

const (
	urlKindLink = iota
	urlKindImage
	urlKindIframe
)

// Sanitizer is an HTML sanitizer.
type Sanitizer struct {
	logger    zerolog.Logger
	blockList *BlockList
	policy    *bluemonday.Policy
	minifier  *minify.M
}

// NewSanitizer creates a new HTML sanitizer.
func NewSanitizer(blockList *BlockList) *Sanitizer {
	log := logger.With().
		Str("component", "sanitizer").
		Logger()

	policy := bluemonday.UGCPolicy()

	policy.AddTargetBlankToFullyQualifiedLinks(true)

	policy.AllowAttrs(
		"width",
		"height",
		"src",
		"allowfullscreen",
		"sandbox",
	).OnElements("iframe")

	policy.AllowAttrs(
		"srcset",
		"sizes",
	).OnElements("img", "source")

	policy.AllowElements(
		"picture",
		"source",
	)

	if blockList != nil {
		log.Info().
			Str("location", blockList.Location()).
			Uint32("size", blockList.Size()).
			Msg("using block-list")
	}

	return &Sanitizer{
		logger:    log,
		blockList: blockList,
		policy:    policy,
		minifier:  newDefaultMinifier(),
	}
}

// cleanURLs removes URL attributes that are malformed, unsafe, or point to
// hosts present in the blocklist.
func (s *Sanitizer) cleanURLs(doc *html.Node) {
	var visit func(*html.Node)

	visit = func(node *html.Node) {
		if node.Type == html.ElementNode {
			node.Attr = s.cleanNodeAttributes(node)
		}

		for child := node.FirstChild; child != nil; child = child.NextSibling {
			visit(child)
		}
	}

	visit(doc)
}

// cleanNodeAttributes cleans URL-bearing attributes on one HTML element.
func (s *Sanitizer) cleanNodeAttributes(node *html.Node) []html.Attribute {
	var cleaned []html.Attribute

	for _, attr := range node.Attr {
		switch attr.Key {
		case "src":
			var kind int

			switch node.DataAtom {
			case atom.Img, atom.Source:
				kind = urlKindImage
			case atom.Iframe:
				kind = urlKindIframe
			default:
				// Do not apply media rules to unknown elements.
				cleaned = append(cleaned, attr)
				continue
			}

			if s.allowedURL(attr.Val, kind) {
				cleaned = append(cleaned, attr)
			}

		case "href":
			if node.DataAtom == atom.A &&
				s.allowedURL(attr.Val, urlKindLink) {
				cleaned = append(cleaned, attr)
			}

		case "srcset":
			if (node.DataAtom == atom.Img || node.DataAtom == atom.Source) &&
				s.allowedSrcSet(attr.Val) {
				cleaned = append(cleaned, attr)
			}

		default:
			cleaned = append(cleaned, attr)
		}
	}

	return cleaned
}

// allowedURL validates one URL and checks its hostname against the blocklist.
//
// The function deliberately does not try to emulate browser URL parsing.
// Ambiguous URLs are rejected rather than interpreted optimistically.
func (s *Sanitizer) allowedURL(raw string, kind int) bool {
	value := strings.TrimSpace(raw)
	if value == "" {
		return false
	}

	if containsUnsafeURLCharacters(value) {
		return false
	}

	parsed, err := url.Parse(value)
	if err != nil {
		// Fail closed. Never retain an unparseable URL.
		return false
	}

	scheme := strings.ToLower(parsed.Scheme)

	switch kind {
	case urlKindIframe:
		// Iframes should use explicit HTTP(S) URLs only.
		if scheme != "http" && scheme != "https" {
			return false
		}

	case urlKindImage:
		// Images may use relative URLs or explicit HTTP(S) URLs.
		// Data URLs are intentionally rejected because they are not needed
		// to protect against the host-based bypass and may carry active
		// content in other contexts.
		if scheme != "" && scheme != "http" && scheme != "https" {
			return false
		}

	case urlKindLink:
		// Permit relative links and ordinary HTTP(S), mailto, and tel links.
		// bluemonday performs an additional protocol allowlist check later.
		if scheme != "" &&
			scheme != "http" &&
			scheme != "https" &&
			scheme != "mailto" &&
			scheme != "tel" {
			return false
		}
	}

	// A URL with an HTTP(S) scheme must contain an authority and hostname.
	// This rejects:
	//
	//   http:////test/a.png
	//   http:///test/a.png
	//   http:/test/a.png
	//
	// These values can be interpreted as network URLs by browsers despite
	// net/url returning an empty hostname.
	if scheme == "http" || scheme == "https" {
		if parsed.Host == "" || parsed.Hostname() == "" {
			return false
		}

		// Userinfo can create confusing display and parsing behavior.
		if parsed.User != nil {
			return false
		}

		return !s.blockedHostname(parsed.Hostname())
	}

	// Protocol-relative URLs, such as //example.com/image.png, have no
	// scheme but do have an authority. They must still be checked.
	if parsed.Host != "" {
		if parsed.Hostname() == "" || parsed.User != nil {
			return false
		}

		return !s.blockedHostname(parsed.Hostname())
	}

	// mailto and tel URLs do not have a host.
	if scheme == "mailto" || scheme == "tel" {
		return parsed.Opaque != "" || parsed.Path != ""
	}

	// Relative URLs are allowed for links and images. They have no host to
	// compare against. Iframes were rejected above unless they use HTTP(S).
	if scheme == "" && parsed.Host == "" {
		return kind == urlKindLink || kind == urlKindImage
	}

	return false
}

// blockedHostname normalizes a hostname before querying the blocklist.
func (s *Sanitizer) blockedHostname(hostname string) bool {
	if s.blockList == nil {
		return false
	}

	host := strings.TrimSpace(strings.ToLower(hostname))
	host = strings.TrimSuffix(host, ".")

	if host == "" {
		return true
	}

	// Normalize IP addresses where possible. This makes equivalent textual
	// forms compare consistently with blocklist entries.
	if ip := net.ParseIP(host); ip != nil {
		host = ip.String()
	}

	return s.blockList.Contains(host)
}

// allowedSrcSet validates every candidate in a srcset attribute.
//
// This parser intentionally rejects data URLs and malformed candidates.
// A comma is safe as a separator here because data URLs are not allowed.
func (s *Sanitizer) allowedSrcSet(raw string) bool {
	value := strings.TrimSpace(raw)
	if value == "" || containsUnsafeURLCharacters(value) {
		return false
	}

	candidates := strings.Split(value, ",")
	if len(candidates) == 0 {
		return false
	}

	for _, candidate := range candidates {
		fields := strings.Fields(strings.TrimSpace(candidate))

		switch len(fields) {
		case 1:
			// URL without a descriptor.

		case 2:
			// URL with one valid width or pixel-density descriptor.
			if !validSrcSetDescriptor(fields[1]) {
				return false
			}

		default:
			return false
		}

		if !s.allowedURL(fields[0], urlKindImage) {
			return false
		}
	}

	return true
}

// validSrcSetDescriptor validates descriptors such as:
//
//	1x
//	2.0x
//	320w
func validSrcSetDescriptor(value string) bool {
	if len(value) < 2 {
		return false
	}

	last := value[len(value)-1]
	number := value[:len(value)-1]

	if number == "" {
		return false
	}

	switch last {
	case 'x':
		f, err := strconv.ParseFloat(number, 64)
		return err == nil && f > 0

	case 'w':
		n, err := strconv.ParseUint(number, 10, 64)
		return err == nil && n > 0

	default:
		return false
	}
}

// containsUnsafeURLCharacters rejects characters that can be normalized or
// interpreted inconsistently by HTML parsers and browsers.
func containsUnsafeURLCharacters(value string) bool {
	for _, r := range value {
		if unicode.IsControl(r) || r == '\\' {
			return true
		}
	}

	return false
}

// Sanitize sanitizes HTML content.
func (s *Sanitizer) Sanitize(content string) string {
	if s.blockList != nil {
		reader := strings.NewReader(content)

		doc, err := dom.Parse(reader)
		if err == nil {
			s.cleanURLs(doc)
			content = dom.InnerHTML(doc)
		}
	}

	minified, err := s.minifier.String("text/html", content)
	if err == nil {
		content = minified
	}

	return s.policy.Sanitize(content)
}
