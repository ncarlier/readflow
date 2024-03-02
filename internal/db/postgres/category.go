package postgres

import (
	"database/sql"
	"errors"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/ncarlier/readflow/internal/model"
)

var categoryColumns = []string{
	"id",
	"user_id",
	"title",
	"created_at",
	"updated_at",
}

func mapRowToCategory(row *sql.Row) (*model.Category, error) {
	cat := &model.Category{}

	err := row.Scan(
		&cat.ID,
		&cat.UserID,
		&cat.Title,
		&cat.CreatedAt,
		&cat.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, mapError(err)
	}
	return cat, nil
}

// CreateCategoryForUser create a category for an user
func (pg *DB) CreateCategoryForUser(uid uint, form model.CategoryCreateForm) (*model.Category, error) {
	query, args, _ := pg.psql.Insert(
		"categories",
	).Columns(
		"user_id", "title",
	).Values(
		uid,
		strings.TrimSpace(form.Title),
	).Suffix(
		"RETURNING " + strings.Join(categoryColumns, ","),
	).ToSql()
	row := pg.db.QueryRow(query, args...)
	return mapRowToCategory(row)
}

// UpdateCategoryForUser update a category for an user
func (pg *DB) UpdateCategoryForUser(uid uint, form model.CategoryUpdateForm) (*model.Category, error) {
	update := map[string]interface{}{
		"updated_at": "NOW()",
	}
	if form.Title != nil {
		update["title"] = strings.TrimSpace(*form.Title)
	}
	query, args, _ := pg.psql.Update(
		"categories",
	).SetMap(update).Where(
		sq.Eq{"id": form.ID},
	).Where(
		sq.Eq{"user_id": uid},
	).Suffix(
		"RETURNING " + strings.Join(categoryColumns, ","),
	).ToSql()

	row := pg.db.QueryRow(query, args...)
	return mapRowToCategory(row)
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

// GetCategoryByUserAndTitle returns a category of an user form the DB
func (pg *DB) GetCategoryByUserAndTitle(uid uint, title string) (*model.Category, error) {
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

// GetCategoriesByUser returns categories of an user from DB
func (pg *DB) GetCategoriesByUser(uid uint) ([]model.Category, error) {
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

// CountCategoriesByUser returns total nb of categories of an user from the DB
func (pg *DB) CountCategoriesByUser(uid uint) (uint, error) {
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

// DeleteCategoryByUser removes an category from the DB
func (pg *DB) DeleteCategoryByUser(uid, id uint) error {
	query, args, _ := pg.psql.Delete("categories").Where(
		sq.Eq{"id": id},
	).Where(
		sq.Eq{"user_id": uid},
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

// DeleteCategoriesByUser removes categories from the DB
func (pg *DB) DeleteCategoriesByUser(uid uint, ids []uint) (int64, error) {
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
