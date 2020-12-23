package user

import (
	"errors"

	"github.com/graphql-go/graphql"
	"github.com/ncarlier/readflow/pkg/helper"
	"github.com/ncarlier/readflow/pkg/model"
	schema "github.com/ncarlier/readflow/pkg/schema/admin"
	"github.com/ncarlier/readflow/pkg/service"
)

var updateUserMutationField = &graphql.Field{
	Type:        userType,
	Description: "delete account and all relative data",
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
	return service.Lookup().UpdateUser(p.Context, form)
}

func init() {
	schema.AddMutationField("updateUser", updateUserMutationField)
}
