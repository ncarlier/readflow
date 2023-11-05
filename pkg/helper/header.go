package helper

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

// AddXForwardHeader add X-Forwarded-For header
func AddXForwardHeader(header *http.Header, host string) {
	if prior := header.Values("X-Forwarded-For"); len(prior) > 0 {
		host = strings.Join(prior, ", ") + ", " + host
	}
	header.Set("X-Forwarded-For", host)
}

// AddCacheHeader add cache headers
func AddCacheHeader(header *http.Header, maxAge int) {
	expires := time.Now().Add(time.Duration(maxAge) * time.Second)
	header.Set("Pragma", "public")
	header.Set("Cache-Control", fmt.Sprintf("max-age=%d, immutable", maxAge))
	header.Set("Expires", expires.Local().String())
}
