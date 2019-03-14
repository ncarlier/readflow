package db

import "github.com/ncarlier/reader/pkg/model"

// ArticleRepository is the repository interface to manage Articles
type ArticleRepository interface {
	GetArticles() ([]model.Article, error)
	GetArticleByID(id uint32) (*model.Article, error)
	CreateOrUpdateArticle(article model.Article) (*model.Article, error)
	DeleteArticle(article model.Article) (*model.Article, error)
}
