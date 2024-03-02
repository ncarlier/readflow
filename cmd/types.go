package cmd

import "github.com/ncarlier/readflow/internal/config"

type Cmd interface {
	Exec(args []string, conf *config.Config) error
	Usage()
}
