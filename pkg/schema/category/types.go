package category

import (
	"errors"

	"github.com/graphql-go/graphql"
	"github.com/ncarlier/readflow/pkg/model"
	"github.com/ncarlier/readflow/pkg/service"
)

// Type of a category
var Type = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Category",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"title": &graphql.Field{
				Type: graphql.String,
			},
			"rule": &graphql.Field{
				Type: graphql.String,
			},
			"unread": &graphql.Field{
				Type: graphql.Int,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					var category = model.Category{}
					switch p.Source.(type) {
					case model.Category:
						category, _ = p.Source.(model.Category)
						break
					case *model.Category:
						cat, _ := p.Source.(*model.Category)
						category = *cat
						break
					default:
						return nil, errors.New("no category received by unread resolver")
					}
					status := "unread"
					req := model.ArticlesPageRequest{
						Category: category.ID,
						Status:   &status,
					}
					return service.Lookup().CountCurrentUserArticles(p.Context, req)
				},
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

// ListResponseType is a list of categories
var ListResponseType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Categories",
		Fields: graphql.Fields{
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
			"entries": &graphql.Field{
				Type: graphql.NewList(Type),
			},
		},
	},
)
