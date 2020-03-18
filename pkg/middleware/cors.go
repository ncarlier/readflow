package middleware

import (
	"net/http"
)

// Cors is a middleware to enabling CORS on HTTP requests
func Cors(allowOrigin string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", allowOrigin)
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization")

			if r.Method != "OPTIONS" {
				next.ServeHTTP(w, r)
			}
			return
		})
	}
}
