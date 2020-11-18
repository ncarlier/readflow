package inboundservice

import (
	"errors"

	"github.com/ncarlier/readflow/pkg/helper"
	"github.com/ncarlier/readflow/pkg/model"
	"github.com/ncarlier/readflow/pkg/schema"

	"github.com/graphql-go/graphql"
	"github.com/ncarlier/readflow/pkg/service"
)

var createOrUpdateInboundServiceMutationField = &graphql.Field{
	Type:        inboundServiceType,
	Description: "create or update an inbound service (use the ID parameter to update)",
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
	},
	Resolve: createOrUpdateInboundServiceResolver,
}

func createOrUpdateInboundServiceResolver(p graphql.ResolveParams) (interface{}, error) {
	alias := helper.GetGQLStringParameter("alias", p.Args)
	provider := helper.GetGQLStringParameter("provider", p.Args)
	config := helper.GetGQLStringParameter("config", p.Args)
	if id, ok := helper.ConvGQLStringToUint(p.Args["id"]); ok {
		form := model.InboundServiceUpdateForm{
			ID:       id,
			Alias:    alias,
			Provider: provider,
			Config:   config,
		}
		return service.Lookup().UpdateInboundService(p.Context, form)
	}
	builder := model.NewInboundServiceCreateFormBuilder()
	form := builder.Alias(*alias).Provider(*provider).Config(*config).Build()
	return service.Lookup().CreateInboundService(p.Context, *form)
}

var deleteInboundServicesMutationField = &graphql.Field{
	Type:        graphql.Int,
	Description: "delete inbound services",
	Args: graphql.FieldConfigArgument{
		"ids": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.NewList(graphql.ID)),
		},
	},
	Resolve: deleteInboundServicesResolver,
}

func deleteInboundServicesResolver(p graphql.ResolveParams) (interface{}, error) {
	idsArg, ok := p.Args["ids"].([]interface{})
	if !ok {
		return nil, errors.New("invalid inbound service ID")
	}
	var ids []uint
	for _, v := range idsArg {
		if id, ok := helper.ConvGQLStringToUint(v); ok {
			ids = append(ids, id)
		}
	}

	return service.Lookup().DeleteInboundServices(p.Context, ids)
}

func init() {
	schema.AddMutationField("createOrUpdateInboundService", createOrUpdateInboundServiceMutationField)
	schema.AddMutationField("deleteInboundServices", deleteInboundServicesMutationField)
}
