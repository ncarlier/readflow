package service

import (
	"github.com/ncarlier/readflow/pkg/config"
	"github.com/ncarlier/readflow/pkg/db"
	"github.com/ncarlier/readflow/pkg/model"
	ruleengine "github.com/ncarlier/readflow/pkg/rule-engine"
	"github.com/ncarlier/readflow/pkg/scraper"
	userplan "github.com/ncarlier/readflow/pkg/user-plan"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var instance *Registry

// Registry is the structure definition of the service registry
type Registry struct {
	db              db.DB
	UserPlans       userplan.UserPlans
	logger          zerolog.Logger
	ruleEngineCache *ruleengine.Cache
	properties      *model.Properties
	webScraper      scraper.WebScraper
}

// GetProperties retrieve service properties
func (reg *Registry) GetProperties() model.Properties {
	return *reg.properties
}

// Configure the global service registry
func Configure(conf config.Config, database db.DB, plans userplan.UserPlans) error {
	webScraper, err := scraper.NewWebScraper(conf.WebScraping)
	if err != nil {
		return err
	}
	instance = &Registry{
		db:              database,
		UserPlans:       plans,
		logger:          log.With().Str("component", "service").Logger(),
		ruleEngineCache: ruleengine.NewRuleEngineCache(1024),
		webScraper:      webScraper,
	}
	return instance.initProperties()
}

// Lookup returns the global service registry
func Lookup() *Registry {
	if instance != nil {
		return instance
	}
	log.Fatal().Msg("service registry not configured")
	return nil
}
