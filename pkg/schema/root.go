package schema

import (
	"github.com/graphql-go/graphql"
	"github.com/rs/zerolog/log"
)

// Root is the root schema
var Root graphql.Schema

var rootQuery = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"me":         meQueryField,
			"categories": categoriesQueryField,
			"articles":   articlesQueryField,
			"article":    articleQueryField,
			"apiKeys":    apiKeysQueryField,
		},
	},
)

var rootMutation = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"createOrUpdateCategory": createOrUpdateCategoryMutationField,
			"deleteCategory":         deleteCategoryMutationField,
			"updateArticleStatus":    updateArticleStatusMutationField,
			"createOrUpdateAPIKey":   createOrUpdateAPIKeyMutationField,
			"deleteAPIKey":           deleteAPIKeyMutationField,
		},
	},
)

func init() {
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:    rootQuery,
		Mutation: rootMutation,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create schema")
	}
	Root = schema
}
