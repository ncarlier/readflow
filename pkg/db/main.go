package db

import (
	"fmt"
	"net/url"

	"github.com/ncarlier/reader/pkg/db/postgres"
	"github.com/rs/zerolog/log"
)

// DB is the global database structure
type DB interface {
	Close() error
	UserRepository
	CategoryRepository
	ArticleRepository
}

// Configure the data store regarding the datasource URI
func Configure(conn string) (DB, error) {
	u, err := url.ParseRequestURI(conn)
	if err != nil {
		return nil, fmt.Errorf("invalid connection URL: %s", conn)
	}
	datastore := u.Scheme
	var db DB

	switch datastore {
	case "postgres":
		db, err = postgres.NewPostgreSQL(u)
		if err != nil {
			return nil, err
		}
		log.Info().Str("component", "store").Str("uri", u.String()).Msg("using PostgreSQL database")
	default:
		return nil, fmt.Errorf("unsuported database: %s", datastore)
	}
	return db, nil
}
