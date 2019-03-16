package postgres

import (
	"database/sql"
	"errors"

	"github.com/ncarlier/reader/pkg/model"
)

func (pg *DB) createAPIKey(apiKey model.APIKey) (*model.APIKey, error) {
	row := pg.db.QueryRow(`
		INSERT INTO api_keys
			(user_id, alias, token)
			VALUES
			($1, $2, $3)
			RETURNING id, user_id, alias, token, created_at
		`,
		apiKey.UserID, apiKey.Alias, apiKey.Token,
	)
	result := model.APIKey{}

	err := row.Scan(
		&result.ID,
		&result.UserID,
		&result.Alias,
		&result.Token,
		&result.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (pg *DB) updateAPIKey(apiKey model.APIKey) (*model.APIKey, error) {
	row := pg.db.QueryRow(`
		UPDATE api_keys SET
			alias=$3,
			updated_at=NOW()
			WHERE id=$1 AND user_id=$2
			RETURNING id, user_id, alias, token, last_usage_at, created_at, updated_at
		`,
		apiKey.ID, apiKey.UserID, apiKey.Alias,
	)

	result := model.APIKey{}

	err := row.Scan(
		&result.ID,
		&result.UserID,
		&result.Alias,
		&result.Token,
		&result.LastUsageAt,
		&result.CreatedAt,
		&result.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateOrUpdateAPIKey creates or updates an apiKey into the DB
func (pg *DB) CreateOrUpdateAPIKey(apiKey model.APIKey) (*model.APIKey, error) {
	if apiKey.ID != nil {
		return pg.updateAPIKey(apiKey)
	}
	return pg.createAPIKey(apiKey)
}

// GetAPIKeyByToken find an apiKey by token form the DB (last usage is updated!)
func (pg *DB) GetAPIKeyByToken(token string) (*model.APIKey, error) {
	row := pg.db.QueryRow(`
		UPDATE api_keys SET
			last_usage_at=NOW()
			WHERE token=$1
			RETURNING id, user_id, alias, token, last_usage_at, created_at, updated_at
		`,
		token,
	)

	result := model.APIKey{}

	err := row.Scan(
		&result.ID,
		&result.UserID,
		&result.Alias,
		&result.Token,
		&result.LastUsageAt,
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

// GetAPIKeyByUserIDAndAlias returns API key of an user by its alias
func (pg *DB) GetAPIKeyByUserIDAndAlias(userID uint32, alias string) (*model.APIKey, error) {
	row := pg.db.QueryRow(`
		SELECT
			id,
			user_id,
			alias,
			token,
			last_usage_at,
			created_at,
			updated_at
		FROM api_keys
		WHERE user_id = $1 AND alias = $2`,
		userID, alias,
	)

	result := model.APIKey{}

	err := row.Scan(
		&result.ID,
		&result.UserID,
		&result.Alias,
		&result.Token,
		&result.LastUsageAt,
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

// GetAPIKeysByUserID returns api-keys of an user from DB
func (pg *DB) GetAPIKeysByUserID(userID uint32) ([]model.APIKey, error) {
	rows, err := pg.db.Query(`
		SELECT
			id,
			user_id,
			alias,
			token,
			last_usage_at,
			created_at,
			updated_at
		FROM api_keys
		WHERE user_id=$1
		ORDER BY alias ASC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.APIKey

	for rows.Next() {
		apiKey := model.APIKey{}
		err = rows.Scan(
			&apiKey.ID,
			&apiKey.UserID,
			&apiKey.Alias,
			&apiKey.Token,
			&apiKey.LastUsageAt,
			&apiKey.CreatedAt,
			&apiKey.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, apiKey)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return result, nil
}

// DeleteAPIKey removes an apiKey from the DB
func (pg *DB) DeleteAPIKey(apiKey model.APIKey) error {
	result, err := pg.db.Exec(`
		DELETE FROM api_keys
			WHERE ID=$1
		`,
		apiKey.ID,
	)
	if err != nil {
		return err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return errors.New("no apiKey has been removed")
	}

	return nil
}
