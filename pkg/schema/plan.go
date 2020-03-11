package schema

import (
	"github.com/graphql-go/graphql"
	"github.com/ncarlier/readflow/pkg/service"
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
		},
	},
)

// QUERIES

var plansQueryField = &graphql.Field{
	Type:    graphql.NewList(planType),
	Resolve: plansResolver,
}

func plansResolver(p graphql.ResolveParams) (interface{}, error) {
	return service.Lookup().UserPlans.GetPlans(), nil
}
