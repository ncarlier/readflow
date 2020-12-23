package incomingwebhook

import (
	"errors"

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
	id, ok := helper.ConvGQLStringToUint(p.Args["id"])
	if !ok {
		return nil, errors.New("invalid incoming webhook ID")
	}
	return service.Lookup().GetIncomingWebhook(p.Context, id)
}

func init() {
	schema.AddQueryField("incomingWebhook", incomingWebhookQueryField)
	schema.AddQueryField("incomingWebhooks", incomingWebhooksQueryField)
}
