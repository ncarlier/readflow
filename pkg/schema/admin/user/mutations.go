package user

import (
	"github.com/graphql-go/graphql"
	"github.com/ncarlier/readflow/pkg/helper"
	"github.com/ncarlier/readflow/pkg/model"
	schema "github.com/ncarlier/readflow/pkg/schema/admin"
	"github.com/ncarlier/readflow/pkg/service"
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
	username := helper.ParseGraphQLArgument[string](p.Args, "username")
	if username == nil {
		return nil, helper.RequireParameterError("username")
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
	uid := helper.ParseGraphQLID(p.Args, "uid")
	if uid == nil {
		return nil, helper.InvalidParameterError("uid")
	}
	form := model.UserForm{
		ID:         *uid,
		Enabled:    helper.ParseGraphQLArgument[bool](p.Args, "enabled"),
		Plan:       helper.ParseGraphQLArgument[string](p.Args, "plan"),
		CustomerID: helper.ParseGraphQLArgument[string](p.Args, "customer_id"),
	}
	return service.Lookup().UpdateUser(p.Context, form)
}

func init() {
	schema.AddMutationField("updateUser", updateUserMutationField)
	schema.AddMutationField("registerUser", registerUserMutationField)
}
