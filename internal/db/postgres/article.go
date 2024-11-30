package postgres

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/ncarlier/readflow/internal/model"
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
	"thumbhash",
	"hash",
	"status",
	"stars",
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
		&article.ThumbHash,
		&article.Hash,
		&article.Status,
		&article.Stars,
		&article.PublishedAt,
		&article.CreatedAt,
		&article.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, mapError(err)
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
		&article.ThumbHash,
		&article.Hash,
		&article.Status,
		&article.Stars,
		&article.PublishedAt,
		&article.CreatedAt,
		&article.UpdatedAt,
	)
}

// CreateArticleForUser creates an article into the DB
func (pg *DB) CreateArticleForUser(uid uint, form model.ArticleCreateForm) (*model.Article, error) {
	query, args, _ := pg.psql.Insert(
		"articles",
	).Columns(
		"user_id",
		"category_id",
		"title",
		"text",
		"html",
		"url",
		"image",
		"hash",
		"status",
		"tags",
		"published_at",
		"updated_at",
	).Values(
		uid,
		form.CategoryID,
		form.Title,
		form.Text,
		form.HTML,
		form.URL,
		form.Image,
		form.Hash(),
		"inbox",
		hashtagsExpr(form.Hashtags()),
		form.PublishedAt,
		"NOW()",
	).Suffix(
		"RETURNING " + strings.Join(articleColumns, ","),
	).ToSql()
	row := pg.db.QueryRow(query, args...)
	return mapRowToArticle(row)
}

// UpdateArticleForUser updates an article into the DB
func (pg *DB) UpdateArticleForUser(uid uint, form model.ArticleUpdateForm) (*model.Article, error) {
	update := map[string]interface{}{
		"updated_at": "NOW()",
	}
	if form.CategoryID != nil {
		update["category_id"] = *form.CategoryID
	}
	if form.Title != nil {
		update["title"] = *form.Title
	}
	if form.Image != nil {
		update["image"] = *form.Image
	}
	if form.HTML != nil {
		update["html"] = *form.HTML
	}
	if form.Text != nil {
		update["text"] = *form.Text
	}
	if form.Status != nil {
		update["status"] = *form.Status
	}
	if form.Stars != nil {
		update["stars"] = *form.Stars
	}
	if tags, found := form.Hashtags(); found {
		update["tags"] = hashtagsExpr(tags)
	}
	query, args, _ := pg.psql.Update(
		"articles",
	).SetMap(update).Where(
		sq.Eq{"id": form.ID},
	).Where(
		sq.Eq{"user_id": uid},
	).Suffix(
		"RETURNING " + strings.Join(articleColumns, ","),
	).ToSql()

	row := pg.db.QueryRow(query, args...)
	return mapRowToArticle(row)
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
func (pg *DB) DeleteArticle(id uint) error {
	result, err := pg.db.Exec(`
		DELETE FROM articles
			WHERE ID=$1
		`,
		id,
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

// MarkAllArticlesAsReadByUser set status to read for all user's articles of a specific status and category
func (pg *DB) MarkAllArticlesAsReadByUser(uid uint, status string, categoryID *uint) (int64, error) {
	update := map[string]interface{}{
		"status":     "read",
		"updated_at": "NOW()",
	}
	queryBuilder := pg.psql.Update(
		"articles",
	).SetMap(update).Where(
		sq.Eq{"user_id": uid},
	)
	queryBuilder = queryBuilder.Where(sq.Eq{"status": status})

	if categoryID != nil {
		queryBuilder = queryBuilder.Where(sq.Eq{"category_id": *categoryID})
	}

	query, args, _ := queryBuilder.ToSql()

	result, err := pg.db.Exec(query, args...)
	if err != nil {
		return 0, err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return count, err
}

// DeleteReadArticlesOlderThan remove old articles from the DB
func (pg *DB) DeleteReadArticlesOlderThan(delay time.Duration) (int64, error) {
	maxAge := time.Now().Add(-delay)
	query, args, _ := pg.psql.Delete(
		"articles",
	).Where(
		sq.Eq{"status": "read"},
	).Where(
		sq.Eq{"stars": 0},
	).Where(
		sq.Lt{"updated_at": maxAge},
	).ToSql()

	result, err := pg.db.Exec(query, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// DeleteAllReadArticlesByUser remove all read articles from the DB
func (pg *DB) DeleteAllReadArticlesByUser(uid uint) (int64, error) {
	query, args, _ := pg.psql.Delete(
		"articles",
	).Where(
		sq.Eq{"status": "read"},
	).Where(
		sq.Eq{"stars": 0},
	).Where(
		sq.Eq{"user_id": uid},
	).ToSql()

	result, err := pg.db.Exec(query, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// SetArticleThumbHash set article thumb hash value
func (pg *DB) SetArticleThumbHash(id uint, hash string) (*model.Article, error) {
	update := map[string]interface{}{
		"updated_at": "NOW()",
		"thumbhash":  hash,
	}
	query, args, _ := pg.psql.Update(
		"articles",
	).SetMap(update).Where(
		sq.Eq{"id": id},
	).Suffix(
		"RETURNING " + strings.Join(articleColumns, ","),
	).ToSql()

	row := pg.db.QueryRow(query, args...)
	return mapRowToArticle(row)
}
