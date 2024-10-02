package test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/go-shiori/dom"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/html"

	htmlrewriter "github.com/ncarlier/readflow/pkg/scraper/html-rewriter"
)

var testCases = []struct {
	input    string
	expected string
}{
	{
		input:    "<html><body><img data-src=\"foo\" /></body></html>",
		expected: "<html><head></head><body><img data-src=\"foo\" src=\"foo\"/></body></html>",
	},
	{
		input:    "<html><body><img data-src=\"foo\" src=\"bar\" /></body></html>",
		expected: "<html><head></head><body><img data-src=\"foo\" src=\"bar\"/></body></html>",
	},
	{
		input: `<picture>
			<source srcSet="https://miro.medium.com/v2/resize:fit:640/format:webp/1*a-eyadQGauD3NbTHtAzudw.png 640w" sizes="" type="image/webp"/>
			<img alt="" width="1000" height="512" role="presentation"/>
		</picture>`,
		expected: `<html><head></head><body><picture>
			<source srcset="https://miro.medium.com/v2/resize:fit:640/format:webp/1*a-eyadQGauD3NbTHtAzudw.png 640w" sizes="" type="image/webp"/>
			<img alt="" width="1000" height="512" role="presentation" src="https://miro.medium.com/v2/resize:fit:640/format:webp/1*a-eyadQGauD3NbTHtAzudw.png"/>
		</picture></body></html>`,
	},
}

func TestRewrite(t *testing.T) {
	for _, testCase := range testCases {
		doc, err := dom.Parse(strings.NewReader(testCase.input))
		require.Nil(t, err)
		htmlrewriter.Rewrite(doc, []htmlrewriter.HTMLRewriterFunc{
			htmlrewriter.RewriteDataSrcToSrcAttribute,
			htmlrewriter.RewritePictureWithoutImgSrcAttribute,
		})
		var b bytes.Buffer
		err = html.Render(&b, doc)
		require.Nil(t, err)
		require.Equal(t, testCase.expected, b.String())
	}
}
