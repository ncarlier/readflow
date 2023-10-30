package api

import (
	"net/http"

	"github.com/ncarlier/readflow/pkg/config"
	"github.com/ncarlier/readflow/pkg/middleware"
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
	authnMiddleware := middleware.Auth(conf.AuthN)
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
			middleware.RateLimiting("webhook", conf.RateLimiting.Webhook),
			middleware.IncomingWebhookAuth,
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
			middleware.IsAdmin,
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
			middleware.IsAdmin,
			authnMiddleware,
			middleware.Methods(http.MethodGet),
			middleware.Cors(origin),
		),
	}
}
