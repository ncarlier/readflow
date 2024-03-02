package user

import (
	"github.com/graphql-go/graphql"
	"github.com/ncarlier/readflow/internal/global"
	"github.com/ncarlier/readflow/internal/model"
	schema "github.com/ncarlier/readflow/internal/schema/admin"
	"github.com/ncarlier/readflow/internal/service"
	"github.com/ncarlier/readflow/pkg/utils"
)

var registerUserMutationField = &graphql.Field{
	Type:        userType,
	Description: "register new user or return it if already exists",
	Args: graphql.FieldConfigArgument{
		"username": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
	},
	Resolve: registerUserResolver,
}

func registerUserResolver(p graphql.ResolveParams) (interface{}, error) {
	username := utils.ParseGraphQLArgument[string](p.Args, "username")
	if username == nil {
		return nil, global.RequireParameterError("username")
	}
	return service.Lookup().GetOrRegisterUser(p.Context, *username)
}

var updateUserMutationField = &graphql.Field{
	Type:        userType,
	Description: "update user account",
	Args: graphql.FieldConfigArgument{
		"uid": &graphql.ArgumentConfig{
			Type: graphql.ID,
		},
		"enabled": &graphql.ArgumentConfig{
			Type: graphql.Boolean,
		},
		"plan": &graphql.ArgumentConfig{
			Type: graphql.String,
		},
		"customer_id": &graphql.ArgumentConfig{
			Type: graphql.String,
		},
	},
	Resolve: updateUserResolver,
}

func updateUserResolver(p graphql.ResolveParams) (interface{}, error) {
	uid := utils.ParseGraphQLID(p.Args, "uid")
	if uid == nil {
		return nil, global.InvalidParameterError("uid")
	}
	form := model.UserForm{
		ID:         *uid,
		Enabled:    utils.ParseGraphQLArgument[bool](p.Args, "enabled"),
		Plan:       utils.ParseGraphQLArgument[string](p.Args, "plan"),
		CustomerID: utils.ParseGraphQLArgument[string](p.Args, "customer_id"),
	}
	return service.Lookup().UpdateUser(p.Context, form)
}

func init() {
	schema.AddMutationField("updateUser", updateUserMutationField)
	schema.AddMutationField("registerUser", registerUserMutationField)
}
