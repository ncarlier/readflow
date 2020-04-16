package postgres

import (
	"strings"

	sq "github.com/Masterminds/squirrel"

	"github.com/ncarlier/readflow/pkg/model"
)

// CountArticles count articles
func (pg *DB) CountArticles(status string) (uint, error) {
	counter := pg.psql.Select("count(*)").From(
		"articles",
	)

	if status != "" {
		counter = counter.Where(sq.Eq{"status": status})
	}

	query, args, _ := counter.ToSql()

	var count uint
	if err := pg.db.QueryRow(query, args...).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

// CountArticlesByUser returns total nb of articles of an user from the DB
func (pg *DB) CountArticlesByUser(uid uint, req model.ArticlesPageRequest) (uint, error) {
	counter := pg.psql.Select("count(*)").From(
		"articles",
	).Where(sq.Eq{"user_id": uid})

	if req.Category != nil {
		counter = counter.Where(sq.Eq{"category_id": *req.Category})
	}

	if req.Status != nil {
		counter = counter.Where(sq.Eq{"status": *req.Status})
	}

	if req.Starred != nil {
		counter = counter.Where(sq.Eq{"starred": *req.Starred})
	}

	query, args, _ := counter.ToSql()

	var count uint
	if err := pg.db.QueryRow(query, args...).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

// GetPaginatedArticlesByUser returns a paginated list of user's articles from the DB
func (pg *DB) GetPaginatedArticlesByUser(uid uint, req model.ArticlesPageRequest) (*model.ArticlesPageResponse, error) {
	// Set defaults
	limit := uint(20)
	if req.Limit != nil {
		limit = *req.Limit
	}
	sortOrder := "asc"
	if req.SortOrder != nil {
		sortOrder = *req.SortOrder
	}

	// Count articles
	total, err := pg.CountArticlesByUser(uid, req)
	if err != nil {
		return nil, err
	}

	result := model.ArticlesPageResponse{
		TotalCount: total,
		HasNext:    false,
	}

	selectBuilder := pg.psql.Select(articleColumns...).From(
		"articles",
	).Where(sq.Eq{"user_id": uid})

	if req.Category != nil {
		selectBuilder = selectBuilder.Where(sq.Eq{"category_id": *req.Category})
	}

	if req.Status != nil {
		selectBuilder = selectBuilder.Where(sq.Eq{"status": *req.Status})
	}

	if req.Starred != nil {
		selectBuilder = selectBuilder.Where(sq.Eq{"starred": *req.Starred})
	}

	if req.AfterCursor != nil {
		if sortOrder == "asc" {
			selectBuilder = selectBuilder.Where(sq.Gt{"id": *req.AfterCursor})
		} else {
			selectBuilder = selectBuilder.Where(sq.Lt{"id": *req.AfterCursor})
		}
	}

	selectBuilder = selectBuilder.OrderBy("id " + strings.ToUpper(sortOrder))
	selectBuilder = selectBuilder.Limit(uint64(limit + 1))

	query, args, _ := selectBuilder.ToSql()
	rows, err := pg.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var index uint
	for rows.Next() {
		index++
		article := &model.Article{}
		if err := mapRowsToArticle(rows, article); err != nil {
			return nil, err
		}
		if index <= limit {
			result.Entries = append(result.Entries, article)
			result.EndCursor = article.ID
		} else {
			result.HasNext = true
		}
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return &result, nil
}
