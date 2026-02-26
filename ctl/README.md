# readflow-ctl

POSIX CLI client for interacting with the readflow GraphQL API.

## Overview

`readflow-ctl` is a command-line tool that communicates with a readflow server through its GraphQL API. It supports managing articles, categories, incoming/outgoing webhooks, and other resources from the terminal. It also includes an MCP server mode for AI tool integration.

## Build

```bash
make build-ctl
# outputs: release/readflow-ctl
```

Cross-compilation is supported via standard `GOOS`/`GOARCH` variables:

```bash
GOOS=darwin GOARCH=arm64 make build-ctl
```

## Usage

```
readflow-ctl [flags] <command> [subcommand] [args]
```

### Global flags

| Flag | Env var | Description |
|---|---|---|
| `--endpoint <url>` | `READFLOW_CTL_ENDPOINT` | Readflow API endpoint URL |
| `--username <str>` | `READFLOW_CTL_USERNAME` | Basic auth username |
| `--password <str>` | `READFLOW_CTL_PASSWORD` | Basic auth password |
| `--client-id <str>` | `READFLOW_CTL_CLIENT_ID` | OIDC client ID |

The `--endpoint` flag is required for all commands except `version`.

## Authentication

The auth method is auto-detected by querying `GET <endpoint>/info`. If the server advertises an OIDC authority, OIDC Device Code Flow is used; otherwise Basic Auth is used.

### Basic Auth

Credentials are sent as `Authorization: Basic ...` on every request.

```bash
readflow-ctl --endpoint https://myserver.local \
             --username john --password secret \
             articles list
```

### OIDC Device Code Flow (RFC 8628)

Designed for CLI tools that cannot open a browser inline:

1. The CLI requests a device code from the OIDC provider.
2. A verification URL and user code are displayed in the terminal.
3. The user opens the URL in a browser and enters the code.
4. The CLI polls the token endpoint until authorization is granted.
5. The access token and refresh token are cached to disk.

On subsequent runs the cached token is reused. When expired, the refresh token is used transparently. Delete the cache file to force re-authentication.

**Token cache location:** `$XDG_CACHE_HOME/readflow-ctl/tokens.json`  
Default: `~/.cache/readflow-ctl/tokens.json` (created with `0600` permissions)

```bash
export READFLOW_CTL_ENDPOINT=https://api.readflow.app
export READFLOW_CTL_CLIENT_ID=readflow-ctl
readflow-ctl me
```

## Commands

### `version`

Print the binary version. Does not require `--endpoint` or authentication.

```bash
readflow-ctl version
```

### `me`

Display the current authenticated user's profile.

```bash
readflow-ctl me
```

### `articles`

Manage articles.

#### `articles list`

```bash
readflow-ctl articles list [flags]
```

| Flag | Default | Description |
|---|---|---|
| `--limit <int>` | `10` | Max number of results |
| `--after <int>` | `0` | Pagination cursor (article ID) |
| `--category <int>` | — | Filter by category ID |
| `--status <str>` | — | Filter by status: `inbox`, `read`, `to_read` |
| `--starred` | `false` | Filter starred articles only |
| `--query <str>` | — | Full-text search query |
| `--sort-order <str>` | `asc` | Sort order: `asc` or `desc` |
| `--sort-by <str>` | `key` | Sort field: `key` or `stars` |

#### `articles get`

```bash
readflow-ctl articles get [--html] [--markdown] <id>
```

| Flag | Description |
|---|---|
| `--html` | Include HTML content in response |
| `--markdown` | Include Markdown content in response |

#### `articles add`

```bash
readflow-ctl articles add [--category <id>] <url>
```

#### `articles update`

```bash
readflow-ctl articles update [flags] <id>
```

| Flag | Description |
|---|---|
| `--title <str>` | New title |
| `--text <str>` | New text |
| `--category <int>` | Category ID |
| `--status <str>` | New status: `inbox`, `read`, `to_read` |
| `--stars <int>` | Stars count |
| `--refresh` | Refresh content from source URL |

#### `articles mark-all-read`

```bash
readflow-ctl articles mark-all-read [--status <status>] [--category <id>]
```

| Flag | Default | Description |
|---|---|---|
| `--status <str>` | `inbox` | Status to mark read: `inbox` or `to_read` |
| `--category <int>` | — | Limit to specific category |

#### `articles clean-history`

Remove all read articles from history.

```bash
readflow-ctl articles clean-history
```

### `categories`

Manage categories.

| Subcommand | Description |
|---|---|
| `categories list` | List all categories with article counts |
| `categories get <id>` | Get a category by ID |
| `categories create --title <title>` | Create a new category |
| `categories update [--title <title>] <id>` | Update a category |
| `categories delete <id> [id...]` | Delete one or more categories |

### `incoming-webhooks`

Manage incoming webhooks (for pushing articles into readflow from external systems).

| Subcommand | Description |
|---|---|
| `incoming-webhooks list` | List all incoming webhooks |
| `incoming-webhooks get <id>` | Get a webhook by ID |
| `incoming-webhooks create --alias <alias> --script <script>` | Create a webhook |
| `incoming-webhooks update [--alias <alias>] [--script <script>] <id>` | Update a webhook |
| `incoming-webhooks delete <id> [id...]` | Delete one or more webhooks |

### `outgoing-webhooks`

Manage outgoing webhooks (for sending articles to external services).

| Subcommand | Description |
|---|---|
| `outgoing-webhooks list` | List all outgoing webhooks |
| `outgoing-webhooks get <id>` | Get a webhook by ID |
| `outgoing-webhooks create` | Create a webhook (see flags below) |
| `outgoing-webhooks update` | Update a webhook (see flags below) |
| `outgoing-webhooks delete <id> [id...]` | Delete one or more webhooks |
| `outgoing-webhooks send [--alias <alias>] <article-id>` | Send an article to a webhook |

**Create/update flags:**

| Flag | Description |
|---|---|
| `--alias <str>` | Webhook alias |
| `--provider <str>` | Provider: `generic`, `keeper`, `wallabag`, `shaarli`, `pocket`, `readflow`, `s3` |
| `--config <json>` | Provider configuration as JSON |
| `--secrets <json>` | Provider secrets as JSON |
| `--default` | Set as default outgoing webhook |

### `mcp`

Start a [Model Context Protocol](https://modelcontextprotocol.io) server over stdio, exposing readflow data to AI tools and agents (e.g. Claude Desktop).

```bash
readflow-ctl --endpoint https://api.readflow.app mcp
```

This starts a long-running stdio server. Exposed capabilities:

**Tools:** `list_articles`, `get_article`, `list_categories`, `get_category`

**Resources:** `readflow://categories`, `readflow://articles/{id}`, `readflow://categories/{id}`

**Prompts:** `summarize_article`, `suggest_categories`

## Examples

```bash
# List inbox articles
readflow-ctl articles list --status inbox --limit 20

# Add an article to a category
readflow-ctl articles add --category 5 https://example.com/article

# Get article with Markdown content
readflow-ctl articles get --markdown 42

# Refresh article content from its source URL
readflow-ctl articles update --refresh 42

# Mark all inbox articles in category 3 as read
readflow-ctl articles mark-all-read --status inbox --category 3

# Create a category
readflow-ctl categories create --title "Technology"

# Create an outgoing webhook
readflow-ctl outgoing-webhooks create \
  --alias my-hook \
  --provider generic \
  --config '{"endpoint":"https://hook.example.com"}' \
  --secrets '{"token":"abc123"}' \
  --default

# Send article 42 to a named outgoing webhook
readflow-ctl outgoing-webhooks send --alias my-hook 42

# Start MCP server for AI tool integration
readflow-ctl --endpoint https://api.readflow.app mcp
```
