package schema

import (
	"github.com/graphql-go/graphql"
	"github.com/ncarlier/reader/pkg/constant"
	"github.com/rs/zerolog/log"
)

var Schema graphql.Schema

var userType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

var queryType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"me": &graphql.Field{
				Type: userType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					user := struct {
						ID   string `json:"id"`
						Name string `json:"name"`
					}{"1", p.Context.Value(constant.Username).(string)}
					return user, nil
				},
			},
		},
	})

func init() {
	s, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: queryType,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create schema")
	}
	Schema = s
}
