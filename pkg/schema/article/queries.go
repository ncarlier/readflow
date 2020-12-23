package article

import (
	"errors"

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
	},
	Resolve: articlesResolver,
}

func articlesResolver(p graphql.ResolveParams) (interface{}, error) {
	pageRequest := model.ArticlesPageRequest{
		Limit:       helper.GetGQLUintParameter("limit", p.Args),
		SortOrder:   helper.GetGQLStringParameter("sortOrder", p.Args),
		AfterCursor: helper.GetGQLUintParameter("afterCursor", p.Args),
		Category:    helper.GetGQLUintParameter("category", p.Args),
		Status:      helper.GetGQLStringParameter("status", p.Args),
		Starred:     helper.GetGQLBoolParameter("starred", p.Args),
		Query:       helper.GetGQLStringParameter("query", p.Args),
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
	id, ok := helper.ConvGQLStringToUint(p.Args["id"])
	if !ok {
		return nil, errors.New("invalid article ID")
	}

	return service.Lookup().GetArticle(p.Context, id)
}

func init() {
	schema.AddQueryField("articles", articlesQueryField)
	schema.AddQueryField("article", articleQueryField)
}
