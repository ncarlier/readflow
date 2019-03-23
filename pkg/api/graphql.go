package api

import (
	"encoding/json"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/ncarlier/reader/pkg/config"
	"github.com/ncarlier/reader/pkg/schema"
)

// graphqlHandler is the handler for GraphQL requets.
func graphqlHandler(conf *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		opts := handler.NewRequestOptions(r)
		params := graphql.Params{
			Schema:         schema.Root,
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
