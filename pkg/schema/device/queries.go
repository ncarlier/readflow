package device

import (
	"errors"

	"github.com/graphql-go/graphql"
	"github.com/ncarlier/readflow/pkg/helper"
	"github.com/ncarlier/readflow/pkg/schema"
	"github.com/ncarlier/readflow/pkg/service"
)

var deviceQueryField = &graphql.Field{
	Type: deviceType,
	Args: graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{
			Type: graphql.ID,
		},
	},
	Resolve: deviceResolver,
}

func deviceResolver(p graphql.ResolveParams) (interface{}, error) {
	id, ok := helper.ConvGQLStringToUint(p.Args["id"])
	if !ok {
		return nil, errors.New("invalid device ID")
	}
	return service.Lookup().GetDevice(p.Context, id)
}

func init() {
	schema.AddQueryField("device", deviceQueryField)
}
