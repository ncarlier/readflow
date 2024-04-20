package downloader

import (
	"context"
	"net/http"

	"github.com/ncarlier/readflow/pkg/cache"
	"github.com/ncarlier/readflow/pkg/defaults"
)

// Downloader is a service used to download assets.
type Downloader interface {
	Get(ctx context.Context, url string, header *http.Header) (*WebAsset, *http.Response, error)
}

// NewDefaultDownloader create new downloader with defaults
func NewDefaultDownloader(downloadCache cache.Cache) Downloader {
	return NewInternalDownloader(defaults.HTTPClient, defaults.UserAgent, downloadCache, defaultMaxConcurentDownload)
}
