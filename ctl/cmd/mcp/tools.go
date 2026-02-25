package mcp

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/ncarlier/readflow/ctl/graphql"
)

func registerTools(s *server.MCPServer, client *graphql.Client) {
	s.AddTools(
		listArticlesTool(client),
		getArticleTool(client),
		listCategoriesTool(client),
		getCategoryTool(client),
	)
}

func listArticlesTool(client *graphql.Client) server.ServerTool {
	tool := mcp.NewTool("list_articles",
		mcp.WithDescription("List articles with optional filters and pagination"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithNumber("limit",
			mcp.Description("Max number of entries to return (default 10)"),
			mcp.Min(1),
			mcp.Max(100),
		),
		mcp.WithNumber("afterCursor",
			mcp.Description("Cursor for pagination (article ID to start after)"),
		),
		mcp.WithNumber("category",
			mcp.Description("Filter by category ID"),
		),
		mcp.WithString("status",
			mcp.Description("Filter by status"),
			mcp.Enum("inbox", "read", "to_read"),
		),
		mcp.WithBoolean("starred",
			mcp.Description("Filter starred articles only"),
		),
		mcp.WithString("query",
			mcp.Description("Full-text search query"),
		),
		mcp.WithString("sortOrder",
			mcp.Description("Sort order (default asc)"),
			mcp.Enum("asc", "desc"),
		),
		mcp.WithString("sortBy",
			mcp.Description("Sort field (default key)"),
			mcp.Enum("key", "stars"),
		),
	)

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		vars := map[string]any{}

		if v, err := request.RequireInt("limit"); err == nil {
			vars["limit"] = v
		}
		if v, err := request.RequireInt("afterCursor"); err == nil {
			vars["afterCursor"] = v
		}
		if v, err := request.RequireInt("category"); err == nil {
			vars["category"] = v
		}
		if v, err := request.RequireString("status"); err == nil {
			vars["status"] = v
		}
		if v, err := request.RequireBool("starred"); err == nil && v {
			vars["starred"] = true
		}
		if v, err := request.RequireString("query"); err == nil {
			vars["query"] = v
		}
		if v, err := request.RequireString("sortOrder"); err == nil {
			vars["sortOrder"] = v
		}
		if v, err := request.RequireString("sortBy"); err == nil {
			vars["sortBy"] = v
		}

		q := `query articles($limit: Int, $afterCursor: Int, $category: Int, $status: status, $starred: Boolean, $query: String, $sortOrder: sortOrder, $sortBy: sortBy) {
			articles(limit: $limit, afterCursor: $afterCursor, category: $category, status: $status, starred: $starred, query: $query, sortOrder: $sortOrder, sortBy: $sortBy) {
				totalCount
				endCursor
				hasNext
				entries {
					id
					title
					status
					stars
					url
					category { id title }
					created_at
				}
			}
		}`

		var result struct {
			Articles any `json:"articles"`
		}
		if err := queryGraphQL(client, q, vars, &result); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		json, err := toJSON(result.Articles)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(json), nil
	}

	return server.ServerTool{Tool: tool, Handler: handler}
}

func getArticleTool(client *graphql.Client) server.ServerTool {
	tool := mcp.NewTool("get_article",
		mcp.WithDescription("Get a single article by ID with full content"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithNumber("id",
			mcp.Description("Article ID"),
			mcp.Required(),
		),
		mcp.WithBoolean("html",
			mcp.Description("Include HTML content in response"),
		),
		mcp.WithBoolean("markdown",
			mcp.Description("Include Markdown content in response"),
		),
	)

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id, err := request.RequireInt("id")
		if err != nil {
			return mcp.NewToolResultError("missing required parameter: id"), nil
		}

		fields := "id title text url image status stars category { id title } published_at created_at updated_at"
		if v, err := request.RequireBool("html"); err == nil && v {
			fields += " html"
		}
		if v, err := request.RequireBool("markdown"); err == nil && v {
			fields += " markdown"
		}
		q := fmt.Sprintf("query article($id: ID!) { article(id: $id) { %s } }", fields)

		var result struct {
			Article any `json:"article"`
		}
		if err := queryGraphQL(client, q, map[string]any{"id": id}, &result); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		json, err := toJSON(result.Article)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(json), nil
	}

	return server.ServerTool{Tool: tool, Handler: handler}
}

func listCategoriesTool(client *graphql.Client) server.ServerTool {
	tool := mcp.NewTool("list_categories",
		mcp.WithDescription("List all categories with article counts"),
		mcp.WithReadOnlyHintAnnotation(true),
	)

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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
			return mcp.NewToolResultError(err.Error()), nil
		}

		json, err := toJSON(result.Categories)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(json), nil
	}

	return server.ServerTool{Tool: tool, Handler: handler}
}

func getCategoryTool(client *graphql.Client) server.ServerTool {
	tool := mcp.NewTool("get_category",
		mcp.WithDescription("Get a single category by ID"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithNumber("id",
			mcp.Description("Category ID"),
			mcp.Required(),
		),
	)

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id, err := request.RequireInt("id")
		if err != nil {
			return mcp.NewToolResultError("missing required parameter: id"), nil
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
			return mcp.NewToolResultError(err.Error()), nil
		}

		json, err := toJSON(result.Category)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(json), nil
	}

	return server.ServerTool{Tool: tool, Handler: handler}
}
