package postgres

import (
	"database/sql"
	"errors"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/ncarlier/readflow/pkg/model"
)

var apiKeyColumns = []string{
	"id",
	"user_id",
	"alias",
	"token",
	"last_usage_at",
	"created_at",
	"updated_at",
}

func mapRowToAPIKey(row *sql.Row) (*model.APIKey, error) {
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

// CreateAPIKeyForUser creates an API key for a user
func (pg *DB) CreateAPIKeyForUser(uid uint, form model.APIKeyCreateForm) (*model.APIKey, error) {
	query, args, _ := pg.psql.Insert(
		"api_keys",
	).Columns(
		"user_id", "alias", "token",
	).Values(
		uid,
		form.Alias,
		form.Token,
	).Suffix(
		"RETURNING " + strings.Join(apiKeyColumns, ","),
	).ToSql()
	row := pg.db.QueryRow(query, args...)
	return mapRowToAPIKey(row)
}

// UpdateAPIKeyForUser updates an API key for a user
func (pg *DB) UpdateAPIKeyForUser(uid uint, form model.APIKeyUpdateForm) (*model.APIKey, error) {
	update := map[string]interface{}{
		"alias":      form.Alias,
		"updated_at": "NOW()",
	}
	query, args, _ := pg.psql.Update(
		"api_keys",
	).SetMap(update).Where(
		sq.Eq{"id": form.ID},
	).Where(
		sq.Eq{"user_id": uid},
	).Suffix(
		"RETURNING " + strings.Join(apiKeyColumns, ","),
	).ToSql()

	row := pg.db.QueryRow(query, args...)
	return mapRowToAPIKey(row)
}

// GetAPIKeyByID get an apiKey from the DB
func (pg *DB) GetAPIKeyByID(id uint) (*model.APIKey, error) {
	query, args, _ := pg.psql.Select(apiKeyColumns...).From(
		"api_keys",
	).Where(
		sq.Eq{"id": id},
	).ToSql()
	row := pg.db.QueryRow(query, args...)
	return mapRowToAPIKey(row)
}

// GetAPIKeyByToken find an apiKey by token form the DB (last usage is updated!)
func (pg *DB) GetAPIKeyByToken(token string) (*model.APIKey, error) {
	query, args, _ := pg.psql.Update(
		"api_keys",
	).Set("last_usage_at", "NOW()").Where(
		sq.Eq{"token": token},
	).Suffix(
		"RETURNING " + strings.Join(apiKeyColumns, ","),
	).ToSql()

	row := pg.db.QueryRow(query, args...)
	return mapRowToAPIKey(row)
}

// GetAPIKeyByUserIDAndAlias returns API key of an user by its alias
func (pg *DB) GetAPIKeyByUserAndAlias(uid uint, alias string) (*model.APIKey, error) {
	query, args, _ := pg.psql.Select(apiKeyColumns...).From(
		"api_keys",
	).Where(
		sq.Eq{"user_id": uid},
	).Where(
		sq.Eq{"alias": alias},
	).ToSql()

	row := pg.db.QueryRow(query, args...)
	return mapRowToAPIKey(row)
}

// GetAPIKeysByUser returns api-keys of an user from DB
func (pg *DB) GetAPIKeysByUser(uid uint) ([]model.APIKey, error) {
	query, args, _ := pg.psql.Select(apiKeyColumns...).From(
		"api_keys",
	).Where(
		sq.Eq{"user_id": uid},
	).OrderBy("alias ASC").ToSql()
	rows, err := pg.db.Query(query, args...)
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

// DeleteAPIKeyByUser removes an apiKey from the DB
func (pg *DB) DeleteAPIKeyByUser(uid uint, id uint) error {
	query, args, _ := pg.psql.Delete("api_keys").Where(
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
		return errors.New("no API key has been removed")
	}

	return nil
}

// DeleteAPIKeysByUser removes API keys from the DB
func (pg *DB) DeleteAPIKeysByUser(uid uint, ids []uint) (int64, error) {
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
