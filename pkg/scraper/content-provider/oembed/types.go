package oembed

//go:generate go run autogen/generate.go
//go:generate gofmt -s -w providers.go

// Endpoint of an oEmbed provider
type Endpoint struct {
	Discovery bool     `json:"discovery,omitempty"`
	URL       string   `json:"url"`
	Schemes   []string `json:"schemes"`
}

// Provider of an oEmbed format
type Provider struct {
	Name      string     `json:"provider_name"`
	URL       string     `json:"provider_url"`
	Endpoints []Endpoint `json:"endpoints"`
}
