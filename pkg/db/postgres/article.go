package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/ncarlier/reader/pkg/model"
)

var articleColumns = []string{
	"id",
	"user_id",
	"category_id",
	"title",
	"text",
	"html",
	"url",
	"image",
	"hash",
	"status",
	"published_at",
	"created_at",
	"updated_at",
}

func mapRowToArticle(row *sql.Row) (*model.Article, error) {
	article := &model.Article{}

	err := row.Scan(
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
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return article, nil
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
		`, strings.Join(articleColumns, ",")),
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
	return mapRowToArticle(row)
}

func (pg *DB) updateArticle(article model.Article) (*model.Article, error) {
	row := pg.db.QueryRow(fmt.Sprintf(`
		UPDATE articles SET
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
		`, strings.Join(articleColumns, ",")),
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
	return mapRowToArticle(row)
}

// CreateOrUpdateArticle creates or updates a article into the DB
func (pg *DB) CreateOrUpdateArticle(article model.Article) (*model.Article, error) {
	if article.ID != nil {
		return pg.updateArticle(article)
	}
	return pg.createArticle(article)
}

// GetArticleByID returns an article by its ID from DB
func (pg *DB) GetArticleByID(id uint) (*model.Article, error) {
	query, args, _ := pg.psql.Select(articleColumns...).From(
		"articles",
	).Where(sq.Eq{"id": id}).ToSql()
	row := pg.db.QueryRow(query, args...)

	return mapRowToArticle(row)
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
