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
			"category":   categoryQueryField,
			"articles":   articlesQueryField,
			"article":    articleQueryField,
			"apiKeys":    apiKeysQueryField,
			"apiKey":     apiKeyQueryField,
			"archivers":  archiversQueryField,
			"archiver":   archiverQueryField,
			"device":     deviceQueryField,
			"plans":      plansQueryField,
		},
	},
)

var rootMutation = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"createOrUpdateCategory": createOrUpdateCategoryMutationField,
			"deleteCategories":       deleteCategoriesMutationField,
			"addArticle":             addArticleMutationField,
			"updateArticle":          updateArticleMutationField,
			"markAllArticlesAsRead":  markAllArticlesAsReadMutationField,
			"cleanHistory":           cleanHistoryMutationField,
			"createOrUpdateAPIKey":   createOrUpdateAPIKeyMutationField,
			"deleteAPIKeys":          deleteAPIKeysMutationField,
			"createOrUpdateArchiver": createOrUpdateArchiverMutationField,
			"deleteArchivers":        deleteArchiversMutationField,
			"archiveArticle":         archiveArticleMutationField,
			"createPushSubscription": createPushSubscriptionMutationField,
			"deletePushSubscription": deletePushSubscriptionMutationField,
			"deleteAccount":          deleteAccountMutationField,
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
