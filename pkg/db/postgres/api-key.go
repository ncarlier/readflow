package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/ncarlier/readflow/pkg/model"
)

const apiKeyColumns = `
	id,
	user_id,
	alias,
	token,
	last_usage_at,
	created_at,
	updated_at
`

func scanAPIKeyRow(row *sql.Row) (*model.APIKey, error) {
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

func (pg *DB) createAPIKey(apiKey model.APIKey) (*model.APIKey, error) {
	row := pg.db.QueryRow(fmt.Sprintf(`
		INSERT INTO api_keys
			(user_id, alias, token)
			VALUES
			($1, $2, $3)
			RETURNING %s
		`, apiKeyColumns),
		apiKey.UserID, apiKey.Alias, apiKey.Token,
	)
	return scanAPIKeyRow(row)
}

func (pg *DB) updateAPIKey(apiKey model.APIKey) (*model.APIKey, error) {
	row := pg.db.QueryRow(fmt.Sprintf(`
		UPDATE api_keys SET
			alias=$3,
			updated_at=NOW()
			WHERE id=$1 AND user_id=$2
			RETURNING %s
		`, apiKeyColumns),
		apiKey.ID, apiKey.UserID, apiKey.Alias,
	)

	return scanAPIKeyRow(row)
}

// CreateOrUpdateAPIKey creates or updates an apiKey into the DB
func (pg *DB) CreateOrUpdateAPIKey(apiKey model.APIKey) (*model.APIKey, error) {
	if apiKey.ID != nil {
		return pg.updateAPIKey(apiKey)
	}
	return pg.createAPIKey(apiKey)
}

// GetAPIKeyByID get an apiKey from the DB
func (pg *DB) GetAPIKeyByID(id uint) (*model.APIKey, error) {
	row := pg.db.QueryRow(fmt.Sprintf(`
		SELECT %s
		FROM api_keys
		WHERE id = $1`, apiKeyColumns),
		id,
	)
	return scanAPIKeyRow(row)
}

// GetAPIKeyByToken find an apiKey by token form the DB (last usage is updated!)
func (pg *DB) GetAPIKeyByToken(token string) (*model.APIKey, error) {
	row := pg.db.QueryRow(fmt.Sprintf(`
		UPDATE api_keys SET
			last_usage_at=NOW()
			WHERE token=$1
			RETURNING %s
		`, apiKeyColumns),
		token,
	)

	return scanAPIKeyRow(row)
}

// GetAPIKeyByUserIDAndAlias returns API key of an user by its alias
func (pg *DB) GetAPIKeyByUserIDAndAlias(userID uint, alias string) (*model.APIKey, error) {
	row := pg.db.QueryRow(fmt.Sprintf(`
		SELECT %s
		FROM api_keys
		WHERE user_id = $1 AND alias = $2`, apiKeyColumns),
		userID, alias,
	)

	return scanAPIKeyRow(row)
}

// GetAPIKeysByUserID returns api-keys of an user from DB
func (pg *DB) GetAPIKeysByUserID(userID uint) ([]model.APIKey, error) {
	rows, err := pg.db.Query(fmt.Sprintf(`
		SELECT %s
		FROM api_keys
		WHERE user_id=$1
		ORDER BY alias ASC`, apiKeyColumns),
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

// DeleteAPIKeys removes API keys from the DB
func (pg *DB) DeleteAPIKeys(uid uint, ids []uint) (int64, error) {
	query, args, _ := pg.psql.Delete("api_keys").Where(
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
