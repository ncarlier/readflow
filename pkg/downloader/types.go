package downloader

import (
	"context"

	"github.com/ncarlier/readflow/pkg/cache"
	"github.com/ncarlier/readflow/pkg/constant"
)

// Downloader is a service used to download assets.
type Downloader interface {
	Download(ctx context.Context, url string) (*WebAsset, error)
}

// NewDefaultDownloader create new downloader with defaults
func NewDefaultDownloader(downloadCache cache.Cache) Downloader {
	return NewInternalDownloader(constant.DefaultClient, downloadCache, defaultMaxConcurentDownload)
}
