package middleware

import (
	"github.com/rs/cors"
)

// Cors is a middleware to enabling CORS on HTTP requests
func Cors(allowOrigin string) Middleware {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{allowOrigin},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "Authorization"},
		ExposedHeaders:   []string{"Content-Disposition", "X-Content-Length", "X-Request-Id"},
		AllowCredentials: true,
	})
	return c.Handler
}
