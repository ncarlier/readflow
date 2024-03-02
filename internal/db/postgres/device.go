package postgres

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/ncarlier/readflow/internal/model"
)

var deviceColumns = []string{
	"id",
	"user_id",
	"key",
	"subscription",
	"last_seen_at",
	"created_at",
}

func mapRowToDevice(row *sql.Row) (*model.Device, error) {
	device := &model.Device{}

	sub := ""

	err := row.Scan(
		&device.ID,
		&device.UserID,
		&device.Key,
		&sub,
		&device.LastSeenAt,
		&device.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, mapError(err)
	}
	if err := device.SetSubscription(sub); err != nil {
		return nil, err
	}
	return device, nil
}

// CreateDevice create a device
func (pg *DB) CreateDevice(device model.Device) (*model.Device, error) {
	dev, err := pg.GetDeviceByUserAndKey(*device.UserID, device.Key)
	if err != nil || dev != nil {
		return dev, err
	}
	sub, err := device.GetSubscription()
	if err != nil {
		return nil, err
	}
	query, args, _ := pg.psql.Insert(
		"devices",
	).Columns(
		"user_id", "key", "subscription", "last_seen_at",
	).Values(
		device.UserID,
		device.Key,
		sub,
		"NOW()",
	).Suffix(
		"RETURNING " + strings.Join(deviceColumns, ","),
	).ToSql()

	row := pg.db.QueryRow(query, args...)
	return mapRowToDevice(row)
}

// GetDeviceByID get a device from the DB
func (pg *DB) GetDeviceByID(id uint) (*model.Device, error) {
	// Update last seen attribute then return the device
	query, args, _ := pg.psql.Update(
		"devices",
	).Set(
		"last_seen_at", "now()",
	).Where(
		sq.Eq{"id": id},
	).Suffix(
		"RETURNING " + strings.Join(deviceColumns, ","),
	).ToSql()
	row := pg.db.QueryRow(query, args...)
	return mapRowToDevice(row)
}

// GetDeviceByUserAndKey get an device from the DB
// Only exposed for testing purpose!
func (pg *DB) GetDeviceByUserAndKey(uid uint, key string) (*model.Device, error) {
	query, args, _ := pg.psql.Select(deviceColumns...).From(
		"devices",
	).Where(
		sq.Eq{"user_id": uid},
	).Where(
		sq.Eq{"key": key},
	).ToSql()

	row := pg.db.QueryRow(query, args...)
	return mapRowToDevice(row)
}

// GetDevicesByUser returns devices of an user from DB
func (pg *DB) GetDevicesByUser(uid uint) ([]model.Device, error) {
	query, args, _ := pg.psql.Select(deviceColumns...).From(
		"devices",
	).Where(
		sq.Eq{"user_id": uid},
	).ToSql()
	rows, err := pg.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.Device

	for rows.Next() {
		device := model.Device{}
		sub := ""
		err = rows.Scan(
			&device.ID,
			&device.UserID,
			&device.Key,
			&sub,
			&device.LastSeenAt,
			&device.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		// Ignore bad subscriptions
		if err = device.SetSubscription(sub); err == nil {
			result = append(result, device)
		}
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return result, nil
}

// CountDevicesByUser returns total nb of devices of an user from the DB
func (pg *DB) CountDevicesByUser(uid uint) (uint, error) {
	counter := pg.psql.Select("count(*)").From(
		"devices",
	).Where(sq.Eq{"user_id": uid})
	query, args, _ := counter.ToSql()

	var count uint
	if err := pg.db.QueryRow(query, args...).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

// DeleteDevice removes an device from the DB
func (pg *DB) DeleteDevice(id uint) error {
	query, args, _ := pg.psql.Delete("devices").Where(
		sq.Eq{"id": id},
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
		return errors.New("no device has been removed")
	}

	return nil
}

// DeleteDevicesByUser removes devices from the DB
func (pg *DB) DeleteDevicesByUser(uid uint, ids []uint) (int64, error) {
	query, args, _ := pg.psql.Delete("devices").Where(
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

// DeleteInactiveDevicesOlderThan remove inactive devices from the DB
func (pg *DB) DeleteInactiveDevicesOlderThan(delay time.Duration) (int64, error) {
	maxAge := time.Now().Add(-delay)
	query, args, _ := pg.psql.Delete(
		"devices",
	).Where(
		sq.Lt{"last_seen_at": maxAge},
	).ToSql()

	result, err := pg.db.Exec(query, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
