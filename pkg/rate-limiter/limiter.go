package ratelimiter

import (
	"context"

	"github.com/ncarlier/readflow/pkg/config"
	"github.com/rs/zerolog/log"
	"github.com/sethvargo/go-limiter/memorystore"
	"github.com/sethvargo/go-limiter/noopstore"
)

// RateLimiter is an interface use to apply rate limiting
type RateLimiter interface {
	Take(ctx context.Context, key string) (tokens, remaining, reset uint64, ok bool, err error)
	Close(ctx context.Context) error
}

// NewWebScraper create new Web Scraping service
func NewRateLimiter(conf config.RateLimiting) (RateLimiter, error) {
	switch conf.Provider {
	case "memory":
		store, err := memorystore.New(&memorystore.Config{
			Tokens:   uint64(conf.Tokens),
			Interval: conf.Interval.Duration,
		})
		if err != nil {
			return nil, err
		}
		log.Info().Str("component", "rate-limiter").Msg("using in memory rate limiter")
		return store, nil
	default:
		store, err := noopstore.New()
		if err != nil {
			return nil, err
		}
		return store, nil
	}
}
