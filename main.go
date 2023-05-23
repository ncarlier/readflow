package main

//go:generate go run generate.go
//go:generate gofmt -s -w autogen/db/postgres/db_sql_migration.go

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ncarlier/readflow/pkg/api"
	"github.com/ncarlier/readflow/pkg/cache"
	"github.com/ncarlier/readflow/pkg/config"
	configflag "github.com/ncarlier/readflow/pkg/config/flag"
	"github.com/ncarlier/readflow/pkg/db"
	"github.com/ncarlier/readflow/pkg/job"
	"github.com/ncarlier/readflow/pkg/logger"
	"github.com/ncarlier/readflow/pkg/metric"
	"github.com/ncarlier/readflow/pkg/service"
	"github.com/ncarlier/readflow/pkg/version"
	"github.com/rs/zerolog/log"
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: readflow OPTIONS\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}
}

func main() {
	// Get executable flags
	flags := config.Flags{}
	configflag.Bind(&flags, "READFLOW")

	// Parse command line (and environment variables)
	flag.Parse()

	// Show version if asked
	if *version.ShowVersion {
		version.Print()
		os.Exit(0)
	}

	// Init config file
	if config.InitConfigFile != nil && *config.InitConfigFile != "" {
		if err := config.WriteConfigFile(*config.InitConfigFile); err != nil {
			log.Fatal().Err(err).Msg("unable to init configuration file")
		}
		os.Exit(0)
	}

	conf := config.NewConfig()
	if flags.Config != "" {
		if err := conf.LoadFile(flags.Config); err != nil {
			log.Fatal().Err(err).Msg("unable to load configuration file")
		}
	}

	// Export configurations vars
	config.ExportVars(conf)

	// Configure the logger
	logger.Configure(flags.LogLevel, flags.LogPretty, conf.Integration.Sentry.DSN)

	log.Debug().Msg("starting readflow server...")

	// Configure the DB
	database, err := db.NewDB(conf.Global.DatabaseURI)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to configure the database")
	}

	// Configure download cache
	downloadCache, err := cache.NewDefault("readflow-downloads")
	if err != nil {
		log.Fatal().Err(err).Msg("unable to configure the cache storage")
	}

	// Configure the service registry
	err = service.Configure(*conf, database, downloadCache)
	if err != nil {
		database.Close()
		log.Fatal().Err(err).Msg("unable to configure the service registry")
	}

	// Start job scheduler
	scheduler := job.StartNewScheduler(database)

	server := &http.Server{
		Addr:    conf.Global.ListenAddr,
		Handler: api.NewRouter(conf),
	}

	var metricsServer *http.Server
	if conf.Global.MetricsListenAddr != "" {
		metricsServer = &http.Server{
			Addr:    conf.Global.MetricsListenAddr,
			Handler: api.NewMetricsRouter(),
		}
		metric.StartCollectors(database)
		go func() {
			log.Info().Str("listen", conf.Global.MetricsListenAddr).Msg("metrics server started")
			if err := metricsServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatal().Err(err).Str("listen", conf.Global.MetricsListenAddr).Msg("unable to start the metrics server")
			}
		}()
	}

	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-quit
		log.Debug().Msg("server is shutting down...")
		scheduler.Shutdown()
		api.Shutdown()

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			log.Fatal().Err(err).Msg("unable to gracefully shutdown the server")
		}
		if metricsServer != nil {
			metric.StopCollectors()
			if err := metricsServer.Shutdown(ctx); err != nil {
				log.Fatal().Err(err).Msg("unable to gracefully shutdown the metrics server")
			}
		}

		if err := downloadCache.Close(); err != nil {
			log.Error().Err(err).Msg("unable to gracefully shutdown the cache storage")
		}

		if err := database.Close(); err != nil {
			log.Fatal().Err(err).Msg("could not gracefully shutdown database connection")
		}

		close(done)
	}()

	api.Start()

	log.Info().Str("listen", conf.Global.ListenAddr).Msg("server started")

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal().Err(err).Str("listen", conf.Global.ListenAddr).Msg("unable to start the server")
	}

	<-done
	log.Debug().Msg("server stopped")
}
