package server

import (
	"context"
	"net/http"

	"github.com/ncarlier/readflow/internal/api"
	"github.com/ncarlier/readflow/internal/config"
	"github.com/ncarlier/readflow/pkg/logger"
	"github.com/rs/zerolog"
)

// MetricsServer is a HTTP server for metrics
type MetricsServer struct {
	self   *http.Server
	logger zerolog.Logger
}

// ListenAndServe start metrics server
func (s *MetricsServer) ListenAndServe() {
	s.logger.Info().Msg("starting the server...")
	if err := s.self.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		s.logger.Fatal().Err(err).Msg("unable to start the server")
	}
}

// Shutdown stop metrics server
func (s *MetricsServer) Shutdown(ctx context.Context) error {
	s.self.SetKeepAlivesEnabled(false)
	return s.self.Shutdown(ctx)
}

// NewMetricsServer create new metrics server
func NewMetricsServer(cfg *config.Config) *MetricsServer {
	addr := cfg.Metrics.ListenAddr
	if addr == "" {
		return nil
	}
	logger := logger.With().Str("component", "metrics").Str("addr", addr).Logger()
	handler := api.NewMetricsRouter()
	return &MetricsServer{
		self: &http.Server{
			Addr:    addr,
			Handler: handler,
		},
		logger: logger,
	}
}
