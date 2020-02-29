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
	authnMiddleware := middleware.Auth(conf.AuthN)
	return Routes{
		route(
			"/",
			index(conf),
			middleware.Methods("GET"),
		),
		route(
			"/articles",
			articles(),
			middleware.APIKeyAuth,
			middleware.Methods("POST"),
		),
		route(
			"/graphql",
			graphqlHandler(),
			authnMiddleware,
			middleware.Methods("GET", "POST"),
		),
		route(
			"/admin",
			adminHandler(),
			middleware.IsAdmin,
			authnMiddleware,
			middleware.Methods("GET", "POST"),
		),
		route(
			"/img",
			imgProxyHandler(conf),
			middleware.Methods("GET"),
		),
		route(
			"/healthz",
			healthz(),
			middleware.Methods("GET"),
		),
		route(
			"/varz",
			varz(),
			middleware.IsAdmin,
			authnMiddleware,
			middleware.Methods("GET"),
		),
	}
}
