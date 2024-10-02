package htmlrewriter

import (
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func RewriteDataSrcToSrcAttribute(node *html.Node) {
	if node.DataAtom != atom.Img && node.DataAtom != atom.Source {
		return
	}

	keys := []string{"src", "srcset"}
	for _, key := range keys {
		realAttr := findHTMLAttribute(node, key)
		if realAttr == nil || realAttr.Val == "" {
			dataAttr := findHTMLAttribute(node, "data-"+key)
			if dataAttr != nil && dataAttr.Val != "" {
				dataAttr.Key = key
				replaceHTMLAttribute(node, *dataAttr)
			}
		}
	}
}
