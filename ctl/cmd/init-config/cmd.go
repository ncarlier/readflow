package initconfig

import (
	"flag"
	"fmt"

	"github.com/ncarlier/readflow/ctl/cmd"
	"github.com/ncarlier/readflow/ctl/config"
	"github.com/ncarlier/readflow/ctl/graphql"
)

const cmdName = "init-config"

type initConfigCmd struct {
	filename string
	flagSet  *flag.FlagSet
}

func (c *initConfigCmd) Exec(args []string, client *graphql.Client, conf *config.Config) error {
	if err := c.flagSet.Parse(args); err != nil {
		return err
	}
	if err := conf.WriteDefaultConfigFile(c.filename); err != nil {
		return err
	}
	fmt.Printf("Configuration file written to %s\n", c.filename)
	return nil
}

func (c *initConfigCmd) Usage() {
	fmt.Fprintf(c.flagSet.Output(), "  %s\tInit configuration file\n", cmdName)
}

func newInitConfigCmd() cmd.Cmd {
	c := &initConfigCmd{
		flagSet: flag.NewFlagSet(cmdName, flag.ExitOnError),
	}
	c.flagSet.StringVar(&c.filename, "f", "config.toml", "Configuration file to create")
	return c
}

func init() {
	cmd.Add(cmdName, newInitConfigCmd)
}
