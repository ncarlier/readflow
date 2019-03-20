package postgres

import (
	"fmt"
	"strings"

	"github.com/ncarlier/reader/pkg/model"
)

func orderByQueryPart(sortOrder string) string {
	return fmt.Sprintf("\nORDER BY id %s", strings.ToUpper(sortOrder))
}

func whereQueryPart(req *model.ArticlesPageRequest) string {
	result := "user_id=$1"
	if req.Category != nil {
		result += " AND category_id=$2"
	}
	if req.AfterCursor != nil {
		var part string
		if req.SortOrder == "asc" {
			part = " AND id > $%d"
		} else {
			part = " AND id < $%d"
		}
		idx := 2
		if req.Category != nil {
			idx = 3
		}
		result += fmt.Sprintf(part, idx)
	}
	return result
}

// CountArticlesByUserID returns total nb of articles of an user from the DB
func (pg *DB) CountArticlesByUserID(uid uint, req model.ArticlesPageRequest) (uint, error) {
	query := fmt.Sprintf(
		"SELECT count(*) FROM articles WHERE %s",
		whereQueryPart(&model.ArticlesPageRequest{
			Category: req.Category,
		}),
	)
	var count uint
	args := []interface{}{uid}
	if req.Category != nil {
		args = append(args, req.Category)
	}
	err := pg.db.QueryRow(
		query,
		args...,
	).Scan(&count)
	if err != nil {
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

	query := fmt.Sprintf(
		"SELECT %s FROM articles WHERE %s ORDER BY id %s LIMIT %d",
		articleColumns,
		whereQueryPart(&req),
		strings.ToUpper(req.SortOrder),
		req.Limit+1,
	)

	args := []interface{}{uid}
	if req.Category != nil {
		args = append(args, req.Category)
	}
	if req.AfterCursor != nil {
		args = append(args, req.AfterCursor)
	}

	rows, err := pg.db.Query(
		query,
		args...,
	)
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
