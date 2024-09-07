package downloader

import (
	"context"
	"net/http"
)

// Downloader is a service used to download assets.
type Downloader interface {
	Get(ctx context.Context, url string, header *http.Header) (*WebAsset, *http.Response, error)
}
