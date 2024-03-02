package ratelimiter

import (
	"github.com/ncarlier/readflow/pkg/logger"
	"github.com/ncarlier/readflow/pkg/types"
	"github.com/sethvargo/go-limiter"
	"github.com/sethvargo/go-limiter/memorystore"
	"github.com/sethvargo/go-limiter/noopstore"
)

// RateLimiterConfig for rate-limiter configuration section
type RateLimiterConfig struct {
	// Provider of the rate limiting store
	Provider string `toml:"provider"`
	// Tokens allowed per interval
	Tokens int `toml:"tokens"`
	// Interval until tokens reset
	Interval types.Duration `toml:"interval"`
}

// RateLimiter is an interface use to apply rate limiting
type RateLimiter limiter.Store

// NewWebScraper create new Web Scraping service
func NewRateLimiter(name string, conf *RateLimiterConfig) (RateLimiter, error) {
	switch conf.Provider {
	case "memory":
		store, err := memorystore.New(&memorystore.Config{
			Tokens:   uint64(conf.Tokens),
			Interval: conf.Interval.Duration,
		})
		if err != nil {
			return nil, err
		}
		logger.Info().Str("name", name).Str("component", "rate-limiter").Msg("using in memory rate limiter")
		return store, nil
	default:
		store, err := noopstore.New()
		if err != nil {
			return nil, err
		}
		return store, nil
	}
}
