package schema

import (
	"errors"

	"github.com/ncarlier/readflow/pkg/tooling"

	"github.com/graphql-go/graphql"
	"github.com/ncarlier/readflow/pkg/service"
)

var apiKeyType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "APIKey",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"alias": &graphql.Field{
				Type: graphql.String,
			},
			"token": &graphql.Field{
				Type: graphql.String,
			},
			"last_usage_at": &graphql.Field{
				Type: graphql.DateTime,
			},
			"created_at": &graphql.Field{
				Type: graphql.DateTime,
			},
			"updated_at": &graphql.Field{
				Type: graphql.DateTime,
			},
		},
	},
)

// QUERIES

var apiKeysQueryField = &graphql.Field{
	Type:    graphql.NewList(apiKeyType),
	Resolve: apiKeysResolver,
}

func apiKeysResolver(p graphql.ResolveParams) (interface{}, error) {
	return service.Lookup().GetAPIKeys(p.Context)
}

var apiKeyQueryField = &graphql.Field{
	Type: apiKeyType,
	Args: graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{
			Type: graphql.ID,
		},
	},
	Resolve: apiKeyResolver,
}

func apiKeyResolver(p graphql.ResolveParams) (interface{}, error) {
	id, ok := tooling.ConvGQLStringToUint(p.Args["id"])
	if !ok {
		return nil, errors.New("invalid API key ID")
	}
	return service.Lookup().GetAPIKey(p.Context, id)
}

// MUTATIONS

var createOrUpdateAPIKeyMutationField = &graphql.Field{
	Type:        apiKeyType,
	Description: "create or update an API key (use the ID parameter to update)",
	Args: graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{
			Type: graphql.ID,
		},
		"alias": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
	},
	Resolve: createOrUpdateAPIKeyResolver,
}

func createOrUpdateAPIKeyResolver(p graphql.ResolveParams) (interface{}, error) {
	var id *uint
	if val, ok := tooling.ConvGQLStringToUint(p.Args["id"]); ok {
		id = &val
	}
	alias, _ := p.Args["alias"].(string)

	return service.Lookup().CreateOrUpdateAPIKey(p.Context, id, alias)
}

var deleteAPIKeysMutationField = &graphql.Field{
	Type:        graphql.Int,
	Description: "delete API keys",
	Args: graphql.FieldConfigArgument{
		"ids": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.NewList(graphql.ID)),
		},
	},
	Resolve: deleteAPIKeysResolver,
}

func deleteAPIKeysResolver(p graphql.ResolveParams) (interface{}, error) {
	idsArg, ok := p.Args["ids"].([]interface{})
	if !ok {
		return nil, errors.New("invalid API Key ID")
	}
	var ids []uint
	for _, v := range idsArg {
		if id, ok := tooling.ConvGQLStringToUint(v); ok {
			ids = append(ids, id)
		}
	}

	return service.Lookup().DeleteAPIKeys(p.Context, ids)
}
