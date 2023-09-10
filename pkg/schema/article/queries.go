package article

import (
	"github.com/graphql-go/graphql"

	"github.com/ncarlier/readflow/pkg/helper"
	"github.com/ncarlier/readflow/pkg/model"
	"github.com/ncarlier/readflow/pkg/schema"
	"github.com/ncarlier/readflow/pkg/service"
)

var articlesQueryField = &graphql.Field{
	Type: articlesType,
	Args: graphql.FieldConfigArgument{
		"limit": &graphql.ArgumentConfig{
			Description:  "max number of entries to returns",
			Type:         graphql.Int,
			DefaultValue: 10,
		},
		"afterCursor": &graphql.ArgumentConfig{
			Description: "retrive entries after this cursor",
			Type:        graphql.Int,
		},
		"category": &graphql.ArgumentConfig{
			Description: "filter entries by this category",
			Type:        graphql.Int,
		},
		"status": &graphql.ArgumentConfig{
			Description: "filter entries by this status",
			Type:        articleStatus,
		},
		"starred": &graphql.ArgumentConfig{
			Description: "filter entries by this starred value",
			Type:        graphql.Boolean,
		},
		"query": &graphql.ArgumentConfig{
			Description: "filter entries by full-text search",
			Type:        graphql.String,
		},
		"sortOrder": &graphql.ArgumentConfig{
			Description:  "sorting order of the entries",
			Type:         sortOrder,
			DefaultValue: "asc",
		},
		"sortBy": &graphql.ArgumentConfig{
			Description:  "sorting attribute of the entries",
			Type:         sortBy,
			DefaultValue: "key",
		},
	},
	Resolve: articlesResolver,
}

func articlesResolver(p graphql.ResolveParams) (interface{}, error) {
	pageRequest := model.ArticlesPageRequest{
		Limit:       helper.ParseGraphQLArgument[int](p.Args, "limit"),
		SortOrder:   helper.ParseGraphQLArgument[string](p.Args, "sortOrder"),
		SortBy:      helper.ParseGraphQLArgument[string](p.Args, "sortBy"),
		AfterCursor: helper.ParseGraphQLID(p.Args, "afterCursor"),
		Category:    helper.ParseGraphQLID(p.Args, "category"),
		Status:      helper.ParseGraphQLArgument[string](p.Args, "status"),
		Starred:     helper.ParseGraphQLArgument[bool](p.Args, "starred"),
		Query:       helper.ParseGraphQLArgument[string](p.Args, "query"),
	}

	return service.Lookup().GetArticles(p.Context, pageRequest)
}

var articleQueryField = &graphql.Field{
	Type: articleType,
	Args: graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.ID),
		},
	},
	Resolve: articleResolver,
}

func articleResolver(p graphql.ResolveParams) (interface{}, error) {
	id := helper.ParseGraphQLID(p.Args, "id")
	if id == nil {
		return nil, helper.InvalidParameterError("id")
	}

	return service.Lookup().GetArticle(p.Context, *id)
}

func init() {
	schema.AddQueryField("articles", articlesQueryField)
	schema.AddQueryField("article", articleQueryField)
}
