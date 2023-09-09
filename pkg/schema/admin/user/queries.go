package user

import (
	"errors"

	"github.com/graphql-go/graphql"
	"github.com/ncarlier/readflow/pkg/helper"
	schema "github.com/ncarlier/readflow/pkg/schema/admin"
	"github.com/ncarlier/readflow/pkg/service"
)

var userQueryField = &graphql.Field{
	Type: userType,
	Args: graphql.FieldConfigArgument{
		"uid": &graphql.ArgumentConfig{
			Type: graphql.ID,
		},
		"username": &graphql.ArgumentConfig{
			Type: graphql.String,
		},
	},
	Resolve: userResolver,
}

func userResolver(p graphql.ResolveParams) (interface{}, error) {
	uid := helper.ParseGraphQLID(p.Args, "uid")
	username := helper.ParseGraphQLArgument[string](p.Args, "username")
	if uid != nil {
		return service.Lookup().GetUserByID(p.Context, *uid)
	} else if username != nil {
		return service.Lookup().GetUserByUsername(p.Context, *username)
	}
	return nil, errors.New("missing uid or username parameter")
}

func init() {
	schema.AddQueryField("user", userQueryField)
}
