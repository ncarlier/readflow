package utils

import (
	"encoding/json"
	"net/http"
)

type problemObject struct {
	Title  string `json:"title"`
	Detail string `json:"detail"`
	Status int    `json:"status"`
}

// WriteJSONProblem write error as JSON Problem Details format
func WriteJSONProblem(w http.ResponseWriter, title, detail string, status int) {
	if title == "" {
		title = http.StatusText(status)
	}
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
