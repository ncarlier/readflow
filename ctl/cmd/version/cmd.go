package version

import (
	"flag"
	"fmt"

	"github.com/ncarlier/readflow/ctl/cmd"
	"github.com/ncarlier/readflow/ctl/graphql"
	"github.com/ncarlier/readflow/internal/version"
)

const cmdName = "version"

type versionCmd struct {
	flagSet *flag.FlagSet
}

func (c *versionCmd) Exec(args []string, client *graphql.Client) error {
	version.Print()
	return nil
}

func (c *versionCmd) Usage() {
	fmt.Fprintf(c.flagSet.Output(), "  %s\tDisplay version\n", cmdName)
}

func newVersionCmd() cmd.Cmd {
	c := &versionCmd{
		flagSet: flag.NewFlagSet(cmdName, flag.ExitOnError),
	}
	return c
}

func init() {
	cmd.Add(cmdName, newVersionCmd)
}
