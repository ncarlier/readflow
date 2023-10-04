package downloader

import (
	"context"
	"errors"
	"io"
	"net/http"
	nurl "net/url"
	"strings"

	"golang.org/x/sync/semaphore"

	"github.com/ncarlier/readflow/pkg/cache"
	"github.com/ncarlier/readflow/pkg/constant"
	"github.com/ncarlier/readflow/pkg/helper"
)

const (
	defaultMaxConcurentDownload = 10
)

var errInvalidURL = errors.New("invalid URL")

// InternalDownloader interface
type InternalDownloader struct {
	cache                 cache.Cache
	maxConcurrentDownload int64
	httpClient            *http.Client
	dlSemaphore           *semaphore.Weighted
}

// NewInternalDownloader create new downloader instance
func NewInternalDownloader(client *http.Client, downloadCache cache.Cache, maxConcurrentDownload int64) Downloader {
	if maxConcurrentDownload <= 0 {
		maxConcurrentDownload = defaultMaxConcurentDownload
	}
	return &InternalDownloader{
		cache:                 downloadCache,
		httpClient:            client,
		maxConcurrentDownload: maxConcurrentDownload,
		dlSemaphore:           semaphore.NewWeighted(maxConcurrentDownload),
	}
}

// Download web asset from its URL
func (dl *InternalDownloader) Download(ctx context.Context, url string) (*WebAsset, error) {
	// Ignore special URLs
	url = strings.TrimSpace(url)
	if url == "" || strings.HasPrefix(url, "data:") || strings.HasPrefix(url, "#") {
		return nil, errInvalidURL
	}
	// Validate URL
	parsedURL, err := nurl.ParseRequestURI(url)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Hostname() == "" {
		return nil, errInvalidURL
	}

	// Get the asset from cache
	hurl := helper.Hash(url)
	data, _ := dl.cache.Get(hurl)
	if data != nil {
		return NewWebAsset(data)
	}

	// Download the asset, use semaphore to limit concurrent downloads
	err = dl.dlSemaphore.Acquire(ctx, 1)
	if err != nil {
		return nil, err
	}
	resp, err := dl.download(url)
	dl.dlSemaphore.Release(1)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Get content type
	contentType := resp.Header.Get("Content-Type")
	contentType = strings.TrimSpace(contentType)
	if contentType == "" {
		contentType = "text/plain"
	}

	// Get response body
	bodyContent, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Put asset into the cache
	asset := &WebAsset{
		Data:        bodyContent,
		ContentType: contentType,
		Name:        url,
	}
	data, err = asset.Encode()
	if err != nil {
		return nil, err
	}
	if err := dl.cache.Put(hurl, data); err != nil {
		return nil, err
	}

	return asset, nil
}

func (dl *InternalDownloader) download(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, http.NoBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", constant.UserAgent)

	resp, err := dl.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
