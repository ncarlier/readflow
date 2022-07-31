package html

import (
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// Meta is a HTML meta tag
type Meta struct {
	Name     string
	Property string
	Content  string
}

// MetaSet is the set of meta tags
type MetaSet map[string]*Meta

// GetContent get first content form keys
func (m MetaSet) GetContent(keys ...string) string {
	for _, key := range keys {
		if m[key] != nil {
			return m[key].Content
		}
	}
	return ""
}

// ExtractMetaFromDOM extracts meta tags from a HTML document.
func ExtractMetaFromDOM(doc *html.Node) MetaSet {
	metaSet := make(map[string]*Meta)

	var parser func(*html.Node)
	parser = func(node *html.Node) {
		if node.Type == html.ElementNode {
			// Don't process body
			if node.DataAtom == atom.Body {
				return
			}

			// Extract title
			if node.DataAtom == atom.Title && node.FirstChild != nil && node.FirstChild.Type == html.NodeType(html.TextToken) {
				meta := Meta{
					Name:    "title",
					Content: node.FirstChild.Data,
				}
				metaSet["title"] = &meta
			}

			// Extract meta
			if node.DataAtom == atom.Meta {
				meta := Meta{}
				for _, a := range node.Attr {
					switch a.Key {
					case "property":
						meta.Property = a.Val
					case "name":
						meta.Name = a.Val
					case "content":
						meta.Content = a.Val
					case "charset":
						meta.Name = "charset"
						meta.Content = a.Val
					}
				}
				key := meta.Name
				if meta.Property != "" {
					key = meta.Property
				}
				metaSet[key] = &meta
			}
		}

		for child := node.FirstChild; child != nil; child = child.NextSibling {
			parser(child)
		}
	}

	parser(doc)
	return metaSet
}
