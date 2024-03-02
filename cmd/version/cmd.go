package initconfig

import (
	"flag"
	"fmt"

	"github.com/ncarlier/readflow/cmd"
	"github.com/ncarlier/readflow/internal/config"
	"github.com/ncarlier/readflow/internal/version"
)

const cmdName = "version"

type VersionCmd struct {
	flagSet *flag.FlagSet
}

func (c *VersionCmd) Exec(args []string, conf *config.Config) error {
	// no args
	version.Print()
	return nil
}

func (c *VersionCmd) Usage() {
	fmt.Fprintf(c.flagSet.Output(), "  %s\tDisplay version\n", cmdName)
}

func newVersionCmd() cmd.Cmd {
	c := &VersionCmd{
		flagSet: flag.NewFlagSet(cmdName, flag.ExitOnError),
	}
	return c
}

func init() {
	cmd.Add("version", newVersionCmd)
}
