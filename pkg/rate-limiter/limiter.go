package ratelimiter

import (
	"github.com/ncarlier/readflow/pkg/config"
	"github.com/rs/zerolog/log"
	"github.com/sethvargo/go-limiter"
	"github.com/sethvargo/go-limiter/memorystore"
	"github.com/sethvargo/go-limiter/noopstore"
)

// RateLimiter is an interface use to apply rate limiting
type RateLimiter limiter.Store

// NewWebScraper create new Web Scraping service
func NewRateLimiter(name string, conf config.RateLimiting) (RateLimiter, error) {
	switch conf.Provider {
	case "memory":
		store, err := memorystore.New(&memorystore.Config{
			Tokens:   uint64(conf.Tokens),
			Interval: conf.Interval.Duration,
		})
		if err != nil {
			return nil, err
		}
		log.Info().Str("name", name).Str("component", "rate-limiter").Msg("using in memory rate limiter")
		return store, nil
	default:
		store, err := noopstore.New()
		if err != nil {
			return nil, err
		}
		return store, nil
	}
}
