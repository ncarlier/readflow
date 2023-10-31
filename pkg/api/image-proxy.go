package api

import (
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/ncarlier/readflow/pkg/config"
	"github.com/ncarlier/readflow/pkg/constant"
	"github.com/ncarlier/readflow/pkg/helper"
)

var hopHeaders = []string{
	"Connection",
	"Keep-Alive",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"Te", // canonicalized version of "TE"
	"Trailers",
	"Transfer-Encoding",
	"Upgrade",
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			if dst.Get(k) != "" {
				dst.Set(k, v)
			} else {
				dst.Add(k, v)
			}
		}
	}
}

func delHopHeaders(header http.Header) {
	for _, h := range hopHeaders {
		header.Del(h)
	}
}

func appendHostToXForwardHeader(header http.Header, host string) {
	// If we aren't the first proxy retain prior
	// X-Forwarded-For information as a comma+space
	// separated list and fold multiple headers into one.
	if prior, ok := header["X-Forwarded-For"]; ok {
		host = strings.Join(prior, ", ") + ", " + host
	}
	header.Set("X-Forwarded-For", host)
}

// imgProxyHandler is the handler for proxying images.
func imgProxyHandler(conf *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/img")
		q := r.URL.Query()
		// Redirect if image proxy service not configured or using old UI
		if q.Has("url") && q.Has("size") {
			img := q.Get("url")
			// legacy UI, redirect
			http.Redirect(w, r, img, http.StatusMovedPermanently)
			return
		}
		if conf.Image.ProxyURL == "" {
			http.Error(w, http.StatusText(http.StatusNoContent), http.StatusNotFound)
			return
		}

		// Build image proxy client
		req, err := http.NewRequest("GET", conf.Image.ProxyURL+path, http.NoBody)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}

		// Manage request headers: copy, add x-forward, del hop
		copyHeader(req.Header, r.Header)
		if clientIP, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
			appendHostToXForwardHeader(req.Header, clientIP)
		}
		delHopHeaders(req.Header)

		// Do proxy request
		resp, err := constant.DefaultClient.Do(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		// Redirect if image proxy failed
		if resp.StatusCode >= 400 {
			// decode image URL from proxy path
			if img, err := helper.DecodeImageProxyPath(path); err != nil {
				http.Error(w, err.Error(), http.StatusBadGateway)
			} else {
				http.Redirect(w, r, string(img), http.StatusTemporaryRedirect)
			}
			return
		}

		// Create proxy response
		delHopHeaders(resp.Header)
		copyHeader(w.Header(), resp.Header)
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	})
}
