package downloader

import (
	"net/http"
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

// delHopByHopheaders remove
func delHopByHopheaders(header *http.Header) {
	for _, h := range hopHeaders {
		header.Del(h)
	}
}

// mergeHeader merge HTTP headers
func mergeHeader(dst, src *http.Header) {
	if src == nil {
		return
	}
	for k, vv := range *src {
		for _, v := range vv {
			if dst.Get(k) != "" {
				dst.Set(k, v)
			} else {
				dst.Add(k, v)
			}
		}
	}
}
