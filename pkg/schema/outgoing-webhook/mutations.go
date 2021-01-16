package outboundservice

import (
	"errors"

	"github.com/ncarlier/readflow/pkg/helper"
	"github.com/ncarlier/readflow/pkg/model"
	"github.com/ncarlier/readflow/pkg/schema"

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
		"is_default": &graphql.ArgumentConfig{
			Type:         graphql.Boolean,
			DefaultValue: false,
		},
	},
	Resolve: createOrUpdateOutgoingWebhookResolver,
}

func createOrUpdateOutgoingWebhookResolver(p graphql.ResolveParams) (interface{}, error) {
	alias := helper.GetGQLStringParameter("alias", p.Args)
	provider := helper.GetGQLStringParameter("provider", p.Args)
	config := helper.GetGQLStringParameter("config", p.Args)
	isDefault := helper.GetGQLBoolParameter("is_default", p.Args)
	if id, ok := helper.ConvGQLStringToUint(p.Args["id"]); ok {
		form := model.OutgoingWebhookUpdateForm{
			ID:        id,
			Alias:     alias,
			Provider:  provider,
			Config:    config,
			IsDefault: isDefault,
		}
		return service.Lookup().UpdateOutgoingWebhook(p.Context, form)
	}
	builder := model.NewOutgoingWebhookCreateFormBuilder()
	builder.Alias(*alias).Provider(*provider).Config(*config)
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
		return nil, errors.New("invalid outgoing webhook ID")
	}
	var ids []uint
	for _, v := range idsArg {
		if id, ok := helper.ConvGQLStringToUint(v); ok {
			ids = append(ids, id)
		}
	}

	return service.Lookup().DeleteOutgoingWebhooks(p.Context, ids)
}

var sendArticleToOutgoingWebhookMutationField = &graphql.Field{
	Type:        graphql.ID,
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
	id, ok := helper.ConvGQLStringToUint(p.Args["id"])
	if !ok {
		return nil, errors.New("invalid article ID")
	}
	var alias *string
	if val, ok := p.Args["alias"]; ok {
		sVal := val.(string)
		alias = &sVal
	}

	err := service.Lookup().SendArticle(p.Context, id, alias)
	return id, err
}

func init() {
	schema.AddMutationField("createOrUpdateOutgoingWebhook", createOrUpdateOutgoingWebhookMutationField)
	schema.AddMutationField("deleteOutgoingWebhooks", deleteOutgoingWebhooksMutationField)
	schema.AddMutationField("sendArticleToOutgoingWebhook", sendArticleToOutgoingWebhookMutationField)
}
