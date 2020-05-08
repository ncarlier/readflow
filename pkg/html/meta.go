package html

import (
	"bytes"
	"io"

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

// ExtractMeta extracts meta tags from a HTML document.
func ExtractMeta(doc io.Reader) (MetaSet, error) {
	var buf bytes.Buffer
	tee := io.TeeReader(doc, &buf)

	metaSet := make(map[string]*Meta)
	tokenizer := html.NewTokenizer(tee)
	for {
		tokenType := tokenizer.Next()
		if tokenType == html.ErrorToken {
			if tokenizer.Err() == io.EOF {
				return metaSet, nil
			}
			return nil, tokenizer.Err()
		}

		token := tokenizer.Token()

		// Header end
		if token.DataAtom == atom.Head && token.Type == html.EndTagToken {
			return metaSet, nil
		}

		// Title tag
		if token.DataAtom == atom.Title && token.Type == html.StartTagToken {
			tokenType = tokenizer.Next()
			if tokenType == html.TextToken {
				meta := Meta{
					Name:    "title",
					Content: tokenizer.Token().Data,
				}
				metaSet["title"] = &meta
			}
			continue
		}

		// Meta tag
		if token.DataAtom == atom.Meta {
			meta := Meta{}
			for _, a := range token.Attr {
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
}
