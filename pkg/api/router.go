package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ncarlier/readflow/pkg/config"
	"github.com/ncarlier/readflow/pkg/middleware"
)

// NewRouter creates router with declared routes
func NewRouter(conf *config.Config) *http.ServeMux {
	router := http.NewServeMux()

	nextRequestID := func() string {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}

	authMiddleware := middleware.Auth(*conf.AuthN)

	for _, route := range routes {
		var handler http.Handler

		handler = route.HandlerFunc(conf)
		if route.AuthNRequired {
			handler = authMiddleware(handler)
		}
		handler = middleware.Method(route.Methods)(handler)
		handler = middleware.Cors(handler)
		handler = middleware.Logger(handler)
		handler = middleware.Tracing(nextRequestID)(handler)
		router.Handle(route.Path, handler)
	}

	return router
}
