package api

import (
	"encoding/base64"
	"errors"
	"regexp"
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
