package db

import (
	"time"

	"github.com/ncarlier/readflow/pkg/model"
)

// ArticleRepository is the repository interface to manage Articles
type ArticleRepository interface {
	CountArticles(status string) (uint, error)
	CountArticlesByUserID(uid uint, req model.ArticlesPageRequest) (uint, error)
	GetPaginatedArticlesByUserID(uid uint, req model.ArticlesPageRequest) (*model.ArticlesPageResponse, error)
	GetArticleByID(id uint) (*model.Article, error)
	CreateOrUpdateArticle(article model.Article) (*model.Article, error)
	DeleteArticle(article model.Article) error
	MarkAllArticlesAsRead(uid uint, categoryID *uint) (int64, error)
	DeleteReadArticlesOlderThan(delay time.Duration) (int64, error)
	DeleteAllReadArticles(uid uint) (int64, error)
}
