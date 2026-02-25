package mcp

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/ncarlier/readflow/ctl/graphql"
)

func registerPrompts(s *server.MCPServer, client *graphql.Client) {
	s.AddPrompt(summarizeArticlePrompt(client))
	s.AddPrompt(suggestCategoriesPrompt(client))
}

func summarizeArticlePrompt(client *graphql.Client) (mcp.Prompt, server.PromptHandlerFunc) {
	prompt := mcp.NewPrompt("summarize_article",
		mcp.WithPromptDescription("Summarize an article's content"),
		mcp.WithArgument("id",
			mcp.ArgumentDescription("Article ID"),
			mcp.RequiredArgument(),
		),
		mcp.WithArgument("style",
			mcp.ArgumentDescription("Summary style: brief or detailed (default: brief)"),
		),
	)

	handler := func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		id, ok := request.Params.Arguments["id"]
		if !ok || id == "" {
			return nil, fmt.Errorf("missing required argument: id")
		}

		style := request.Params.Arguments["style"]
		if style == "" {
			style = "brief"
		}

		q := `query article($id: ID!) {
			article(id: $id) {
				id
				title
				text
				url
				category { id title }
			}
		}`

		var result struct {
			Article struct {
				ID       int    `json:"id"`
				Title    string `json:"title"`
				Text     string `json:"text"`
				URL      string `json:"url"`
				Category *struct {
					ID    int    `json:"id"`
					Title string `json:"title"`
				} `json:"category"`
			} `json:"article"`
		}
		if err := queryGraphQL(client, q, map[string]any{"id": id}, &result); err != nil {
			return nil, fmt.Errorf("unable to fetch article: %w", err)
		}

		a := result.Article
		categoryInfo := ""
		if a.Category != nil {
			categoryInfo = fmt.Sprintf("\nCategory: %s", a.Category.Title)
		}

		var styleInstruction string
		if style == "detailed" {
			styleInstruction = "Provide a detailed summary covering the main points, key arguments, and conclusions."
		} else {
			styleInstruction = "Provide a brief summary in 2-3 sentences capturing the main idea."
		}

		content := fmt.Sprintf(`Please summarize the following article.

Title: %s
URL: %s%s

Content:
%s

%s`, a.Title, a.URL, categoryInfo, a.Text, styleInstruction)

		return mcp.NewGetPromptResult(
			fmt.Sprintf("Summarize article: %s", a.Title),
			[]mcp.PromptMessage{
				mcp.NewPromptMessage(mcp.RoleUser, mcp.NewTextContent(content)),
			},
		), nil
	}

	return prompt, handler
}

func suggestCategoriesPrompt(client *graphql.Client) (mcp.Prompt, server.PromptHandlerFunc) {
	prompt := mcp.NewPrompt("suggest_categories",
		mcp.WithPromptDescription("Suggest the best category for an article based on its content"),
		mcp.WithArgument("id",
			mcp.ArgumentDescription("Article ID"),
			mcp.RequiredArgument(),
		),
	)

	handler := func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		id, ok := request.Params.Arguments["id"]
		if !ok || id == "" {
			return nil, fmt.Errorf("missing required argument: id")
		}

		// Fetch article
		articleQuery := `query article($id: ID!) {
			article(id: $id) {
				id
				title
				text
				url
				category { id title }
			}
		}`

		var articleResult struct {
			Article struct {
				ID       int    `json:"id"`
				Title    string `json:"title"`
				Text     string `json:"text"`
				URL      string `json:"url"`
				Category *struct {
					ID    int    `json:"id"`
					Title string `json:"title"`
				} `json:"category"`
			} `json:"article"`
		}
		if err := queryGraphQL(client, articleQuery, map[string]any{"id": id}, &articleResult); err != nil {
			return nil, fmt.Errorf("unable to fetch article: %w", err)
		}

		// Fetch categories
		categoriesQuery := `query {
			categories {
				entries {
					id
					title
				}
			}
		}`

		var categoriesResult struct {
			Categories struct {
				Entries []struct {
					ID    int    `json:"id"`
					Title string `json:"title"`
				} `json:"entries"`
			} `json:"categories"`
		}
		if err := queryGraphQL(client, categoriesQuery, nil, &categoriesResult); err != nil {
			return nil, fmt.Errorf("unable to fetch categories: %w", err)
		}

		a := articleResult.Article
		currentCategory := "none"
		if a.Category != nil {
			currentCategory = a.Category.Title
		}

		categoriesList := ""
		for _, c := range categoriesResult.Categories.Entries {
			categoriesList += fmt.Sprintf("- [%d] %s\n", c.ID, c.Title)
		}

		content := fmt.Sprintf(`Based on the following article, suggest the most appropriate category from the available categories.

Article Title: %s
Article URL: %s
Current Category: %s

Article Content:
%s

Available Categories:
%s
Please suggest the best category for this article and explain your reasoning. If none of the existing categories fit well, suggest a new category name.`, a.Title, a.URL, currentCategory, a.Text, categoriesList)

		return mcp.NewGetPromptResult(
			fmt.Sprintf("Suggest category for article: %s", a.Title),
			[]mcp.PromptMessage{
				mcp.NewPromptMessage(mcp.RoleUser, mcp.NewTextContent(content)),
			},
		), nil
	}

	return prompt, handler
}
