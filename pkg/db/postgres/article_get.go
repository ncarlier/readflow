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

// CountArticlesByUserID returns total nb of articles of an user from the DB
func (pg *DB) CountArticlesByUserID(uid uint, req model.ArticlesPageRequest) (uint, error) {
	counter := pg.psql.Select("count(*)").From(
		"articles",
	).Where(sq.Eq{"user_id": uid})

	if req.Category != nil {
		counter = counter.Where(sq.Eq{"category_id": *req.Category})
	}

	if req.Status != nil {
		counter = counter.Where(sq.Eq{"status": *req.Status})
	}

	query, args, _ := counter.ToSql()

	var count uint
	if err := pg.db.QueryRow(query, args...).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

// GetPaginatedArticlesByUserID returns a paginated list of user's articles from the DB
func (pg *DB) GetPaginatedArticlesByUserID(uid uint, req model.ArticlesPageRequest) (*model.ArticlesPageResponse, error) {
	total, err := pg.CountArticlesByUserID(uid, req)
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

	if req.AfterCursor != nil {
		if req.SortOrder == "asc" {
			selectBuilder = selectBuilder.Where(sq.Gt{"id": *req.AfterCursor})
		} else {
			selectBuilder = selectBuilder.Where(sq.Lt{"id": *req.AfterCursor})
		}
	}

	selectBuilder = selectBuilder.OrderBy("id " + strings.ToUpper(req.SortOrder))
	selectBuilder = selectBuilder.Limit(uint64(req.Limit + 1))

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
		if index <= req.Limit {
			result.Entries = append(result.Entries, article)
			result.EndCursor = *article.ID
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
