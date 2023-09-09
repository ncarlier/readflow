package incomingwebhook

import (
	"github.com/ncarlier/readflow/pkg/helper"
	"github.com/ncarlier/readflow/pkg/schema"

	"github.com/graphql-go/graphql"
	"github.com/ncarlier/readflow/pkg/service"
)

var incomingWebhooksQueryField = &graphql.Field{
	Type:    graphql.NewList(incomingWebhookType),
	Resolve: incomingWebhooksResolver,
}

func incomingWebhooksResolver(p graphql.ResolveParams) (interface{}, error) {
	return service.Lookup().GetIncomingWebhooks(p.Context)
}

var incomingWebhookQueryField = &graphql.Field{
	Type: incomingWebhookType,
	Args: graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{
			Type: graphql.ID,
		},
	},
	Resolve: incomingWebhookResolver,
}

func incomingWebhookResolver(p graphql.ResolveParams) (interface{}, error) {
	id := helper.ParseGraphQLID(p.Args, "id")
	if id == nil {
		return nil, helper.InvalidParameterError("id")
	}
	return service.Lookup().GetIncomingWebhook(p.Context, *id)
}

func init() {
	schema.AddQueryField("incomingWebhook", incomingWebhookQueryField)
	schema.AddQueryField("incomingWebhooks", incomingWebhooksQueryField)
}
