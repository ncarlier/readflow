package user

import (
	"errors"

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
	username := helper.GetGQLStringParameter("username", p.Args)
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
	uid, ok := helper.ConvGQLStringToUint(p.Args["uid"])
	if !ok {
		return nil, errors.New("invalid user ID")
	}
	form := model.UserForm{
		ID: uid,
	}
	if val, ok := p.Args["enabled"]; ok {
		b := val.(bool)
		form.Enabled = &b
	}
	if val, ok := p.Args["plan"]; ok {
		s := val.(string)
		form.Plan = &s
	}
	if val, ok := p.Args["customer_id"]; ok {
		s := val.(string)
		form.CustomerID = &s
	}
	return service.Lookup().UpdateUser(p.Context, form)
}

func init() {
	schema.AddMutationField("updateUser", updateUserMutationField)
	schema.AddMutationField("registerUser", registerUserMutationField)
}
