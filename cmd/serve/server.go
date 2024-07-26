package serve

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ncarlier/readflow/internal/api"
	"github.com/ncarlier/readflow/internal/config"
	"github.com/ncarlier/readflow/internal/db"
	"github.com/ncarlier/readflow/internal/exporter"
	"github.com/ncarlier/readflow/internal/exporter/pdf"
	"github.com/ncarlier/readflow/internal/metric"
	"github.com/ncarlier/readflow/internal/server"
	"github.com/ncarlier/readflow/internal/service"
	"github.com/ncarlier/readflow/pkg/logger"
)

func startServer(conf *config.Config) error {
	logger.Debug().Msg("starting readflow...")

	// configure the DB
	database, err := db.NewDB(conf.Database.URI)
	if err != nil {
		return fmt.Errorf("unable to configure the database: %w", err)
	}

	// configure the service registry
	err = service.Configure(*conf, database)
	if err != nil {
		database.Close()
		return fmt.Errorf("unable to configure the service registry: %w", err)
	}

	// register external exporters...
	if conf.PDF.ServiceProvider != "" {
		logger.Info().Str("provider", conf.PDF.ServiceProvider).Msg("using PDF generator service")
		exporter.Register("pdf", pdf.NewPDFExporter(conf.PDF.ServiceProvider))
	}

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
		logger.Debug().Msg("shutting down readflow...")
		api.Shutdown()

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := httpServer.Shutdown(ctx); err != nil {
			logger.Fatal().Err(err).Msg("unable to gracefully shutdown the HTTP server")
		}
		if smtpServer != nil {
			if err := smtpServer.Shutdown(ctx); err != nil {
				logger.Fatal().Err(err).Msg("unable to gracefully shutdown the SMTP server")
			}
		}
		if metricsServer != nil {
			metric.StopCollectors()
			if err := metricsServer.Shutdown(ctx); err != nil {
				logger.Fatal().Err(err).Msg("unable to gracefully shutdown the metrics server")
			}
		}

		service.Shutdown()

		if err := database.Close(); err != nil {
			logger.Fatal().Err(err).Msg("could not gracefully shutdown database connection")
		}

		close(done)
	}()

	// set API health check as started
	api.Start()

	// start HTTP server
	httpServer.ListenAndServe()

	<-done
	logger.Debug().Msg("readflow stopped")

	return nil
}
