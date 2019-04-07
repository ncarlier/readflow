package schema

import (
	"errors"

	"github.com/ncarlier/readflow/pkg/tooling"

	"github.com/graphql-go/graphql"
	"github.com/ncarlier/readflow/pkg/service"
)

var categoryType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Category",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"title": &graphql.Field{
				Type: graphql.String,
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

var categoriesQueryField = &graphql.Field{
	Type:    graphql.NewList(categoryType),
	Resolve: categoriesResolver,
}

func categoriesResolver(p graphql.ResolveParams) (interface{}, error) {
	categories, err := service.Lookup().GetCategories(p.Context)
	if err != nil {
		return nil, err
	}
	return categories, nil
}

var categoryQueryField = &graphql.Field{
	Type: categoryType,
	Args: graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{
			Type: graphql.ID,
		},
	},
	Resolve: categoryResolver,
}

func categoryResolver(p graphql.ResolveParams) (interface{}, error) {
	id, ok := tooling.ConvGQLStringToUint(p.Args["id"])
	if !ok {
		return nil, errors.New("invalid category ID")
	}
	category, err := service.Lookup().GetCategory(p.Context, id)
	if err != nil {
		return nil, err
	}
	return category, nil
}

// MUTATIONS

var createOrUpdateCategoryMutationField = &graphql.Field{
	Type:        categoryType,
	Description: "create or update a category (use the ID parameter to update)",
	Args: graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{
			Type: graphql.ID,
		},
		"title": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
	},
	Resolve: createOrUpdateCategoryResolver,
}

func createOrUpdateCategoryResolver(p graphql.ResolveParams) (interface{}, error) {
	var id *uint
	if val, ok := tooling.ConvGQLStringToUint(p.Args["id"]); ok {
		id = &val
	}
	title, _ := p.Args["title"].(string)

	category, err := service.Lookup().CreateOrUpdateCategory(p.Context, id, title)
	if err != nil {
		return nil, err
	}
	return category, nil
}

var deleteCategoriesMutationField = &graphql.Field{
	Type:        graphql.Int,
	Description: "delete categories",
	Args: graphql.FieldConfigArgument{
		"ids": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.NewList(graphql.ID)),
		},
	},
	Resolve: deleteCategoriesResolver,
}

func deleteCategoriesResolver(p graphql.ResolveParams) (interface{}, error) {
	idsArg, ok := p.Args["ids"].([]interface{})
	if !ok {
		return nil, errors.New("invalid category ID")
	}
	var ids []uint
	for _, v := range idsArg {
		if id, ok := tooling.ConvGQLStringToUint(v); ok {
			ids = append(ids, id)
		}
	}

	nb, err := service.Lookup().DeleteCategories(p.Context, ids)
	if err != nil {
		return nil, err
	}
	return nb, nil
}
