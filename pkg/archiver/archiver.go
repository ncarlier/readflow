package archiver

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	nurl "net/url"
	"strings"
	"time"

	"github.com/go-shiori/dom"
	"golang.org/x/net/html"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"

	"github.com/ncarlier/readflow/pkg/cache"
	"github.com/ncarlier/readflow/pkg/helper"
)

var errSkippedURL = errors.New("skip processing url")

type WebArchiver struct {
	cache                 cache.Cache
	requestTimeout        time.Duration
	maxConcurrentDownload int64
	httpClient            *http.Client
	dlSemaphore           *semaphore.Weighted
}

func NewWebArchiver(_cache cache.Cache, maxConcurrentDownload int64, requestTimeout time.Duration) *WebArchiver {
	if maxConcurrentDownload <= 0 {
		maxConcurrentDownload = 10
	}
	httpClient := &http.Client{
		Timeout: requestTimeout,
	}
	return &WebArchiver{
		cache:                 _cache,
		httpClient:            httpClient,
		requestTimeout:        requestTimeout,
		maxConcurrentDownload: maxConcurrentDownload,
		dlSemaphore:           semaphore.NewWeighted(maxConcurrentDownload),
	}
}

func (arc *WebArchiver) Archive(ctx context.Context, input io.Reader, baseURL string) ([]byte, error) {
	url, err := nurl.ParseRequestURI(baseURL)
	if err != nil || url.Scheme == "" || url.Hostname() == "" {
		return nil, fmt.Errorf("url \"%s\" is not valid", baseURL)
	}
	doc, err := html.Parse(input)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}
	nodes := make(map[*html.Node]struct{})
	for _, node := range dom.GetElementsByTagName(doc, "img") {
		nodes[node] = struct{}{}
	}
	g, ctx := errgroup.WithContext(ctx)
	for node := range nodes {
		node := node
		g.Go(func() error {
			return arc.processNode(ctx, node, url)
		})
	}

	// Wait until all nodes processed
	if err = g.Wait(); err != nil {
		return nil, err
	}

	var buffer bytes.Buffer
	err = html.Render(&buffer, doc)
	return buffer.Bytes(), err
}

func (arc *WebArchiver) processNode(ctx context.Context, node *html.Node, baseURL *nurl.URL) error {
	err := arc.processURLAttribute(ctx, node, "src", baseURL)
	if err != nil {
		return err
	}
	return nil
}

func (arc *WebArchiver) processURLAttribute(ctx context.Context, node *html.Node, attrName string, baseURL *nurl.URL) error {
	if !dom.HasAttribute(node, attrName) {
		return nil
	}

	url := dom.GetAttribute(node, attrName)
	content, contentType, err := arc.processURL(ctx, url, baseURL.String())
	if err != nil && err != errSkippedURL {
		return err
	}

	newURL := url
	if err == nil {
		newURL = createDataURL(content, contentType)
	}

	dom.SetAttribute(node, attrName, newURL)
	return nil
}

func (arc *WebArchiver) processURL(ctx context.Context, url string, parentURL string) ([]byte, string, error) {
	// Ignore special URLs
	url = strings.TrimSpace(url)
	if url == "" || strings.HasPrefix(url, "data:") || strings.HasPrefix(url, "#") {
		return nil, "", errSkippedURL
	}
	// Validate URL
	parsedURL, err := nurl.ParseRequestURI(url)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Hostname() == "" {
		return nil, "", errSkippedURL
	}

	// Get the asset from cache
	hurl := helper.Hash(url)
	b, _ := arc.cache.Get(hurl)
	if b != nil {
		asset, _ := DecodeWebAsset(b)
		if asset != nil {
			return asset.Data, asset.ContentType, nil
		}
	}

	// Download the asset, use semaphore to limit concurrent downloads
	err = arc.dlSemaphore.Acquire(ctx, 1)
	if err != nil {
		return nil, "", nil
	}
	resp, err := arc.downloadFile(url, parentURL)
	arc.dlSemaphore.Release(1)
	if err != nil {
		return nil, "", errSkippedURL
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
		return nil, "", err
	}

	// Put asset into the cache
	asset := WebAsset{
		Data:        bodyContent,
		ContentType: contentType,
	}
	b, _ = asset.Encode()
	if b != nil {
		arc.cache.Put(hurl, b)
	}

	return bodyContent, contentType, nil
}
