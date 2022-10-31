package outboundservice

import (
	"github.com/graphql-go/graphql"
)

var providerEnum = graphql.NewEnum(
	graphql.EnumConfig{
		Name:        "outgoingWebhookProvider",
		Description: "Outgoing webhook provider",
		Values: graphql.EnumValueConfigMap{
			"generic": &graphql.EnumValueConfig{
				Value:       "generic",
				Description: "Use a generic webhook as outgoing webhook provider",
			},
			"keeper": &graphql.EnumValueConfig{
				Value:       "keeper",
				Description: "Use Nunux Keeper as outgoing webhook provider",
			},
			"wallabag": &graphql.EnumValueConfig{
				Value:       "wallabag",
				Description: "Use Wallabag as outgoing webhook provider",
			},
			"shaarli": &graphql.EnumValueConfig{
				Value:       "shaarli",
				Description: "Use Shaarli as outgoing webhook provider",
			},
			"pocket": &graphql.EnumValueConfig{
				Value:       "pocket",
				Description: "Use Pocket as outgoing webhook provider",
			},
			"readflow": &graphql.EnumValueConfig{
				Value:       "readflow",
				Description: "Use Readflow as outgoing webhook provider",
			},
			"s3": &graphql.EnumValueConfig{
				Value:       "s3",
				Description: "Use S3 bucket as outgoing webhook provider",
			},
		},
	},
)

var outgoingWebhookType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "OutgoingWebhook",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"alias": &graphql.Field{
				Type: graphql.String,
			},
			"provider": &graphql.Field{
				Type: providerEnum,
			},
			"config": &graphql.Field{
				Type: graphql.String,
			},
			"is_default": &graphql.Field{
				Type: graphql.Boolean,
			},
			"created_at": &graphql.Field{
				Type: graphql.DateTime,
			},
			"updated_at": &graphql.Field{
				Type: graphql.DateTime,
			},
		},
	},
)
