package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ncarlier/readflow/pkg/config"
	"github.com/ncarlier/readflow/pkg/middleware"
)

func nextRequestID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

var commonMiddlewares = []middleware.Middleware{
	middleware.Cors,
	middleware.Logger,
	middleware.Tracing(nextRequestID),
}

// NewRouter creates router with declared routes
func NewRouter(conf *config.Config) *http.ServeMux {
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
