package postgres

import (
	"fmt"

	"github.com/ncarlier/reader/pkg/model"
)

// GetArticles returns articles from DB
func (pg *DB) GetArticles() ([]model.Article, error) {
	rows, err := pg.db.Query(`
		SELECT
			id,
			category_id,
			title,
			text,
			html,
			url,
			image,
			hash,
			status,
			published_at,
			created_at,
			updated_at
		FROM articles
		ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.Article

	for rows.Next() {
		article := model.Article{}
		err = rows.Scan(
			&article.ID,
			&article.CategoryID,
			&article.Title,
			&article.Text,
			&article.HTML,
			&article.URL,
			&article.Image,
			&article.Hash,
			&article.Status,
			&article.PublishedAt,
			&article.CreatedAt,
			&article.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, article)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetArticleByID returns an article by its ID from DB
func (pg *DB) GetArticleByID(id uint32) (*model.Article, error) {
	return nil, fmt.Errorf("Not yet implemented")
}

// CreateOrUpdateArticle creates an article into the DB
func (pg *DB) CreateOrUpdateArticle(article model.Article) (*model.Article, error) {
	return nil, fmt.Errorf("Not yet implemented")
}

// DeleteArticle remove an article from the DB
func (pg *DB) DeleteArticle(article model.Article) (*model.Article, error) {
	return nil, fmt.Errorf("Not yet implemented")
}
