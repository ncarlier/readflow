package txtpaper

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ncarlier/readflow/pkg/downloader"
	"github.com/ncarlier/readflow/pkg/exporter"
	"github.com/ncarlier/readflow/pkg/model"
)

const txtpaperURL = "https://txtpaper.com/api/v1/"

// TxtpaperExporter convert an article using txtpaper service
type TxtpaperExporter struct {
	format string
}

func newTxtpaperExporter(format string) func(dl downloader.Downloader) (exporter.ArticleExporter, error) {
	return func(dl downloader.Downloader) (exporter.ArticleExporter, error) {
		return &TxtpaperExporter{
			format: format,
		}, nil
	}
}

// Export an article using txtpaper service
func (exp *TxtpaperExporter) Export(ctx context.Context, article *model.Article) (*downloader.WebAsset, error) {
	form := url.Values{}
	form.Add("byline", *article.URL)
	form.Add("content", *article.HTML)
	form.Add("format", exp.format)

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "POST", txtpaperURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("invalid txtpaper response: %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return &downloader.WebAsset{
		Data:        body,
		ContentType: res.Header.Get("Content-Type"),
		Name:        strings.TrimRight(article.Title, ". ") + "." + exp.format,
	}, nil
}

func init() {
	exporter.Register("pdf", newTxtpaperExporter("pdf"))
	exporter.Register("md", newTxtpaperExporter("md"))
	exporter.Register("mobi", newTxtpaperExporter("mobi"))
}
