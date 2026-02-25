package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ncarlier/readflow/ctl/auth"
	"github.com/ncarlier/readflow/ctl/cmd"
	"github.com/ncarlier/readflow/ctl/graphql"

	// activate all ctl commands
	_ "github.com/ncarlier/readflow/ctl/cmd/all"
)

var (
	endpoint string
	username string
	password string
	clientID string
)

func envDefault(envKey, fallback string) string {
	if v, ok := os.LookupEnv(envKey); ok {
		return v
	}
	return fallback
}

func init() {
	flag.StringVar(&endpoint, "endpoint", envDefault("READFLOW_CTL_ENDPOINT", ""), "Readflow API endpoint URL [ENV: READFLOW_CTL_ENDPOINT]")
	flag.StringVar(&username, "username", envDefault("READFLOW_CTL_USERNAME", ""), "Basic auth username [ENV: READFLOW_CTL_USERNAME]")
	flag.StringVar(&password, "password", envDefault("READFLOW_CTL_PASSWORD", ""), "Basic auth password [ENV: READFLOW_CTL_PASSWORD]")
	flag.StringVar(&clientID, "client-id", envDefault("READFLOW_CTL_CLIENT_ID", ""), "OIDC client ID [ENV: READFLOW_CTL_CLIENT_ID]")

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
	if commandLabel == "version" {
		if err := command.Exec(args[idx+1:], nil); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// validate required endpoint
	if endpoint == "" {
		fmt.Fprintf(os.Stderr, "error: --endpoint is required [ENV: READFLOW_CTL_ENDPOINT]\n")
		os.Exit(1)
	}

	// discover auth method from /info
	method, issuer, err := auth.DiscoverAuthMethod(endpoint)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: unable to discover auth method: %v\n", err)
		os.Exit(1)
	}

	// validate method-specific required parameters
	authConf := &auth.Config{
		Method:   method,
		Issuer:   issuer,
		Username: username,
		Password: password,
		ClientID: clientID,
	}
	switch method {
	case "basic":
		if authConf.Username == "" {
			fmt.Fprintf(os.Stderr, "error: --username is required for basic auth [ENV: READFLOW_CTL_USERNAME]\n")
			os.Exit(1)
		}
		if authConf.Password == "" {
			fmt.Fprintf(os.Stderr, "error: --password is required for basic auth [ENV: READFLOW_CTL_PASSWORD]\n")
			os.Exit(1)
		}
	case "oidc":
		if authConf.ClientID == "" {
			fmt.Fprintf(os.Stderr, "error: --client-id is required for OIDC auth [ENV: READFLOW_CTL_CLIENT_ID]\n")
			os.Exit(1)
		}
	}

	// create authenticated HTTP client
	httpClient, err := auth.NewHTTPClient(authConf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	// create GraphQL client
	client := graphql.NewClient(endpoint, httpClient)

	if err := command.Exec(args[idx+1:], client); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
