package postgres

import (
	"database/sql"
	"errors"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/ncarlier/readflow/internal/model"
)

var inboundServiceColumns = []string{
	"id",
	"user_id",
	"alias",
	"token",
	"script",
	"last_usage_at",
	"created_at",
	"updated_at",
}

func mapRowToIncomingWebhook(row *sql.Row) (*model.IncomingWebhook, error) {
	result := model.IncomingWebhook{}

	err := row.Scan(
		&result.ID,
		&result.UserID,
		&result.Alias,
		&result.Token,
		&result.Script,
		&result.LastUsageAt,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, mapError(err)
	}
	return &result, nil
}

// CreateIncomingWebhookForUser creates an incoming webhook for a user
func (pg *DB) CreateIncomingWebhookForUser(uid uint, form model.IncomingWebhookCreateForm) (*model.IncomingWebhook, error) {
	query, args, _ := pg.psql.Insert(
		"incoming_webhooks",
	).Columns(
		"user_id", "alias", "token", "script",
	).Values(
		uid,
		form.Alias,
		form.Token,
		form.Script,
	).Suffix(
		"RETURNING " + strings.Join(inboundServiceColumns, ","),
	).ToSql()
	row := pg.db.QueryRow(query, args...)
	return mapRowToIncomingWebhook(row)
}

// UpdateIncomingWebhookForUser updates an incoming webhook for a user
func (pg *DB) UpdateIncomingWebhookForUser(uid uint, form model.IncomingWebhookUpdateForm) (*model.IncomingWebhook, error) {
	update := map[string]interface{}{
		"updated_at": "NOW()",
	}
	if form.Alias != nil {
		update["alias"] = *form.Alias
	}
	if form.Script != nil {
		update["script"] = *form.Script
	}
	query, args, _ := pg.psql.Update(
		"incoming_webhooks",
	).SetMap(update).Where(
		sq.Eq{"id": form.ID},
	).Where(
		sq.Eq{"user_id": uid},
	).Suffix(
		"RETURNING " + strings.Join(inboundServiceColumns, ","),
	).ToSql()

	row := pg.db.QueryRow(query, args...)
	return mapRowToIncomingWebhook(row)
}

// GetIncomingWebhookByID get an incoming webhook from the DB
func (pg *DB) GetIncomingWebhookByID(id uint) (*model.IncomingWebhook, error) {
	query, args, _ := pg.psql.Select(inboundServiceColumns...).From(
		"incoming_webhooks",
	).Where(
		sq.Eq{"id": id},
	).ToSql()
	row := pg.db.QueryRow(query, args...)
	return mapRowToIncomingWebhook(row)
}

// GetIncomingWebhookByToken find an incoming webhook by token form the DB (last usage is updated!)
func (pg *DB) GetIncomingWebhookByToken(token string) (*model.IncomingWebhook, error) {
	query, args, _ := pg.psql.Update(
		"incoming_webhooks",
	).Set("last_usage_at", "NOW()").Where(
		sq.Eq{"token": token},
	).Suffix(
		"RETURNING " + strings.Join(inboundServiceColumns, ","),
	).ToSql()

	row := pg.db.QueryRow(query, args...)
	return mapRowToIncomingWebhook(row)
}

// GetIncomingWebhookByUserAndAlias returns incoming webhook of an user by its alias
func (pg *DB) GetIncomingWebhookByUserAndAlias(uid uint, alias string) (*model.IncomingWebhook, error) {
	query, args, _ := pg.psql.Select(inboundServiceColumns...).From(
		"incoming_webhooks",
	).Where(
		sq.Eq{"user_id": uid},
	).Where(
		sq.Eq{"alias": alias},
	).ToSql()

	row := pg.db.QueryRow(query, args...)
	return mapRowToIncomingWebhook(row)
}

// GetIncomingWebhooksByUser returns incoming webhooks of an user from DB
func (pg *DB) GetIncomingWebhooksByUser(uid uint) ([]model.IncomingWebhook, error) {
	query, args, _ := pg.psql.Select(inboundServiceColumns...).From(
		"incoming_webhooks",
	).Where(
		sq.Eq{"user_id": uid},
	).OrderBy("alias ASC").ToSql()
	rows, err := pg.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.IncomingWebhook

	for rows.Next() {
		item := model.IncomingWebhook{}
		err = rows.Scan(
			&item.ID,
			&item.UserID,
			&item.Alias,
			&item.Token,
			&item.Script,
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

// CountIncomingWebhooksByUser returns total nb of incoming webhooks of an user from the DB
func (pg *DB) CountIncomingWebhooksByUser(uid uint) (uint, error) {
	counter := pg.psql.Select("count(*)").From(
		"incoming_webhooks",
	).Where(sq.Eq{"user_id": uid})
	query, args, _ := counter.ToSql()

	var count uint
	if err := pg.db.QueryRow(query, args...).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

// DeleteIncomingWebhookByUser removes an inboundService from the DB
func (pg *DB) DeleteIncomingWebhookByUser(uid, id uint) error {
	query, args, _ := pg.psql.Delete("incoming_webhooks").Where(
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
		return errors.New("no incoming webhook been removed")
	}

	return nil
}

// DeleteIncomingWebhooksByUser removes incoming webhooks from the DB
func (pg *DB) DeleteIncomingWebhooksByUser(uid uint, ids []uint) (int64, error) {
	query, args, _ := pg.psql.Delete("incoming_webhooks").Where(
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
