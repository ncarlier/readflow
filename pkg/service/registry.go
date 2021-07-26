package service

import (
	"github.com/ncarlier/readflow/pkg/archiver"
	"github.com/ncarlier/readflow/pkg/cache"
	"github.com/ncarlier/readflow/pkg/config"
	"github.com/ncarlier/readflow/pkg/constant"
	"github.com/ncarlier/readflow/pkg/db"
	"github.com/ncarlier/readflow/pkg/helper"
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
	conf            config.Config
	db              db.DB
	UserPlans       userplan.UserPlans
	logger          zerolog.Logger
	ruleEngineCache *ruleengine.Cache
	downloadCache   cache.Cache
	properties      *model.Properties
	webScraper      scraper.WebScraper
	webArchiver     *archiver.WebArchiver
	hashid          *helper.HashIDHandler
}

// Configure the global service registry
func Configure(conf config.Config, database db.DB, downloadCache cache.Cache, plans userplan.UserPlans) error {
	webArchiver := archiver.NewWebArchiver(downloadCache, 10, constant.DefaultTimeout)
	webScraper, err := scraper.NewWebScraper(conf.WebScraping)
	if err != nil {
		return err
	}
	hashid, err := helper.NewHashIDHandler(conf.SecretSalt)
	if err != nil {
		return err
	}
	instance = &Registry{
		conf:            conf,
		db:              database,
		UserPlans:       plans,
		logger:          log.With().Str("component", "service").Logger(),
		ruleEngineCache: ruleengine.NewRuleEngineCache(1024),
		downloadCache:   downloadCache,
		webScraper:      webScraper,
		webArchiver:     webArchiver,
		hashid:          hashid,
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
