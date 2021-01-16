package scraper

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	nurl "net/url"
	"strings"
	"time"

	read "github.com/go-shiori/go-readability"
	"github.com/ncarlier/readflow/pkg/constant"
	"github.com/ncarlier/readflow/pkg/helper"
	"github.com/ncarlier/readflow/pkg/html"
	"golang.org/x/net/html/charset"
)

type internalWebScraper struct {
	httpClient *http.Client
}

// NewInternalWebScraper create an internal web scrapping service
func NewInternalWebScraper() WebScraper {
	return &internalWebScraper{
		httpClient: &http.Client{
			Timeout: constant.DefaultTimeout,
		},
	}
}

func (ws internalWebScraper) Scrap(ctx context.Context, url string) (*WebPage, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	// Validate URL
	_, err := nurl.ParseRequestURI(url)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %v", err)
	}

	// Get URL content type
	contentType, err := ws.getContentType(ctx, url)
	if err != nil {
		return nil, err
	}

	if !strings.HasPrefix(contentType, "text/html") {
		return nil, fmt.Errorf("invalid content-type: %s", contentType)
	}

	// Get URL content
	res, err := ws.get(ctx, url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := charset.NewReader(res.Body, contentType)
	if err != nil {
		return nil, err
	}

	// Extract meta
	meta, err := html.ExtractMeta(body)
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

	var buffer bytes.Buffer
	tee := io.TeeReader(body, &buffer)

	// Test if the HTML page is readable by Shiori readability
	if !read.IsReadable(tee) {
		return result, fmt.Errorf("unable to extract content from HTML page")
	}

	// Extract content from the HTML page
	article, err := read.FromReader(&buffer, result.URL)
	if err != nil {
		return result, err
	}

	// Complete result with extracted properties
	result.HTML = article.Content
	result.Favicon = article.Favicon
	result.Length = article.Length
	result.SiteName = article.SiteName
	// FIXME: readability excerpt don't well support UTF8
	result.Excerpt = helper.ToUTF8(article.Excerpt)

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

func (ws internalWebScraper) getContentType(ctx context.Context, url string) (string, error) {
	req, err := http.NewRequest("HEAD", url, nil)
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

func (ws internalWebScraper) get(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", constant.UserAgent)
	req = req.WithContext(ctx)
	return ws.httpClient.Do(req)
}
