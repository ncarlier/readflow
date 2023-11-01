package helper

import (
	"encoding/base64"
	"regexp"
)

var proxyPathRe = regexp.MustCompile(`^/([^/]+)/([^/]+)/(.+)`)

// Encode image URL to Image Proxy path
func EncodeImageProxyPath(url, size string) string {
	return "/resize:fit:" + size + "/" + base64.StdEncoding.EncodeToString([]byte(url))
}

// Decode image URL from Image Proxy Path
func DecodeImageProxyPath(path string) (url string, err error) {
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
