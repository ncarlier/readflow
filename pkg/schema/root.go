package schema

import (
	"github.com/graphql-go/graphql"
	"github.com/rs/zerolog/log"
)

// Root is the root schema
var Root graphql.Schema

var queryType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"me": meQueryField,
		},
	})

func init() {
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: queryType,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create schema")
	}
	Root = schema
}
