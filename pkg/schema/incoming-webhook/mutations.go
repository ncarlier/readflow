package incomingwebhook

import (
	"errors"

	"github.com/ncarlier/readflow/pkg/helper"
	"github.com/ncarlier/readflow/pkg/model"
	"github.com/ncarlier/readflow/pkg/schema"

	"github.com/graphql-go/graphql"
	"github.com/ncarlier/readflow/pkg/service"
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
	},
	Resolve: createOrUpdateIncomingWebhookResolver,
}

func createOrUpdateIncomingWebhookResolver(p graphql.ResolveParams) (interface{}, error) {
	alias := helper.GetGQLStringParameter("alias", p.Args)
	if id, ok := helper.ConvGQLStringToUint(p.Args["id"]); ok {
		form := model.IncomingWebhookUpdateForm{
			ID:    id,
			Alias: alias,
		}
		return service.Lookup().UpdateIncomingWebhook(p.Context, form)
	}
	builder := model.NewIncomingWebhookCreateFormBuilder()
	form := builder.Alias(*alias).Build()
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
		return nil, errors.New("invalid incoming webhook ID")
	}
	var ids []uint
	for _, v := range idsArg {
		if id, ok := helper.ConvGQLStringToUint(v); ok {
			ids = append(ids, id)
		}
	}

	return service.Lookup().DeleteIncomingWebhooks(p.Context, ids)
}

func init() {
	schema.AddMutationField("createOrUpdateIncomingWebhook", createOrUpdateIncomingWebhookMutationField)
	schema.AddMutationField("deleteIncomingWebhooks", deleteIncomingWebhooksMutationField)
}
