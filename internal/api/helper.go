package api

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"
)

var proxyPathRe = regexp.MustCompile(`^/([^/]+)/([^/]+)/(.+)`)

// Decode image URL from Image Proxy Path
func decodeImageProxyPath(path string) (signature, options, url string, err error) {
	parts := proxyPathRe.FindStringSubmatch(path)
	if len(parts) != 4 {
		err = errors.New("invalid image proxy path")
		return
	}
	signature = parts[1]
	options = parts[2]
	encoded := parts[3]
	var decoded []byte
	decoded, err = base64.StdEncoding.DecodeString(encoded)
	if err == nil {
		url = string(decoded)
	}

	return
}

// addXForwardHeader add X-Forwarded-For header
func addXForwardHeader(header *http.Header, host string) {
	if prior := header.Values("X-Forwarded-For"); len(prior) > 0 {
		host = strings.Join(prior, ", ") + ", " + host
	}
	header.Set("X-Forwarded-For", host)
}

// addCacheHeader add cache headers
func addCacheHeader(header *http.Header, maxAge int) {
	expires := time.Now().Add(time.Duration(maxAge) * time.Second)
	header.Set("Pragma", "public")
	header.Set("Cache-Control", fmt.Sprintf("max-age=%d, immutable", maxAge))
	header.Set("Expires", expires.Local().String())
}
