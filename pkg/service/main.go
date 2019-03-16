package service

import (
	"github.com/ncarlier/reader/pkg/db"
	"github.com/rs/zerolog/log"
)

var instance *Registry

// Registry is the structure definition of the service registry
type Registry struct {
	db db.DB
}

// InitRegistry init the service registry
func InitRegistry(_db db.DB) {
	instance = &Registry{
		db: _db,
	}
}

// Lookup returns the global service registry
func Lookup() *Registry {
	if instance != nil {
		return instance
	}
	log.Fatal().Msg("Service registry not initialized")
	return nil
}
