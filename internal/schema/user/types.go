package user

import (
	"errors"
	"strings"

	"github.com/graphql-go/graphql"

	"github.com/ncarlier/readflow/internal/model"
	"github.com/ncarlier/readflow/internal/service"
	"github.com/ncarlier/readflow/pkg/utils"
)

var userType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"username": &graphql.Field{
				Type: graphql.String,
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
		},
	},
)
