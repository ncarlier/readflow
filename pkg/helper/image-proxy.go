package helper

import (
	"encoding/base64"
	"strings"
)

// Encode image URL to Image Proxy path
func EncodeImageProxyPath(url, size string) string {
	return "/resize:fit:" + size + "/" + base64.StdEncoding.EncodeToString([]byte(url))
}

// Decode image URL from Image Proxy Path
func DecodeImageProxyPath(path string) (url string, err error) {
	_, encoded, ok := strings.Cut(path[1:], "/")
	if !ok {
		return
	}
	var decoded []byte
	decoded, err = base64.StdEncoding.DecodeString(encoded)
	if err == nil {
		url = string(decoded)
	}

	return
}
