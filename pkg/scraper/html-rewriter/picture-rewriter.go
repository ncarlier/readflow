package htmlrewriter

import (
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func RewritePictureWithoutImgSrcAttribute(node *html.Node) {
	if node.DataAtom != atom.Picture {
		return
	}

	var imgNode, sourceNode *html.Node
	for n := node.FirstChild; n != nil; n = n.NextSibling {
		switch n.DataAtom {
		case atom.Source:
			sourceNode = n
		case atom.Img:
			imgNode = n
		}
	}
	if imgNode == nil || sourceNode == nil {
		return
	}

	srcAttr := findHTMLAttribute(imgNode, "src")
	if srcAttr != nil && srcAttr.Val != "" {
		return
	}
	srcsetAttr := findHTMLAttribute(sourceNode, "srcset")
	if srcsetAttr == nil || srcsetAttr.Val == "" {
		return
	}

	if srcAttr == nil {
		attr := newHTMLAttribute("src")
		srcAttr = &attr
	}
	srcAttr.Val = strings.Split(srcsetAttr.Val, " ")[0]
	replaceHTMLAttribute(imgNode, *srcAttr)
}
