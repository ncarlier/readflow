package inboundservice

import (
	"github.com/graphql-go/graphql"
)

var providerEnum = graphql.NewEnum(
	graphql.EnumConfig{
		Name:        "provider",
		Description: "Inbound service provider",
		Values: graphql.EnumValueConfigMap{
			"webhook": &graphql.EnumValueConfig{
				Value:       "webhook",
				Description: "Use a webhook as inbound service provider",
			},
		},
	},
)

var inboundServiceType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "InboundService",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"alias": &graphql.Field{
				Type: graphql.String,
			},
			"token": &graphql.Field{
				Type: graphql.String,
			},
			"provider": &graphql.Field{
				Type: providerEnum,
			},
			"config": &graphql.Field{
				Type: graphql.String,
			},
			"last_usage_at": &graphql.Field{
				Type: graphql.DateTime,
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
