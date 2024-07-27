package service

import (

	// load all exporters
	"github.com/ncarlier/readflow/internal/config"
	"github.com/ncarlier/readflow/internal/db"
	_ "github.com/ncarlier/readflow/internal/exporter/all"
	"github.com/ncarlier/readflow/internal/model"
	"github.com/ncarlier/readflow/internal/scripting"
	"github.com/ncarlier/readflow/pkg/cache"
	"github.com/ncarlier/readflow/pkg/defaults"
	"github.com/ncarlier/readflow/pkg/downloader"
	"github.com/ncarlier/readflow/pkg/event"
	"github.com/ncarlier/readflow/pkg/event/dispatcher"
	"github.com/ncarlier/readflow/pkg/hashid"
	"github.com/ncarlier/readflow/pkg/job"
	"github.com/ncarlier/readflow/pkg/logger"
	ratelimiter "github.com/ncarlier/readflow/pkg/rate-limiter"
	"github.com/ncarlier/readflow/pkg/sanitizer"
	"github.com/ncarlier/readflow/pkg/scraper"
	"github.com/ncarlier/readflow/pkg/secret"
	"github.com/rs/zerolog"
)

var instance *Registry

// Registry is the structure definition of the service registry
type Registry struct {
	conf                    config.Config
	db                      db.DB
	logger                  zerolog.Logger
	downloadCache           cache.Cache
	properties              *model.Properties
	webScraper              scraper.WebScraper
	dl                      downloader.Downloader
	hashid                  *hashid.HashIDHandler
	notificationRateLimiter ratelimiter.RateLimiter
	scriptEngine            *scripting.ScriptEngine
	sanitizer               *sanitizer.Sanitizer
	events                  *event.Manager
	scheduler               *job.Scheduler
	dispatcher              dispatcher.Dispatcher
	secretsEngineProvider   secret.EngineProvider
}

// Configure the global service registry
func Configure(conf config.Config, database db.DB) error {
	// configure download cache
	downloadCache, err := cache.New(conf.Downloader.Cache)
	if err != nil {
		return err
	}
	webScraper, err := scraper.NewWebScraper(defaults.HTTPClient, defaults.UserAgent, conf.Scraping.ServiceProvider)
	if err != nil {
		return err
	}
	hid, err := hashid.NewHashIDHandler(conf.Hash.SecretSalt.Value)
	if err != nil {
		return err
	}
	notificationRateLimiter, err := ratelimiter.NewRateLimiter("notification", &conf.RateLimiting.Notification)
	if err != nil {
		return err
	}
	blockList, err := sanitizer.NewBlockList(conf.Scraping.BlockList, sanitizer.DefaultBlockList)
	if err != nil {
		return err
	}
	_dispatcher, err := dispatcher.NewDispatcher(conf.Event.BrokerURI)
	if err != nil {
		return err
	}
	secretsEngineProvider, err := secret.NewSecretsEngineProvider(conf.Secrets.ServiceProvider)
	if err != nil {
		return err
	}
	scheduler := job.NewScheduler(
		db.NewCleanupDatabaseJob(database),
	)

	dl := downloader.NewInternalDownloader(
		defaults.HTTPClient,
		conf.Downloader.UserAgent,
		downloadCache,
		conf.Downloader.MaxConcurentDownloads,
		conf.Downloader.Timeout.Duration,
	)

	instance = &Registry{
		conf:                    conf,
		db:                      database,
		logger:                  logger.With().Str("component", "service").Logger(),
		downloadCache:           downloadCache,
		webScraper:              webScraper,
		dl:                      dl,
		hashid:                  hid,
		notificationRateLimiter: notificationRateLimiter,
		sanitizer:               sanitizer.NewSanitizer(blockList),
		scriptEngine:            scripting.NewScriptEngine(128),
		dispatcher:              _dispatcher,
		events:                  event.NewEventManager(),
		scheduler:               scheduler,
		secretsEngineProvider:   secretsEngineProvider,
	}
	instance.registerEventHandlers()
	return instance.initProperties()
}

func (reg *Registry) GetConfig() config.Config {
	return reg.conf
}

// Shutdown service internals jobs
func Shutdown() {
	if instance == nil {
		return
	}
	instance.scheduler.Shutdown()
	if err := instance.downloadCache.Close(); err != nil {
		instance.logger.Error().Err(err).Msg("unable to gracefully shutdown the cache storage")
	}
}

// Lookup returns the global service registry
func Lookup() *Registry {
	if instance != nil {
		return instance
	}
	logger.Fatal().Msg("service registry not configured")
	return nil
}
