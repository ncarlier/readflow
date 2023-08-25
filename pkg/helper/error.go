package helper

import (
	"encoding/json"
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
