package outboundservice

import (
	"github.com/graphql-go/graphql"
)

var providerEnum = graphql.NewEnum(
	graphql.EnumConfig{
		Name:        "outboundProvider",
		Description: "Outbound service provider",
		Values: graphql.EnumValueConfigMap{
			"keeper": &graphql.EnumValueConfig{
				Value:       "keeper",
				Description: "Use Nunux Keeper as outbound service provider",
			},
			"webhook": &graphql.EnumValueConfig{
				Value:       "webhook",
				Description: "Use a webhook as outbound service provider",
			},
			"wallabag": &graphql.EnumValueConfig{
				Value:       "wallabag",
				Description: "Use Wallabag as outbound service provider",
			},
		},
	},
)

var outboundServiceType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "OutboundService",
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
