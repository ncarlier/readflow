package helper

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type problemObject struct {
	Title  string `json:"title"`
	Detail string `json:"detail"`
	Status int    `json:"status"`
}

// WriteProblemDetail write error as JSON Problem Details format
func WriteProblemDetail(w http.ResponseWriter, title string, detail string, status int) {
	err := problemObject{
		Title:  title,
		Detail: detail,
		Status: status,
	}
	w.Header().Set("Content-Type", "application/problem+json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(err)
}

// InvalidParameterError create invalid parameter error
func InvalidParameterError(name string) error {
	return fmt.Errorf("a parameter is not valid: %s", name)
}

// RequireParameterError create require parameter error
func RequireParameterError(name string) error {
	return fmt.Errorf("a required parameter is missing: %s", name)
}
