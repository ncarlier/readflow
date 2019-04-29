package main

//go:generate go run generate.go
//go:generate gofmt -s -w autogen/db/postgres/db_sql_migration.go

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ncarlier/readflow/pkg/job"

	eventbroker "github.com/ncarlier/readflow/pkg/event-broker"
	_ "github.com/ncarlier/readflow/pkg/event-listener"
	"github.com/ncarlier/readflow/pkg/service"

	"github.com/ncarlier/readflow/pkg/db"

	"github.com/ncarlier/readflow/pkg/api"
	"github.com/ncarlier/readflow/pkg/config"
	"github.com/ncarlier/readflow/pkg/logger"
	"github.com/ncarlier/readflow/pkg/version"
	"github.com/rs/zerolog/log"
)

func main() {
	flag.Parse()

	conf := config.Get()

	if *conf.Version {
		version.Print()
		return
	}

	// Configure the logger
	logger.Configure(*conf.LogLevel, *conf.LogPretty, *conf.SentryDSN)

	log.Debug().Msg("starting Nunux Reader server...")

	// Configure the DB
	_db, err := db.Configure(*conf.DB)
	if err != nil {
		log.Fatal().Err(err).Msg("could not configure Database")
	}

	// Configure Event Broker
	_, err = eventbroker.Configure(*conf.Broker)
	if err != nil {
		log.Fatal().Err(err).Msg("could not configure Event Broker")
	}

	// Init service registry
	err = service.InitRegistry(_db)
	if err != nil {
		log.Fatal().Err(err).Msg("could not init service registry")
	}

	// Start job scheduler
	scheduler := job.StartNewScheduler(_db)

	server := &http.Server{
		Addr:    *conf.ListenAddr,
		Handler: api.NewRouter(conf),
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
			log.Fatal().Err(err).Msg("could not gracefully shutdown the server")
		}
		close(done)
	}()

	api.Start()

	log.Info().Str("listen", *conf.ListenAddr).Msg("server started")

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal().Err(err).Str("listen", *conf.ListenAddr).Msg("could not start the server")
	}

	<-done
	log.Debug().Msg("server stopped")
}
