package admin

import (
	"github.com/graphql-go/graphql"
	"github.com/rs/zerolog/log"
)

// Root is the Admin root schema
var Root graphql.Schema

var rootQuery = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"user": userQueryField,
		},
	},
)

func init() {
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: rootQuery,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create schema")
	}
	Root = schema
}
