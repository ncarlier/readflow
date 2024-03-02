package ratelimiter

import (
	"errors"
	"net/http"
	"strings"

	"github.com/ncarlier/readflow/pkg/logger"
	"github.com/ncarlier/readflow/pkg/middleware"
	"github.com/sethvargo/go-limiter/httplimit"
)

const unabbleToCreateMiddleware = "unable to create rate-limiting middleware"

func customKeyFunc() httplimit.KeyFunc {
	return func(r *http.Request) (string, error) {
		s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
		if len(s) != 2 {
			return "", errors.New("rate-limiting key not found")
		}
		return s[1], nil
	}
}

// RateLimiting is a middleware to limite usage of an endpoint
func Middleware(name string, conf *RateLimiterConfig) middleware.Middleware {
	store, err := NewRateLimiter(name, conf)
	if err != nil {
		logger.Fatal().Err(err).Msg(unabbleToCreateMiddleware)
	}

	middleware, err := httplimit.NewMiddleware(store, customKeyFunc())
	if err != nil {
		logger.Fatal().Err(err).Msg(unabbleToCreateMiddleware)
	}
	return middleware.Handle
}
