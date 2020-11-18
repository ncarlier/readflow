package device

import (
	"github.com/graphql-go/graphql"
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
