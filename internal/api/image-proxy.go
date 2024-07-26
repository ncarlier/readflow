package api

import (
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/ncarlier/readflow/internal/config"
	"github.com/ncarlier/readflow/pkg/cache"
	"github.com/ncarlier/readflow/pkg/defaults"
	"github.com/ncarlier/readflow/pkg/downloader"
	"github.com/ncarlier/readflow/pkg/logger"
)

// imgProxyHandler is the handler for proxying images.
func imgProxyHandler(conf *config.Config) http.Handler {
	if conf.Image.ProxyURL == "" {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, http.StatusText(http.StatusNoContent), http.StatusNotFound)
		})
	}
	c, err := cache.New(conf.Image.Cache)
	if err != nil {
		logger.Fatal().Err(err).Msg("unable to setup Image Proxy cache")
	}
	down := downloader.NewInternalDownloader(defaults.HTTPClient, defaults.UserAgent, c, 0, defaults.Timeout)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		img := strings.TrimPrefix(r.URL.Path, "/img")
		_, opts, src, err := decodeImageProxyPath(img)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		logger := logger.With().Str("src", src).Str("opts", opts).Logger()

		if host, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
			addXForwardHeader(&r.Header, host)
		}
		logger.Debug().Msg("getting image via proxy")
		asset, resp, err := down.Get(r.Context(), conf.Image.ProxyURL+img, &r.Header)
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
