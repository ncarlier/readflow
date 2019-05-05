# readflow

[![Build Status](https://travis-ci.org/ncarlier/readflow.svg?branch=master)](https://travis-ci.org/ncarlier/readflow)
[![Image size](https://images.microbadger.com/badges/image/ncarlier/readflow.svg)](https://microbadger.com/images/ncarlier/readflow)
[![Docker pulls](https://img.shields.io/docker/pulls/ncarlier/readflow.svg)](https://hub.docker.com/r/ncarlier/readflow/)

Read your Internet article flow in one place with complete peace of mind and freedom.

![Logo](readflow.svg)

## Features

- Read articles from anywhere in one place.
- Save articles for offline reading.
- Save articles to external services ([Keeper][keeper], [Wallabag][wallabag], Webhooks...).
- Create categories and classify new articles automatically thanks to a customizable rule engine.
- Receive notifications when new articles are to be read.
- Good user experience on mobile devices.
- No tracking and no publicity.

## Installation

Run the following command:

```bash
$ go get -v github.com/ncarlier/readflow
```

**Or** download the binary regarding your architecture:

```bash
$ sudo curl -s https://raw.githubusercontent.com/ncarlier/readflow/master/install.sh | bash
```

**Or** use Docker:

```bash
$ docker run -d --name=readflow ncarlier/readflow
```

## Configuration

You can configure the server by setting environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `APP_LISTEN` | `:8080` | Service listen address |
| `APP_LISTEN_METRICS` | none | Metrics listen address |
| `APP_DB` | `postgres://postgres:testpwd@localhost/reader_test` | Database connection string |
| `APP_AUTHN` | `https://login.nunux.org/auth/realms/readflow` | Authentication method ("mock", "proxy" or OIDC if URL) |
| `APP_BROKER` | none | External event broker URI for outgoing events |
| `APP_PUBLIC_URL` | `https://api.readflow.app` | Public URL |
| `APP_SENTRY_DSN` | none | Sentry DSN URL for error reporting |
| `APP_LOG_LEVEL` | `info` | Logging level (`debug`, `info`, `warn` or `error`) |
| `APP_LOG_PRETTY` | `false` | Plain text log output format if true (JSON otherwise) |
| `APP_LOG_OUTPUT` | `stdout` | Log output target (`stdout` or `file://sample.log`) |

You can also override these settings using program parameters.
Type `readflow --help` to see options.

## UI

You can access Web UI on http://localhost:8080/ui

![Screenshot](screenshot.png)

## Documentation

The documentation can be found here: https://about.readlow.app/docs

## GraphQL API

You can explore the server API using GraphiQL endpoint: http://localhost:8080/graphiql

## For development

To be able to build the project you will need to:

- Install `makefiles` external helpers:
  ```bash
  $ git submodule init
  $ git submodule update
  ```

Then you can build the project using make:

```bash
$ make
```

Type `make help` to see other possibilities.

---

[keeper]: https://keeper.nunux.org
[wallabag]: https://www.wallabag.org
