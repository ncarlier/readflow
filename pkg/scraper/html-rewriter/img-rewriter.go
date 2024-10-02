package htmlrewriter

import (
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func RewriteDataSrcToSrcAttribute(node *html.Node) {
	if node.DataAtom != atom.Img {
		return
	}

	attrs := []html.Attribute{}
	srcAttr := newHTMLAttribute("src")
	dataSrcAttr := newHTMLAttribute("data-src")
	for _, attr := range node.Attr {
		switch attr.Key {
		case srcAttr.Key:
			srcAttr = attr
		case dataSrcAttr.Key:
			dataSrcAttr = attr
		default:
			attrs = append(attrs, attr)
		}
	}
	if dataSrcAttr.Val != "" && srcAttr.Val == "" {
		srcAttr.Val = dataSrcAttr.Val
		attrs = append(attrs, srcAttr)
		node.Attr = attrs
	}
}
