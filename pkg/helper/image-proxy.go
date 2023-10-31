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
	url = path[strings.LastIndex(path, "/")+1:]
	var decoded []byte
	decoded, err = base64.StdEncoding.DecodeString(url)
	if err == nil {
		url = string(decoded)
	}

	return
}
