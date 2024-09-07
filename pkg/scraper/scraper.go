package scraper

import (
	"context"
	"net/http"
)

// WebPage is the result of a web scraping
type WebPage struct {
	URL      string `json:"url,omitempty"`
	Title    string `json:"title,omitempty"`
	HTML     string `json:"html,omitempty"`
	Text     string `json:"text,omitempty"`
	Length   int    `json:"length,omitempty"`
	Excerpt  string `json:"excerpt,omitempty"`
	SiteName string `json:"sitename,omitempty"`
	Image    string `json:"image,omitempty"`
	Favicon  string `json:"favicon,omitempty"`
}

// WebScraper is an interface with Web Scrapping provider
type WebScraper interface {
	Scrap(ctx context.Context, rawurl string) (*WebPage, error)
}

// WebScraperConfiguration to configure a Web scraper
type WebScraperConfiguration struct {
	HttpClient              *http.Client
	UserAgent               string
	ExternalServiceEndpoint string
}

// NewWebScraper create new Web Scraping service
func NewWebScraper(conf *WebScraperConfiguration) (WebScraper, error) {
	if conf.ExternalServiceEndpoint == "" {
		return NewInternalWebScraper(conf), nil
	}
	return NewExternalWebScraper(conf)
}
