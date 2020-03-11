package schema

import (
	"errors"

	"github.com/ncarlier/readflow/pkg/tooling"

	"github.com/graphql-go/graphql"
	"github.com/ncarlier/readflow/pkg/service"
)

var pushKeysType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "PushKeys",
		Fields: graphql.Fields{
			"auth": &graphql.Field{
				Type: graphql.String,
			},
			"p256dh": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

var pushSubscriptionType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "PushSubscription",
		Fields: graphql.Fields{
			"enpoint": &graphql.Field{
				Type: graphql.String,
			},
			"keys": &graphql.Field{
				Type: pushKeysType,
			},
		},
	},
)

var deviceType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Device",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"key": &graphql.Field{
				Type: graphql.String,
			},
			"created_at": &graphql.Field{
				Type: graphql.DateTime,
			},
		},
	},
)

// QUERIES

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
	id, ok := tooling.ConvGQLStringToUint(p.Args["id"])
	if !ok {
		return nil, errors.New("invalid device ID")
	}
	return service.Lookup().GetDevice(p.Context, id)
}

// MUTATIONS

var createPushSubscriptionMutationField = &graphql.Field{
	Type:        deviceType,
	Description: "create push subscription for a device",
	Args: graphql.FieldConfigArgument{
		"sub": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
	},
	Resolve: createPushSubscriptionResolver,
}

func createPushSubscriptionResolver(p graphql.ResolveParams) (interface{}, error) {
	sub, _ := p.Args["sub"].(string)
	return service.Lookup().CreateDevice(p.Context, sub)
}

var deletePushSubscriptionMutationField = &graphql.Field{
	Type:        deviceType,
	Description: "remove device push subscription",
	Args: graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.ID),
		},
	},
	Resolve: deletePushSubscriptionResolver,
}

func deletePushSubscriptionResolver(p graphql.ResolveParams) (interface{}, error) {
	id, ok := tooling.ConvGQLStringToUint(p.Args["id"])
	if !ok {
		return nil, errors.New("invalid device ID")
	}
	return service.Lookup().DeleteDevice(p.Context, id)
}
