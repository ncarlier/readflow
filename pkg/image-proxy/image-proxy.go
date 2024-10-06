package imageproxy

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"regexp"
	"strings"
)

var proxyPathRe = regexp.MustCompile(`^/([^/]+)/([^/]+)/(.+)`)

type ImageProxyHashSet struct {
	Size string
	Hash string
}

type ImageProxyConfiguration struct {
	URL        string
	Sizes      string
	SecretKey  []byte
	SecretSalt []byte
}

type ImageProxy struct {
	url        string
	sizes      []string
	secretKey  []byte
	secretSalt []byte
}

func NewImageProxy(config *ImageProxyConfiguration) *ImageProxy {
	return &ImageProxy{
		url:        config.URL,
		sizes:      strings.Split(config.Sizes, ","),
		secretKey:  config.SecretKey,
		secretSalt: config.SecretSalt,
	}
}

// Decode image URL from Image Proxy Path
func Decode(path string) (signature, options, url string, err error) {
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

func (ip *ImageProxy) URL() string {
	return ip.url
}

// Encode image URL to Image Proxy path
func (ip *ImageProxy) Encode(url, size string) string {
	if size == "" {
		// use last size by default (aka: highest resolution)
		size = ip.sizes[len(ip.sizes)-1]
	}
	return ip.URL() + "/" + ip.getHash(url, size) + "/resize:fit:" + size + "/" + base64.StdEncoding.EncodeToString([]byte(url))
}

// getHash from image URL and size
func (ip *ImageProxy) getHash(url, size string) string {
	path := "/resize:fit:" + size + "/" + base64.StdEncoding.EncodeToString([]byte(url))
	mac := hmac.New(sha256.New, ip.secretKey)
	mac.Write(ip.secretSalt)
	mac.Write([]byte(path))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

// GetHashSet return image proxy URLs hashset
func (ip *ImageProxy) GetHashSet(url string) *[]ImageProxyHashSet {
	result := make([]ImageProxyHashSet, len(ip.sizes))
	for i, size := range ip.sizes {
		result[i] = ImageProxyHashSet{
			Size: size,
			Hash: ip.getHash(url, size),
		}
	}
	return &result
}
