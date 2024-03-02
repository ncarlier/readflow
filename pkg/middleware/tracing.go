package middleware

import (
	"context"
	"net/http"
)

const requestIDHeader = "X-Request-Id"

// Tracing is a middleware to trace HTTP request
func Tracing(nextRequestID func() string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(requestIDHeader)
			if requestID == "" {
				requestID = nextRequestID()
			}
			ctx := context.WithValue(r.Context(), ContextRequestID, requestID)
			w.Header().Set(requestIDHeader, requestID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
