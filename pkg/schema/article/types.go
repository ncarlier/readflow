package article

import (
	"errors"

	"github.com/graphql-go/graphql"

	"github.com/ncarlier/readflow/pkg/model"
	"github.com/ncarlier/readflow/pkg/schema/category"
	"github.com/ncarlier/readflow/pkg/service"
)

var sortOrder = graphql.NewEnum(
	graphql.EnumConfig{
		Name:        "sortOrder",
		Description: "Sorting order",
		Values: graphql.EnumValueConfigMap{
			"asc": &graphql.EnumValueConfig{
				Value:       "asc",
				Description: "from older to newer",
			},
			"desc": &graphql.EnumValueConfig{
				Value:       "desc",
				Description: "from newer to older",
			},
		},
	},
)

var sortBy = graphql.NewEnum(
	graphql.EnumConfig{
		Name:        "sortBy",
		Description: "Sorting by",
		Values: graphql.EnumValueConfigMap{
			"key": &graphql.EnumValueConfig{
				Value:       "key",
				Description: "sort by key",
			},
			"stars": &graphql.EnumValueConfig{
				Value:       "stars",
				Description: "sort by stars",
			},
		},
	},
)

var articleStatus = graphql.NewEnum(
	graphql.EnumConfig{
		Name:        "status",
		Description: "Article status",
		Values: graphql.EnumValueConfigMap{
			"read": &graphql.EnumValueConfig{
				Value:       "read",
				Description: "article is read",
			},
			"unread": &graphql.EnumValueConfig{
				Value:       "unread",
				Description: "article is not read",
			},
		},
	},
)

var articleType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Article",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"title": &graphql.Field{
				Type: graphql.String,
			},
			"text": &graphql.Field{
				Type: graphql.String,
			},
			"html": &graphql.Field{
				Type: graphql.String,
			},
			"url": &graphql.Field{
				Type: graphql.String,
			},
			"image": &graphql.Field{
				Type: graphql.String,
			},
			"status": &graphql.Field{
				Type: articleStatus,
			},
			"stars": &graphql.Field{
				Type: graphql.Int,
			},
			"category": &graphql.Field{
				Type: category.Type,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					article, ok := p.Source.(*model.Article)
					if !ok {
						return nil, errors.New("no article received by category resolver")
					}
					if article.CategoryID != nil {
						return service.Lookup().GetCategory(p.Context, *article.CategoryID)
					}
					return nil, nil
				},
			},
			"published_at": &graphql.Field{
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

var articlesType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Articles",
		Fields: graphql.Fields{
			"totalCount": &graphql.Field{
				Type: graphql.Int,
			},
			"endCursor": &graphql.Field{
				Type: graphql.Int,
			},
			"hasNext": &graphql.Field{
				Type: graphql.Boolean,
			},
			"entries": &graphql.Field{
				Type: graphql.NewList(articleType),
			},
		},
	},
)

var articleUpdateResponseType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "ArticleUpdateResponseType",
		Fields: graphql.Fields{
			"article": &graphql.Field{
				Type: articleType,
			},
			"_all": &graphql.Field{
				Type: graphql.Int,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					status := "unread"
					req := model.ArticlesPageRequest{
						Status: &status,
					}
					return service.Lookup().CountCurrentUserArticles(p.Context, req)
				},
			},
			"_starred": &graphql.Field{
				Type: graphql.Int,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					starred := true
					req := model.ArticlesPageRequest{
						Starred: &starred,
					}
					return service.Lookup().CountCurrentUserArticles(p.Context, req)
				},
			},
		},
	},
)
