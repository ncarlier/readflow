package api

import (
	"encoding/json"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/ncarlier/readflow/internal/schema"
	"github.com/ncarlier/readflow/pkg/logger"

	// import all GraphQl schema
	_ "github.com/ncarlier/readflow/internal/schema/all"
)

// graphqlHandler is the handler for GraphQL requets.
func graphqlHandler() http.Handler {
	root, err := schema.BuildRootSchema()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to create schema")
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		opts := handler.NewRequestOptions(r)
		params := graphql.Params{
			Schema:         root,
			RequestString:  opts.Query,
			VariableValues: opts.Variables,
			OperationName:  opts.OperationName,
			Context:        ctx,
		}

		result := graphql.Do(params)
		if len(result.Errors) > 0 {
			w.WriteHeader(http.StatusBadRequest)
		}
		json.NewEncoder(w).Encode(result)
	})
}
