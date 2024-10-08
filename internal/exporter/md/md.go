package md

import (
	"context"

	"github.com/ncarlier/readflow/internal/exporter"
	"github.com/ncarlier/readflow/internal/exporter/html"
	"github.com/ncarlier/readflow/internal/model"

	"github.com/ncarlier/readflow/pkg/downloader"
	"github.com/ncarlier/readflow/pkg/mediatype"
	"github.com/ncarlier/readflow/pkg/utils"

	md "github.com/JohannesKaufmann/html-to-markdown"
)

// MarkdownExporter convert an article to Markdown format
type MarkdownExporter struct {
	htmlExporter *html.HTMLExporter
	converter    *md.Converter
}

func newMarkdownExporter(dl downloader.Downloader) (exporter.ArticleExporter, error) {
	return &MarkdownExporter{
		htmlExporter: &html.HTMLExporter{},
		converter:    md.NewConverter("", true, nil),
	}, nil
}

// Export an article to Markdown format
func (exp *MarkdownExporter) Export(ctx context.Context, article *model.Article) (*downloader.WebAsset, error) {
	asset, err := exp.htmlExporter.Export(ctx, article)
	if err != nil {
		return nil, err
	}

	data, err := exp.converter.ConvertBytes(asset.Data)
	if err != nil {
		return nil, err
	}

	return &downloader.WebAsset{
		Data:        data,
		ContentType: mediatype.HTML,
		Name:        utils.SanitizeFilename(article.Title) + ".md",
	}, nil
}

func init() {
	exporter.Register("md", newMarkdownExporter)
}
