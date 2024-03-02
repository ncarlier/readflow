package schema

import (
	"github.com/graphql-go/graphql"
	"github.com/ncarlier/readflow/internal/schema"
	"github.com/ncarlier/readflow/internal/service"
)

var plansQueryField = &graphql.Field{
	Type:    graphql.NewList(planType),
	Resolve: plansResolver,
}

func plansResolver(p graphql.ResolveParams) (interface{}, error) {
	return service.Lookup().GetUserPlans(), nil
}

func init() {
	schema.AddQueryField("plans", plansQueryField)
}
