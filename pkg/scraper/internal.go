package scraper

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-shiori/dom"
	read "github.com/go-shiori/go-readability"
	"github.com/ncarlier/readflow/pkg/constant"
	"github.com/ncarlier/readflow/pkg/html"
	"golang.org/x/net/html/charset"
)

type internalWebScraper struct {
	httpClient *http.Client
}

// NewInternalWebScraper create an internal web scrapping service
func NewInternalWebScraper(httpClient *http.Client) WebScraper {
	return &internalWebScraper{
		httpClient: httpClient,
	}
}

func (ws internalWebScraper) Scrap(ctx context.Context, rawurl string) (*WebPage, error) {
	// Validate URL
	_, err := url.ParseRequestURI(rawurl)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %v", err)
	}

	// Get content provider for this URL
	contentProvider := GetContentProvider(rawurl)
	if contentProvider != nil {
		return contentProvider.Get(ctx, rawurl)
	}

	// Get URL content type
	contentType, err := ws.getContentType(ctx, rawurl)
	if err != nil {
		return nil, err
	}

	if !strings.HasPrefix(contentType, "text/html") {
		return nil, fmt.Errorf("invalid content-type: %s", contentType)
	}

	// Get URL content
	res, err := ws.get(ctx, rawurl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := charset.NewReader(res.Body, contentType)
	if err != nil {
		return nil, err
	}

	// Parse DOM
	doc, err := dom.Parse(body)
	if err != nil {
		return nil, err
	}

	// Extract meta
	meta := html.ExtractMetaFromDOM(doc)
	if err != nil {
		return nil, err
	}

	// Create article with Open Graph attributes
	result := &WebPage{
		Title: meta.GetContent("og:title", "twitter:title", "title"),
		Text:  meta.GetContent("og:description", "twitter:description", "description"),
		Image: meta.GetContent("og:image", "twitter:image"),
	}

	// Set canonical URL
	result.URL = res.Request.URL.String()

	// Extract content from the HTML page
	article, err := read.FromDocument(doc, res.Request.URL)
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

func (ws internalWebScraper) getContentType(ctx context.Context, rawurl string) (string, error) {
	req, err := http.NewRequest("HEAD", rawurl, http.NoBody)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", constant.UserAgent)
	req = req.WithContext(ctx)
	res, err := ws.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	return res.Header.Get("Content-type"), nil
}

func (ws internalWebScraper) get(ctx context.Context, rawurl string) (*http.Response, error) {
	req, err := http.NewRequest("GET", rawurl, http.NoBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", constant.UserAgent)
	req = req.WithContext(ctx)
	return ws.httpClient.Do(req)
}
