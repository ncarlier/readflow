package sanitizer

import (
	"net/url"
	"strings"

	"github.com/go-shiori/dom"
	"github.com/microcosm-cc/bluemonday"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// Sanitizer is a HTML sanitizer
type Sanitizer struct {
	logger    zerolog.Logger
	blockList *BlockList
	policy    *bluemonday.Policy
}

// NewSanitizer create new HTML sanitizer
func NewSanitizer(blockList *BlockList) *Sanitizer {
	logger := log.With().Str("component", "sanitizer").Logger()
	policy := bluemonday.UGCPolicy()
	policy.AddTargetBlankToFullyQualifiedLinks(true)
	policy.AllowAttrs("width", "height", "src", "allowfullscreen", "sandbox").OnElements("iframe")
	policy.AllowAttrs("srcset", "sizes", "data-src").OnElements("img")

	if blockList != nil {
		logger.Info().
			Str("location", blockList.Location()).
			Uint32("size", blockList.Size()).
			Msg("using block-list")
	}

	return &Sanitizer{
		logger:    logger,
		blockList: blockList,
		policy:    policy,
	}
}

func (s *Sanitizer) cleanURLs(doc *html.Node) {
	var parser func(*html.Node)
	parser = func(node *html.Node) {
		if node.Type == html.ElementNode {
			switch node.DataAtom {
			case atom.A, atom.Img, atom.Iframe:
				var attrs []html.Attribute
				for _, a := range node.Attr {
					switch a.Key {
					case "src", "href":
						u, err := url.Parse(a.Val)
						if err != nil {
							attrs = append(attrs, a)
							continue
						}
						if !s.blockList.Contains(u.Hostname()) {
							attrs = append(attrs, a)
						}
					default:
						attrs = append(attrs, a)
					}
				}
				node.Attr = attrs
			}
		}

		for child := node.FirstChild; child != nil; child = child.NextSibling {
			parser(child)
		}
	}
	parser(doc)
}

// Sanitize HTML content
func (s *Sanitizer) Sanitize(content string) string {
	if s.blockList != nil {
		// convert content to DOM
		r := strings.NewReader(content)
		doc, err := dom.Parse(r)
		if err == nil {
			// clean URLs
			s.cleanURLs(doc)
			content = dom.InnerHTML(doc)
		}
	}

	return s.policy.Sanitize(content)
}
