package oembed

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/ncarlier/readflow/pkg/defaults"
	"github.com/ncarlier/readflow/pkg/scraper"
)

const queryParams = "?maxheight=600&maxwidth=800&format=json&url="

type oEmbedResponse struct {
	Title        string `json:"title,omitempty"`
	Type         string `json:"type,omitempty"`
	HTML         string `json:"html,omitempty"`
	ThumbnailURL string `json:"thumbnail_url,omitempty"`
	AuthorName   string `json:"author_name,omitempty"`
	ProviderName string `json:"provider_name,omitempty"`
}

type oEmbedContentProvider struct {
	name       string
	endpoint   string
	re         *regexp.Regexp
	httpClient *http.Client
}

func (cp oEmbedContentProvider) Get(ctx context.Context, rawurl string) (*scraper.WebPage, error) {
	oembedURL := cp.endpoint + queryParams + rawurl

	req, err := http.NewRequest("GET", oembedURL, http.NoBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", defaults.UserAgent)
	req = req.WithContext(ctx)
	res, err := cp.httpClient.Do(req)
	if err != nil || res.StatusCode >= 300 {
		if err == nil {
			err = fmt.Errorf("bad status code: %d", res.StatusCode)
		}
		return nil, err
	}
	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	payload := oEmbedResponse{}
	err = json.Unmarshal(body, &payload)
	if err != nil {
		return nil, err
	}

	return &scraper.WebPage{
		Title:    payload.Title,
		HTML:     payload.HTML,
		Image:    payload.ThumbnailURL,
		URL:      rawurl,
		Text:     fmt.Sprintf("%s %s from %s", payload.ProviderName, payload.Type, payload.AuthorName),
		SiteName: payload.ProviderName,
	}, nil
}

func (cp oEmbedContentProvider) Match(rawurl string) bool {
	return cp.re.MatchString(rawurl)
}

func init() {
	for _, provider := range Providers {
		for _, endpoint := range provider.Endpoints {
			for _, scheme := range endpoint.Schemes {
				re, err := Scheme2Regexp(scheme)
				if err != nil {
					log.Fatal(err)
				}
				scraper.ContentProviders = append(scraper.ContentProviders, &oEmbedContentProvider{
					name:       provider.Name,
					endpoint:   strings.ReplaceAll(endpoint.URL, "{format}", "json"),
					re:         re,
					httpClient: http.DefaultClient,
				})
			}
		}
	}
}
