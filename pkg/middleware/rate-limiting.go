package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/ncarlier/readflow/pkg/config"
	ratelimiter "github.com/ncarlier/readflow/pkg/rate-limiter"
	"github.com/sethvargo/go-limiter/httplimit"

	"github.com/rs/zerolog/log"
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
func RateLimiting(name string, conf config.RateLimiting) Middleware {
	store, err := ratelimiter.NewRateLimiter(name, conf)
	if err != nil {
		log.Fatal().Err(err).Msg(unabbleToCreateMiddleware)
	}

	middleware, err := httplimit.NewMiddleware(store, customKeyFunc())
	if err != nil {
		log.Fatal().Err(err).Msg(unabbleToCreateMiddleware)
	}
	return middleware.Handle
}
