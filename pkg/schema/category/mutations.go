package category

import (
	"errors"

	"github.com/graphql-go/graphql"
	"github.com/ncarlier/readflow/pkg/helper"
	"github.com/ncarlier/readflow/pkg/model"
	"github.com/ncarlier/readflow/pkg/schema"
	"github.com/ncarlier/readflow/pkg/service"
)

var createOrUpdateCategoryMutationField = &graphql.Field{
	Type:        Type,
	Description: "create or update a category (use the ID parameter to update)",
	Args: graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{
			Type: graphql.ID,
		},
		"title": &graphql.ArgumentConfig{
			Type: graphql.String,
		},
		"rule": &graphql.ArgumentConfig{
			Type: graphql.String,
		},
	},
	Resolve: createOrUpdateCategoryResolver,
}

func createOrUpdateCategoryResolver(p graphql.ResolveParams) (interface{}, error) {
	title := helper.GetGQLStringParameter("title", p.Args)
	rule := helper.GetGQLStringParameter("rule", p.Args)
	if id, ok := helper.ConvGQLStringToUint(p.Args["id"]); ok {
		form := model.CategoryUpdateForm{
			ID:    id,
			Title: title,
			Rule:  rule,
		}
		return service.Lookup().UpdateCategory(p.Context, form)
	}
	form := model.CategoryCreateForm{
		Title: *title,
		Rule:  rule,
	}
	return service.Lookup().CreateCategory(p.Context, form)
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
		if id, ok := helper.ConvGQLStringToUint(v); ok {
			ids = append(ids, id)
		}
	}

	return service.Lookup().DeleteCategories(p.Context, ids)
}

func init() {
	schema.AddMutationField("createOrUpdateCategory", createOrUpdateCategoryMutationField)
	schema.AddMutationField("deleteCategories", deleteCategoriesMutationField)
}
