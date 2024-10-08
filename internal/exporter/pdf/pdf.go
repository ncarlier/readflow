package pdf

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/ncarlier/readflow/internal/exporter"
	"github.com/ncarlier/readflow/internal/exporter/html"
	"github.com/ncarlier/readflow/internal/model"

	"github.com/ncarlier/readflow/pkg/downloader"
	"github.com/ncarlier/readflow/pkg/utils"
)

// PDFExporter convert an article to PDF file
type PDFExporter struct {
	endpoint     string
	htmlExporter *html.HTMLExporter
}

// NewPDFExporter create new PDF exporter
func NewPDFExporter(endpoint string) func(dl downloader.Downloader) (exporter.ArticleExporter, error) {
	return func(dl downloader.Downloader) (exporter.ArticleExporter, error) {
		return &PDFExporter{
			endpoint:     endpoint,
			htmlExporter: &html.HTMLExporter{},
		}, nil
	}
}

// Export an article using txtpaper service
func (exp *PDFExporter) Export(ctx context.Context, article *model.Article) (*downloader.WebAsset, error) {
	asset, err := exp.htmlExporter.Export(ctx, article)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	part, err := writer.CreateFormFile("files", "index.html")
	if err != nil {
		return nil, err
	}
	if _, err := part.Write(asset.Data); err != nil {
		return nil, err
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", exp.endpoint, &buf)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", writer.FormDataContentType())
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		cause, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("invalid PDF generator service response: %d - %s", res.StatusCode, cause)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return &downloader.WebAsset{
		Data:        body,
		ContentType: res.Header.Get("Content-Type"),
		Name:        utils.SanitizeFilename(article.Title) + ".pdf",
	}, nil
}
