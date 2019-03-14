package api

import (
	"encoding/json"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/ncarlier/reader/pkg/config"
	"github.com/ncarlier/reader/pkg/schema"
	"github.com/rs/zerolog/log"
)

// graphqlHandler is the handler for GraphQL requets.
func graphqlHandler(conf *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		query := r.URL.Query().Get("query")
		log.Debug().Str("query", query).Msg("GraphQL request")

		result := graphql.Do(graphql.Params{
			Schema:        schema.Schema,
			RequestString: query,
			Context:       ctx,
		})

		if len(result.Errors) > 0 {
			http.Error(w, result.Errors[0].Error(), http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(result)
	})
}
