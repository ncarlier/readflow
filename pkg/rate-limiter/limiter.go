package ratelimiter

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/sethvargo/go-limiter/memorystore"
)

// RateLimiter is an interface use to apply rate limiting
type RateLimiter interface {
	Take(ctx context.Context, key string) (tokens, remaining, reset uint64, ok bool, err error)
	Close(ctx context.Context) error
}

// NewWebScraper create new Web Scraping service
func NewRateLimiter(uri string) (RateLimiter, error) {
	if uri == "" {
		return nil, nil
	}
	u, err := url.ParseRequestURI(uri)
	if err != nil {
		return nil, fmt.Errorf("invalid configuration URI: %s", uri)
	}

	tokens := uint64(1)
	if val := u.Query().Get("tokens"); val != "" {
		tokens, _ = strconv.ParseUint(string(val), 10, 64)
	}
	interval := time.Hour
	if val := u.Query().Get("interval"); val != "" {
		interval, _ = time.ParseDuration(val)
	}

	switch u.Scheme {
	case "memory":
		store, err := memorystore.New(&memorystore.Config{
			// Number of tokens allowed per interval.
			Tokens: tokens,
			// Interval until tokens reset.
			Interval: interval,
		})
		if err != nil {
			return nil, err
		}
		log.Info().Str("component", "rate-limiter").Str("uri", u.Redacted()).Msg("using in memory rate limiter")
		return store, nil
	default:
		return nil, fmt.Errorf("unsupported rate limiter: %s", u.Scheme)
	}
}
