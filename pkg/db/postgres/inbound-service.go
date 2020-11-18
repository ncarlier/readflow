package postgres

import (
	"database/sql"
	"errors"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/ncarlier/readflow/pkg/model"
)

var inboundServiceColumns = []string{
	"id",
	"user_id",
	"alias",
	"token",
	"provider",
	"config",
	"last_usage_at",
	"created_at",
	"updated_at",
}

func mapRowToInboundService(row *sql.Row) (*model.InboundService, error) {
	result := model.InboundService{}

	err := row.Scan(
		&result.ID,
		&result.UserID,
		&result.Alias,
		&result.Token,
		&result.Provider,
		&result.Config,
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

// CreateInboundServiceForUser creates an inbound service for a user
func (pg *DB) CreateInboundServiceForUser(uid uint, form model.InboundServiceCreateForm) (*model.InboundService, error) {
	query, args, _ := pg.psql.Insert(
		"inbound_services",
	).Columns(
		"user_id", "alias", "token", "provider", "config",
	).Values(
		uid,
		form.Alias,
		form.Token,
		form.Provider,
		form.Config,
	).Suffix(
		"RETURNING " + strings.Join(inboundServiceColumns, ","),
	).ToSql()
	row := pg.db.QueryRow(query, args...)
	return mapRowToInboundService(row)
}

// UpdateInboundServiceForUser updates an inbound service for a user
func (pg *DB) UpdateInboundServiceForUser(uid uint, form model.InboundServiceUpdateForm) (*model.InboundService, error) {
	update := map[string]interface{}{
		"updated_at": "NOW()",
	}
	if form.Alias != nil {
		update["alias"] = *form.Alias
	}
	if form.Provider != nil {
		update["provider"] = *form.Provider
	}
	if form.Config != nil {
		update["config"] = *form.Config
	}
	query, args, _ := pg.psql.Update(
		"inbound_services",
	).SetMap(update).Where(
		sq.Eq{"id": form.ID},
	).Where(
		sq.Eq{"user_id": uid},
	).Suffix(
		"RETURNING " + strings.Join(inboundServiceColumns, ","),
	).ToSql()

	row := pg.db.QueryRow(query, args...)
	return mapRowToInboundService(row)
}

// GetInboundServiceByID get an inbound service from the DB
func (pg *DB) GetInboundServiceByID(id uint) (*model.InboundService, error) {
	query, args, _ := pg.psql.Select(inboundServiceColumns...).From(
		"inbound_services",
	).Where(
		sq.Eq{"id": id},
	).ToSql()
	row := pg.db.QueryRow(query, args...)
	return mapRowToInboundService(row)
}

// GetInboundServiceByToken find an inbound service by token form the DB (last usage is updated!)
func (pg *DB) GetInboundServiceByToken(token string) (*model.InboundService, error) {
	query, args, _ := pg.psql.Update(
		"inbound_services",
	).Set("last_usage_at", "NOW()").Where(
		sq.Eq{"token": token},
	).Suffix(
		"RETURNING " + strings.Join(inboundServiceColumns, ","),
	).ToSql()

	row := pg.db.QueryRow(query, args...)
	return mapRowToInboundService(row)
}

// GetInboundServiceByUserAndAlias returns inbound service of an user by its alias
func (pg *DB) GetInboundServiceByUserAndAlias(uid uint, alias string) (*model.InboundService, error) {
	query, args, _ := pg.psql.Select(inboundServiceColumns...).From(
		"inbound_services",
	).Where(
		sq.Eq{"user_id": uid},
	).Where(
		sq.Eq{"alias": alias},
	).ToSql()

	row := pg.db.QueryRow(query, args...)
	return mapRowToInboundService(row)
}

// GetInboundServicesByUser returns inbound services of an user from DB
func (pg *DB) GetInboundServicesByUser(uid uint) ([]model.InboundService, error) {
	query, args, _ := pg.psql.Select(inboundServiceColumns...).From(
		"inbound_services",
	).Where(
		sq.Eq{"user_id": uid},
	).OrderBy("alias ASC").ToSql()
	rows, err := pg.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.InboundService

	for rows.Next() {
		item := model.InboundService{}
		err = rows.Scan(
			&item.ID,
			&item.UserID,
			&item.Alias,
			&item.Token,
			&item.Provider,
			&item.Config,
			&item.LastUsageAt,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return result, nil
}

// DeleteInboundServiceByUser removes an inboundService from the DB
func (pg *DB) DeleteInboundServiceByUser(uid uint, id uint) error {
	query, args, _ := pg.psql.Delete("inbound_services").Where(
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
		return errors.New("no inbound service been removed")
	}

	return nil
}

// DeleteInboundServicesByUser removes inbound services from the DB
func (pg *DB) DeleteInboundServicesByUser(uid uint, ids []uint) (int64, error) {
	query, args, _ := pg.psql.Delete("inbound_services").Where(
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
