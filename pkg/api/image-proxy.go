
package api

import (
	"net/http"
	"strings"
	"net"
	"io"

	"github.com/ncarlier/readflow/pkg/config"
)

const userAgent = "Mozilla/5.0 (X11; Linux x86_64; rv:10.0) Gecko/20100101 Firefox/10.0"

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
			dst.Add(k, v)
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
		q := r.URL.Query()

		// Extract and validate url parameter
		img := q.Get("url")
		if img == "" {
			http.Error(w, "bad parameter", http.StatusBadRequest)
			return
		}

		// Redirect if image proxy service not configured
		if conf.ImageProxy == nil || *conf.ImageProxy == "" {
			http.Redirect(w, r, img, 301)
			return
		}
		
		// Build image proxy client
		client := &http.Client{}
		req, err := http.NewRequest("GET", *conf.ImageProxy + "/resize?" + q.Encode(), nil)
        if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
        }

		// Manage request headers
        req.Header.Set("User-Agent", userAgent)
		if clientIP, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
			appendHostToXForwardHeader(req.Header, clientIP)
		}
		delHopHeaders(r.Header)

		// Do proxy request
        resp, err := client.Do(req)
        if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
        }
        defer resp.Body.Close()

		// Create proxy response
		delHopHeaders(resp.Header)
		copyHeader(w.Header(), resp.Header)
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	})
}
