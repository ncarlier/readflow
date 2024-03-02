package main

//go:generate go run generate.go
//go:generate gofmt -s -w autogen/db/postgres/db_sql_migration.go

import (
	"flag"
	"fmt"
	"os"

	"github.com/ncarlier/readflow/cmd"
	"github.com/ncarlier/readflow/internal/config"
	"github.com/ncarlier/readflow/pkg/logger"

	_ "github.com/ncarlier/readflow/cmd/all"
)

func main() {
	// parse command line
	flag.Parse()

	// load configuration
	conf := config.NewConfig()
	if cmd.ConfigFlag != "" {
		if err := conf.LoadFile(cmd.ConfigFlag); err != nil {
			logger.Fatal().Err(err).Str("filename", cmd.ConfigFlag).Msg("unable to load configuration file")
		}
	}

	// export configurations vars
	config.ExportVars(conf)

	// configure the logger
	logger.Configure(conf.Log.Level, conf.Log.Format, conf.Integration.Sentry.DSN)

	args := flag.Args()
	commandLabel, idx := cmd.GetFirstCommand(args)

	if command, ok := cmd.Commands[commandLabel]; ok {
		if err := command.Exec(args[idx+1:], conf); err != nil {
			logger.Fatal().Err(err).Str("command", commandLabel).Msg("error during command execution")
		}
	} else {
		if commandLabel != "" {
			fmt.Fprintf(os.Stderr, "undefined command: %s\n", commandLabel)
		}
		flag.Usage()
		os.Exit(0)
	}
}
