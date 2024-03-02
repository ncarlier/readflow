package postgres

import (
	"database/sql"
	"errors"

	"github.com/ncarlier/readflow/internal/model"
)

func (pg *DB) createUser(user model.User) (*model.User, error) {
	row := pg.db.QueryRow(`
		INSERT INTO users
			(username, enabled, plan)
			VALUES
			($1, $2, $3)
			RETURNING id, username, enabled, plan, created_at
		`,
		user.Username, user.Enabled, user.Plan,
	)
	result := model.User{}

	err := row.Scan(
		&result.ID,
		&result.Username,
		&result.Enabled,
		&result.Plan,
		&result.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (pg *DB) updateUser(user model.User) (*model.User, error) {
	row := pg.db.QueryRow(`
		UPDATE users SET
			enabled=$2,
			plan=$3,
			customer_id=$4,
			last_login_at=$5,
			updated_at=NOW()
			WHERE id=$1
			RETURNING id, username, enabled, plan, customer_id, last_login_at, created_at, updated_at
		`,
		user.ID, user.Enabled, user.Plan, user.CustomerID, user.LastLoginAt,
	)

	result := model.User{}

	err := row.Scan(
		&result.ID,
		&result.Username,
		&result.Enabled,
		&result.Plan,
		&result.CustomerID,
		&result.LastLoginAt,
		&result.CreatedAt,
		&result.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateOrUpdateUser creates or updates a user into the DB
func (pg *DB) CreateOrUpdateUser(user model.User) (*model.User, error) {
	if user.ID != nil {
		return pg.updateUser(user)
	}
	return pg.createUser(user)
}

// GetUserByID returns a user by its ID from DB
func (pg *DB) GetUserByID(id uint) (*model.User, error) {
	row := pg.db.QueryRow(`
		SELECT
			id,
			username,
			enabled,
			plan,
			customer_id,
			last_login_at,
			created_at,
			updated_at
		FROM users
		WHERE id = $1`,
		id,
	)

	result := model.User{}

	err := row.Scan(
		&result.ID,
		&result.Username,
		&result.Enabled,
		&result.Plan,
		&result.CustomerID,
		&result.LastLoginAt,
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

// GetUserByUsername returns user by its username from DB
func (pg *DB) GetUserByUsername(username string) (*model.User, error) {
	row := pg.db.QueryRow(`
		SELECT
			id,
			username,
			enabled,
			plan,
			customer_id,
			last_login_at,
			created_at,
			updated_at
		FROM users
		WHERE username = $1`,
		username,
	)

	result := model.User{}

	err := row.Scan(
		&result.ID,
		&result.Username,
		&result.Enabled,
		&result.Plan,
		&result.CustomerID,
		&result.LastLoginAt,
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

// DeleteUser removes an user from the DB
func (pg *DB) DeleteUser(user model.User) error {
	result, err := pg.db.Exec(`
		DELETE FROM users
			WHERE username=$1
		`,
		user.Username,
	)
	if err != nil {
		return err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return errors.New("no user has been removed")
	}

	return nil
}

// CountUsers count users
func (pg *DB) CountUsers() (uint, error) {
	counter := pg.psql.Select("count(*)").From("users")
	query, args, _ := counter.ToSql()

	var count uint
	if err := pg.db.QueryRow(query, args...).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}
