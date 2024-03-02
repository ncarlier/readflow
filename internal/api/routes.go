package api

import (
	"net/http"

	"github.com/ncarlier/readflow/internal/auth"
	_ "github.com/ncarlier/readflow/internal/auth/methods"
	"github.com/ncarlier/readflow/internal/config"
	"github.com/ncarlier/readflow/pkg/logger"
	"github.com/ncarlier/readflow/pkg/middleware"
	ratelimiter "github.com/ncarlier/readflow/pkg/rate-limiter"
)

// Route is the structure of an HTTP route definition
type Route struct {
	Path        string
	Handler     http.Handler
	Middlewares []middleware.Middleware
}

func route(path string, handler http.Handler, middlewares ...middleware.Middleware) Route {
	return Route{
		Path:        path,
		Handler:     handler,
		Middlewares: middlewares,
	}
}

// Routes is a list of Route
type Routes []Route

func routes(conf *config.Config) Routes {
	origin := conf.UI.PublicURL
	if origin == "" {
		origin = "*"
	}
	authnMiddleware, err := auth.NewAuthMiddleware(&conf.AuthN)
	if err != nil {
		logger.Fatal().Err(err).Msg("unable to create authentication middleware")
	}
	return Routes{
		route(
			"/",
			index(conf),
			middleware.Methods(http.MethodGet),
			middleware.Cors("*"),
		),
		route(
			"/info",
			info(conf),
			middleware.Methods(http.MethodGet),
			middleware.Cors(origin),
		),
		route(
			"/articles",
			articles(),
			ratelimiter.Middleware("webhook", &conf.RateLimiting.Webhook),
			auth.IncomingWebhookAuth,
			middleware.Methods(http.MethodPost),
			middleware.Cors("*"),
		),
		route(
			"/articles/",
			download(),
			authnMiddleware,
			middleware.Methods(http.MethodGet),
			middleware.Cors(origin),
		),
		route(
			"/graphql",
			graphqlHandler(),
			authnMiddleware,
			middleware.Methods(http.MethodGet, http.MethodPost),
			middleware.Cors(origin),
		),
		route(
			"/linking/",
			linking(conf),
			authnMiddleware,
			middleware.Methods(http.MethodGet),
			middleware.Cors(origin),
		),
		route(
			"/admin",
			adminHandler(),
			auth.IsAdmin,
			authnMiddleware,
			middleware.Methods(http.MethodGet, http.MethodPost),
			middleware.Cors(origin),
		),
		route(
			"/img/",
			imgProxyHandler(conf),
			middleware.Methods(http.MethodGet),
			middleware.Cors(origin),
		),
		route(
			"/qr",
			qrcodeHandler(conf),
			authnMiddleware,
			middleware.Methods(http.MethodGet),
			middleware.Cors(origin),
		),
		route(
			"/avatar/",
			avatarHandler(conf),
			middleware.Methods(http.MethodGet),
			middleware.Cors(origin),
		),
		route(
			"/healthz",
			healthz(),
			middleware.Methods(http.MethodGet, http.MethodHead),
			middleware.Cors("*"),
		),
		route(
			"/varz",
			varz(),
			auth.IsAdmin,
			authnMiddleware,
			middleware.Methods(http.MethodGet),
			middleware.Cors(origin),
		),
	}
}
