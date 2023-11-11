package api

import (
	"net"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/ncarlier/readflow/pkg/cache"
	"github.com/ncarlier/readflow/pkg/config"
	"github.com/ncarlier/readflow/pkg/constant"
	"github.com/ncarlier/readflow/pkg/downloader"
	"github.com/ncarlier/readflow/pkg/helper"
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
		log.Fatal().Err(err).Msg("unable to setup Image Proxy cache")
	}
	down := downloader.NewInternalDownloader(constant.DefaultClient, c, 0)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		img := strings.TrimPrefix(r.URL.Path, "/img")
		q := r.URL.Query()
		// Redirect if image proxy service not configured or using old UI
		if q.Has("url") && q.Has("size") {
			img = q.Get("url")
			// legacy UI, redirect
			http.Redirect(w, r, img, http.StatusMovedPermanently)
			return
		}

		if host, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
			helper.AddXForwardHeader(&r.Header, host)
		}
		asset, resp, err := down.Get(r.Context(), conf.Image.ProxyURL+img, &r.Header)
		if err != nil {
			// Redirect if image proxy failed
			if decoded, err := decodeImageProxyPath(img); err != nil {
				http.Error(w, err.Error(), http.StatusBadGateway)
			} else {
				http.Redirect(w, r, strings.Replace(decoded, "http://", "https://", 1), http.StatusTemporaryRedirect)
			}
			return
		}

		header := http.Header{}
		if resp != nil {
			header = resp.Header
		}

		// Write response
		w.WriteHeader(http.StatusOK)
		helper.AddCacheHeader(&header, constant.CacheMaxAge)
		asset.Write(w, header)
	})
}
