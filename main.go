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

	eventbroker "github.com/ncarlier/reader/pkg/event-broker"
	"github.com/ncarlier/reader/pkg/service"

	"github.com/ncarlier/reader/pkg/db"

	"github.com/ncarlier/reader/pkg/api"
	"github.com/ncarlier/reader/pkg/config"
	"github.com/ncarlier/reader/pkg/logger"
	"github.com/ncarlier/reader/pkg/version"
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
	level := "info"
	if *conf.Debug {
		level = "debug"
	}
	logger.Configure(level, true, *conf.SentryDSN)

	// Configure the DB
	_db, err := db.Configure(*conf.DB)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not configure Database")
	}

	// Configure Event Broker
	_, err = eventbroker.Configure(*conf.Broker)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not configure Event Broker")
	}

	// Init service registry
	service.InitRegistry(_db)

	log.Debug().Msg("Starting Nunux Reader server...")

	server := &http.Server{
		Addr:    *conf.ListenAddr,
		Handler: api.NewRouter(conf),
	}

	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-quit
		log.Debug().Msg("Server is shutting down...")
		api.Shutdown()

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			log.Fatal().Err(err).Msg("Could not gracefully shutdown the server")
		}
		close(done)
	}()

	api.Start()

	log.Info().Str("listen", *conf.ListenAddr).Msg("Server started")

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal().Err(err).Str("listen", *conf.ListenAddr).Msg("Could not start the server")
	}

	<-done
	log.Debug().Msg("Server stopped")
}
