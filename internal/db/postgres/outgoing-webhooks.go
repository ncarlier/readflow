package postgres

import (
	"database/sql"
	"errors"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/ncarlier/readflow/internal/model"
)

var outgoingWebhookColumns = []string{
	"id",
	"user_id",
	"alias",
	"is_default",
	"provider",
	"config",
	"secrets",
	"created_at",
	"updated_at",
}

func mapRowToOutgoingWebhook(row *sql.Row) (*model.OutgoingWebhook, error) {
	result := &model.OutgoingWebhook{}

	err := row.Scan(
		&result.ID,
		&result.UserID,
		&result.Alias,
		&result.IsDefault,
		&result.Provider,
		&result.Config,
		&result.Secrets,
		&result.CreatedAt,
		&result.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, mapError(err)
	}
	return result, nil
}

// CreateOutgoingWebhookForUser creates an outgoing webhook into the DB
func (pg *DB) CreateOutgoingWebhookForUser(uid uint, form model.OutgoingWebhookCreateForm) (*model.OutgoingWebhook, error) {
	tx, err := pg.db.Begin()
	if err != nil {
		return nil, err
	}
	query, args, _ := pg.psql.Insert(
		"outgoing_webhooks",
	).Columns(
		"user_id", "alias", "is_default", "provider", "config", "secrets",
	).Values(
		uid,
		form.Alias,
		form.IsDefault,
		form.Provider,
		form.Config,
		form.Secrets,
	).Suffix(
		"RETURNING " + strings.Join(outgoingWebhookColumns, ","),
	).ToSql()

	row := pg.db.QueryRow(query, args...)
	result, err := mapRowToOutgoingWebhook(row)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if result != nil && result.IsDefault {
		// Unset default for other outgoing webhooks
		err = pg.unsetDefaultForOtherOutgoingWebhooks(result)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	return result, tx.Commit()
}

// UpdateOutgoingWebhookForUser update an outgoing webhook of the DB
func (pg *DB) UpdateOutgoingWebhookForUser(uid uint, form model.OutgoingWebhookUpdateForm) (*model.OutgoingWebhook, error) {
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
		"outgoing_webhooks",
	).SetMap(update).Where(
		sq.Eq{"id": form.ID},
	).Where(
		sq.Eq{"user_id": uid},
	).Suffix(
		"RETURNING " + strings.Join(outgoingWebhookColumns, ","),
	).ToSql()

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	row := pg.db.QueryRow(query, args...)
	result, err := mapRowToOutgoingWebhook(row)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if result != nil {
		if form.Secrets != nil {
			// Update secrets
			err = pg.updateOutgoingWebhookSecrets(result, form.Secrets)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
		}
		if result.IsDefault {
			// Unset default for other outgoing webhooks
			err = pg.unsetDefaultForOtherOutgoingWebhooks(result)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}
	return result, tx.Commit()
}

func (pg *DB) unsetDefaultForOtherOutgoingWebhooks(webhook *model.OutgoingWebhook) error {
	update := map[string]interface{}{
		"is_default": false,
	}
	query, args, _ := pg.psql.Update(
		"outgoing_webhooks",
	).SetMap(update).Where(
		sq.NotEq{"id": webhook.ID},
	).Where(
		sq.Eq{"user_id": webhook.UserID},
	).ToSql()

	_, err := pg.db.Exec(query, args...)
	return err
}

// GetOutgoingWebhookByID get an outgoing webhook from the DB
func (pg *DB) GetOutgoingWebhookByID(id uint) (*model.OutgoingWebhook, error) {
	query, args, _ := pg.psql.Select(outgoingWebhookColumns...).From(
		"outgoing_webhooks",
	).Where(
		sq.Eq{"id": id},
	).ToSql()
	row := pg.db.QueryRow(query, args...)
	return mapRowToOutgoingWebhook(row)
}

// GetOutgoingWebhookByUserAndAlias get an outgoing webhook from the DB.
// Default outgoing webhook is returned if alias is nil.
func (pg *DB) GetOutgoingWebhookByUserAndAlias(uid uint, alias *string) (*model.OutgoingWebhook, error) {
	selectBuilder := pg.psql.Select(outgoingWebhookColumns...).From(
		"outgoing_webhooks",
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
	return mapRowToOutgoingWebhook(row)
}

// GetOutgoingWebhooksByUser returns outgoing webhooks of an user from DB
func (pg *DB) GetOutgoingWebhooksByUser(uid uint) ([]model.OutgoingWebhook, error) {
	query, args, _ := pg.psql.Select(outgoingWebhookColumns...).From(
		"outgoing_webhooks",
	).Where(
		sq.Eq{"user_id": uid},
	).ToSql()
	rows, err := pg.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.OutgoingWebhook

	for rows.Next() {
		OutgoingWebhook := model.OutgoingWebhook{}
		err = rows.Scan(
			&OutgoingWebhook.ID,
			&OutgoingWebhook.UserID,
			&OutgoingWebhook.Alias,
			&OutgoingWebhook.IsDefault,
			&OutgoingWebhook.Provider,
			&OutgoingWebhook.Config,
			&OutgoingWebhook.Secrets,
			&OutgoingWebhook.CreatedAt,
			&OutgoingWebhook.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, OutgoingWebhook)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return result, nil
}

// CountOutgoingWebhooksByUser returns total nb of outgoing webhooks of an user from the DB
func (pg *DB) CountOutgoingWebhooksByUser(uid uint) (uint, error) {
	counter := pg.psql.Select("count(*)").From(
		"outgoing_webhooks",
	).Where(sq.Eq{"user_id": uid})
	query, args, _ := counter.ToSql()

	var count uint
	if err := pg.db.QueryRow(query, args...).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

// DeleteOutgoingWebhookByUser removes an outgoing webhook from the DB
func (pg *DB) DeleteOutgoingWebhookByUser(uid, id uint) error {
	query, args, _ := pg.psql.Delete("outgoing_webhooks").Where(
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
		return errors.New("no outgoing webhook has been removed")
	}

	return nil
}

// DeleteOutgoingWebhooksByUser removes outgoing webhooks from the DB
func (pg *DB) DeleteOutgoingWebhooksByUser(uid uint, ids []uint) (int64, error) {
	query, args, _ := pg.psql.Delete("outgoing_webhooks").Where(
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
