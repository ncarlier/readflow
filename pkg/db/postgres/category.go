package postgres

import (
	"database/sql"
	"errors"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/ncarlier/readflow/pkg/model"
)

var categoryColumns = []string{
	"id",
	"user_id",
	"title",
	"rule",
	"created_at",
	"updated_at",
}

func mapRowToCategory(row *sql.Row) (*model.Category, error) {
	cat := &model.Category{}

	err := row.Scan(
		&cat.ID,
		&cat.UserID,
		&cat.Title,
		&cat.Rule,
		&cat.CreatedAt,
		&cat.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return cat, nil
}
func (pg *DB) createCategory(category model.Category) (*model.Category, error) {
	query, args, _ := pg.psql.Insert(
		"categories",
	).Columns(
		"user_id", "title", "rule",
	).Values(
		category.UserID,
		category.Title,
		category.Rule,
	).Suffix(
		"RETURNING " + strings.Join(categoryColumns, ","),
	).ToSql()
	row := pg.db.QueryRow(query, args...)
	return mapRowToCategory(row)
}

func (pg *DB) updateCategory(category model.Category) (*model.Category, error) {
	update := map[string]interface{}{
		"title":      category.Title,
		"rule":       category.Rule,
		"updated_at": "NOW()",
	}
	query, args, _ := pg.psql.Update(
		"categories",
	).SetMap(update).Where(
		sq.Eq{"id": category.ID},
	).Where(
		sq.Eq{"user_id": category.UserID},
	).Suffix(
		"RETURNING " + strings.Join(categoryColumns, ","),
	).ToSql()

	row := pg.db.QueryRow(query, args...)
	return mapRowToCategory(row)
}

// CreateOrUpdateCategory creates or updates a category into the DB
func (pg *DB) CreateOrUpdateCategory(category model.Category) (*model.Category, error) {
	if category.ID != nil {
		return pg.updateCategory(category)
	}
	return pg.createCategory(category)
}

// GetCategoryByID returns a category from the DB
func (pg *DB) GetCategoryByID(id uint) (*model.Category, error) {
	query, args, _ := pg.psql.Select(categoryColumns...).From(
		"categories",
	).Where(
		sq.Eq{"id": id},
	).ToSql()
	row := pg.db.QueryRow(query, args...)
	return mapRowToCategory(row)
}

// GetCategoryByUserIDAndTitle returns a category of an user form the DB
func (pg *DB) GetCategoryByUserIDAndTitle(uid uint, title string) (*model.Category, error) {
	query, args, _ := pg.psql.Select(categoryColumns...).From(
		"categories",
	).Where(
		sq.Eq{"user_id": uid},
	).Where(
		sq.Eq{"title": title},
	).ToSql()

	row := pg.db.QueryRow(query, args...)
	return mapRowToCategory(row)
}

// GetCategoriesByUserID returns categories of an user from DB
func (pg *DB) GetCategoriesByUserID(uid uint) ([]model.Category, error) {
	query, args, _ := pg.psql.Select(categoryColumns...).From(
		"categories",
	).Where(
		sq.Eq{"user_id": uid},
	).OrderBy("title ASC").ToSql()
	rows, err := pg.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.Category

	for rows.Next() {
		cat := model.Category{}
		err = rows.Scan(
			&cat.ID,
			&cat.UserID,
			&cat.Title,
			&cat.Rule,
			&cat.CreatedAt,
			&cat.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, cat)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return result, nil
}

// CountCategoriesByUserID returns total nb of categories of an user from the DB
func (pg *DB) CountCategoriesByUserID(uid uint) (uint, error) {
	counter := pg.psql.Select("count(*)").From(
		"categories",
	).Where(sq.Eq{"user_id": uid})
	query, args, _ := counter.ToSql()

	var count uint
	if err := pg.db.QueryRow(query, args...).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

// DeleteCategory removes an category from the DB
func (pg *DB) DeleteCategory(category model.Category) error {
	query, args, _ := pg.psql.Delete("categories").Where(
		sq.Eq{"id": category.ID},
	).Where(
		sq.Eq{"user_id": category.UserID},
	).ToSql()
	result, err := pg.db.Exec(query, args...)
	if err != nil {
		return err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return errors.New("no category has been removed")
	}

	return nil
}

// DeleteCategories removes categories from the DB
func (pg *DB) DeleteCategories(uid uint, ids []uint) (int64, error) {
	query, args, _ := pg.psql.Delete("categories").Where(
		sq.Eq{"user_id": uid},
	).Where(
		sq.Eq{"id": ids},
	).ToSql()
	result, err := pg.db.Exec(query, args...)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}
