package cmd

import (
	"github.com/ncarlier/readflow/ctl/graphql"
)

// Cmd is the interface for a CLI command
type Cmd interface {
	Exec(args []string, client *graphql.Client) error
	Usage()
}
