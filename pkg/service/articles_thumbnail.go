package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"

	"github.com/ncarlier/readflow/pkg/model"
)

// encode image URL to Image Proxy path
func encodeImageProxyPath(url, size string) string {
	return "/resize:fit:" + size + "/" + base64.StdEncoding.EncodeToString([]byte(url))
}

// GetArticleThumbnail return article thumbnail URL
func (reg *Registry) GetArticleThumbnailHash(article *model.Article, size string) string {
	if article.Image == nil || *article.Image == "" {
		return ""
	}
	path := encodeImageProxyPath(*article.Image, size)

	mac := hmac.New(sha256.New, reg.conf.Hash.SecretKey.Value)
	mac.Write(reg.conf.Hash.SecretSalt.Value)
	mac.Write([]byte(path))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}
