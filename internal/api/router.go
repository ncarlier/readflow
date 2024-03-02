package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ncarlier/readflow/internal/config"
	"github.com/ncarlier/readflow/pkg/middleware"
)

func nextRequestID() string {
	now := time.Now().UnixNano()
	return strconv.FormatInt(now, 32)
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
