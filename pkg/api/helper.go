package api

import (
	"encoding/base64"
	"regexp"
)


var proxyPathRe = regexp.MustCompile(`^/([^/]+)/([^/]+)/(.+)`)

// Decode image URL from Image Proxy Path
func decodeImageProxyPath(path string) (url string, err error) {
	parts := proxyPathRe.FindStringSubmatch(path)
	if len(parts) != 4 {
		return
	}
	encoded := parts[3]
	var decoded []byte
	decoded, err = base64.StdEncoding.DecodeString(encoded)
	if err == nil {
		url = string(decoded)
	}

	return
}
