# readflow

[![Build Status](https://travis-ci.org/ncarlier/readflow.svg?branch=master)](https://travis-ci.org/ncarlier/readflow)
[![Image size](https://images.microbadger.com/badges/image/ncarlier/readflow.svg)](https://microbadger.com/images/ncarlier/readflow)
[![Docker pulls](https://img.shields.io/docker/pulls/ncarlier/readflow.svg)](https://hub.docker.com/r/ncarlier/readflow/)

Read the Internet with complete peace of mind and freedom.

![Logo](readflow.svg)

## Features

- Read articles from anywhere in one place.
- Save articles for offline reading.
- Save articles to external services ([Keeper][keeper], [Wallabag][wallabag], ...).
- Create categories and classify new articles automatically thanks to a customizable rule engine.
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
| `APP_LISTEN_ADDR` | `:8080` | HTTP server address |
| `APP_DB` | `postgres://postgres:testpwd@localhost/reader_test` | Database connection string |
| `APP_AUTHN` | `https://login.nunux.org/auth/realms/readflow` | Authentication method ("mock", "proxy" or OIDC if URL) |
| `APP_BROKER` | none | External event broker URI for outgoing events |
| `APP_SENTRY_DSN` | none | Sentry DSN URL for error reporting |
| `APP_LOG_LEVEL` | `info` | Logging level (`debug`, `info`, `warn` or `error`) |
| `APP_LOG_PRETTY` | `false` | Plain text log output format if true (JSON otherwise) |
| `APP_LOG_OUTPUT` | `stdout` | Log output target (`stdout` or `file://sample.log`) |

You can also override these settings using program parameters.
Type `readflow --help` to see options.

## UI

You can access Web UI on http://localhost:8080/ui

![Screenshot](screenshot.png)

## New articles

The only way to create a new article is to use the integration API.

This API requires an API key.
Use the user interface to obtain an API key.
Then use the API key with your requests.

```bash
$ cat payload.json | http \
  -a api:89b5700d-e4da-407e-94a0-7303417189c5 \
  :8080/articles
```

The JSON payload is an array of article to create.
It must comply this structure:

```js
[
  {
    "title": "foo",          // Article title
    "html": "<p>foo</>",     // Article HTML content
    "text": "foo",           // (optional) Article excerpt
    "url": "http://foo.com", // (optional) Article URL
    "image": "http://...",   // (optional) Article illustration
    "tags": "test,foo",      // (optional) Article tags
    "category": "Test",      // (optional) Target category title
    "published_at": "2019-04-07T13:04:44.247Z", // (optional) Article publication date
  },
  {
    "title": "bar",
    "html": "<p>bar</p>"
  }
]
```

If the article `URL` is provided and `image` or `text` is missing, readflow will try to retrieve this information.

`Tags` can be used by the rule engine to put an article into the relevant category.

If the `category` is set then the rule engine is bypassed.

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
