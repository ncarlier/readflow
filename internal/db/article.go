package db

import (
	"time"

	"github.com/ncarlier/readflow/internal/model"
)

// ArticleRepository is the repository interface to manage Articles
type ArticleRepository interface {
	CountArticles(status string) (uint, error)
	CountArticlesByUser(uid uint, req model.ArticlesPageRequest) (uint, error)
	GetPaginatedArticlesByUser(uid uint, req model.ArticlesPageRequest) (*model.ArticlesPageResponse, error)
	GetArticleByID(id uint) (*model.Article, error)
	CreateArticleForUser(uid uint, form model.ArticleCreateForm) (*model.Article, error)
	UpdateArticleForUser(uid uint, form model.ArticleUpdateForm) (*model.Article, error)
	MarkAllArticlesAsReadByUser(uid uint, status string, categoryID *uint) (int64, error)
	DeleteArticle(id uint) error
	DeleteReadArticlesOlderThan(delay time.Duration) (int64, error)
	DeleteAllReadArticlesByUser(uid uint) (int64, error)
	SetArticleThumbHash(id uint, hash string) (*model.Article, error)
}
