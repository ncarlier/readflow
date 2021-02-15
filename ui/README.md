# UI

User interface of readflow.

## Configuration

You can configure the webapp build by setting environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `REACT_APP_API_ROOT` | `https://api.readflow.app` | API base URL. |
| `REACT_APP_AUTHORITY` | `https://login.readflow.app/auth/realms/readflow` | OpenID Connect authority provider URL. Set to `mock` for both proxy and mock authentication. |
| `REACT_APP_CLIENT_ID` | `webapp` | OpenID Connect client ID. |
| `REACT_APP_REDIR_URL` | `https://about.readflow.app` | Page to redirect unauthenticated clients to. Set to `/login` for selfhosting.

Example:

```bash
$ export REACT_APP_API_ROOT=http://localhost:8080
```

## Development server

Use `npm start` command to start the development server.

The website will be available here: http://localhost:3000

## Production build

Use `npm run build` to build the UI.

The result is stored into the `./build` directory.
This directory can be served by any web server.

---

