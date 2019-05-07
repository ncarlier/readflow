+++
title = "Flux RSS avec feedpushr"
description = "Utilisation de feedpushr pour intégrer des flux RSS"
weight = 2
+++

![](images/feedpushr.png)

[Feedpushr](https://github.com/ncarlier/feedpushr) est un puissant aggrégateur Open Source de flux RSS capable de traiter et envoyer les articles vers des serices tiers.
Et parmi ces services, nous avons readflow.

## En utilisant l'exécutable

Pour utiliser readflow avec feedpushr, vous devez utiliser et configurer son plugin:

```bash
$ # Configurer l'URL de readflow
$ export APP_READFLOW_URL=https://api.readflow.app
$ # Configurer la clé d'API
$ export APP_READFLOW_API_KEY=d70*************456
$ # Lancer feedpushr
$ feedpushr --log-pretty --plugin ./feedpushr-readflow.so --output readflow://
```

## En utilisant Docker

Créer le fichier de configuration suivant (`config.env`):

```bash
APP_PLUGINS=feedpushr-readflow.so                                           
APP_READFLOW_URL=https://api.readflow.app
APP_READFLOW_API_KEY=<YOUR API KEY>
APP_OUTPUTS=readflow://
APP_FILTERS=fetch://#fetch,minify://#minify
APP_LOG_PRETTY=true
```

Lancer Docker:

```bash
$ docker run -d --name feedpushr -p 8080:8080 --env-file=conf.env ncarlier/feedpushr-contrib
```

## En utilisant Docker Compose

Créer le fichier `docker-compose.yml` :

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

Personaliser la configuration avec un fichier `.env` :

```bash
PORT=8080
READFLOW_API_KEY=<YOUR API KEY>
```

Lancer Compose:

```bash
$ docker-compose up -d
```

## The UI

Vous devriez voir ceci sur [l'IHM de feedpushr](http://localhost:8080/ui) :

![](images/feedpushr-ui.png)

Vous pouvez ensuite importer vos abonnements OPML dans feedpushr et voir vos articles arriver dans readflow.
