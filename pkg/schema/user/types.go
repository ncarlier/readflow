package user

import (
	"github.com/graphql-go/graphql"
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
			"plan": &graphql.Field{
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
