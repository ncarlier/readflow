package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ncarlier/readflow/ctl/auth"
	"github.com/ncarlier/readflow/ctl/cmd"
	ctlconfig "github.com/ncarlier/readflow/ctl/config"
	"github.com/ncarlier/readflow/ctl/graphql"

	// activate all ctl commands
	_ "github.com/ncarlier/readflow/ctl/cmd/all"
)

var (
	configFlag string
)

func init() {
	defaultValue := ""
	if value, ok := os.LookupEnv("READFLOW_CTL_CONFIG"); ok {
		defaultValue = value
	}
	flag.StringVar(&configFlag, "c", defaultValue, "Configuration file to load [ENV: READFLOW_CTL_CONFIG]")
	flag.Usage = func() {
		out := flag.CommandLine.Output()
		fmt.Fprintf(out, "readflow-ctl is a CLI client for the readflow GraphQL API.\n")
		fmt.Fprintf(out, "\nUsage:\n  readflow-ctl [flags] <command> [subcommand] [args]\n")
		fmt.Fprintf(out, "\nAvailable Commands:\n")
		for _, c := range cmd.Commands {
			c.Usage()
		}
		fmt.Fprintf(out, "\nFlags:\n")
		flag.PrintDefaults()
		fmt.Fprintf(out, "\nUse \"readflow-ctl <command> --help\" for more information about a command.\n\n")
	}
}

func main() {
	flag.Parse()

	conf := ctlconfig.NewConfig()
	if configFlag != "" {
		if err := conf.LoadFile(configFlag); err != nil {
			fmt.Fprintf(os.Stderr, "error: unable to load configuration file %q: %v\n", configFlag, err)
			os.Exit(1)
		}
	}

	args := flag.Args()
	commandLabel, idx := cmd.GetFirstCommand(args)

	command, ok := cmd.Commands[commandLabel]
	if !ok {
		if commandLabel != "" {
			fmt.Fprintf(os.Stderr, "undefined command: %s\n", commandLabel)
		}
		flag.Usage()
		os.Exit(0)
	}

	// commands that do not require authentication
	if commandLabel == "version" || commandLabel == "init-config" {
		if err := command.Exec(args[idx+1:], nil, conf); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// validate required configuration
	if conf.API.Endpoint == "" {
		fmt.Fprintf(os.Stderr, "error: API endpoint is not configured. Set [api] endpoint in config file or READFLOW_CTL_API_ENDPOINT env var.\n")
		os.Exit(1)
	}
	if conf.AuthN.Method == "" {
		fmt.Fprintf(os.Stderr, "error: authentication method is not configured. Set [authn] method in config file.\n")
		os.Exit(1)
	}

	// create authenticated HTTP client
	httpClient, err := auth.NewHTTPClient(conf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	// create GraphQL client
	client := graphql.NewClient(conf.API.Endpoint, httpClient)

	if err := command.Exec(args[idx+1:], client, conf); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
