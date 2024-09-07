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

type internalWebScraper struct {
	httpClient *http.Client
	userAgent  string
}

// NewInternalWebScraper create an internal web scrapping service
func NewInternalWebScraper(conf *WebScraperConfiguration) WebScraper {
	return &internalWebScraper{
		httpClient: utils.If(conf.HttpClient == nil, defaults.HTTPClient, conf.HttpClient),
		userAgent:  utils.If(conf.UserAgent == "", defaults.UserAgent, conf.UserAgent),
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

	return ReadWebPage(body, res.Request.URL)
}

func (ws internalWebScraper) getContentType(ctx context.Context, rawurl string) (string, error) {
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

func (ws internalWebScraper) get(ctx context.Context, rawurl string) (*http.Response, error) {
	req, err := http.NewRequest("GET", rawurl, http.NoBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", ws.userAgent)
	req = req.WithContext(ctx)
	return ws.httpClient.Do(req)
}
