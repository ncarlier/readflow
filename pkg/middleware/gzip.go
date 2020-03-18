package middleware

import (
	"net/http"

	"github.com/NYTimes/gziphandler"
)

// Gzip is a middleware to enabling GZIP on HTTP requests
func Gzip(inner http.Handler) http.Handler {
	return gziphandler.GzipHandler(inner)
}
