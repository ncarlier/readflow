package postgres

import (
	"database/sql"
	"errors"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/ncarlier/reader/pkg/model"
)

var archiverColumns = []string{
	"id",
	"user_id",
	"alias",
	"is_default",
	"provider",
	"config",
	"created_at",
	"updated_at",
}

func mapRowToArchiver(row *sql.Row) (*model.Archiver, error) {
	archiver := &model.Archiver{}

	err := row.Scan(
		&archiver.ID,
		&archiver.UserID,
		&archiver.Alias,
		&archiver.IsDefault,
		&archiver.Provider,
		&archiver.Config,
		&archiver.CreatedAt,
		&archiver.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return archiver, nil
}

func (pg *DB) createArchiver(archiver model.Archiver) (*model.Archiver, error) {
	query, args, _ := pg.psql.Insert(
		"archivers",
	).Columns(
		"user_id", "alias", "is_default", "provider", "config",
	).Values(
		archiver.UserID,
		archiver.Alias,
		archiver.IsDefault,
		archiver.Provider,
		archiver.Config,
	).Suffix(
		"RETURNING " + strings.Join(archiverColumns, ","),
	).ToSql()

	row := pg.db.QueryRow(query, args...)
	return mapRowToArchiver(row)
}

func (pg *DB) updateArchiver(archiver model.Archiver) (*model.Archiver, error) {
	update := map[string]interface{}{
		"alias":      archiver.Alias,
		"is_default": archiver.IsDefault,
		"provider":   archiver.Provider,
		"config":     archiver.Config,
		"updated_at": "NOW()",
	}
	query, args, _ := pg.psql.Update(
		"archivers",
	).SetMap(update).Where(
		sq.Eq{"id": archiver.ID},
	).Where(
		sq.Eq{"user_id": archiver.UserID},
	).Suffix(
		"RETURNING " + strings.Join(archiverColumns, ","),
	).ToSql()

	row := pg.db.QueryRow(query, args...)
	return mapRowToArchiver(row)
}

// CreateOrUpdateArchiver creates or updates an archiver into the DB
func (pg *DB) CreateOrUpdateArchiver(archiver model.Archiver) (*model.Archiver, error) {
	if archiver.ID != nil {
		return pg.updateArchiver(archiver)
	}
	return pg.createArchiver(archiver)
}

// GetArchiverByID get an archiver from the DB
func (pg *DB) GetArchiverByID(id uint) (*model.Archiver, error) {
	query, args, _ := pg.psql.Select(archiverColumns...).From(
		"archivers",
	).Where(
		sq.Eq{"id": id},
	).ToSql()
	row := pg.db.QueryRow(query, args...)
	return mapRowToArchiver(row)
}

// GetArchiverByUserIDAndAlias get an archiver from the DB
func (pg *DB) GetArchiverByUserIDAndAlias(uid uint, alias string) (*model.Archiver, error) {
	query, args, _ := pg.psql.Select(archiverColumns...).From(
		"archivers",
	).Where(
		sq.Eq{"user_id": uid},
		sq.Eq{"alias": alias},
	).ToSql()
	row := pg.db.QueryRow(query, args...)
	return mapRowToArchiver(row)
}

// GetArchiversByUserID returns archivers of an user from DB
func (pg *DB) GetArchiversByUserID(uid uint) ([]model.Archiver, error) {
	query, args, _ := pg.psql.Select(archiverColumns...).From(
		"archivers",
	).Where(
		sq.Eq{"user_id": uid},
	).ToSql()
	rows, err := pg.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.Archiver

	for rows.Next() {
		archiver := model.Archiver{}
		err = rows.Scan(
			&archiver.ID,
			&archiver.UserID,
			&archiver.Alias,
			&archiver.IsDefault,
			&archiver.Provider,
			&archiver.Config,
			&archiver.CreatedAt,
			&archiver.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, archiver)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return result, nil
}

// DeleteArchiver removes an archiver from the DB
func (pg *DB) DeleteArchiver(archiver model.Archiver) error {
	query, args, _ := pg.psql.Delete("archivers").Where(
		sq.Eq{"id": archiver.ID},
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
		return errors.New("no archiver has been removed")
	}

	return nil
}
