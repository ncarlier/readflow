package readability

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

// Metas is the set of meta tags
type Metas map[string]*Meta

// GetContent get first content form keys
func (m Metas) GetContent(keys ...string) *string {
	for _, key := range keys {
		if m[key] != nil {
			return &m[key].Content
		}
	}
	return nil
}

// ExtractMetas extracts meta tags from a HTML document.
func ExtractMetas(doc io.Reader) (Metas, error) {
	var buf bytes.Buffer
	tee := io.TeeReader(doc, &buf)

	metas := make(map[string]*Meta)
	z := html.NewTokenizer(tee)
	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			if z.Err() == io.EOF {
				return metas, nil
			}
			return nil, z.Err()
		}

		t := z.Token()

		if t.DataAtom == atom.Head && t.Type == html.EndTagToken {
			return metas, nil
		}

		if t.DataAtom == atom.Meta {
			meta := Meta{}
			for _, a := range t.Attr {
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
			metas[key] = &meta
		}
	}
}
