package postgres

import (
	"database/sql"
	"errors"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/ncarlier/readflow/pkg/model"
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

// CreateArchiverForUser creates an archiver into the DB
func (pg *DB) CreateArchiverForUser(uid uint, form model.ArchiverCreateForm) (*model.Archiver, error) {
	tx, err := pg.db.Begin()
	if err != nil {
		return nil, err
	}
	query, args, _ := pg.psql.Insert(
		"archivers",
	).Columns(
		"user_id", "alias", "is_default", "provider", "config",
	).Values(
		uid,
		form.Alias,
		form.IsDefault,
		form.Provider,
		form.Config,
	).Suffix(
		"RETURNING " + strings.Join(archiverColumns, ","),
	).ToSql()

	row := pg.db.QueryRow(query, args...)
	result, err := mapRowToArchiver(row)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if result != nil && result.IsDefault {
		// Unset previous archiver default
		err = pg.setDefaultArchiver(result)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	return result, tx.Commit()
}

// UpdateArchiverForUser update an archiver of the DB
func (pg *DB) UpdateArchiverForUser(uid uint, form model.ArchiverUpdateForm) (*model.Archiver, error) {
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
		"archivers",
	).SetMap(update).Where(
		sq.Eq{"id": form.ID},
	).Where(
		sq.Eq{"user_id": uid},
	).Suffix(
		"RETURNING " + strings.Join(archiverColumns, ","),
	).ToSql()

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	row := pg.db.QueryRow(query, args...)
	result, err := mapRowToArchiver(row)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if result != nil && result.IsDefault {
		// Unset previous archiver default
		err = pg.setDefaultArchiver(result)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	return result, tx.Commit()
}

func (pg *DB) setDefaultArchiver(archiver *model.Archiver) error {
	update := map[string]interface{}{
		"is_default": false,
	}
	query, args, _ := pg.psql.Update(
		"archivers",
	).SetMap(update).Where(
		sq.NotEq{"id": archiver.ID},
	).Where(
		sq.Eq{"user_id": archiver.UserID},
	).ToSql()

	_, err := pg.db.Exec(query, args...)
	return err
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

// GetArchiverByUserIDAndAlias get an archiver from the DB.
// Default archiver is returned if alias is nil.
func (pg *DB) GetArchiverByUserAndAlias(uid uint, alias *string) (*model.Archiver, error) {
	selectBuilder := pg.psql.Select(archiverColumns...).From(
		"archivers",
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
	return mapRowToArchiver(row)
}

// GetArchiversByUser returns archivers of an user from DB
func (pg *DB) GetArchiversByUser(uid uint) ([]model.Archiver, error) {
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

// DeleteArchiverByUser removes an archiver from the DB
func (pg *DB) DeleteArchiverByUser(uid uint, id uint) error {
	query, args, _ := pg.psql.Delete("archivers").Where(
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
		return errors.New("no archiver has been removed")
	}

	return nil
}

// DeleteArchiversByUser removes archivers from the DB
func (pg *DB) DeleteArchiversByUser(uid uint, ids []uint) (int64, error) {
	query, args, _ := pg.psql.Delete("archivers").Where(
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
