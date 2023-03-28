package exporter

import (
	"context"

	"github.com/ncarlier/readflow/pkg/downloader"
	"github.com/ncarlier/readflow/pkg/model"
)

// ArticleExporter is service used to export an article to a specific format.
type ArticleExporter interface {
	// Export an article to a specific format
	Export(ctx context.Context, article *model.Article) (*downloader.WebAsset, error)
}
