package initconfig

import (
	"flag"
	"fmt"

	"github.com/ncarlier/readflow/cmd"
	"github.com/ncarlier/readflow/internal/config"
)

const cmdName = "init-config"

type InitConfigCmd struct {
	filename string
	flagSet  *flag.FlagSet
}

func (c *InitConfigCmd) Exec(args []string, conf *config.Config) error {
	if err := c.flagSet.Parse(args); err != nil {
		return err
	}
	return conf.WriteDefaultConfigFile(c.filename)
}

func (c *InitConfigCmd) Usage() {
	fmt.Fprintf(c.flagSet.Output(), "  %s\tInit configuration file\n", cmdName)
}

func newInitConfigCmd() cmd.Cmd {
	c := &InitConfigCmd{
		flagSet: flag.NewFlagSet(cmdName, flag.ExitOnError),
	}
	c.flagSet.StringVar(&c.filename, "f", "config.toml", "Configuration file to create")
	return c
}

func init() {
	cmd.Add(cmdName, newInitConfigCmd)
}
