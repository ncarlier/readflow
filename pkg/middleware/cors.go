package middleware

import (
	"github.com/rs/cors"
)

// Cors is a middleware to enabling CORS on HTTP requests
func Cors(allowOrigin string) Middleware {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{allowOrigin},
		AllowCredentials: true,
	})
	return c.Handler
}
