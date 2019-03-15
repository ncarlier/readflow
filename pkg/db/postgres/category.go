package postgres

import (
	"database/sql"
	"errors"

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

// GetCategoryByUserIDAndTitle returns a category of an user form the DB
func (pg *DB) GetCategoryByUserIDAndTitle(userID uint32, title string) (*model.Category, error) {
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
func (pg *DB) GetCategoriesByUserID(userID uint32) ([]model.Category, error) {
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
			WHERE ID=$1
		`,
		category.ID,
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
