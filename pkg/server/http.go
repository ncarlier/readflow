package server

import (
	"context"
	"net/http"

	"github.com/ncarlier/readflow/pkg/api"
	"github.com/ncarlier/readflow/pkg/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// HTTPServer is a HTTP server wrapper
type HTTPServer struct {
	self   *http.Server
	logger zerolog.Logger
}

// ListenAndServe start HTTP server
func (s *HTTPServer) ListenAndServe() {
	s.logger.Info().Msg("starting the server...")
	if err := s.self.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		s.logger.Fatal().Err(err).Msg("unable to start the server")
	}
}

// Shutdown stop HTTP server
func (s *HTTPServer) Shutdown(ctx context.Context) error {
	s.self.SetKeepAlivesEnabled(false)
	return s.self.Shutdown(ctx)
}

// NewHTTPServer create new HTTP server
func NewHTTPServer(cfg *config.Config) *HTTPServer {
	addr := cfg.HTTP.ListenAddr
	handler := api.NewRouter(cfg)
	logger := log.With().Str("component", "http").Str("addr", addr).Logger()
	return &HTTPServer{
		self: &http.Server{
			Addr:    addr,
			Handler: handler,
		},
		logger: logger,
	}
}
