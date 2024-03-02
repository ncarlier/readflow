package user

import (
	"github.com/graphql-go/graphql"
	"github.com/ncarlier/readflow/internal/schema"
	"github.com/ncarlier/readflow/internal/service"
)

var meQueryField = &graphql.Field{
	Type:    userType,
	Resolve: meResolver,
}

func meResolver(p graphql.ResolveParams) (interface{}, error) {
	return service.Lookup().GetCurrentUser(p.Context)
}

func init() {
	schema.AddQueryField("me", meQueryField)
}
