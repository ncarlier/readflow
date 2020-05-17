package api

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ncarlier/readflow/pkg/config"
	"github.com/ncarlier/readflow/pkg/middleware"
)

func nextRequestID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func buildOriginFromPublicURL(publicURL string) string {
	if publicURL == "" {
		return "*"
	}
	u, err := url.Parse(publicURL)
	if err != nil {
		return "*"
	}
	return fmt.Sprintf("%s://%s", u.Scheme, strings.TrimPrefix(u.Host, "api."))
}

// NewRouter creates router with declared routes
func NewRouter(conf *config.Config) *http.ServeMux {
	commonMiddlewares := []middleware.Middleware{
		middleware.Gzip,
		middleware.Logger,
		middleware.Tracing(nextRequestID),
	}
	router := http.NewServeMux()
	for _, route := range routes(conf) {
		handler := route.Handler
		for _, mw := range route.Middlewares {
			handler = mw(handler)
		}
		for _, mw := range commonMiddlewares {
			handler = mw(handler)
		}
		router.Handle(route.Path, handler)
	}

	return router
}
