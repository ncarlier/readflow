# ctl

POSIX CLI client for interacting with the readflow GraphQL API.

## Overview

`readflow-ctl` is a command-line tool that communicates with a readflow server instance through its GraphQL API. It supports managing articles, categories, webhooks, and other resources from the terminal.

## Authentication

The CLI supports two authentication methods:

### Basic Auth

Authenticate using a username and password. Credentials are sent with each request via HTTP Basic Authentication.

```toml
[authn]
method = "basic"
[authn.basic]
username = "john"
password = "secret"
```

### OIDC Device Code Flow

Authenticate using [OAuth 2.0 Device Authorization Grant](https://datatracker.ietf.org/doc/html/rfc8628) (RFC 8628). This flow is designed for CLI tools that cannot open a browser inline:

1. The CLI requests a device code from the OIDC provider.
2. A verification URL and user code are displayed in the terminal.
3. The user opens the URL in a browser and enters the code.
4. The CLI polls the token endpoint until authorization is granted.
5. The access token (and refresh token) are stored in a local cache file.

On subsequent invocations, the cached access token is used directly. When the access token expires, the refresh token is used to obtain a new one transparently.

```toml
[authn]
method = "oidc"
[authn.oidc]
issuer = "https://accounts.readflow.app"
client_id = "readflow-ctl"
```

#### Token cache

Tokens are stored in `$XDG_CACHE_HOME/readflow-ctl/tokens.json` (defaults to `~/.cache/readflow-ctl/tokens.json`). The file is created with `0600` permissions. Delete this file to force re-authentication.

## Configuration

Similar to the server, configuration is loaded from a TOML file specified with the `-c` flag or the `READFLOW_CTL_CONFIG` environment variable.

```bash
readflow-ctl -c config.toml articles list
```

Environment variables can be used anywhere in the config file using `${VAR}` syntax.

### Configuration file

```toml
# readflow-ctl configuration file

[api]
## Readflow server URL
# Example: "https://api.readflow.app"
endpoint = "${READFLOW_CTL_API_ENDPOINT}"

[authn]
## Authentication method
# Values: "basic" or "oidc"
method = "${READFLOW_CTL_AUTHN_METHOD}"

[authn.basic]
## Username for Basic Authentication
username = "${READFLOW_CTL_AUTHN_BASIC_USERNAME}"
## Password for Basic Authentication
password = "${READFLOW_CTL_AUTHN_BASIC_PASSWORD}"

[authn.oidc]
## OpenID Connect issuer URL
issuer = "${READFLOW_CTL_AUTHN_OIDC_ISSUER}"
## OpenID Connect client ID
client_id = "${READFLOW_CTL_AUTHN_OIDC_CLIENT_ID}"
```

---
