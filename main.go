package main

//go:generate go run generate.go
//go:generate gofmt -s -w autogen/db/postgres/db_sql_migration.go

import (
	"context"
	"flag"
	"fmt"
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
	"github.com/ncarlier/readflow/pkg/server"
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
	// get executable flags
	flags := config.Flags{}
	configflag.Bind(&flags, "READFLOW")

	// parse command line (and environment variables)
	flag.Parse()

	// show version if asked
	if *version.ShowVersion {
		version.Print()
		os.Exit(0)
	}

	// init config file
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

	// export configurations vars
	config.ExportVars(conf)

	// configure the logger
	logger.Configure(flags.LogLevel, flags.LogPretty, conf.Integration.Sentry.DSN)

	log.Debug().Msg("starting readflow...")

	// configure the DB
	database, err := db.NewDB(conf.Global.DatabaseURI)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to configure the database")
	}

	// configure download cache
	downloadCache, err := cache.NewDefault("readflow-downloads")
	if err != nil {
		log.Fatal().Err(err).Msg("unable to configure the cache storage")
	}

	// configure the service registry
	err = service.Configure(*conf, database, downloadCache)
	if err != nil {
		database.Close()
		log.Fatal().Err(err).Msg("unable to configure the service registry")
	}

	// ctart job scheduler
	scheduler := job.StartNewScheduler(database)

	// create HTTP server
	httpServer := server.NewHTTPServer(conf)

	// create and start metrics server
	metricsServer := server.NewMetricsServer(conf)
	if metricsServer != nil {
		metric.StartCollectors(database)
		go metricsServer.ListenAndServe()
	}

	// create and start SMTP server
	smtpServer := server.NewSMTPHTTPServer(conf)
	if smtpServer != nil {
		go smtpServer.ListenAndServe()
	}

	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-quit
		log.Debug().Msg("shutting down readflow...")
		scheduler.Shutdown()
		api.Shutdown()

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := httpServer.Shutdown(ctx); err != nil {
			log.Fatal().Err(err).Msg("unable to gracefully shutdown the HTTP server")
		}
		if smtpServer != nil {
			if err := smtpServer.Shutdown(ctx); err != nil {
				log.Fatal().Err(err).Msg("unable to gracefully shutdown the SMTP server")
			}
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

	// set API health check as started
	api.Start()

	// start HTTP server
	httpServer.ListenAndServe()

	<-done
	log.Debug().Msg("readflow stopped")
}
