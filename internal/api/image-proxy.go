package api

import (
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/ncarlier/readflow/internal/config"
	"github.com/ncarlier/readflow/internal/service"
	"github.com/ncarlier/readflow/pkg/defaults"
	imageproxy "github.com/ncarlier/readflow/pkg/image-proxy"
	"github.com/ncarlier/readflow/pkg/logger"
)

// imgProxyHandler is the handler for proxying images.
func imgProxyHandler(conf *config.Config) http.Handler {
	if conf.ImageProxy.URL == "" {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, http.StatusText(http.StatusNoContent), http.StatusNotFound)
		})
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		img := strings.TrimPrefix(r.URL.Path, "/img")
		_, opts, src, err := imageproxy.Decode(img)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		logger := logger.With().Str("src", src).Str("opts", opts).Logger()

		if host, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
			addXForwardHeader(&r.Header, host)
		}
		logger.Debug().Msg("getting image via proxy")
		asset, resp, err := service.Lookup().Download(r.Context(), conf.ImageProxy.URL+img, &r.Header)
		if err != nil {
			logger.Info().Err(err).Dur("took", time.Since(start)).Msg("unable to get image via proxy")
			// Redirect if image proxy failed
			http.Redirect(w, r, strings.Replace(src, "http://", "https://", 1), http.StatusTemporaryRedirect)
			return
		}

		header := http.Header{}
		if resp != nil {
			header = resp.Header
		}

		// Write response
		w.WriteHeader(http.StatusOK)
		addCacheHeader(&header, defaults.CacheMaxAge)
		asset.Write(w, header)
		logger.Info().Dur("took", time.Since(start)).Msg("got image via proxy")
	})
}
