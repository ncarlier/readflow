package server

import (
	"context"

	"github.com/emersion/go-smtp"
	"github.com/ncarlier/readflow/pkg/config"
	"github.com/ncarlier/readflow/pkg/constant"
	"github.com/ncarlier/readflow/pkg/helper"
	"github.com/ncarlier/readflow/pkg/mail"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// SMTPServer is a SMTP server
type SMTPServer struct {
	self   *smtp.Server
	logger zerolog.Logger
}

// ListenAndServe start SMTP server
func (s *SMTPServer) ListenAndServe() {
	s.logger.Info().Msg("starting the server...")
	if err := s.self.ListenAndServe(); err != nil && err != smtp.ErrServerClosed {
		s.logger.Fatal().Err(err).Msg("unable to start the server")
	}
}

// Shutdown stop SMTP server
func (s *SMTPServer) Shutdown(ctx context.Context) error {
	return s.self.Shutdown(ctx)
}

// NewSMTPServer create new SMTP server
func NewSMTPHTTPServer(cfg *config.Config) *SMTPServer {
	addr := cfg.Global.SMTPListenAddr
	if addr == "" {
		return nil
	}
	backend := mail.NewBackend()

	s := smtp.NewServer(backend)

	s.Addr = addr
	s.Domain = helper.GetMailHostname()
	s.ReadTimeout = constant.DefaultTimeout
	s.WriteTimeout = constant.DefaultTimeout
	s.MaxMessageBytes = 1024 * 1024
	s.MaxRecipients = 1
	s.AllowInsecureAuth = true

	logger := log.With().Str("component", "smtp").Str("addr", s.Addr).Logger()
	return &SMTPServer{
		self:   s,
		logger: logger,
	}
}
