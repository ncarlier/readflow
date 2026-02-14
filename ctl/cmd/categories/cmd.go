package categories

import (
	"flag"
	"fmt"
	"strconv"
	"strings"

	"github.com/ncarlier/readflow/ctl/cmd"
	"github.com/ncarlier/readflow/ctl/cmd/output"
	"github.com/ncarlier/readflow/ctl/config"
	"github.com/ncarlier/readflow/ctl/graphql"
)

const cmdName = "categories"

type categoriesCmd struct {
	flagSet *flag.FlagSet
}

func (c *categoriesCmd) Exec(args []string, client *graphql.Client, conf *config.Config) error {
	subcmd, idx := cmd.GetFirstCommand(args)
	switch subcmd {
	case "list":
		return c.list(client)
	case "get":
		return c.get(args[idx+1:], client)
	case "create":
		return c.create(args[idx+1:], client)
	case "update":
		return c.update(args[idx+1:], client)
	case "delete":
		return c.delete(args[idx+1:], client)
	default:
		c.Usage()
		return nil
	}
}

func (c *categoriesCmd) Usage() {
	fmt.Fprintf(c.flagSet.Output(), "  %s\tManage categories\n", cmdName)
	fmt.Fprintf(c.flagSet.Output(), "    Subcommands: list, get, create, update, delete\n")
}

func (c *categoriesCmd) list(client *graphql.Client) error {
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
		Categories struct {
			Inbox   int `json:"_inbox"`
			ToRead  int `json:"_to_read"`
			Starred int `json:"_starred"`
			Entries []struct {
				ID        int    `json:"id"`
				Title     string `json:"title"`
				Inbox     int    `json:"inbox"`
				CreatedAt string `json:"created_at"`
				UpdatedAt string `json:"updated_at"`
			} `json:"entries"`
		} `json:"categories"`
	}

	if err := client.Execute(q, nil, &result); err != nil {
		return err
	}

	headers := []string{"ID", "TITLE", "INBOX", "CREATED"}
	var rows [][]string
	for _, cat := range result.Categories.Entries {
		rows = append(rows, []string{
			strconv.Itoa(cat.ID),
			cat.Title,
			strconv.Itoa(cat.Inbox),
			cat.CreatedAt,
		})
	}
	output.PrintTable(headers, rows)
	fmt.Printf("\nInbox: %d, ToRead: %d, Starred: %d\n", result.Categories.Inbox, result.Categories.ToRead, result.Categories.Starred)
	return nil
}

func (c *categoriesCmd) get(args []string, client *graphql.Client) error {
	fs := flag.NewFlagSet("categories get", flag.ExitOnError)
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() < 1 {
		return fmt.Errorf("usage: categories get <id>")
	}
	id, err := strconv.Atoi(fs.Arg(0))
	if err != nil {
		return fmt.Errorf("invalid category ID: %s", fs.Arg(0))
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
	if err := client.Execute(q, map[string]any{"id": id}, &result); err != nil {
		return err
	}
	return output.PrintJSON(result.Category)
}

func (c *categoriesCmd) create(args []string, client *graphql.Client) error {
	fs := flag.NewFlagSet("categories create", flag.ExitOnError)
	title := fs.String("title", "", "Category title (required)")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *title == "" {
		return fmt.Errorf("usage: categories create -title <title>")
	}

	q := `mutation createOrUpdateCategory($title: String!) {
		createOrUpdateCategory(title: $title) {
			id
			title
			created_at
		}
	}`

	var result struct {
		CreateOrUpdateCategory any `json:"createOrUpdateCategory"`
	}
	if err := client.Execute(q, map[string]any{"title": *title}, &result); err != nil {
		return err
	}
	return output.PrintJSON(result.CreateOrUpdateCategory)
}

func (c *categoriesCmd) update(args []string, client *graphql.Client) error {
	fs := flag.NewFlagSet("categories update", flag.ExitOnError)
	title := fs.String("title", "", "New title")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() < 1 {
		return fmt.Errorf("usage: categories update -title <title> <id>")
	}
	id, err := strconv.Atoi(fs.Arg(0))
	if err != nil {
		return fmt.Errorf("invalid category ID: %s", fs.Arg(0))
	}

	vars := map[string]any{
		"id": id,
	}
	if *title != "" {
		vars["title"] = *title
	}

	q := `mutation createOrUpdateCategory($id: ID, $title: String) {
		createOrUpdateCategory(id: $id, title: $title) {
			id
			title
			updated_at
		}
	}`

	var result struct {
		CreateOrUpdateCategory any `json:"createOrUpdateCategory"`
	}
	if err := client.Execute(q, vars, &result); err != nil {
		return err
	}
	return output.PrintJSON(result.CreateOrUpdateCategory)
}

func (c *categoriesCmd) delete(args []string, client *graphql.Client) error {
	fs := flag.NewFlagSet("categories delete", flag.ExitOnError)
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() < 1 {
		return fmt.Errorf("usage: categories delete <id> [id...]")
	}

	var ids []int
	for _, arg := range fs.Args() {
		id, err := strconv.Atoi(arg)
		if err != nil {
			return fmt.Errorf("invalid category ID: %s", arg)
		}
		ids = append(ids, id)
	}

	q := `mutation deleteCategories($ids: [ID]!) {
		deleteCategories(ids: $ids)
	}`

	var result struct {
		DeleteCategories int `json:"deleteCategories"`
	}
	if err := client.Execute(q, map[string]any{"ids": ids}, &result); err != nil {
		return err
	}
	fmt.Printf("Deleted %d categor%s\n", result.DeleteCategories, pluralize(result.DeleteCategories, "y", "ies"))
	return nil
}

func pluralize(count int, singular, plural string) string {
	if count == 1 {
		return singular
	}
	return plural
}

// parseIDList parses a comma-separated list of IDs
func parseIDList(s string) ([]int, error) {
	parts := strings.Split(s, ",")
	var ids []int
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		id, err := strconv.Atoi(p)
		if err != nil {
			return nil, fmt.Errorf("invalid ID: %s", p)
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func newCategoriesCmd() cmd.Cmd {
	c := &categoriesCmd{
		flagSet: flag.NewFlagSet(cmdName, flag.ExitOnError),
	}
	return c
}

func init() {
	cmd.Add(cmdName, newCategoriesCmd)
}
