package db

import (
	"fmt"
	"net/url"

	"github.com/ncarlier/readflow/pkg/db/postgres"
	"github.com/rs/zerolog/log"
)

// DB is the global database structure
type DB interface {
	Close() error
	UserRepository
	CategoryRepository
	ArticleRepository
	APIKeyRepository
	ArchiverRepository
	DeviceRepository
	PropertiesRepository
}

// NewDB create new database provider regarding the datasource URI
func NewDB(conn string) (DB, error) {
	u, err := url.ParseRequestURI(conn)
	if err != nil {
		return nil, fmt.Errorf("invalid connection URL: %s", conn)
	}
	provider := u.Scheme
	var db DB

	switch provider {
	case "postgres":
		db, err = postgres.NewPostgreSQL(u)
		if err != nil {
			return nil, err
		}
		log.Info().Str("component", "database").Str("uri", u.String()).Msg("using PostgreSQL database")
	default:
		return nil, fmt.Errorf("unsupported database: %s", provider)
	}
	return db, nil
}
