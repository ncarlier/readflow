package middleware

import (
	"net/http"
	"net/url"
)

// Cors is a middleware to enabling CORS on HTTP requests
func Cors(allowOrigin string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			u, err := url.Parse(origin)
			if !(err == nil && u.Hostname() == "localhost") {
				// Use configured origin for non localhost requests
				origin = allowOrigin
			}

			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization")
			w.Header().Set("Access-Control-Expose-Headers", "*")

			if r.Method != "OPTIONS" {
				next.ServeHTTP(w, r)
			}
		})
	}
}
