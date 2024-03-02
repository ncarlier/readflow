package postgres

import (
	"database/sql"
	"strings"

	"github.com/ncarlier/readflow/internal/model"
)

var propertiesColumns = []string{
	"rev",
	"vapid_public_key",
	"vapid_private_key",
	"created_at",
}

func mapRowToProperties(row *sql.Row) (*model.Properties, error) {
	properties := &model.Properties{}

	err := row.Scan(
		&properties.Rev,
		&properties.VAPIDPublicKey,
		&properties.VAPIDPrivateKey,
		&properties.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return properties, nil
}

// CreateProperties update properties (create a new revision)
func (pg *DB) CreateProperties(properties model.Properties) (*model.Properties, error) {
	query, args, _ := pg.psql.Insert(
		"properties",
	).Columns(
		"vapid_public_key", "vapid_private_key",
	).Values(
		properties.VAPIDPublicKey,
		properties.VAPIDPrivateKey,
	).Suffix(
		"RETURNING " + strings.Join(propertiesColumns, ","),
	).ToSql()

	row := pg.db.QueryRow(query, args...)
	return mapRowToProperties(row)
}

// GetProperties get last revision of properties from the DB
func (pg *DB) GetProperties() (*model.Properties, error) {
	query, args, _ := pg.psql.Select(propertiesColumns...).From(
		"properties",
	).OrderBy("rev DESC").Limit(1).ToSql()
	row := pg.db.QueryRow(query, args...)
	return mapRowToProperties(row)
}
