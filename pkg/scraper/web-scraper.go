package scraper

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/ncarlier/readflow/pkg/defaults"
	"github.com/ncarlier/readflow/pkg/utils"
	"golang.org/x/net/html/charset"
)

type WebScraper struct {
	httpClient   *http.Client
	userAgent    string
	forwardProxy *ForwardProxyConfiguration
}

// NewWebScraper create an internal web scrapping service
func NewWebScraper(conf *WebScraperConfiguration) *WebScraper {
	return &WebScraper{
		httpClient:   utils.If(conf.HttpClient == nil, defaults.HTTPClient, conf.HttpClient),
		userAgent:    utils.If(conf.UserAgent == "", defaults.UserAgent, conf.UserAgent),
		forwardProxy: conf.ForwardProxy,
	}
}

func (ws WebScraper) Scrap(ctx context.Context, rawurl string) (*WebPage, error) {
	// Validate URL
	pageURL, err := url.ParseRequestURI(rawurl)
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
	res, err := ws.get(ctx, pageURL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := charset.NewReader(res.Body, contentType)
	if err != nil {
		return nil, err
	}

	return ReadWebPage(body, pageURL)
}

func (ws WebScraper) getContentType(ctx context.Context, rawurl string) (string, error) {
	req, err := http.NewRequest("HEAD", rawurl, http.NoBody)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", ws.userAgent)
	req = req.WithContext(ctx)
	res, err := ws.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	return res.Header.Get("Content-type"), nil
}

func (ws WebScraper) get(ctx context.Context, pageURL *url.URL) (*http.Response, error) {
	rawurl := pageURL.String()
	if ws.forwardProxy != nil && ws.forwardProxy.Endpoint != "" && ws.forwardProxy.Match(pageURL.Hostname()) {
		rawurl = strings.ReplaceAll(ws.forwardProxy.Endpoint, "{url}", rawurl)
	}
	req, err := http.NewRequest("GET", rawurl, http.NoBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", ws.userAgent)
	req = req.WithContext(ctx)
	return ws.httpClient.Do(req)
}
