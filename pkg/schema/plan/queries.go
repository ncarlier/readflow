package schema

import (
	"github.com/graphql-go/graphql"
	"github.com/ncarlier/readflow/pkg/schema"
	"github.com/ncarlier/readflow/pkg/service"
)

var plansQueryField = &graphql.Field{
	Type:    graphql.NewList(planType),
	Resolve: plansResolver,
}

func plansResolver(p graphql.ResolveParams) (interface{}, error) {
	return service.Lookup().UserPlans.Plans, nil
}

func init() {
	schema.AddQueryField("plans", plansQueryField)
}
