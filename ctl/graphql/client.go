package graphql

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Request represents a GraphQL request
type Request struct {
	Query     string         `json:"query"`
	Variables map[string]any `json:"variables,omitempty"`
}

// Response represents a GraphQL response
type Response struct {
	Data   json.RawMessage `json:"data"`
	Errors []Error         `json:"errors,omitempty"`
}

// Error represents a GraphQL error
type Error struct {
	Message string `json:"message"`
}

// Client is a simple GraphQL HTTP client
type Client struct {
	endpoint   string
	httpClient *http.Client
}

// NewClient creates a new GraphQL client
func NewClient(endpoint string, httpClient *http.Client) *Client {
	return &Client{
		endpoint:   endpoint,
		httpClient: httpClient,
	}
}

// Execute sends a GraphQL request and decodes the response
func (c *Client) Execute(query string, variables map[string]any, target any) error {
	reqBody := Request{
		Query:     query,
		Variables: variables,
	}
	body, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("unable to marshal GraphQL request: %w", err)
	}

	req, err := http.NewRequest("POST", c.endpoint+"/graphql", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("unable to create HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("unable to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("unable to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected HTTP status: %s (body: %s)", resp.Status, string(respBody))
	}

	var gqlResp Response
	if err := json.Unmarshal(respBody, &gqlResp); err != nil {
		return fmt.Errorf("unable to decode GraphQL response: %w", err)
	}

	if len(gqlResp.Errors) > 0 {
		return fmt.Errorf("GraphQL error: %s", gqlResp.Errors[0].Message)
	}

	if target != nil {
		if err := json.Unmarshal(gqlResp.Data, target); err != nil {
			return fmt.Errorf("unable to decode GraphQL data: %w", err)
		}
	}

	return nil
}
