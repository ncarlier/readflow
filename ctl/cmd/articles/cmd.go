package articles

import (
	"flag"
	"fmt"
	"strconv"

	"github.com/ncarlier/readflow/ctl/cmd"
	"github.com/ncarlier/readflow/ctl/cmd/output"
	"github.com/ncarlier/readflow/ctl/graphql"
)

const cmdName = "articles"

type articlesCmd struct {
	flagSet *flag.FlagSet
}

func (c *articlesCmd) Exec(args []string, client *graphql.Client) error {
	subcmd, idx := cmd.GetFirstCommand(args)
	switch subcmd {
	case "list":
		return c.list(args[idx+1:], client)
	case "get":
		return c.get(args[idx+1:], client)
	case "add":
		return c.add(args[idx+1:], client)
	case "update":
		return c.update(args[idx+1:], client)
	case "mark-all-read":
		return c.markAllRead(args[idx+1:], client)
	case "clean-history":
		return c.cleanHistory(client)
	default:
		c.Usage()
		return nil
	}
}

func (c *articlesCmd) Usage() {
	fmt.Fprintf(c.flagSet.Output(), "  %s\tManage articles\n", cmdName)
	fmt.Fprintf(c.flagSet.Output(), "    Subcommands: list, get, add, update, mark-all-read, clean-history\n")
}

func (c *articlesCmd) list(args []string, client *graphql.Client) error {
	fs := flag.NewFlagSet("articles list", flag.ExitOnError)
	limit := fs.Int("limit", 10, "Max number of entries")
	afterCursor := fs.Int("after", 0, "Cursor for pagination")
	category := fs.Int("category", 0, "Filter by category ID")
	status := fs.String("status", "", "Filter by status (inbox, read, to_read)")
	starred := fs.Bool("starred", false, "Filter starred articles")
	query := fs.String("query", "", "Full-text search query")
	sortOrder := fs.String("sort-order", "asc", "Sort order (asc, desc)")
	sortBy := fs.String("sort-by", "key", "Sort by (key, stars)")
	if err := fs.Parse(args); err != nil {
		return err
	}

	vars := map[string]any{
		"limit":     *limit,
		"sortOrder": *sortOrder,
		"sortBy":    *sortBy,
	}
	if *afterCursor > 0 {
		vars["afterCursor"] = *afterCursor
	}
	if *category > 0 {
		vars["category"] = *category
	}
	if *status != "" {
		vars["status"] = *status
	}
	if *starred {
		vars["starred"] = true
	}
	if *query != "" {
		vars["query"] = *query
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
				created_at
			}
		}
	}`

	var result struct {
		Articles struct {
			TotalCount int  `json:"totalCount"`
			EndCursor  int  `json:"endCursor"`
			HasNext    bool `json:"hasNext"`
			Entries    []struct {
				ID        int    `json:"id"`
				Title     string `json:"title"`
				Status    string `json:"status"`
				Stars     int    `json:"stars"`
				URL       string `json:"url"`
				CreatedAt string `json:"created_at"`
			} `json:"entries"`
		} `json:"articles"`
	}

	if err := client.Execute(q, vars, &result); err != nil {
		return err
	}

	headers := []string{"ID", "TITLE", "STATUS", "STARS", "CREATED"}
	var rows [][]string
	for _, a := range result.Articles.Entries {
		title := a.Title
		if len(title) > 60 {
			title = title[:57] + "..."
		}
		rows = append(rows, []string{
			strconv.Itoa(a.ID),
			title,
			a.Status,
			strconv.Itoa(a.Stars),
			a.CreatedAt,
		})
	}
	output.PrintTable(headers, rows)
	fmt.Printf("\nTotal: %d, HasNext: %v\n", result.Articles.TotalCount, result.Articles.HasNext)
	return nil
}

func (c *articlesCmd) get(args []string, client *graphql.Client) error {
	fs := flag.NewFlagSet("articles get", flag.ExitOnError)
	html := fs.Bool("html", false, "Include HTML content in response")
	markdown := fs.Bool("markdown", false, "Include Markdown content in response")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() < 1 {
		return fmt.Errorf("usage: articles get [--html] [--markdown] <id>")
	}
	id, err := strconv.Atoi(fs.Arg(0))
	if err != nil {
		return fmt.Errorf("invalid article ID: %s", fs.Arg(0))
	}

	fields := "id title text url image status stars category { id title } published_at created_at updated_at"
	if *html {
		fields += " html"
	}
	if *markdown {
		fields += " markdown"
	}
	q := fmt.Sprintf("query article($id: ID!) { article(id: $id) { %s } }", fields)

	var result struct {
		Article any `json:"article"`
	}
	if err := client.Execute(q, map[string]any{"id": id}, &result); err != nil {
		return err
	}
	return output.PrintJSON(result.Article)
}

func (c *articlesCmd) add(args []string, client *graphql.Client) error {
	fs := flag.NewFlagSet("articles add", flag.ExitOnError)
	category := fs.Int("category", 0, "Category ID")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() < 1 {
		return fmt.Errorf("usage: articles add [-category <id>] <url>")
	}
	articleURL := fs.Arg(0)

	vars := map[string]any{
		"url": articleURL,
	}
	if *category > 0 {
		vars["category"] = *category
	}

	q := `mutation addArticle($url: String!, $category: ID) {
		addArticle(url: $url, category: $category) {
			id
			title
			url
			status
			created_at
		}
	}`

	var result struct {
		AddArticle any `json:"addArticle"`
	}
	if err := client.Execute(q, vars, &result); err != nil {
		return err
	}
	return output.PrintJSON(result.AddArticle)
}

func (c *articlesCmd) update(args []string, client *graphql.Client) error {
	fs := flag.NewFlagSet("articles update", flag.ExitOnError)
	title := fs.String("title", "", "New title")
	text := fs.String("text", "", "New text")
	categoryID := fs.Int("category", 0, "Category ID")
	status := fs.String("status", "", "Status (inbox, read, to_read)")
	stars := fs.Int("stars", -1, "Stars count")
	refresh := fs.Bool("refresh", false, "Refresh article content")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() < 1 {
		return fmt.Errorf("usage: articles update [flags] <id>")
	}
	id, err := strconv.Atoi(fs.Arg(0))
	if err != nil {
		return fmt.Errorf("invalid article ID: %s", fs.Arg(0))
	}

	vars := map[string]any{
		"id": id,
	}
	if *title != "" {
		vars["title"] = *title
	}
	if *text != "" {
		vars["text"] = *text
	}
	if *categoryID > 0 {
		vars["category_id"] = *categoryID
	}
	if *status != "" {
		vars["status"] = *status
	}
	if *stars >= 0 {
		vars["stars"] = *stars
	}
	if *refresh {
		vars["refresh"] = true
	}

	q := `mutation updateArticle($id: ID!, $title: String, $text: String, $category_id: ID, $status: status, $stars: Int, $refresh: Boolean) {
		updateArticle(id: $id, title: $title, text: $text, category_id: $category_id, status: $status, stars: $stars, refresh: $refresh) {
			article {
				id
				title
				status
				stars
				updated_at
			}
		}
	}`

	var result struct {
		UpdateArticle struct {
			Article any `json:"article"`
		} `json:"updateArticle"`
	}
	if err := client.Execute(q, vars, &result); err != nil {
		return err
	}
	return output.PrintJSON(result.UpdateArticle.Article)
}

func (c *articlesCmd) markAllRead(args []string, client *graphql.Client) error {
	fs := flag.NewFlagSet("articles mark-all-read", flag.ExitOnError)
	status := fs.String("status", "inbox", "Status to mark as read (inbox, to_read)")
	category := fs.Int("category", 0, "Category ID")
	if err := fs.Parse(args); err != nil {
		return err
	}

	vars := map[string]any{
		"status": *status,
	}
	if *category > 0 {
		vars["category"] = *category
	}

	q := `mutation markAllArticlesAsRead($status: status, $category: ID) {
		markAllArticlesAsRead(status: $status, category: $category) {
			_inbox
			_to_read
			_starred
			entries { id title }
		}
	}`

	var result struct {
		MarkAllArticlesAsRead any `json:"markAllArticlesAsRead"`
	}
	if err := client.Execute(q, vars, &result); err != nil {
		return err
	}
	return output.PrintJSON(result.MarkAllArticlesAsRead)
}

func (c *articlesCmd) cleanHistory(client *graphql.Client) error {
	q := `mutation {
		cleanHistory {
			_inbox
			_to_read
			_starred
			entries { id title }
		}
	}`

	var result struct {
		CleanHistory any `json:"cleanHistory"`
	}
	if err := client.Execute(q, nil, &result); err != nil {
		return err
	}
	return output.PrintJSON(result.CleanHistory)
}

func newArticlesCmd() cmd.Cmd {
	c := &articlesCmd{
		flagSet: flag.NewFlagSet(cmdName, flag.ExitOnError),
	}
	return c
}

func init() {
	cmd.Add(cmdName, newArticlesCmd)
}
