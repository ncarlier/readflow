package outboundservice

import (
	"errors"

	"github.com/graphql-go/graphql"
	"github.com/ncarlier/readflow/internal/model"
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

var outgoingWebhookSecretsField = &graphql.Field{
	Type: graphql.NewList(graphql.String),
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		obj, ok := p.Source.(*model.OutgoingWebhook)
		if !ok {
			return nil, errors.New("invalid object received by secrets resolver")
		}
		keys := make([]string, 0, len(obj.Secrets))
		for k := range obj.Secrets {
			keys = append(keys, k)
		}
		return keys, nil
	},
}

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
			"secrets": outgoingWebhookSecretsField,
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

var outgoingWebhookResponseType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "OutgoingWebhookResponse",
		Fields: graphql.Fields{
			"url": &graphql.Field{
				Type: graphql.String,
			},
			"text": &graphql.Field{
				Type: graphql.String,
			},
			"json": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)
