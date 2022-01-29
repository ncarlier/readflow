package service

import (
	"github.com/ncarlier/readflow/pkg/cache"
	"github.com/ncarlier/readflow/pkg/config"
	"github.com/ncarlier/readflow/pkg/constant"
	"github.com/ncarlier/readflow/pkg/db"
	"github.com/ncarlier/readflow/pkg/exporter"
	_ "github.com/ncarlier/readflow/pkg/exporter/all"
	"github.com/ncarlier/readflow/pkg/helper"
	"github.com/ncarlier/readflow/pkg/model"
	ratelimiter "github.com/ncarlier/readflow/pkg/rate-limiter"
	ruleengine "github.com/ncarlier/readflow/pkg/rule-engine"
	"github.com/ncarlier/readflow/pkg/scraper"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var instance *Registry

// Registry is the structure definition of the service registry
type Registry struct {
	conf                    config.Config
	db                      db.DB
	logger                  zerolog.Logger
	ruleEngineCache         *ruleengine.Cache
	downloadCache           cache.Cache
	properties              *model.Properties
	webScraper              scraper.WebScraper
	downloader              exporter.Downloader
	hashid                  *helper.HashIDHandler
	notificationRateLimiter ratelimiter.RateLimiter
}

// Configure the global service registry
func Configure(conf config.Config, database db.DB, downloadCache cache.Cache) error {
	downloader := exporter.NewInternalDownloader(downloadCache, 10, constant.DefaultTimeout)
	webScraper, err := scraper.NewWebScraper(conf.Integration.ExternalWebScraperURL)
	if err != nil {
		return err
	}
	hashid, err := helper.NewHashIDHandler(conf.Global.SecretSalt)
	if err != nil {
		return err
	}
	notificationRateLimiter, err := ratelimiter.NewRateLimiter("notification", conf.RateLimiting.Notification)
	if err != nil {
		return err
	}
	instance = &Registry{
		conf:                    conf,
		db:                      database,
		logger:                  log.With().Str("component", "service").Logger(),
		ruleEngineCache:         ruleengine.NewRuleEngineCache(1024),
		downloadCache:           downloadCache,
		webScraper:              webScraper,
		downloader:              downloader,
		hashid:                  hashid,
		notificationRateLimiter: notificationRateLimiter,
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
