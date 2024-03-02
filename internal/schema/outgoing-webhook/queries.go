package outboundservice

import (
	"github.com/ncarlier/readflow/internal/global"
	"github.com/ncarlier/readflow/internal/schema"
	"github.com/ncarlier/readflow/pkg/utils"

	"github.com/graphql-go/graphql"
	"github.com/ncarlier/readflow/internal/service"
)

var outgoingWebhooksQueryField = &graphql.Field{
	Type:    graphql.NewList(outgoingWebhookType),
	Resolve: outgoingWebhooksResolver,
}

func outgoingWebhooksResolver(p graphql.ResolveParams) (interface{}, error) {
	return service.Lookup().GetOutgoingWebhooks(p.Context)
}

var outgoingWebhookQueryField = &graphql.Field{
	Type: outgoingWebhookType,
	Args: graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{
			Type: graphql.ID,
		},
	},
	Resolve: outgoingWebhookResolver,
}

func outgoingWebhookResolver(p graphql.ResolveParams) (interface{}, error) {
	id := utils.ParseGraphQLID(p.Args, "id")
	if id == nil {
		return nil, global.InvalidParameterError("id")
	}
	return service.Lookup().GetOutgoingWebhook(p.Context, *id)
}

func init() {
	schema.AddQueryField("outgoingWebhooks", outgoingWebhooksQueryField)
	schema.AddQueryField("outgoingWebhook", outgoingWebhookQueryField)
}
