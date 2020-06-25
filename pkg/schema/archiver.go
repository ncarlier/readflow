package schema

import (
	"errors"

	"github.com/ncarlier/readflow/pkg/model"

	"github.com/ncarlier/readflow/pkg/tooling"

	"github.com/graphql-go/graphql"
	"github.com/ncarlier/readflow/pkg/service"
)

var providerEnum = graphql.NewEnum(
	graphql.EnumConfig{
		Name:        "provider",
		Description: "Archive provider",
		Values: graphql.EnumValueConfigMap{
			"keeper": &graphql.EnumValueConfig{
				Value:       "keeper",
				Description: "Use Nunux Keeper as archiver provider",
			},
			"webhook": &graphql.EnumValueConfig{
				Value:       "webhook",
				Description: "Use a webhook as archiver provider",
			},
			"wallabag": &graphql.EnumValueConfig{
				Value:       "wallabag",
				Description: "Use Wallabag as archiver provider",
			},
		},
	},
)

var archiverType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Archiver",
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

// QUERIES

var archiversQueryField = &graphql.Field{
	Type:    graphql.NewList(archiverType),
	Resolve: archiversResolver,
}

func archiversResolver(p graphql.ResolveParams) (interface{}, error) {
	return service.Lookup().GetArchivers(p.Context)
}

var archiverQueryField = &graphql.Field{
	Type: archiverType,
	Args: graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{
			Type: graphql.ID,
		},
	},
	Resolve: archiverResolver,
}

func archiverResolver(p graphql.ResolveParams) (interface{}, error) {
	id, ok := tooling.ConvGQLStringToUint(p.Args["id"])
	if !ok {
		return nil, errors.New("invalid archiver ID")
	}
	return service.Lookup().GetArchiver(p.Context, id)
}

// MUTATIONS

var createOrUpdateArchiverMutationField = &graphql.Field{
	Type:        archiverType,
	Description: "create or update a archiver (use the ID parameter to update)",
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
	Resolve: createOrUpdateArchiverResolver,
}

func createOrUpdateArchiverResolver(p graphql.ResolveParams) (interface{}, error) {
	alias := tooling.GetGQLStringParameter("alias", p.Args)
	provider := tooling.GetGQLStringParameter("provider", p.Args)
	config := tooling.GetGQLStringParameter("config", p.Args)
	isDefault := tooling.GetGQLBoolParameter("is_default", p.Args)
	if id, ok := tooling.ConvGQLStringToUint(p.Args["id"]); ok {
		form := model.ArchiverUpdateForm{
			ID:        id,
			Alias:     alias,
			Provider:  provider,
			Config:    config,
			IsDefault: isDefault,
		}
		return service.Lookup().UpdateArchiver(p.Context, form)
	}
	builder := model.NewArchiverCreateFormBuilder()
	builder.Alias(*alias).Provider(*provider).Config(*config)
	if isDefault != nil && *isDefault {
		builder.IsDefault(true)
	}
	form := builder.Build()

	return service.Lookup().CreateArchiver(p.Context, *form)
}

var deleteArchiversMutationField = &graphql.Field{
	Type:        graphql.Int,
	Description: "delete archivers",
	Args: graphql.FieldConfigArgument{
		"ids": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.NewList(graphql.ID)),
		},
	},
	Resolve: deleteArchiversResolver,
}

func deleteArchiversResolver(p graphql.ResolveParams) (interface{}, error) {
	idsArg, ok := p.Args["ids"].([]interface{})
	if !ok {
		return nil, errors.New("invalid archiver ID")
	}
	var ids []uint
	for _, v := range idsArg {
		if id, ok := tooling.ConvGQLStringToUint(v); ok {
			ids = append(ids, id)
		}
	}

	return service.Lookup().DeleteArchivers(p.Context, ids)
}

var archiveArticleMutationField = &graphql.Field{
	Type:        graphql.ID,
	Description: "archive an article",
	Args: graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{
			Description: "article ID",
			Type:        graphql.NewNonNull(graphql.ID),
		},
		"archiver": &graphql.ArgumentConfig{
			Description: "archiver alias (using default if missing)",
			Type:        graphql.String,
		},
	},
	Resolve: archiveArticleResolver,
}

func archiveArticleResolver(p graphql.ResolveParams) (interface{}, error) {
	id, ok := tooling.ConvGQLStringToUint(p.Args["id"])
	if !ok {
		return nil, errors.New("invalid article ID")
	}
	var archiver *string
	if val, ok := p.Args["archiver"]; ok {
		sVal := val.(string)
		archiver = &sVal
	}

	err := service.Lookup().ArchiveArticle(p.Context, id, archiver)
	return id, err
}
