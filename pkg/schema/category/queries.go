package category

import (
	"github.com/graphql-go/graphql"
	"github.com/ncarlier/readflow/pkg/helper"
	"github.com/ncarlier/readflow/pkg/model"
	"github.com/ncarlier/readflow/pkg/schema"
	"github.com/ncarlier/readflow/pkg/service"
)

var categoriesQueryField = &graphql.Field{
	Type:    ListResponseType,
	Resolve: CategoriesResolver,
}

// CategoriesResolver is the resolver for retrieve categories
func CategoriesResolver(p graphql.ResolveParams) (interface{}, error) {
	categories, err := service.Lookup().GetCategories(p.Context)
	if err != nil {
		return nil, err
	}
	return struct {
		Entries []model.Category
	}{
		categories,
	}, nil
}

var categoryQueryField = &graphql.Field{
	Type: Type,
	Args: graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{
			Type: graphql.ID,
		},
	},
	Resolve: categoryResolver,
}

func categoryResolver(p graphql.ResolveParams) (interface{}, error) {
	id := helper.ParseGraphQLID(p.Args, "id")
	if id == nil {
		return nil, helper.InvalidParameterError("id")
	}
	return service.Lookup().GetCategory(p.Context, *id)
}

func init() {
	schema.AddQueryField("category", categoryQueryField)
	schema.AddQueryField("categories", categoriesQueryField)
}
