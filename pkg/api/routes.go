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
	origin := buildOriginFromPublicURL(conf.Global.PublicURL)
	authnMiddleware := middleware.Auth(conf.Global.AuthN)
	return Routes{
		route(
			"/",
			index(conf),
			middleware.Methods("GET"),
			middleware.Cors(origin),
		),
		route(
			"/articles",
			articles(),
			middleware.RateLimiting("webhook", conf.RateLimiting.Webhook),
			middleware.APIKeyAuth,
			middleware.Methods("POST"),
			middleware.Cors("*"),
		),
		route(
			"/articles/",
			download(),
			authnMiddleware,
			middleware.Methods("GET"),
			middleware.Cors(origin),
		),
		route(
			"/graphql",
			graphqlHandler(),
			authnMiddleware,
			middleware.Methods("GET", "POST"),
			middleware.Cors(origin),
		),
		route(
			"/linking/",
			linking(conf),
			authnMiddleware,
			middleware.Methods("GET"),
			middleware.Cors(origin),
		),
		route(
			"/admin",
			adminHandler(),
			middleware.IsAdmin,
			authnMiddleware,
			middleware.Methods("GET", "POST"),
			middleware.Cors(origin),
		),
		route(
			"/img",
			imgProxyHandler(conf),
			middleware.Methods("GET"),
			middleware.Cors(origin),
		),
		route(
			"/qr",
			qrcodeHandler(conf),
			authnMiddleware,
			middleware.Methods("GET"),
			middleware.Cors(origin),
		),
		route(
			"/healthz",
			healthz(),
			middleware.Methods("GET"),
			middleware.Cors("*"),
		),
		route(
			"/varz",
			varz(),
			middleware.IsAdmin,
			authnMiddleware,
			middleware.Methods("GET"),
			middleware.Cors(origin),
		),
	}
}
