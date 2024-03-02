package account

import (
	"net/http"
)

// Provider used for account linking
type Provider interface {
	RequestHandler(w http.ResponseWriter, r *http.Request) error
	AuthorizeHandler(w http.ResponseWriter, r *http.Request) error
}
