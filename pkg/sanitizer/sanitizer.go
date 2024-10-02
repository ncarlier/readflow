package sanitizer

import (
	"net/url"
	"strings"

	"github.com/go-shiori/dom"
	"github.com/microcosm-cc/bluemonday"
	"github.com/ncarlier/readflow/pkg/logger"
	"github.com/rs/zerolog"
	"github.com/tdewolff/minify/v2"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// Sanitizer is a HTML sanitizer
type Sanitizer struct {
	logger    zerolog.Logger
	blockList *BlockList
	policy    *bluemonday.Policy
	minifier  *minify.M
}

// NewSanitizer create new HTML sanitizer
func NewSanitizer(blockList *BlockList) *Sanitizer {
	logger := logger.With().Str("component", "sanitizer").Logger()
	policy := bluemonday.UGCPolicy()
	policy.AddTargetBlankToFullyQualifiedLinks(true)
	policy.AllowAttrs("width", "height", "src", "allowfullscreen", "sandbox").OnElements("iframe")
	policy.AllowAttrs("srcset", "sizes").OnElements("img", "source")
	policy.AllowElements("picture", "source")

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
		minifier:  newDefaultMinifier(),
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
					case "src", "href", "srcset":
						// get raw URL (split if srcset)
						value := strings.Split(a.Val, " ")[0]
						u, err := url.Parse(value)
						if err != nil {
							attrs = append(attrs, a)
							continue
						}
						if u.Scheme != "data" {
							value = u.Hostname()
						}
						if !s.blockList.Contains(value) {
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

	min, err := s.minifier.String("text/html", content)
	if err == nil {
		content = min
	}

	return s.policy.Sanitize(content)
}
