package scraper

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type extrenalWebScraper struct {
	uri        string
	httpClient *http.Client
	logger     zerolog.Logger
}

// NewExternalWebScraper create an external web scrapping service
func NewExternalWebScraper(httpClient *http.Client, uri string) (WebScraper, error) {
	if _, err := url.ParseRequestURI(uri); err != nil {
		return nil, fmt.Errorf("invalid Web Scraping service URI: %s", uri)
	}
	logger := log.With().Str("component", "webscraper").Str("uri", uri).Logger()
	logger.Debug().Msg("using external service")

	return &extrenalWebScraper{
		uri:        uri,
		httpClient: httpClient,
		logger:     logger,
	}, nil
}

func (ws extrenalWebScraper) Scrap(ctx context.Context, rawurl string) (*WebPage, error) {
	webPage, err := ws.scrap(ctx, rawurl)
	if err != nil {
		ws.logger.Error().Err(err).Msg("unable to scrap web page with external service, fallback on internal service")
		return NewInternalWebScraper(ws.httpClient).Scrap(ctx, rawurl)
	}
	return webPage, nil
}

func (ws extrenalWebScraper) scrap(ctx context.Context, rawurl string) (*WebPage, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", ws.uri, http.NoBody)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Add("u", rawurl)
	req.URL.RawQuery = q.Encode()

	ws.logger.Debug().Str("url", rawurl).Msg("scraping webpage")
	res, err := ws.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("invalid web scraping response: %d", res.StatusCode)
	}

	if ct := res.Header.Get("Content-Type"); ct != "" {
		if !strings.HasPrefix(ct, "application/json") {
			return nil, fmt.Errorf("invalid web scraping Content-Type response: %s", ct)
		}
	}

	webPage := WebPage{}
	err = json.NewDecoder(res.Body).Decode(&webPage)
	return &webPage, err
}
