package html

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	nurl "net/url"
	"path"
	"strings"

	"github.com/go-shiori/dom"
	"golang.org/x/net/html"

	"github.com/ncarlier/readflow/internal/exporter"
	"github.com/ncarlier/readflow/internal/model"

	"github.com/ncarlier/readflow/pkg/downloader"
	"github.com/ncarlier/readflow/pkg/utils"
)

// ZIPExporter convert an article to a ZIP archive
type ZIPExporter struct {
	dl downloader.Downloader
}

func newZIPExporter(dl downloader.Downloader) (exporter.ArticleExporter, error) {
	return &ZIPExporter{
		dl: dl,
	}, nil
}

// Export an article to ZIP archive
func (exp *ZIPExporter) Export(ctx context.Context, article *model.Article) (*downloader.WebAsset, error) {
	var buffer bytes.Buffer
	if err := articleAsHTMLTpl.Execute(&buffer, article); err != nil {
		return nil, err
	}
	r := bytes.NewReader(buffer.Bytes())

	// Create a buffer to write our archive to.
	buf := new(bytes.Buffer)
	// Create a new zip archive.
	w := zip.NewWriter(buf)

	err := exp.exportZIPArchive(ctx, r, w, *article.URL)
	if err != nil {
		w.Close()
		return nil, err
	}
	w.Close()

	return &downloader.WebAsset{
		Data:        buf.Bytes(),
		ContentType: "application/zip",
		Name:        utils.SanitizeFilename(article.Title) + ".zip",
	}, nil
}

func (exp *ZIPExporter) exportZIPArchive(ctx context.Context, input io.Reader, output *zip.Writer, baseURL string) error {
	url, err := nurl.ParseRequestURI(baseURL)
	if err != nil || url.Scheme == "" || url.Hostname() == "" {
		return fmt.Errorf("url \"%s\" is not valid", baseURL)
	}
	doc, err := html.Parse(input)
	if err != nil {
		return fmt.Errorf("failed to parse HTML: %w", err)
	}
	for _, node := range dom.GetElementsByTagName(doc, "img") {
		if err := exp.processNode(ctx, output, node, url); err != nil {
			return err
		}
	}

	f, err := output.Create("index.html")
	if err != nil {
		return err
	}
	err = html.Render(f, doc)
	if err != nil {
		return err
	}
	return err
}

func (exp *ZIPExporter) processNode(ctx context.Context, output *zip.Writer, node *html.Node, baseURL *nurl.URL) error {
	err := exp.processURLAttribute(ctx, output, node, "src", baseURL)
	if err != nil {
		return err
	}
	return nil
}

func (exp *ZIPExporter) processURLAttribute(ctx context.Context, output *zip.Writer, node *html.Node, attrName string, baseURL *nurl.URL) error {
	if !dom.HasAttribute(node, attrName) {
		return nil
	}

	url := dom.GetAttribute(node, attrName)
	asset, err := exp.processURL(ctx, url, baseURL.String())
	if err != nil {
		if err == errSkippedURL {
			return nil
		}
		return err
	}

	newURL := strings.Split(path.Base(asset.Name), "?")[0]
	f, err := output.Create(newURL)
	if err != nil {
		return err
	}
	_, err = f.Write(asset.Data)
	if err != nil {
		return err
	}
	dom.SetAttribute(node, attrName, newURL)
	return nil
}

func (exp *ZIPExporter) processURL(ctx context.Context, url, parentURL string) (*downloader.WebAsset, error) {
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
	exporter.Register("zip", newZIPExporter)
}
