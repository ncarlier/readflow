package api

import (
	"net/http"

	"github.com/ncarlier/reader/pkg/config"
)

// HandlerFunc custom function handler
type HandlerFunc func(conf *config.Config) http.Handler

// Route is the structure of an HTTP route definition
type Route struct {
	Methods       []string
	Path          string
	AuthNRequired bool
	HandlerFunc   HandlerFunc
}

// Routes is a list of Route
type Routes []Route

var routes = Routes{
	Route{
		[]string{"GET"},
		"/",
		false,
		index,
	},
	Route{
		[]string{"POST"},
		"/articles",
		false,
		articles,
	},
	Route{
		[]string{"GET", "POST"},
		"/graphql",
		true,
		graphqlHandler,
	},
	Route{
		[]string{"GET", "POST"},
		"/graphiql",
		false,
		graphiqlHandler,
	},
	Route{
		[]string{"GET"},
		"/healthz",
		false,
		healthz,
	},
	Route{
		[]string{"GET"},
		"/varz",
		true,
		varz,
	},
	Route{
		[]string{"GET"},
		"/metrics",
		true,
		metrics,
	},
}
