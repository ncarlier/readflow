package scraper

import (
	"net/http"
	"strings"
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

// ForwardProxyConfiguration to configure forward proxy
type ForwardProxyConfiguration struct {
	Endpoint string
	Hosts    []string
}

// Match test if hostname is in the hosts list
func (fpc *ForwardProxyConfiguration) Match(hostname string) bool {
	for _, value := range fpc.Hosts {
		if strings.HasSuffix(hostname, value) {
			return true
		}
	}
	return false
}

// WebScraperConfiguration to configure a Web scraper
type WebScraperConfiguration struct {
	HttpClient   *http.Client
	UserAgent    string
	ForwardProxy *ForwardProxyConfiguration
}
