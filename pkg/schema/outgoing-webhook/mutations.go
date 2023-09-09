package outboundservice

import (
	"fmt"

	"github.com/ncarlier/readflow/pkg/helper"
	"github.com/ncarlier/readflow/pkg/model"
	"github.com/ncarlier/readflow/pkg/schema"
	"github.com/ncarlier/readflow/pkg/secret"

	"github.com/graphql-go/graphql"
	"github.com/ncarlier/readflow/pkg/service"
)

var createOrUpdateOutgoingWebhookMutationField = &graphql.Field{
	Type:        outgoingWebhookType,
	Description: "create or update a outgoing webhook (use the ID parameter to update)",
	Args: graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{
			Type: graphql.ID,
		},
		"alias": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"provider": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(providerEnum),
		},
		"config": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"secrets": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"is_default": &graphql.ArgumentConfig{
			Type:         graphql.Boolean,
			DefaultValue: false,
		},
	},
	Resolve: createOrUpdateOutgoingWebhookResolver,
}

func createOrUpdateOutgoingWebhookResolver(p graphql.ResolveParams) (interface{}, error) {
	alias := helper.ParseGraphQLArgument[string](p.Args, "alias")
	provider := helper.ParseGraphQLArgument[string](p.Args, "provider")
	config := helper.ParseGraphQLArgument[string](p.Args, "config")
	isDefault := helper.ParseGraphQLArgument[bool](p.Args, "is_default")

	// decode secrets
	secretsParams := helper.ParseGraphQLArgument[string](p.Args, "secrets")
	secrets := make(secret.Secrets)
	if secretsParams != nil {
		if err := secrets.Scan(*secretsParams); err != nil {
			return nil, fmt.Errorf("invalid secrets: %v", err)
		}
	}

	if id := helper.ParseGraphQLID(p.Args, "id"); id != nil {
		form := model.OutgoingWebhookUpdateForm{
			ID:        *id,
			Alias:     alias,
			Provider:  provider,
			Config:    config,
			Secrets:   &secrets,
			IsDefault: isDefault,
		}
		return service.Lookup().UpdateOutgoingWebhook(p.Context, form)
	}
	builder := model.NewOutgoingWebhookCreateFormBuilder()
	builder.Alias(*alias).Provider(*provider).Config(*config).Secrets(secrets)
	if isDefault != nil && *isDefault {
		builder.IsDefault(true)
	}
	form := builder.Build()

	return service.Lookup().CreateOutgoingWebhook(p.Context, *form)
}

var deleteOutgoingWebhooksMutationField = &graphql.Field{
	Type:        graphql.Int,
	Description: "delete outgoing webhooks",
	Args: graphql.FieldConfigArgument{
		"ids": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.NewList(graphql.ID)),
		},
	},
	Resolve: deleteOutgoingWebhooksResolver,
}

func deleteOutgoingWebhooksResolver(p graphql.ResolveParams) (interface{}, error) {
	idsArg, ok := p.Args["ids"].([]interface{})
	if !ok {
		return nil, helper.InvalidParameterError("ids")
	}
	var ids []uint
	for _, v := range idsArg {
		if id := helper.ConvGraphQLID(v); id != nil {
			ids = append(ids, *id)
		}
	}

	return service.Lookup().DeleteOutgoingWebhooks(p.Context, ids)
}

var sendArticleToOutgoingWebhookMutationField = &graphql.Field{
	Type:        outgoingWebhookResponseType,
	Description: "send an article to the outgoing webhook",
	Args: graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{
			Description: "article ID",
			Type:        graphql.NewNonNull(graphql.ID),
		},
		"alias": &graphql.ArgumentConfig{
			Description: "outgoing webhook alias (using default if missing)",
			Type:        graphql.String,
		},
	},
	Resolve: sendArticleToOutgoingWebhookResolver,
}

func sendArticleToOutgoingWebhookResolver(p graphql.ResolveParams) (interface{}, error) {
	id := helper.ParseGraphQLID(p.Args, "id")
	if id == nil {
		return nil, helper.InvalidParameterError("id")
	}
	alias := helper.ParseGraphQLArgument[string](p.Args, "alias")
	return service.Lookup().SendArticle(p.Context, *id, alias)
}

func init() {
	schema.AddMutationField("createOrUpdateOutgoingWebhook", createOrUpdateOutgoingWebhookMutationField)
	schema.AddMutationField("deleteOutgoingWebhooks", deleteOutgoingWebhooksMutationField)
	schema.AddMutationField("sendArticleToOutgoingWebhook", sendArticleToOutgoingWebhookMutationField)
}
