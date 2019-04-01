package db

import "github.com/ncarlier/reader/pkg/model"

// ArticleRepository is the repository interface to manage Articles
type ArticleRepository interface {
	GetPaginatedArticlesByUserID(uid uint, req model.ArticlesPageRequest) (*model.ArticlesPageResponse, error)
	GetArticleByID(id uint) (*model.Article, error)
	CreateOrUpdateArticle(article model.Article) (*model.Article, error)
	DeleteArticle(article model.Article) error
	MarkAllArticlesAsRead(uid uint, categoryID *uint) (int64, error)
}
