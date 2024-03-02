package incomingwebhook

import (
	"github.com/ncarlier/readflow/internal/global"
	"github.com/ncarlier/readflow/internal/schema"
	"github.com/ncarlier/readflow/pkg/utils"

	"github.com/graphql-go/graphql"
	"github.com/ncarlier/readflow/internal/service"
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
	id := utils.ParseGraphQLID(p.Args, "id")
	if id == nil {
		return nil, global.InvalidParameterError("id")
	}
	return service.Lookup().GetIncomingWebhook(p.Context, *id)
}

func init() {
	schema.AddQueryField("incomingWebhook", incomingWebhookQueryField)
	schema.AddQueryField("incomingWebhooks", incomingWebhooksQueryField)
}
