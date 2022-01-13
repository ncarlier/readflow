package postgres

import (
	"database/sql"
	"strconv"

	migration "github.com/ncarlier/readflow/autogen/db/postgres"
	"github.com/rs/zerolog/log"
)

const schemaVersion = 11

// Migrate executes database migrations.
func Migrate(db *sql.DB) {
	var currentVersion int
	db.QueryRow(`select version from schema_version`).Scan(&currentVersion)

	log.Debug().Int("current", currentVersion).Int("latest", schemaVersion).Msg("Database version")

	for version := currentVersion + 1; version <= schemaVersion; version++ {
		log.Warn().Int("version", version).Msg("Migrating Database...")

		tx, err := db.Begin()
		if err != nil {
			log.Fatal().Err(err).Msg("Unable to init Database session")
		}

		rawSQL := migration.DatabaseSQLMigration["db_migration_"+strconv.Itoa(version)]
		_, err = tx.Exec(rawSQL)
		if err != nil {
			tx.Rollback()
			log.Fatal().Err(err).Msg("Unable to apply migration")
		}

		if _, err := tx.Exec(`delete from schema_version`); err != nil {
			tx.Rollback()
			log.Fatal().Err(err).Msg("Unable to update schema version number")
		}

		if _, err := tx.Exec(`insert into schema_version (version) values($1)`, version); err != nil {
			tx.Rollback()
			log.Fatal().Err(err).Msg("Unable to update schema version number")
		}

		if err := tx.Commit(); err != nil {
			log.Fatal().Err(err).Msg("Unable to apply migration on COMMIT")
		}
	}
}
