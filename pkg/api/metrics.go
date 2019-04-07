package api

import (
	"net/http"

	"github.com/ncarlier/readflow/pkg/config"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func metrics(conf *config.Config) http.Handler {
	return promhttp.Handler()
}
