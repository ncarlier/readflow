package htmlrewriter

import (
	"golang.org/x/net/html"
)

type HTMLRewriterFunc func(doc *html.Node)

func Rewrite(doc *html.Node, rewriters []HTMLRewriterFunc) {
	var parser func(*html.Node)
	parser = func(node *html.Node) {
		if node.Type == html.ElementNode {
			for _, rewrite := range rewriters {
				rewrite(node)
			}
		}

		for child := node.FirstChild; child != nil; child = child.NextSibling {
			parser(child)
		}
	}

	parser(doc)
}
