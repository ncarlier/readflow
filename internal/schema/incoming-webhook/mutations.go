package incomingwebhook

import (
	"github.com/ncarlier/readflow/internal/global"
	"github.com/ncarlier/readflow/internal/model"
	"github.com/ncarlier/readflow/internal/schema"
	"github.com/ncarlier/readflow/pkg/utils"

	"github.com/graphql-go/graphql"
	"github.com/ncarlier/readflow/internal/service"
)

var createOrUpdateIncomingWebhookMutationField = &graphql.Field{
	Type:        incomingWebhookType,
	Description: "create or update an incoming webhook (use the ID parameter to update)",
	Args: graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{
			Type: graphql.ID,
		},
		"alias": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"script": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
	},
	Resolve: createOrUpdateIncomingWebhookResolver,
}

func createOrUpdateIncomingWebhookResolver(p graphql.ResolveParams) (interface{}, error) {
	alias := utils.ParseGraphQLArgument[string](p.Args, "alias")
	script := utils.ParseGraphQLArgument[string](p.Args, "script")
	if id := utils.ParseGraphQLID(p.Args, "id"); id != nil {
		form := model.IncomingWebhookUpdateForm{
			ID:     *id,
			Alias:  alias,
			Script: script,
		}
		return service.Lookup().UpdateIncomingWebhook(p.Context, form)
	}
	if alias == nil {
		return nil, global.RequireParameterError("alias")
	}
	if script == nil {
		return nil, global.RequireParameterError("script")
	}
	builder := model.NewIncomingWebhookCreateFormBuilder()
	form := builder.Alias(*alias).Script(*script).Build()
	return service.Lookup().CreateIncomingWebhook(p.Context, *form)
}

var deleteIncomingWebhooksMutationField = &graphql.Field{
	Type:        graphql.Int,
	Description: "delete incoming webhooks",
	Args: graphql.FieldConfigArgument{
		"ids": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.NewList(graphql.ID)),
		},
	},
	Resolve: deleteIncomingWebhooksResolver,
}

func deleteIncomingWebhooksResolver(p graphql.ResolveParams) (interface{}, error) {
	idsArg, ok := p.Args["ids"].([]interface{})
	if !ok {
		return nil, global.InvalidParameterError("ids")
	}
	var ids []uint
	for _, v := range idsArg {
		if id := utils.ConvGraphQLID(v); id != nil {
			ids = append(ids, *id)
		}
	}

	return service.Lookup().DeleteIncomingWebhooks(p.Context, ids)
}

func init() {
	schema.AddMutationField("createOrUpdateIncomingWebhook", createOrUpdateIncomingWebhookMutationField)
	schema.AddMutationField("deleteIncomingWebhooks", deleteIncomingWebhooksMutationField)
}
