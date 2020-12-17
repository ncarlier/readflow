package outboundservice

import (
	"errors"

	"github.com/ncarlier/readflow/pkg/helper"
	"github.com/ncarlier/readflow/pkg/schema"

	"github.com/graphql-go/graphql"
	"github.com/ncarlier/readflow/pkg/service"
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
	id, ok := helper.ConvGQLStringToUint(p.Args["id"])
	if !ok {
		return nil, errors.New("invalid outgoing webhook ID")
	}
	return service.Lookup().GetOutgoingWebhook(p.Context, id)
}

func init() {
	schema.AddQueryField("outgoingWebhooks", outgoingWebhooksQueryField)
	schema.AddQueryField("outgoingWebhook", outgoingWebhookQueryField)
}
