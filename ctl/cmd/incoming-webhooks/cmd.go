package incomingwebhooks

import (
	"flag"
	"fmt"
	"strconv"

	"github.com/ncarlier/readflow/ctl/cmd"
	"github.com/ncarlier/readflow/ctl/cmd/output"
	"github.com/ncarlier/readflow/ctl/graphql"
)

const cmdName = "incoming-webhooks"

type incomingWebhooksCmd struct {
	flagSet *flag.FlagSet
}

func (c *incomingWebhooksCmd) Exec(args []string, client *graphql.Client) error {
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

func (c *incomingWebhooksCmd) Usage() {
	fmt.Fprintf(c.flagSet.Output(), "  %s\tManage incoming webhooks\n", cmdName)
	fmt.Fprintf(c.flagSet.Output(), "    Subcommands: list, get, create, update, delete\n")
}

func (c *incomingWebhooksCmd) list(client *graphql.Client) error {
	q := `query {
		incomingWebhooks {
			id
			alias
			token
			script
			last_usage_at
			created_at
		}
	}`

	var result struct {
		IncomingWebhooks []struct {
			ID          int    `json:"id"`
			Alias       string `json:"alias"`
			Token       string `json:"token"`
			Script      string `json:"script"`
			LastUsageAt string `json:"last_usage_at"`
			CreatedAt   string `json:"created_at"`
		} `json:"incomingWebhooks"`
	}

	if err := client.Execute(q, nil, &result); err != nil {
		return err
	}

	headers := []string{"ID", "ALIAS", "TOKEN", "LAST USED", "CREATED"}
	var rows [][]string
	for _, wh := range result.IncomingWebhooks {
		rows = append(rows, []string{
			strconv.Itoa(wh.ID),
			wh.Alias,
			wh.Token,
			wh.LastUsageAt,
			wh.CreatedAt,
		})
	}
	output.PrintTable(headers, rows)
	return nil
}

func (c *incomingWebhooksCmd) get(args []string, client *graphql.Client) error {
	fs := flag.NewFlagSet("incoming-webhooks get", flag.ExitOnError)
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() < 1 {
		return fmt.Errorf("usage: incoming-webhooks get <id>")
	}
	id, err := strconv.Atoi(fs.Arg(0))
	if err != nil {
		return fmt.Errorf("invalid webhook ID: %s", fs.Arg(0))
	}

	q := `query incomingWebhook($id: ID!) {
		incomingWebhook(id: $id) {
			id
			alias
			token
			script
			last_usage_at
			created_at
			updated_at
		}
	}`

	var result struct {
		IncomingWebhook any `json:"incomingWebhook"`
	}
	if err := client.Execute(q, map[string]any{"id": id}, &result); err != nil {
		return err
	}
	return output.PrintJSON(result.IncomingWebhook)
}

func (c *incomingWebhooksCmd) create(args []string, client *graphql.Client) error {
	fs := flag.NewFlagSet("incoming-webhooks create", flag.ExitOnError)
	alias := fs.String("alias", "", "Webhook alias (required)")
	script := fs.String("script", "", "Webhook script (required)")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *alias == "" || *script == "" {
		return fmt.Errorf("usage: incoming-webhooks create -alias <alias> -script <script>")
	}

	q := `mutation createOrUpdateIncomingWebhook($alias: String!, $script: String!) {
		createOrUpdateIncomingWebhook(alias: $alias, script: $script) {
			id
			alias
			token
			script
			created_at
		}
	}`

	var result struct {
		CreateOrUpdateIncomingWebhook any `json:"createOrUpdateIncomingWebhook"`
	}
	if err := client.Execute(q, map[string]any{"alias": *alias, "script": *script}, &result); err != nil {
		return err
	}
	return output.PrintJSON(result.CreateOrUpdateIncomingWebhook)
}

func (c *incomingWebhooksCmd) update(args []string, client *graphql.Client) error {
	fs := flag.NewFlagSet("incoming-webhooks update", flag.ExitOnError)
	alias := fs.String("alias", "", "Webhook alias")
	script := fs.String("script", "", "Webhook script")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() < 1 {
		return fmt.Errorf("usage: incoming-webhooks update -alias <alias> -script <script> <id>")
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
	if *script != "" {
		vars["script"] = *script
	}

	q := `mutation createOrUpdateIncomingWebhook($id: ID, $alias: String!, $script: String!) {
		createOrUpdateIncomingWebhook(id: $id, alias: $alias, script: $script) {
			id
			alias
			token
			script
			updated_at
		}
	}`

	var result struct {
		CreateOrUpdateIncomingWebhook any `json:"createOrUpdateIncomingWebhook"`
	}
	if err := client.Execute(q, vars, &result); err != nil {
		return err
	}
	return output.PrintJSON(result.CreateOrUpdateIncomingWebhook)
}

func (c *incomingWebhooksCmd) delete(args []string, client *graphql.Client) error {
	fs := flag.NewFlagSet("incoming-webhooks delete", flag.ExitOnError)
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() < 1 {
		return fmt.Errorf("usage: incoming-webhooks delete <id> [id...]")
	}

	var ids []int
	for _, arg := range fs.Args() {
		id, err := strconv.Atoi(arg)
		if err != nil {
			return fmt.Errorf("invalid webhook ID: %s", arg)
		}
		ids = append(ids, id)
	}

	q := `mutation deleteIncomingWebhooks($ids: [ID]!) {
		deleteIncomingWebhooks(ids: $ids)
	}`

	var result struct {
		DeleteIncomingWebhooks int `json:"deleteIncomingWebhooks"`
	}
	if err := client.Execute(q, map[string]any{"ids": ids}, &result); err != nil {
		return err
	}
	fmt.Printf("Deleted %d incoming webhook(s)\n", result.DeleteIncomingWebhooks)
	return nil
}

func newIncomingWebhooksCmd() cmd.Cmd {
	c := &incomingWebhooksCmd{
		flagSet: flag.NewFlagSet(cmdName, flag.ExitOnError),
	}
	return c
}

func init() {
	cmd.Add(cmdName, newIncomingWebhooksCmd)
}
