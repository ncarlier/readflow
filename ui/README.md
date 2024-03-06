# UI

User interface of readflow.

## Build configuration

You can configure the UI building process by setting environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `REACT_APP_API_ROOT` | `/` | API base URL to use by default if runtime configuration is not set. |
| `REACT_APP_AUTHORITY` | `none` | OpenID Connect authorithy to use by default if runtime configuration is not set. OpenID Connect authority provider URL or `none` if the authentication is delegated to another method (ex: Basic Auth). |
| `REACT_APP_CLIENT_ID` | '' | OpenID Connect client ID. |
| `REACT_APP_PORTAL_URL` | '' | Redirect page for new visitors. |

Example:

```bash
$ export REACT_APP_API_ROOT=http://localhost:8080
$ export REACT_APP_AUTHORITY=none
```

## Runtime configuration

You can customize the UI at runtime by editing the configuration file: [config.js](./public/config.js)

## Dependencies

Use `npm install --legacy-peer-deps` to install dependencies.

## Development server

Use `npm start` command to start the development server.

The website will be available here: http://localhost:3000

## Production build

Use `npm run build` to build the UI.

The result is stored into the `./build` directory.
This directory can be served by any web server.

---
