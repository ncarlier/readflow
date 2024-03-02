package exporter

import (
	"context"

	"github.com/ncarlier/readflow/internal/model"

	"github.com/ncarlier/readflow/pkg/downloader"
)

// ArticleExporter is service used to export an article to a specific format.
type ArticleExporter interface {
	// Export an article to a specific format
	Export(ctx context.Context, article *model.Article) (*downloader.WebAsset, error)
}
