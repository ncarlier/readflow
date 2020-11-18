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
			"total_articles": &graphql.Field{
				Type: graphql.Int,
			},
			"total_categories": &graphql.Field{
				Type: graphql.Int,
			},
		},
	},
)
