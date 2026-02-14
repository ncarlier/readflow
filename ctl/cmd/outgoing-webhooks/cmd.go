package outgoingwebhooks

import (
	"flag"
	"fmt"
	"strconv"

	"github.com/ncarlier/readflow/ctl/cmd"
	"github.com/ncarlier/readflow/ctl/cmd/output"
	"github.com/ncarlier/readflow/ctl/config"
	"github.com/ncarlier/readflow/ctl/graphql"
)

const cmdName = "outgoing-webhooks"

type outgoingWebhooksCmd struct {
	flagSet *flag.FlagSet
}

func (c *outgoingWebhooksCmd) Exec(args []string, client *graphql.Client, conf *config.Config) error {
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
	case "send":
		return c.send(args[idx+1:], client)
	default:
		c.Usage()
		return nil
	}
}

func (c *outgoingWebhooksCmd) Usage() {
	fmt.Fprintf(c.flagSet.Output(), "  %s\tManage outgoing webhooks\n", cmdName)
	fmt.Fprintf(c.flagSet.Output(), "    Subcommands: list, get, create, update, delete, send\n")
}

func (c *outgoingWebhooksCmd) list(client *graphql.Client) error {
	q := `query {
		outgoingWebhooks {
			id
			alias
			provider
			is_default
			created_at
		}
	}`

	var result struct {
		OutgoingWebhooks []struct {
			ID        int    `json:"id"`
			Alias     string `json:"alias"`
			Provider  string `json:"provider"`
			IsDefault bool   `json:"is_default"`
			CreatedAt string `json:"created_at"`
		} `json:"outgoingWebhooks"`
	}

	if err := client.Execute(q, nil, &result); err != nil {
		return err
	}

	headers := []string{"ID", "ALIAS", "PROVIDER", "DEFAULT", "CREATED"}
	var rows [][]string
	for _, wh := range result.OutgoingWebhooks {
		isDefault := "no"
		if wh.IsDefault {
			isDefault = "yes"
		}
		rows = append(rows, []string{
			strconv.Itoa(wh.ID),
			wh.Alias,
			wh.Provider,
			isDefault,
			wh.CreatedAt,
		})
	}
	output.PrintTable(headers, rows)
	return nil
}

func (c *outgoingWebhooksCmd) get(args []string, client *graphql.Client) error {
	fs := flag.NewFlagSet("outgoing-webhooks get", flag.ExitOnError)
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() < 1 {
		return fmt.Errorf("usage: outgoing-webhooks get <id>")
	}
	id, err := strconv.Atoi(fs.Arg(0))
	if err != nil {
		return fmt.Errorf("invalid webhook ID: %s", fs.Arg(0))
	}

	q := `query outgoingWebhook($id: ID!) {
		outgoingWebhook(id: $id) {
			id
			alias
			provider
			config
			secrets
			is_default
			created_at
			updated_at
		}
	}`

	var result struct {
		OutgoingWebhook any `json:"outgoingWebhook"`
	}
	if err := client.Execute(q, map[string]any{"id": id}, &result); err != nil {
		return err
	}
	return output.PrintJSON(result.OutgoingWebhook)
}

func (c *outgoingWebhooksCmd) create(args []string, client *graphql.Client) error {
	fs := flag.NewFlagSet("outgoing-webhooks create", flag.ExitOnError)
	alias := fs.String("alias", "", "Webhook alias (required)")
	provider := fs.String("provider", "", "Provider: generic, keeper, wallabag, shaarli, pocket, readflow, s3 (required)")
	config := fs.String("config", "", "Provider configuration as JSON (required)")
	secrets := fs.String("secrets", "", "Provider secrets as JSON (required)")
	isDefault := fs.Bool("default", false, "Set as default webhook")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *alias == "" || *provider == "" || *config == "" || *secrets == "" {
		return fmt.Errorf("usage: outgoing-webhooks create -alias <alias> -provider <provider> -config <json> -secrets <json>")
	}

	vars := map[string]any{
		"alias":      *alias,
		"provider":   *provider,
		"config":     *config,
		"secrets":    *secrets,
		"is_default": *isDefault,
	}

	q := `mutation createOrUpdateOutgoingWebhook($alias: String!, $provider: outgoingWebhookProvider!, $config: String!, $secrets: String!, $is_default: Boolean) {
		createOrUpdateOutgoingWebhook(alias: $alias, provider: $provider, config: $config, secrets: $secrets, is_default: $is_default) {
			id
			alias
			provider
			is_default
			created_at
		}
	}`

	var result struct {
		CreateOrUpdateOutgoingWebhook any `json:"createOrUpdateOutgoingWebhook"`
	}
	if err := client.Execute(q, vars, &result); err != nil {
		return err
	}
	return output.PrintJSON(result.CreateOrUpdateOutgoingWebhook)
}

func (c *outgoingWebhooksCmd) update(args []string, client *graphql.Client) error {
	fs := flag.NewFlagSet("outgoing-webhooks update", flag.ExitOnError)
	alias := fs.String("alias", "", "Webhook alias")
	provider := fs.String("provider", "", "Provider")
	config := fs.String("config", "", "Provider configuration as JSON")
	secrets := fs.String("secrets", "", "Provider secrets as JSON")
	isDefault := fs.Bool("default", false, "Set as default webhook")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() < 1 {
		return fmt.Errorf("usage: outgoing-webhooks update [flags] <id>")
	}
	id, err := strconv.Atoi(fs.Arg(0))
	if err != nil {
		return fmt.Errorf("invalid webhook ID: %s", fs.Arg(0))
	}

	vars := map[string]any{
		"id": id,
	}
	if *alias != "" {
		vars["alias"] = *alias
	}
	if *provider != "" {
		vars["provider"] = *provider
	}
	if *config != "" {
		vars["config"] = *config
	}
	if *secrets != "" {
		vars["secrets"] = *secrets
	}
	if *isDefault {
		vars["is_default"] = true
	}

	q := `mutation createOrUpdateOutgoingWebhook($id: ID, $alias: String!, $provider: outgoingWebhookProvider!, $config: String!, $secrets: String!, $is_default: Boolean) {
		createOrUpdateOutgoingWebhook(id: $id, alias: $alias, provider: $provider, config: $config, secrets: $secrets, is_default: $is_default) {
			id
			alias
			provider
			is_default
			updated_at
		}
	}`

	var result struct {
		CreateOrUpdateOutgoingWebhook any `json:"createOrUpdateOutgoingWebhook"`
	}
	if err := client.Execute(q, vars, &result); err != nil {
		return err
	}
	return output.PrintJSON(result.CreateOrUpdateOutgoingWebhook)
}

func (c *outgoingWebhooksCmd) delete(args []string, client *graphql.Client) error {
	fs := flag.NewFlagSet("outgoing-webhooks delete", flag.ExitOnError)
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() < 1 {
		return fmt.Errorf("usage: outgoing-webhooks delete <id> [id...]")
	}

	var ids []int
	for _, arg := range fs.Args() {
		id, err := strconv.Atoi(arg)
		if err != nil {
			return fmt.Errorf("invalid webhook ID: %s", arg)
		}
		ids = append(ids, id)
	}

	q := `mutation deleteOutgoingWebhooks($ids: [ID]!) {
		deleteOutgoingWebhooks(ids: $ids)
	}`

	var result struct {
		DeleteOutgoingWebhooks int `json:"deleteOutgoingWebhooks"`
	}
	if err := client.Execute(q, map[string]any{"ids": ids}, &result); err != nil {
		return err
	}
	fmt.Printf("Deleted %d outgoing webhook(s)\n", result.DeleteOutgoingWebhooks)
	return nil
}

func (c *outgoingWebhooksCmd) send(args []string, client *graphql.Client) error {
	fs := flag.NewFlagSet("outgoing-webhooks send", flag.ExitOnError)
	alias := fs.String("alias", "", "Outgoing webhook alias (uses default if empty)")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() < 1 {
		return fmt.Errorf("usage: outgoing-webhooks send [-alias <alias>] <article-id>")
	}
	articleID, err := strconv.Atoi(fs.Arg(0))
	if err != nil {
		return fmt.Errorf("invalid article ID: %s", fs.Arg(0))
	}

	vars := map[string]any{
		"id": articleID,
	}
	if *alias != "" {
		vars["alias"] = *alias
	}

	q := `mutation sendArticleToOutgoingWebhook($id: ID!, $alias: String) {
		sendArticleToOutgoingWebhook(id: $id, alias: $alias) {
			url
			text
			json
		}
	}`

	var result struct {
		SendArticleToOutgoingWebhook any `json:"sendArticleToOutgoingWebhook"`
	}
	if err := client.Execute(q, vars, &result); err != nil {
		return err
	}
	return output.PrintJSON(result.SendArticleToOutgoingWebhook)
}

func newOutgoingWebhooksCmd() cmd.Cmd {
	c := &outgoingWebhooksCmd{
		flagSet: flag.NewFlagSet(cmdName, flag.ExitOnError),
	}
	return c
}

func init() {
	cmd.Add(cmdName, newOutgoingWebhooksCmd)
}
