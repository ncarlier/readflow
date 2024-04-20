package downloader

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	nurl "net/url"
	"strings"

	"golang.org/x/sync/semaphore"

	"github.com/ncarlier/readflow/pkg/cache"
	"github.com/ncarlier/readflow/pkg/utils"
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
	userAgent             string
	dlSemaphore           *semaphore.Weighted
}

// NewInternalDownloader create new downloader instance
func NewInternalDownloader(client *http.Client, userAgent string, downloadCache cache.Cache, maxConcurrentDownload int64) Downloader {
	if maxConcurrentDownload <= 0 {
		maxConcurrentDownload = defaultMaxConcurentDownload
	}
	return &InternalDownloader{
		cache:                 downloadCache,
		httpClient:            client,
		userAgent:             userAgent,
		maxConcurrentDownload: maxConcurrentDownload,
		dlSemaphore:           semaphore.NewWeighted(maxConcurrentDownload),
	}
}

// Download web asset from its URL
func (dl *InternalDownloader) Get(ctx context.Context, url string, header *http.Header) (*WebAsset, *http.Response, error) {
	// Ignore special URLs
	url = strings.TrimSpace(url)
	if url == "" || strings.HasPrefix(url, "data:") || strings.HasPrefix(url, "#") {
		return nil, nil, errInvalidURL
	}
	// Validate URL
	parsedURL, err := nurl.ParseRequestURI(url)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Hostname() == "" {
		return nil, nil, errInvalidURL
	}

	// Get the asset from cache
	hurl := utils.Hash(url)
	data, _ := dl.cache.Get(hurl)
	if data != nil {
		asset, err := NewWebAsset(data)
		return asset, nil, err
	}

	// Download the asset, use semaphore to limit concurrent downloads
	err = dl.dlSemaphore.Acquire(ctx, 1)
	if err != nil {
		return nil, nil, err
	}
	resp, err := dl.get(url, header)
	dl.dlSemaphore.Release(1)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, resp, fmt.Errorf("invalide HTTP response: %d", resp.StatusCode)
	}

	// Get content type
	contentType := resp.Header.Get("Content-Type")
	contentType = strings.TrimSpace(contentType)
	if contentType == "" {
		contentType = "text/plain"
	}

	// Get response body
	bodyContent, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp, err
	}

	// Put asset into the cache
	asset := &WebAsset{
		Data:        bodyContent,
		ContentType: contentType,
		Name:        url,
	}
	data, err = asset.Encode()
	if err != nil {
		return nil, resp, err
	}
	if err := dl.cache.Put(hurl, data); err != nil {
		return nil, resp, err
	}

	return asset, resp, nil
}

func (dl *InternalDownloader) get(url string, header *http.Header) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, http.NoBody)
	if err != nil {
		return nil, err
	}

	// Manage request headers: defaults, merge, del hop
	req.Header.Set("User-Agent", dl.userAgent)
	mergeHeader(&req.Header, header)
	delHopByHopheaders(&req.Header)

	resp, err := dl.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	// Clear hop headers
	delHopByHopheaders(&req.Header)

	return resp, nil
}
