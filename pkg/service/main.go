package service

import (
	"github.com/ncarlier/readflow/pkg/db"
	"github.com/ncarlier/readflow/pkg/model"
	ruleengine "github.com/ncarlier/readflow/pkg/rule-engine"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var instance *Registry

// Registry is the structure definition of the service registry
type Registry struct {
	db              db.DB
	logger          zerolog.Logger
	ruleEngineCache *ruleengine.Cache
	properties      *model.Properties
}

// GetProperties retieve service properties
func (reg *Registry) GetProperties() model.Properties {
	return *reg.properties
}

// InitRegistry init the service registry
func InitRegistry(_db db.DB) error {
	instance = &Registry{
		db:              _db,
		logger:          log.With().Str("component", "service").Logger(),
		ruleEngineCache: ruleengine.NewRuleEngineCache(1024),
	}
	return instance.initProperties()
}

// Lookup returns the global service registry
func Lookup() *Registry {
	if instance != nil {
		return instance
	}
	log.Fatal().Msg("service registry not initialized")
	return nil
}
