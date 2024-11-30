package postgres

import (
	"strings"

	sq "github.com/Masterminds/squirrel"

	"github.com/ncarlier/readflow/internal/model"
	"github.com/ncarlier/readflow/pkg/utils"
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
		if *req.Starred {
			counter = counter.Where(sq.Gt{"stars": 0})
		} else {
			counter = counter.Where(sq.Eq{"stars": 0})
		}
	}

	if req.Query != nil && strings.TrimSpace(*req.Query) != "" {
		query := utils.ReplaceHashtagsPrefix(*req.Query, "~")
		counter = counter.Where(sq.Expr("search @@ websearch_to_tsquery(?)", query))
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
	limit := 20
	if req.Limit != nil {
		limit = *req.Limit
	}
	sortOrder := "asc"
	if req.SortOrder != nil {
		// Note that sort order is ignored when full-text query is used
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
		if *req.Starred {
			selectBuilder = selectBuilder.Where(sq.Gt{"stars": 0})
		} else {
			selectBuilder = selectBuilder.Where(sq.Eq{"stars": 0})
		}
	}

	var offset uint
	isFullText := req.Query != nil && strings.TrimSpace(*req.Query) != ""
	isSortedBy := req.SortBy != nil && strings.TrimSpace(*req.SortBy) != "key"
	if isFullText {
		// Full-text search query:
		// Classic Limit-Offset pagination (beware of performance issue)
		query := utils.ReplaceHashtagsPrefix(*req.Query, "~")
		selectBuilder = selectBuilder.Where(sq.Expr("search @@ websearch_to_tsquery(?)", query))
		selectBuilder = selectBuilder.OrderByClause("ts_rank(search, websearch_to_tsquery(?)) DESC", query)
		if req.AfterCursor != nil {
			offset = *req.AfterCursor
			selectBuilder = selectBuilder.Offset(uint64(offset))
		}
	} else if isSortedBy {
		// Search query with order by:
		// Classic Limit-Offset pagination (beware of performance issue)
		by := strings.TrimSpace(*req.SortBy)
		selectBuilder = selectBuilder.OrderBy(by + " " + strings.ToUpper(sortOrder))
		if req.AfterCursor != nil {
			offset = *req.AfterCursor
			selectBuilder = selectBuilder.Offset(uint64(offset))
		}
	} else {
		// Standard search:
		// Keyset Pagination using id as ordered column
		if req.AfterCursor != nil {
			if sortOrder == "asc" {
				selectBuilder = selectBuilder.Where(sq.Gt{"id": *req.AfterCursor})
			} else {
				selectBuilder = selectBuilder.Where(sq.Lt{"id": *req.AfterCursor})
			}
		}
		selectBuilder = selectBuilder.OrderBy("id " + strings.ToUpper(sortOrder))
	}

	// Limit the query by "plus-one" in order to known if the query has more results
	selectBuilder = selectBuilder.Limit(uint64(limit + 1))

	query, args, _ := selectBuilder.ToSql()
	rows, err := pg.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Read resultset...
	var index int
	for rows.Next() {
		index++
		article := &model.Article{}
		if err := mapRowsToArticle(rows, article); err != nil {
			return nil, err
		}
		if index <= limit {
			result.Entries = append(result.Entries, article)
			if isFullText || isSortedBy {
				result.EndCursor = offset + uint(index)
			} else {
				result.EndCursor = article.ID
			}
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
