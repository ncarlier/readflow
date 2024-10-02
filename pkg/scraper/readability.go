package scraper

import (
	"io"
	"net/url"

	"github.com/go-shiori/dom"
	"github.com/go-shiori/go-readability"
	"github.com/ncarlier/readflow/pkg/html"
	htmlrewriter "github.com/ncarlier/readflow/pkg/scraper/html-rewriter"
)

const MAX_RESPONSE_SIZE = 2 << 20 // 2Mb

func ReadWebPage(body io.Reader, pageUrl *url.URL) (*WebPage, error) {
	// Set body limit
	body = io.LimitReader(body, MAX_RESPONSE_SIZE)
	// Parse DOM
	doc, err := dom.Parse(body)
	if err != nil {
		return nil, err
	}
	// Extract meta
	meta := html.ExtractMetaFromDOM(doc)

	// Rewrite HTML if needed
	htmlrewriter.Rewrite(doc, []htmlrewriter.HTMLRewriterFunc{
		htmlrewriter.RewriteDataSrcToSrcAttribute,
		htmlrewriter.RewritePictureWithoutImgSrcAttribute,
	})

	// Create article with Open Graph attributes
	result := &WebPage{
		Title: meta.GetContent("og:title", "twitter:title", "title"),
		Text:  meta.GetContent("og:description", "twitter:description", "description"),
		Image: meta.GetContent("og:image", "twitter:image"),
	}

	// Set canonical URL
	result.URL = pageUrl.String()

	// Extract content from the HTML page
	article, err := readability.FromDocument(doc, pageUrl)
	if err != nil {
		return result, err
	}

	// Complete result with extracted properties
	result.HTML = article.Content
	result.Favicon = article.Favicon
	result.Length = article.Length
	result.SiteName = article.SiteName
	// FIXME: readability excerpt don't well support UTF8
	// result.Excerpt = helper.ToUTF8(article.Excerpt)

	// Fill in empty Open Graph attributes
	if result.Title == "" {
		result.Title = article.Title
	}
	if result.Text == "" {
		result.Text = result.Excerpt
	}
	if result.Image == "" {
		result.Image = article.Image
	}

	return result, nil
}
