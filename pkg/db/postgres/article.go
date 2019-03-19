package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/ncarlier/reader/pkg/model"
)

const articleColumns = `
	id,
	user_id,
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
`

func mapRowToArticle(row *sql.Row, article *model.Article) error {
	return row.Scan(
		&article.ID,
		&article.UserID,
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
}

func mapRowsToArticle(rows *sql.Rows, article *model.Article) error {
	return rows.Scan(
		&article.ID,
		&article.UserID,
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
}

func (pg *DB) createArticle(article model.Article) (*model.Article, error) {
	row := pg.db.QueryRow(fmt.Sprintf(`
		INSERT INTO articles (
				user_id,
				category_id,
				title,
				text,
				html,
				url,
				image,
				hash,
				status,
				published_at
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
			RETURNING %s
		`, articleColumns),
		article.UserID,
		article.CategoryID,
		article.Title,
		article.Text,
		article.HTML,
		article.URL,
		article.Image,
		article.Hash,
		article.Status,
		article.PublishedAt,
	)
	result := &model.Article{}

	if err := mapRowToArticle(row, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (pg *DB) updateArticle(article model.Article) (*model.Article, error) {
	row := pg.db.QueryRow(fmt.Sprintf(`
		UPDATE article SET
			category_id = $3,
			title = $4,
			text = $5,
			html = $6,
			url = $7,
			image = $8,
			hash = $9,
			status = $10,
			published_at = $11,
			updated_at=NOW()
			WHERE id=$1 AND user_id=$2
			RETURNING %s
		`, articleColumns),
		article.ID,
		article.UserID,
		article.CategoryID,
		article.Title,
		article.Text,
		article.HTML,
		article.URL,
		article.Image,
		article.Hash,
		article.Status,
		article.PublishedAt,
	)

	result := &model.Article{}

	if err := mapRowToArticle(row, result); err != nil {
		return nil, err
	}
	return result, nil
}

// CreateOrUpdateArticle creates or updates a article into the DB
func (pg *DB) CreateOrUpdateArticle(article model.Article) (*model.Article, error) {
	if article.ID != nil {
		return pg.updateArticle(article)
	}
	return pg.createArticle(article)
}

// GetArticlesByUserID returns user's articles from DB
func (pg *DB) GetArticlesByUserID(userID uint) ([]model.Article, error) {
	rows, err := pg.db.Query(fmt.Sprintf(`
		SELECT %s
		FROM articles
		WHERE user_id=$1
		ORDER BY created_at DESC`, articleColumns), userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.Article

	for rows.Next() {
		article := &model.Article{}
		if err := mapRowsToArticle(rows, article); err != nil {
			return nil, err
		}
		result = append(result, *article)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetArticleByID returns an article by its ID from DB
func (pg *DB) GetArticleByID(id uint) (*model.Article, error) {
	row := pg.db.QueryRow(fmt.Sprintf(`
		SELECT %s
		FROM articles
		WHERE id = $1`, articleColumns),
		id,
	)

	result := &model.Article{}
	err := mapRowToArticle(row, result)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return result, nil
}

// DeleteArticle remove an article from the DB
func (pg *DB) DeleteArticle(article model.Article) error {
	result, err := pg.db.Exec(`
		DELETE FROM articles
			WHERE ID=$1
		`,
		article.ID,
	)
	if err != nil {
		return err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return errors.New("no article has been removed")
	}

	return nil
}
