package mcp

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/ncarlier/readflow/ctl/graphql"
)

func registerResources(s *server.MCPServer, client *graphql.Client) {
	// Static resource: all categories
	s.AddResource(
		mcp.NewResource(
			"readflow://categories",
			"All Categories",
			mcp.WithResourceDescription("Overview of all categories with article counts"),
			mcp.WithMIMEType("application/json"),
		),
		allCategoriesHandler(client),
	)

	// Resource template: single article by ID
	s.AddResourceTemplate(
		mcp.NewResourceTemplate(
			"readflow://articles/{id}",
			"Article",
			mcp.WithTemplateDescription("A single article with full content"),
			mcp.WithTemplateMIMEType("application/json"),
		),
		articleResourceHandler(client),
	)

	// Resource template: single category by ID
	s.AddResourceTemplate(
		mcp.NewResourceTemplate(
			"readflow://categories/{id}",
			"Category",
			mcp.WithTemplateDescription("A single category with details"),
			mcp.WithTemplateMIMEType("application/json"),
		),
		categoryResourceHandler(client),
	)
}

func allCategoriesHandler(client *graphql.Client) server.ResourceHandlerFunc {
	return func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		q := `query {
			categories {
				_inbox
				_to_read
				_starred
				entries {
					id
					title
					inbox
					created_at
					updated_at
				}
			}
		}`

		var result struct {
			Categories any `json:"categories"`
		}
		if err := queryGraphQL(client, q, nil, &result); err != nil {
			return nil, fmt.Errorf("unable to fetch categories: %w", err)
		}

		json, err := toJSON(result.Categories)
		if err != nil {
			return nil, fmt.Errorf("unable to marshal categories: %w", err)
		}

		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      "readflow://categories",
				MIMEType: "application/json",
				Text:     json,
			},
		}, nil
	}
}

func articleResourceHandler(client *graphql.Client) server.ResourceTemplateHandlerFunc {
	return func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		id, err := extractIDFromURI(request.Params.URI, "readflow://articles/")
		if err != nil {
			return nil, err
		}

		q := `query article($id: ID!) {
			article(id: $id) {
				id
				title
				text
				url
				image
				status
				stars
				category { id title }
				published_at
				created_at
				updated_at
			}
		}`

		var result struct {
			Article any `json:"article"`
		}
		if err := queryGraphQL(client, q, map[string]any{"id": id}, &result); err != nil {
			return nil, fmt.Errorf("unable to fetch article: %w", err)
		}

		json, err := toJSON(result.Article)
		if err != nil {
			return nil, fmt.Errorf("unable to marshal article: %w", err)
		}

		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      request.Params.URI,
				MIMEType: "application/json",
				Text:     json,
			},
		}, nil
	}
}

func categoryResourceHandler(client *graphql.Client) server.ResourceTemplateHandlerFunc {
	return func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		id, err := extractIDFromURI(request.Params.URI, "readflow://categories/")
		if err != nil {
			return nil, err
		}

		q := `query category($id: ID!) {
			category(id: $id) {
				id
				title
				inbox
				created_at
				updated_at
			}
		}`

		var result struct {
			Category any `json:"category"`
		}
		if err := queryGraphQL(client, q, map[string]any{"id": id}, &result); err != nil {
			return nil, fmt.Errorf("unable to fetch category: %w", err)
		}

		json, err := toJSON(result.Category)
		if err != nil {
			return nil, fmt.Errorf("unable to marshal category: %w", err)
		}

		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      request.Params.URI,
				MIMEType: "application/json",
				Text:     json,
			},
		}, nil
	}
}

// extractIDFromURI extracts the ID portion from a URI given its prefix.
func extractIDFromURI(uri, prefix string) (string, error) {
	if !strings.HasPrefix(uri, prefix) {
		return "", fmt.Errorf("invalid URI: %s", uri)
	}
	id := strings.TrimPrefix(uri, prefix)
	if id == "" {
		return "", fmt.Errorf("missing ID in URI: %s", uri)
	}
	return id, nil
}
