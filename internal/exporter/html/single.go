package html

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	nurl "net/url"
	"strings"

	"github.com/ncarlier/readflow/internal/exporter"
	"github.com/ncarlier/readflow/internal/model"

	"github.com/ncarlier/readflow/pkg/downloader"
	"github.com/ncarlier/readflow/pkg/mediatype"
	"github.com/ncarlier/readflow/pkg/utils"

	"github.com/go-shiori/dom"
	"golang.org/x/net/html"
	"golang.org/x/sync/errgroup"
)

var errSkippedURL = errors.New("skip processing url")

// SingleHTMLExporter convert an article to HTML offline format
type SingleHTMLExporter struct {
	dl downloader.Downloader
}

func newHTMLOffflineExporter(dl downloader.Downloader) (exporter.ArticleExporter, error) {
	return &SingleHTMLExporter{
		dl: dl,
	}, nil
}

// Export an article to HTML offline format
func (exp *SingleHTMLExporter) Export(ctx context.Context, article *model.Article) (*downloader.WebAsset, error) {
	var buffer bytes.Buffer
	if err := articleAsHTMLTpl.Execute(&buffer, article); err != nil {
		return nil, err
	}
	r := bytes.NewReader(buffer.Bytes())
	data, err := exp.exportWithEmbededAssets(ctx, r, *article.URL)
	if err != nil {
		return nil, err
	}

	return &downloader.WebAsset{
		Data:        data,
		ContentType: mediatype.HTML,
		Name:        utils.SanitizeFilename(article.Title) + ".html",
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
	if err := g.Wait(); err != nil {
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
		dom.RemoveAttribute(node, attrName+"-set")
	}

	dom.SetAttribute(node, attrName, newURL)
	return nil
}

func (exp *SingleHTMLExporter) processURL(ctx context.Context, url, parentURL string) (*downloader.WebAsset, error) {
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
	asset, _, err := exp.dl.Get(ctx, url, nil)
	if err != nil {
		return nil, errSkippedURL
	}
	return asset, nil
}

func init() {
	exporter.Register("html-single", newHTMLOffflineExporter)
}
