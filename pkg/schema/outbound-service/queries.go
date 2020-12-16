package outboundservice

import (
	"errors"

	"github.com/ncarlier/readflow/pkg/helper"
	"github.com/ncarlier/readflow/pkg/schema"

	"github.com/graphql-go/graphql"
	"github.com/ncarlier/readflow/pkg/service"
)

var outboundServicesQueryField = &graphql.Field{
	Type:    graphql.NewList(outboundServiceType),
	Resolve: outboundServicesResolver,
}

func outboundServicesResolver(p graphql.ResolveParams) (interface{}, error) {
	return service.Lookup().GetOutboundServices(p.Context)
}

var outboundServiceQueryField = &graphql.Field{
	Type: outboundServiceType,
	Args: graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{
			Type: graphql.ID,
		},
	},
	Resolve: outboundServiceResolver,
}

func outboundServiceResolver(p graphql.ResolveParams) (interface{}, error) {
	id, ok := helper.ConvGQLStringToUint(p.Args["id"])
	if !ok {
		return nil, errors.New("invalid outbound service ID")
	}
	return service.Lookup().GetOutboundService(p.Context, id)
}

func init() {
	schema.AddQueryField("outboundServices", outboundServicesQueryField)
	schema.AddQueryField("outboundService", outboundServiceQueryField)
}
