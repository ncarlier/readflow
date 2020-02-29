package admin

import (
	"errors"

	"github.com/graphql-go/graphql"
	"github.com/ncarlier/readflow/pkg/model"
	"github.com/ncarlier/readflow/pkg/service"
	"github.com/ncarlier/readflow/pkg/tooling"
)

var userType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"username": &graphql.Field{
				Type: graphql.String,
			},
			"hash": &graphql.Field{
				Type: graphql.String,
			},
			"enabled": &graphql.Field{
				Type: graphql.Boolean,
			},
			"last_login_at": &graphql.Field{
				Type: graphql.DateTime,
			},
			"created_at": &graphql.Field{
				Type: graphql.DateTime,
			},
			"read": &graphql.Field{
				Type: graphql.Int,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					user, ok := p.Source.(*model.User)
					if !ok {
						return nil, errors.New("no user received by unread resolver")
					}
					status := "read"
					req := model.ArticlesPageRequest{
						Status: &status,
					}
					return service.Lookup().CountUserArticles(p.Context, *user.ID, req)
				},
			},
			"unread": &graphql.Field{
				Type: graphql.Int,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					user, ok := p.Source.(*model.User)
					if !ok {
						return nil, errors.New("no user received by unread resolver")
					}
					status := "unread"
					req := model.ArticlesPageRequest{
						Status: &status,
					}
					return service.Lookup().CountUserArticles(p.Context, *user.ID, req)
				},
			},
		},
	},
)

var userQueryField = &graphql.Field{
	Type: userType,
	Args: graphql.FieldConfigArgument{
		"uid": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.ID),
		},
	},
	Resolve: userResolver,
}

func userResolver(p graphql.ResolveParams) (interface{}, error) {
	uid, ok := tooling.ConvGQLStringToUint(p.Args["uid"])
	if !ok {
		return nil, errors.New("invalid category ID")
	}
	user, err := service.Lookup().GetUserByID(p.Context, uid)
	if err != nil {
		return nil, err
	}
	return user, nil
}
