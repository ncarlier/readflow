package mcp

import (
	"encoding/json"
	"flag"
	"fmt"

	"github.com/mark3labs/mcp-go/server"

	"github.com/ncarlier/readflow/ctl/cmd"
	"github.com/ncarlier/readflow/ctl/graphql"
	"github.com/ncarlier/readflow/internal/version"
)

const cmdName = "mcp"

type mcpCmd struct {
	flagSet *flag.FlagSet
}

func (c *mcpCmd) Exec(args []string, client *graphql.Client) error {
	s := server.NewMCPServer(
		"readflow",
		version.Version,
		server.WithToolCapabilities(false),
		server.WithResourceCapabilities(false, false),
		server.WithPromptCapabilities(false),
		server.WithInstructions("Readflow is a read-it-later application. This MCP server provides read-only access to articles and categories."),
	)

	registerTools(s, client)
	registerResources(s, client)
	registerPrompts(s, client)

	return server.ServeStdio(s)
}

func (c *mcpCmd) Usage() {
	fmt.Fprintf(c.flagSet.Output(), "  %s\tStart MCP server (stdio)\n", cmdName)
}

func newMCPCmd() cmd.Cmd {
	return &mcpCmd{
		flagSet: flag.NewFlagSet(cmdName, flag.ExitOnError),
	}
}

func init() {
	cmd.Add(cmdName, newMCPCmd)
}

// queryGraphQL is a helper that executes a GraphQL query and returns JSON bytes.
func queryGraphQL(client *graphql.Client, query string, vars map[string]any, target any) error {
	return client.Execute(query, vars, target)
}

// toJSON marshals a value to indented JSON.
func toJSON(v any) (string, error) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}
