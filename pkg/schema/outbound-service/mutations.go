package outboundservice

import (
	"errors"

	"github.com/ncarlier/readflow/pkg/helper"
	"github.com/ncarlier/readflow/pkg/model"
	"github.com/ncarlier/readflow/pkg/schema"

	"github.com/graphql-go/graphql"
	"github.com/ncarlier/readflow/pkg/service"
)

var createOrUpdateOutboundServiceMutationField = &graphql.Field{
	Type:        outboundServiceType,
	Description: "create or update a outbound service (use the ID parameter to update)",
	Args: graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{
			Type: graphql.ID,
		},
		"alias": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"provider": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(providerEnum),
		},
		"config": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"is_default": &graphql.ArgumentConfig{
			Type:         graphql.Boolean,
			DefaultValue: false,
		},
	},
	Resolve: createOrUpdateOutboundServiceResolver,
}

func createOrUpdateOutboundServiceResolver(p graphql.ResolveParams) (interface{}, error) {
	alias := helper.GetGQLStringParameter("alias", p.Args)
	provider := helper.GetGQLStringParameter("provider", p.Args)
	config := helper.GetGQLStringParameter("config", p.Args)
	isDefault := helper.GetGQLBoolParameter("is_default", p.Args)
	if id, ok := helper.ConvGQLStringToUint(p.Args["id"]); ok {
		form := model.OutboundServiceUpdateForm{
			ID:        id,
			Alias:     alias,
			Provider:  provider,
			Config:    config,
			IsDefault: isDefault,
		}
		return service.Lookup().UpdateOutboundService(p.Context, form)
	}
	builder := model.NewOutboundServiceCreateFormBuilder()
	builder.Alias(*alias).Provider(*provider).Config(*config)
	if isDefault != nil && *isDefault {
		builder.IsDefault(true)
	}
	form := builder.Build()

	return service.Lookup().CreateOutboundService(p.Context, *form)
}

var deleteOutboundServicesMutationField = &graphql.Field{
	Type:        graphql.Int,
	Description: "delete outbound services",
	Args: graphql.FieldConfigArgument{
		"ids": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.NewList(graphql.ID)),
		},
	},
	Resolve: deleteOutboundServicesResolver,
}

func deleteOutboundServicesResolver(p graphql.ResolveParams) (interface{}, error) {
	idsArg, ok := p.Args["ids"].([]interface{})
	if !ok {
		return nil, errors.New("invalid outbound service ID")
	}
	var ids []uint
	for _, v := range idsArg {
		if id, ok := helper.ConvGQLStringToUint(v); ok {
			ids = append(ids, id)
		}
	}

	return service.Lookup().DeleteOutboundServices(p.Context, ids)
}

var sendArticleToOutboundServiceMutationField = &graphql.Field{
	Type:        graphql.ID,
	Description: "send an article to the outbound service",
	Args: graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{
			Description: "article ID",
			Type:        graphql.NewNonNull(graphql.ID),
		},
		"alias": &graphql.ArgumentConfig{
			Description: "outbound service alias (using default if missing)",
			Type:        graphql.String,
		},
	},
	Resolve: sendArticleToOutboundServiceResolver,
}

func sendArticleToOutboundServiceResolver(p graphql.ResolveParams) (interface{}, error) {
	id, ok := helper.ConvGQLStringToUint(p.Args["id"])
	if !ok {
		return nil, errors.New("invalid article ID")
	}
	var alias *string
	if val, ok := p.Args["alias"]; ok {
		sVal := val.(string)
		alias = &sVal
	}

	err := service.Lookup().ArchiveArticle(p.Context, id, alias)
	return id, err
}

func init() {
	schema.AddMutationField("createOrUpdateOutboundService", createOrUpdateOutboundServiceMutationField)
	schema.AddMutationField("deleteOutboundServices", deleteOutboundServicesMutationField)
	schema.AddMutationField("sendArticleToOutboundService", sendArticleToOutboundServiceMutationField)
}
