package device

import (
	"errors"

	"github.com/graphql-go/graphql"
	"github.com/ncarlier/readflow/pkg/helper"
	"github.com/ncarlier/readflow/pkg/schema"
	"github.com/ncarlier/readflow/pkg/service"
)

var devicesQueryField = &graphql.Field{
	Type:    graphql.NewList(deviceType),
	Resolve: DevicesResolver,
}

// DevicesResolver is the resolver for retrieve devices
func DevicesResolver(p graphql.ResolveParams) (interface{}, error) {
	return service.Lookup().GetDevices(p.Context)
}

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
	schema.AddQueryField("devices", devicesQueryField)
	schema.AddQueryField("device", deviceQueryField)
}
