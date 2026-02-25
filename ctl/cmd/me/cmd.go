package me

import (
	"flag"
	"fmt"

	"github.com/ncarlier/readflow/ctl/cmd"
	"github.com/ncarlier/readflow/ctl/cmd/output"
	"github.com/ncarlier/readflow/ctl/graphql"
)

const cmdName = "me"

type meCmd struct {
	flagSet *flag.FlagSet
}

func (c *meCmd) Exec(args []string, client *graphql.Client) error {
	q := `query {
		me {
			username
			hashid
			hash
			plan
			customer_id
			last_login_at
			created_at
		}
	}`

	var result struct {
		Me any `json:"me"`
	}
	if err := client.Execute(q, nil, &result); err != nil {
		return err
	}
	return output.PrintJSON(result.Me)
}

func (c *meCmd) Usage() {
	fmt.Fprintf(c.flagSet.Output(), "  %s\t\tDisplay current user info\n", cmdName)
}

func newMeCmd() cmd.Cmd {
	c := &meCmd{
		flagSet: flag.NewFlagSet(cmdName, flag.ExitOnError),
	}
	return c
}

func init() {
	cmd.Add(cmdName, newMeCmd)
}
