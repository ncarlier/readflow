package html

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	nurl "net/url"
	"strings"

	"github.com/go-shiori/dom"
	"github.com/ncarlier/readflow/pkg/constant"
	"github.com/ncarlier/readflow/pkg/exporter"
	"github.com/ncarlier/readflow/pkg/model"
	"golang.org/x/net/html"
	"golang.org/x/sync/errgroup"
)

var errSkippedURL = errors.New("skip processing url")

// SingleHTMLExporter convert an article to HTML offline format
type SingleHTMLExporter struct {
	downloader exporter.Downloader
}

func newHTMLOffflineExporter(downloader exporter.Downloader) (exporter.ArticleExporter, error) {
	return &SingleHTMLExporter{
		downloader: downloader,
	}, nil
}

// Export an article to HTML offline format
func (exp *SingleHTMLExporter) Export(ctx context.Context, article *model.Article) (*model.FileAsset, error) {
	var buffer bytes.Buffer
	if err := articleAsHTMLTpl.Execute(&buffer, article); err != nil {
		return nil, err
	}
	r := bytes.NewReader(buffer.Bytes())
	data, err := exp.exportWithEmbededAssets(ctx, r, *article.URL)
	if err != nil {
		return nil, err
	}

	return &model.FileAsset{
		Data:        data,
		ContentType: constant.ContentTypeHTML,
		Name:        strings.TrimRight(article.Title, ".") + ".html",
	}, nil
}

func (exp *SingleHTMLExporter) exportWithEmbededAssets(ctx context.Context, input io.Reader, baseURL string) ([]byte, error) {
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
			return exp.processNode(ctx, node, url)
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

func (exp *SingleHTMLExporter) processNode(ctx context.Context, node *html.Node, baseURL *nurl.URL) error {
	err := exp.processURLAttribute(ctx, node, "src", baseURL)
	if err != nil {
		return err
	}
	return nil
}

func (exp *SingleHTMLExporter) processURLAttribute(ctx context.Context, node *html.Node, attrName string, baseURL *nurl.URL) error {
	if !dom.HasAttribute(node, attrName) {
		return nil
	}

	url := dom.GetAttribute(node, attrName)
	asset, err := exp.processURL(ctx, url, baseURL.String())
	if err != nil && err != errSkippedURL {
		return err
	}

	newURL := url
	if err == nil {
		// TODO convert images to webp format
		newURL = asset.ToDataURL()
		dom.SetAttribute(node, "data-"+attrName, url)
	}

	dom.SetAttribute(node, attrName, newURL)
	return nil
}

func (exp *SingleHTMLExporter) processURL(ctx context.Context, url string, parentURL string) (*model.FileAsset, error) {
	// Ignore special URLs
	url = strings.TrimSpace(url)
	if url == "" || strings.HasPrefix(url, "data:") || strings.HasPrefix(url, "#") {
		return nil, errSkippedURL
	}
	// Validate URL
	parsedURL, err := nurl.ParseRequestURI(url)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Hostname() == "" {
		return nil, errSkippedURL
	}

	// Download URL
	asset, err := exp.downloader.Download(ctx, url)
	if err != nil {
		return nil, errSkippedURL
	}
	return asset, nil
}

func init() {
	exporter.Register("html-single", newHTMLOffflineExporter)
}
