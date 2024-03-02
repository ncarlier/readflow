package category

import (
	"github.com/graphql-go/graphql"
	"github.com/ncarlier/readflow/internal/global"
	"github.com/ncarlier/readflow/internal/model"
	"github.com/ncarlier/readflow/internal/schema"
	"github.com/ncarlier/readflow/internal/service"
	"github.com/ncarlier/readflow/pkg/utils"
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
	id := utils.ParseGraphQLID(p.Args, "id")
	if id == nil {
		return nil, global.InvalidParameterError("id")
	}
	return service.Lookup().GetCategory(p.Context, *id)
}

func init() {
	schema.AddQueryField("category", categoryQueryField)
	schema.AddQueryField("categories", categoriesQueryField)
}
