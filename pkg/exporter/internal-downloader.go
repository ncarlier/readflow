package exporter

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	nurl "net/url"
	"strings"
	"time"

	"golang.org/x/sync/semaphore"

	"github.com/ncarlier/readflow/pkg/cache"
	"github.com/ncarlier/readflow/pkg/constant"
	"github.com/ncarlier/readflow/pkg/helper"
	"github.com/ncarlier/readflow/pkg/model"
)

var errInvalidURL = errors.New("invalid URL")

type InternalDownloader struct {
	cache                 cache.Cache
	requestTimeout        time.Duration
	maxConcurrentDownload int64
	httpClient            *http.Client
	dlSemaphore           *semaphore.Weighted
}

func NewInternalDownloader(_cache cache.Cache, maxConcurrentDownload int64, requestTimeout time.Duration) *InternalDownloader {
	if maxConcurrentDownload <= 0 {
		maxConcurrentDownload = 10
	}
	httpClient := &http.Client{
		Timeout: requestTimeout,
	}
	return &InternalDownloader{
		cache:                 _cache,
		httpClient:            httpClient,
		requestTimeout:        requestTimeout,
		maxConcurrentDownload: maxConcurrentDownload,
		dlSemaphore:           semaphore.NewWeighted(maxConcurrentDownload),
	}
}

func (d *InternalDownloader) Download(ctx context.Context, url string) (*model.FileAsset, error) {
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
	asset, _ := d.cache.Get(hurl)
	if asset != nil {
		return asset, nil
	}

	// Download the asset, use semaphore to limit concurrent downloads
	err = d.dlSemaphore.Acquire(ctx, 1)
	if err != nil {
		return nil, err
	}
	resp, err := d.download(url)
	d.dlSemaphore.Release(1)
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
	bodyContent, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Put asset into the cache
	asset = &model.FileAsset{
		Data:        bodyContent,
		ContentType: contentType,
		Name:        url,
	}
	if err := d.cache.Put(hurl, asset); err != nil {
		return nil, err
	}

	return asset, nil
}

func (d *InternalDownloader) download(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", constant.UserAgent)
	// req.Header.Set("Referer", parentURL)

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
