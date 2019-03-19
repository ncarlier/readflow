package schema

import (
	"github.com/graphql-go/graphql"
	"github.com/ncarlier/reader/pkg/service"
)

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
			"last_login_at": &graphql.Field{
				Type: graphql.DateTime,
			},
		},
	},
)

var meQueryField = &graphql.Field{
	Type:    userType,
	Resolve: meResolver,
}

func meResolver(p graphql.ResolveParams) (interface{}, error) {
	user, err := service.Lookup().GetCurrentUser(p.Context)
	if err != nil {
		return nil, err
	}
	return user, nil
}
