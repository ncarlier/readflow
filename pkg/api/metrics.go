package api

import (
	"net/http"

	"github.com/ncarlier/readflow/pkg/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// NewRouter creates metrics router
func NewMetricsRouter() *http.ServeMux {
	router := http.NewServeMux()

	handler := promhttp.Handler()
	handler = middleware.Methods("GET")(handler)
	handler = middleware.Cors("*")(handler)
	handler = middleware.Logger(handler)
	router.Handle("/metrics", handler)

	return router
}
