# UI

User interface of readflow.

## Configuration

By default the API endpoint is `/api`.

You can change this by setting the `REACT_APP_API_BASE_URL` environment variable.

Example:

```bash
$ export REACT_APP_API_BASE_URL=http://localhost:8080
```

## Development server

Use `npm start` command to start the development server.

The website will be available here: http://localhost:3000

## Production build

Use `npm run build` to build the UI.

The result is stored into the `./build` directory.
This directory can be served by any web server.

---

