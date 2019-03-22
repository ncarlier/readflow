package schema

import (
	"errors"

	"github.com/ncarlier/reader/pkg/model"

	"github.com/ncarlier/reader/pkg/tooling"

	"github.com/graphql-go/graphql"
	"github.com/ncarlier/reader/pkg/service"
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
	archivers, err := service.Lookup().GetArchivers(p.Context)
	if err != nil {
		return nil, err
	}
	return archivers, nil
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
	var id *uint
	if val, ok := tooling.ConvGQLStringToUint(p.Args["id"]); ok {
		id = &val
	}
	alias, _ := p.Args["alias"].(string)
	provider, _ := p.Args["provider"].(string)
	config, _ := p.Args["config"].(string)
	isDefault, _ := p.Args["is_default"].(bool)

	form := model.ArchiverForm{
		ID:        id,
		Alias:     alias,
		Provider:  provider,
		Config:    config,
		IsDefault: isDefault,
	}

	archiver, err := service.Lookup().CreateOrUpdateArchiver(p.Context, form)
	if err != nil {
		return nil, err
	}
	return archiver, nil
}

var deleteArchiverMutationField = &graphql.Field{
	Type:        archiverType,
	Description: "delete a archiver",
	Args: graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.ID),
		},
	},
	Resolve: deleteArchiverResolver,
}

func deleteArchiverResolver(p graphql.ResolveParams) (interface{}, error) {
	id, ok := tooling.ConvGQLStringToUint(p.Args["id"])
	if !ok {
		return nil, errors.New("invalid archiver ID")
	}

	archiver, err := service.Lookup().DeleteArchiver(p.Context, id)
	if err != nil {
		return nil, err
	}
	return archiver, nil
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
