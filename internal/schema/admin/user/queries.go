package user

import (
	"errors"

	"github.com/graphql-go/graphql"
	schema "github.com/ncarlier/readflow/internal/schema/admin"
	"github.com/ncarlier/readflow/internal/service"
	"github.com/ncarlier/readflow/pkg/utils"
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
	uid := utils.ParseGraphQLID(p.Args, "uid")
	username := utils.ParseGraphQLArgument[string](p.Args, "username")
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
