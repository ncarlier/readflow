package postgres

import (
	"database/sql"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/ncarlier/reader/pkg/model"
)

func (pg *DB) createCategory(category model.Category) (*model.Category, error) {
	row := pg.db.QueryRow(`
		INSERT INTO categories
			(user_id, title)
			VALUES
			($1, $2)
			RETURNING id, user_id, title, created_at
		`,
		category.UserID, category.Title,
	)
	result := model.Category{}

	err := row.Scan(
		&result.ID,
		&result.UserID,
		&result.Title,
		&result.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (pg *DB) updateCategory(category model.Category) (*model.Category, error) {
	row := pg.db.QueryRow(`
		UPDATE categories SET
			title=$3,
			updated_at=NOW()
			WHERE id=$1 AND user_id=$2
			RETURNING id, user_id, title, created_at, updated_at
		`,
		category.ID, category.UserID, category.Title,
	)

	result := model.Category{}

	err := row.Scan(
		&result.ID,
		&result.UserID,
		&result.Title,
		&result.CreatedAt,
		&result.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &result, nil
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
	row := pg.db.QueryRow(`
		SELECT
			id,
			user_id,
			title,
			created_at,
			updated_at
		FROM categories
		WHERE id = $1`,
		id,
	)

	result := model.Category{}

	err := row.Scan(
		&result.ID,
		&result.UserID,
		&result.Title,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetCategoryByUserIDAndTitle returns a category of an user form the DB
func (pg *DB) GetCategoryByUserIDAndTitle(userID uint, title string) (*model.Category, error) {
	row := pg.db.QueryRow(`
		SELECT
			id,
			user_id,
			title,
			created_at,
			updated_at
		FROM categories
		WHERE user_id = $1 AND title = $2`,
		userID, title,
	)

	result := model.Category{}

	err := row.Scan(
		&result.ID,
		&result.UserID,
		&result.Title,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetCategoriesByUserID returns categories of an user from DB
func (pg *DB) GetCategoriesByUserID(userID uint) ([]model.Category, error) {
	rows, err := pg.db.Query(`
		SELECT
			id,
			user_id,
			title,
			created_at,
			updated_at
		FROM categories
		WHERE user_id=$1
		ORDER BY title ASC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.Category

	for rows.Next() {
		category := model.Category{}
		err = rows.Scan(
			&category.ID,
			&category.UserID,
			&category.Title,
			&category.CreatedAt,
			&category.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, category)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return result, nil
}

// DeleteCategory removes an category from the DB
func (pg *DB) DeleteCategory(category model.Category) error {
	result, err := pg.db.Exec(`
		DELETE FROM categories
			WHERE id=$1 AND user_id=$2
		`,
		category.ID, category.UserID,
	)
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
