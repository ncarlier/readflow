package middleware

import (
	"encoding/json"
	"net/http"
)

type errorObject struct {
	Message string `json:"message"`
}

type errorsObject struct {
	Errors []errorObject `json:"errors"`
}

func jsonErrors(w http.ResponseWriter, message string, code int) {
	err := errorsObject{
		Errors: []errorObject{{Message: message}},
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(err)
}
