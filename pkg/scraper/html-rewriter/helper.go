package htmlrewriter

import "golang.org/x/net/html"

func newHTMLAttribute(key string) html.Attribute {
	return html.Attribute{
		Key: key,
	}
}

func findHTMLAttribute(node *html.Node, key string) *html.Attribute {
	for _, attr := range node.Attr {
		if attr.Key == key {
			return &attr
		}
	}
	return nil
}

func replaceHTMLAttribute(node *html.Node, attribute html.Attribute) {
	attrs := []html.Attribute{}
	for _, attr := range node.Attr {
		if attr.Key != attribute.Key {
			attrs = append(attrs, attr)
		}
	}
	attrs = append(attrs, attribute)
	node.Attr = attrs
}
