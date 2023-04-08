package middleware

import (
	"net/http"

	"github.com/klauspost/compress/gzhttp"
)

// Gzip is a middleware to enabling GZIP on HTTP requests
func Gzip(inner http.Handler) http.Handler {
	return gzhttp.GzipHandler(inner)
}
