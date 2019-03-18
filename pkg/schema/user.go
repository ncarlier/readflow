package schema

import (
	"fmt"

	"github.com/graphql-go/graphql"
	"github.com/ncarlier/reader/pkg/constant"
	"github.com/ncarlier/reader/pkg/service"
)

var userType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"username": &graphql.Field{
				Type: graphql.String,
			},
			"last_login_at": &graphql.Field{
				Type: graphql.DateTime,
			},
		},
	},
)

var meQueryField = &graphql.Field{
	Type:    userType,
	Resolve: meResolver,
}

func meResolver(p graphql.ResolveParams) (interface{}, error) {
	uid := p.Context.Value(constant.UserID).(uint32)
	user, err := service.Lookup().GetUserByID(p.Context, uid)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("User not found: %d", uid)
	}
	return user, nil
}
