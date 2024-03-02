package user

import (
	"errors"
	"strings"

	"github.com/graphql-go/graphql"
	"github.com/ncarlier/readflow/internal/model"
	"github.com/ncarlier/readflow/internal/service"
	"github.com/ncarlier/readflow/pkg/utils"
)

func statusResolver(status string) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		user, ok := p.Source.(*model.User)
		if !ok {
			return nil, errors.New("no user received by read resolver")
		}
		req := model.ArticlesPageRequest{
			Status: &status,
		}
		return service.Lookup().CountUserArticles(p.Context, *user.ID, req)
	}
}

var userType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"username": &graphql.Field{
				Type: graphql.String,
			},
			"enabled": &graphql.Field{
				Type: graphql.Boolean,
			},
			"plan": &graphql.Field{
				Type: graphql.String,
			},
			"customer_id": &graphql.Field{
				Type: graphql.String,
			},
			"last_login_at": &graphql.Field{
				Type: graphql.DateTime,
			},
			"created_at": &graphql.Field{
				Type: graphql.DateTime,
			},
			"updated_at": &graphql.Field{
				Type: graphql.DateTime,
			},
			"inbox": &graphql.Field{
				Type:    graphql.Int,
				Resolve: statusResolver("inbox"),
			},
			"read": &graphql.Field{
				Type:    graphql.Int,
				Resolve: statusResolver("read"),
			},
			"to_read": &graphql.Field{
				Type:    graphql.Int,
				Resolve: statusResolver("to_read"),
			},
			"hash": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					user, ok := p.Source.(*model.User)
					if !ok {
						return nil, errors.New("unsuported type received by hash resolver")
					}
					if user.ID != nil {
						return utils.Hash(strings.ToLower(user.Username)), nil
					}
					return nil, nil
				},
			},
			"hashid": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					user, ok := p.Source.(*model.User)
					if !ok {
						return nil, errors.New("unsuported type received by hashid resolver")
					}
					if user.ID != nil {
						return service.Lookup().GetUserHashID(*user.ID), nil
					}
					return nil, nil
				},
			},
		},
	},
)
