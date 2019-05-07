+++
title = "RSS feed with feedpushr"
description = "Use feedpushr to integrate RSS feeds"
weight = 2
+++

![](images/feedpushr.png)

[Feedpushr](https://github.com/ncarlier/feedpushr) is a powerful open source RSS feed aggregator able to process and send articles to third party services.
And among these services, we have readflow.

## Using the binary

To use readflow with feedpushr, you must use and configure its plugin:

```bash
$ # Set readflow URL
$ export APP_READFLOW_URL=https://api.readflow.app
$ # Set API key
$ export APP_READFLOW_API_KEY=d70*************456
$ # Start feedpushr
$ feedpushr --log-pretty --plugin ./feedpushr-readflow.so --output readflow://
```

## Using Docker

Create the following `config.env` file:

```bash
APP_PLUGINS=feedpushr-readflow.so                                           
APP_READFLOW_URL=https://api.readflow.app
APP_READFLOW_API_KEY=<YOUR API KEY>
APP_OUTPUTS=readflow://
APP_FILTERS=fetch://#fetch,minify://#minify
APP_LOG_PRETTY=true
```

Start Docker:

```bash
$ docker run -d --name feedpushr -p 8080:8080 --env-file=config.env ncarlier/feedpushr-contrib
```

## Using Docker Compose

Create the following `docker-compose.yml` file:

```yaml
version: "3"
services:
  feedpushr:
    image: "ncarlier/feedpushr-contrib:latest"
    restart: always
    ports:
      - "${PORT:-8080}:8080"
    environment:
      - APP_DB=boltdb://var/opt/feedpushr/data.db
      - APP_PLUGINS=feedpushr-readflow.so
      - APP_OUTPUTS=readflow://
      - "APP_FILTERS=fetch://#fetch,minify://#minify"
      - APP_LOG_PRETTY=true
      - APP_DELAY=5m
      - APP_READFLOW_URL=${READFLOW_URL:-https://api.readflow.app}
      - APP_READFLOW_API_KEY=${READFLOW_API_KEY}
    volumes:
      - feedpushr-data:/var/opt/feedpushr

volumes:
  feedpushr-data:
```

Customize the configuration with an `.env` file:

```bash
PORT=8080
READFLOW_API_KEY=<YOUR API KEY>
```

Start Compose:

```bash
$ docker-compose up -d
```

## The UI

You should see this on the [feedpushr UI](http://localhost:8080/ui):

![](images/feedpushr-ui.png)

You can then import your OPML subscriptions into feedpushr and see your articles in readflow.
