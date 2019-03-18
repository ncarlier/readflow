package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ncarlier/reader/pkg/config"
	"github.com/ncarlier/reader/pkg/middleware"
)

// NewRouter creates router with declared routes
func NewRouter(conf *config.Config) *http.ServeMux {
	router := http.NewServeMux()

	nextRequestID := func() string {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}

	for _, route := range routes {
		var handler http.Handler

		handler = route.HandlerFunc(conf)
		if route.AuthNRequired {
			handler = middleware.MockAuth(handler)
		}
		handler = middleware.Method(handler, route.Methods)
		handler = middleware.Cors(handler)
		handler = middleware.Logger(handler)
		handler = middleware.Tracing(nextRequestID)(handler)
		router.Handle(route.Path, handler)
	}

	return router
}
