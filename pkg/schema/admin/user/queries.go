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
			Type: graphql.NewNonNull(graphql.ID),
		},
	},
	Resolve: userResolver,
}

func userResolver(p graphql.ResolveParams) (interface{}, error) {
	uid, ok := helper.ConvGQLStringToUint(p.Args["uid"])
	if !ok {
		return nil, errors.New("invalid user ID")
	}
	user, err := service.Lookup().GetUserByID(p.Context, uid)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func init() {
	schema.AddQueryField("user", userQueryField)
}
