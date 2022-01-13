package category

import (
	"errors"

	"github.com/graphql-go/graphql"
	"github.com/ncarlier/readflow/pkg/model"
	"github.com/ncarlier/readflow/pkg/service"
)

var notificationStrategy = graphql.NewEnum(
	graphql.EnumConfig{
		Name:        "notificationStrategy",
		Description: "Notification strategy",
		Values: graphql.EnumValueConfigMap{
			"none": &graphql.EnumValueConfig{
				Value:       "none",
				Description: "no notification will be sent",
			},
			"individual": &graphql.EnumValueConfig{
				Value:       "individual",
				Description: "a notification will be sent as soon as an article is received",
			},
			"global": &graphql.EnumValueConfig{
				Value:       "global",
				Description: "a notification will be sent using the global strategy",
			},
		},
	},
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
				Type:        graphql.String,
				Description: "title of the category",
			},
			"rule": &graphql.Field{
				Type:        graphql.String,
				Description: "rule definition to put articles into this category",
			},
			"notification_strategy": &graphql.Field{
				Type:        notificationStrategy,
				Description: "notification strategy for this category",
			},
			"inbox": &graphql.Field{
				Type:        graphql.Int,
				Description: "number of received articles for this category",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					var category = model.Category{}
					switch p.Source.(type) {
					case model.Category:
						category, _ = p.Source.(model.Category)
					case *model.Category:
						cat, _ := p.Source.(*model.Category)
						category = *cat
					default:
						return nil, errors.New("no category received by inbox resolver")
					}
					status := "inbox"
					req := model.ArticlesPageRequest{
						Category: category.ID,
						Status:   &status,
					}
					return service.Lookup().CountCurrentUserArticles(p.Context, req)
				},
			},
			"created_at": &graphql.Field{
				Type:        graphql.DateTime,
				Description: "creation date",
			},
			"updated_at": &graphql.Field{
				Type:        graphql.DateTime,
				Description: "modification date",
			},
		},
	},
)

func statusResolver(status string) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		req := model.ArticlesPageRequest{
			Status: &status,
		}
		return service.Lookup().CountCurrentUserArticles(p.Context, req)
	}
}

// ListResponseType is a list of categories
var ListResponseType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Categories",
		Fields: graphql.Fields{
			"_inbox": &graphql.Field{
				Type:    graphql.Int,
				Resolve: statusResolver("inbox"),
			},
			"_to_read": &graphql.Field{
				Type:    graphql.Int,
				Resolve: statusResolver("to_read"),
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
