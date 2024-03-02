package user

import (
	"github.com/graphql-go/graphql"
	"github.com/ncarlier/readflow/internal/schema"
	"github.com/ncarlier/readflow/internal/service"
)

var deleteAccountMutationField = &graphql.Field{
	Type:        graphql.Boolean,
	Description: "delete account and all relative data",
	Resolve:     deleteAccountResolver,
}

func deleteAccountResolver(p graphql.ResolveParams) (interface{}, error) {
	return service.Lookup().DeleteAccount(p.Context)
}

func init() {
	schema.AddMutationField("deleteAccount", deleteAccountMutationField)
}
