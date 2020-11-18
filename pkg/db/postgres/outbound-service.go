package postgres

import (
	"database/sql"
	"errors"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/ncarlier/readflow/pkg/model"
)

var outboundServiceColumns = []string{
	"id",
	"user_id",
	"alias",
	"is_default",
	"provider",
	"config",
	"created_at",
	"updated_at",
}

func mapRowToOutboundService(row *sql.Row) (*model.OutboundService, error) {
	result := &model.OutboundService{}

	err := row.Scan(
		&result.ID,
		&result.UserID,
		&result.Alias,
		&result.IsDefault,
		&result.Provider,
		&result.Config,
		&result.CreatedAt,
		&result.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return result, nil
}

// CreateOutboundServiceForUser creates an outbound service into the DB
func (pg *DB) CreateOutboundServiceForUser(uid uint, form model.OutboundServiceCreateForm) (*model.OutboundService, error) {
	tx, err := pg.db.Begin()
	if err != nil {
		return nil, err
	}
	query, args, _ := pg.psql.Insert(
		"outbound_services",
	).Columns(
		"user_id", "alias", "is_default", "provider", "config",
	).Values(
		uid,
		form.Alias,
		form.IsDefault,
		form.Provider,
		form.Config,
	).Suffix(
		"RETURNING " + strings.Join(outboundServiceColumns, ","),
	).ToSql()

	row := pg.db.QueryRow(query, args...)
	result, err := mapRowToOutboundService(row)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if result != nil && result.IsDefault {
		// Unset previous outboundService default
		err = pg.setDefaultOutboundService(result)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	return result, tx.Commit()
}

// UpdateOutboundServiceForUser update an outbound service of the DB
func (pg *DB) UpdateOutboundServiceForUser(uid uint, form model.OutboundServiceUpdateForm) (*model.OutboundService, error) {
	tx, err := pg.db.Begin()
	if err != nil {
		return nil, err
	}
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
	if form.IsDefault != nil {
		update["is_default"] = *form.IsDefault
	}
	query, args, err := pg.psql.Update(
		"outbound_services",
	).SetMap(update).Where(
		sq.Eq{"id": form.ID},
	).Where(
		sq.Eq{"user_id": uid},
	).Suffix(
		"RETURNING " + strings.Join(outboundServiceColumns, ","),
	).ToSql()

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	row := pg.db.QueryRow(query, args...)
	result, err := mapRowToOutboundService(row)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if result != nil && result.IsDefault {
		// Unset previous outbound service default
		err = pg.setDefaultOutboundService(result)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	return result, tx.Commit()
}

func (pg *DB) setDefaultOutboundService(outboundService *model.OutboundService) error {
	update := map[string]interface{}{
		"is_default": false,
	}
	query, args, _ := pg.psql.Update(
		"outbound_services",
	).SetMap(update).Where(
		sq.NotEq{"id": outboundService.ID},
	).Where(
		sq.Eq{"user_id": outboundService.UserID},
	).ToSql()

	_, err := pg.db.Exec(query, args...)
	return err
}

// GetOutboundServiceByID get an outbound service from the DB
func (pg *DB) GetOutboundServiceByID(id uint) (*model.OutboundService, error) {
	query, args, _ := pg.psql.Select(outboundServiceColumns...).From(
		"outbound_services",
	).Where(
		sq.Eq{"id": id},
	).ToSql()
	row := pg.db.QueryRow(query, args...)
	return mapRowToOutboundService(row)
}

// GetOutboundServiceByUserAndAlias get an outbound service from the DB.
// Default outbound service is returned if alias is nil.
func (pg *DB) GetOutboundServiceByUserAndAlias(uid uint, alias *string) (*model.OutboundService, error) {
	selectBuilder := pg.psql.Select(outboundServiceColumns...).From(
		"outbound_services",
	).Where(
		sq.Eq{"user_id": uid},
	)

	if alias != nil {
		selectBuilder = selectBuilder.Where(sq.Eq{"alias": *alias})
	} else {
		selectBuilder = selectBuilder.Where(sq.Eq{"is_default": true})
	}

	query, args, _ := selectBuilder.ToSql()
	row := pg.db.QueryRow(query, args...)
	return mapRowToOutboundService(row)
}

// GetOutboundServicesByUser returns outbound services of an user from DB
func (pg *DB) GetOutboundServicesByUser(uid uint) ([]model.OutboundService, error) {
	query, args, _ := pg.psql.Select(outboundServiceColumns...).From(
		"outbound_services",
	).Where(
		sq.Eq{"user_id": uid},
	).ToSql()
	rows, err := pg.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.OutboundService

	for rows.Next() {
		outboundService := model.OutboundService{}
		err = rows.Scan(
			&outboundService.ID,
			&outboundService.UserID,
			&outboundService.Alias,
			&outboundService.IsDefault,
			&outboundService.Provider,
			&outboundService.Config,
			&outboundService.CreatedAt,
			&outboundService.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, outboundService)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return result, nil
}

// DeleteOutboundServiceByUser removes an outbound service from the DB
func (pg *DB) DeleteOutboundServiceByUser(uid uint, id uint) error {
	query, args, _ := pg.psql.Delete("outbound_services").Where(
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
		return errors.New("no outbound service has been removed")
	}

	return nil
}

// DeleteOutboundServicesByUser removes outbound services from the DB
func (pg *DB) DeleteOutboundServicesByUser(uid uint, ids []uint) (int64, error) {
	query, args, _ := pg.psql.Delete("outbound_services").Where(
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
