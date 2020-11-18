package inboundservice

import (
	"errors"

	"github.com/ncarlier/readflow/pkg/helper"
	"github.com/ncarlier/readflow/pkg/schema"

	"github.com/graphql-go/graphql"
	"github.com/ncarlier/readflow/pkg/service"
)

var inboundServicesQueryField = &graphql.Field{
	Type:    graphql.NewList(inboundServiceType),
	Resolve: inboundServicesResolver,
}

func inboundServicesResolver(p graphql.ResolveParams) (interface{}, error) {
	return service.Lookup().GetInboundServices(p.Context)
}

var inboundServiceQueryField = &graphql.Field{
	Type: inboundServiceType,
	Args: graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{
			Type: graphql.ID,
		},
	},
	Resolve: inboundServiceResolver,
}

func inboundServiceResolver(p graphql.ResolveParams) (interface{}, error) {
	id, ok := helper.ConvGQLStringToUint(p.Args["id"])
	if !ok {
		return nil, errors.New("invalid inbound service ID")
	}
	return service.Lookup().GetInboundService(p.Context, id)
}

func init() {
	schema.AddQueryField("inboundService", inboundServiceQueryField)
	schema.AddQueryField("inboundServices", inboundServicesQueryField)
}
