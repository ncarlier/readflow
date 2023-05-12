package schema

import (
	"github.com/graphql-go/graphql"
)

var planType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Plan",
		Fields: graphql.Fields{
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"articles_limit": &graphql.Field{
				Type: graphql.Int,
			},
			"categories_limit": &graphql.Field{
				Type: graphql.Int,
			},
			"incoming_webhooks_limit": &graphql.Field{
				Type: graphql.Int,
			},
			"outgoing_webhooks_limit": &graphql.Field{
				Type: graphql.Int,
			},
		},
	},
)
