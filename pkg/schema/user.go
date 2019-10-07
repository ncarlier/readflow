package schema

import (
	"github.com/graphql-go/graphql"
	"github.com/ncarlier/readflow/pkg/service"
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
			"last_login_at": &graphql.Field{
				Type: graphql.DateTime,
			},
			"created_at": &graphql.Field{
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

var deleteAccountMutationField = &graphql.Field{
	Type:        graphql.Boolean,
	Description: "delete account and all relative data",
	Resolve:     deleteAccountResolver,
}

func deleteAccountResolver(p graphql.ResolveParams) (interface{}, error) {
	return service.Lookup().DeleteAccount(p.Context)
}
