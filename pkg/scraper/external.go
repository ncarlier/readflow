package scraper

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/ncarlier/readflow/pkg/logger"
	"github.com/rs/zerolog"
)

type extrenalWebScraper struct {
	uri        string
	httpClient *http.Client
	userAgent  string
	logger     zerolog.Logger
}

// NewExternalWebScraper create an external web scrapping service
func NewExternalWebScraper(httpClient *http.Client, userAgent, uri string) (WebScraper, error) {
	if _, err := url.ParseRequestURI(uri); err != nil {
		return nil, fmt.Errorf("invalid Web Scraping service URI: %s", uri)
	}
	log := logger.With().Str("component", "webscraper").Str("uri", uri).Logger()
	log.Debug().Msg("using external service")

	return &extrenalWebScraper{
		uri:        uri,
		userAgent:  userAgent,
		httpClient: httpClient,
		logger:     log,
	}, nil
}

func (ws extrenalWebScraper) Scrap(ctx context.Context, rawurl string) (*WebPage, error) {
	webPage, err := ws.scrap(ctx, rawurl)
	if err != nil {
		ws.logger.Error().Err(err).Msg("unable to scrap web page with external service, fallback on internal service")
		return NewInternalWebScraper(ws.httpClient, ws.userAgent).Scrap(ctx, rawurl)
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
