package admin

import "github.com/graphql-go/graphql"

// QueryFields registry
var QueryFields = graphql.Fields{}

// AddQueryField to query fields
func AddQueryField(name string, field *graphql.Field) {
	QueryFields[name] = field
}

// MutationFields registry
var MutationFields = graphql.Fields{}

// AddMutationField to query fields
func AddMutationField(name string, field *graphql.Field) {
	MutationFields[name] = field
}

// BuildRootSchema build root GraphQl schema
func BuildRootSchema() (graphql.Schema, error) {
	return graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(
			graphql.ObjectConfig{
				Name:   "Query",
				Fields: QueryFields,
			}),
		Mutation: graphql.NewObject(
			graphql.ObjectConfig{
				Name:   "Mutation",
				Fields: MutationFields,
			},
		),
	})
}
