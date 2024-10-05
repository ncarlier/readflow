package service

import (
	"github.com/ncarlier/readflow/internal/model"
	imageproxy "github.com/ncarlier/readflow/pkg/image-proxy"
)

// GetArticleThumbnailHashSet return article thumbnail hash set
func (reg *Registry) GetArticleThumbnailHashSet(article *model.Article) *[]imageproxy.ImageProxyHashSet {
	if reg.imageProxy.URL() == "" || article.Image == nil || *article.Image == "" {
		return nil
	}
	return reg.imageProxy.GetHashSet(*article.Image)
}
