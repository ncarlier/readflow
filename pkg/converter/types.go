package converter

import (
	"context"

	"github.com/ncarlier/readflow/pkg/model"
)

// ArticleConverter is service used to convert an article to a specific format.
type ArticleConverter interface {
	// Convert an article to a specific format
	Convert(ctx context.Context, article *model.Article) (*model.FileAsset, error)
}
