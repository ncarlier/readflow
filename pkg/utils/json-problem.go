package utils

import (
	"encoding/json"
	"net/http"
)

type JSONProblem struct {
	Title   string                 `json:"title"`
	Detail  string                 `json:"detail"`
	Status  int                    `json:"status"`
	Context map[string]interface{} `json:"context"`
}

// WriteJSONProblem write error as JSON Problem Details format
func WriteJSONProblem(w http.ResponseWriter, problem JSONProblem) {
	if problem.Status == 0 {
		problem.Status = http.StatusInternalServerError
	}

	if problem.Title == "" {
		problem.Title = http.StatusText(problem.Status)
	}
	w.Header().Set("Content-Type", "application/problem+json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(problem.Status)
	json.NewEncoder(w).Encode(problem)
}
