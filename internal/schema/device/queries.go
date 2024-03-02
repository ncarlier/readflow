package device

import (
	"github.com/graphql-go/graphql"
	"github.com/ncarlier/readflow/internal/global"
	"github.com/ncarlier/readflow/internal/schema"
	"github.com/ncarlier/readflow/internal/service"
	"github.com/ncarlier/readflow/pkg/utils"
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
	id := utils.ParseGraphQLID(p.Args, "id")
	if id == nil {
		return nil, global.InvalidParameterError("id")
	}
	return service.Lookup().GetDevice(p.Context, *id)
}

func init() {
	schema.AddQueryField("devices", devicesQueryField)
	schema.AddQueryField("device", deviceQueryField)
}
