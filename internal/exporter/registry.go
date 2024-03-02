package exporter

import (
	"fmt"

	"github.com/ncarlier/readflow/pkg/downloader"
)

// ArticleExporterCreator function for create an article exporter
type ArticleExporterCreator func(dl downloader.Downloader) (ArticleExporter, error)

// Registry of all Exporters
var registry = map[string]ArticleExporterCreator{}

// Register an article Exporter to the registry
func Register(format string, creator ArticleExporterCreator) {
	registry[format] = creator
}

// NewArticleExporter create new article Exporter
func NewArticleExporter(format string, dl downloader.Downloader) (ArticleExporter, error) {
	creator, ok := registry[format]
	if !ok {
		return nil, fmt.Errorf("unsupported export format: %s", format)
	}
	return creator(dl)
}
