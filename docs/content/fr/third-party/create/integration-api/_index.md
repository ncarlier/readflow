+++
title = "L'API d'intégration"
description = "Utilisation de l'API de création de contenu"
weight = 1
+++

L'API d'intégration est une simple URL accessible en POST avec un payload JSON.

Le payload est un tableau d'articles. Un article peut avoir la structure suivant:

- `url`: L'URL d'origine de l'article.
- `title`: Titre de l'article.
- `html`: Contenu HTML de l'article.
- `text`: Contenu textuel de l'article (le plus souvent un résumé).
- `image`: L'illustration de l'article.
- `published_at`: La date de publication de l'article.
- `category`: La catégorie cible dans readflow.

Toutes les propriétés sont facultatives à l'exception du titre.

*Exemple:*

```json
[
    {
        "title": "",
        "url": "http://example.org/foo.html"
    },
    {
        "title": "Lorem ipsum 1",
        "html": "<p>Lorem ipsum <strong>dolor</strong> sit amet, consectetur adipiscing elit, ...</p>"
    },
    {
        "title": "Lorem ipsum 2",
        "html": "<p>Lorem ipsum <strong>dolor</strong> sit amet, consectetur adipiscing elit, ...</p>",
        "text": "Lorem ipsum dolor sit amet, consectetur adipiscing elit, ...",
        "image": "http://example.org/foo.png",
        "published_at": "2019-04-30T14:00:20.786Z",
        "category": "test",
    }
]
```

L'API utilise l'API key via une authentification basic:

```bash
$ cat payload.json | http \
  -a api:89b5700d-e4da-407e-94a0-7303417189c5 \
  https://api.readflow.app/articles
```

Si l'URL est présente et que l'article est incomplet (pas de contenu HTML, ou de texte, etc...) alors l'article original est téléchargé puis traité avant son intégration.

Le traitement consiste à :

- nettoyer le contenu HTML
- extraire le titre de la page ou des balises [OpenGraph][opengraph]
- extraire le texte du contenu HTML ou des balises OpenGraph
- extraire l'illustration des balises OpenGraph

En retour l'API produit un document JSON avec comme propriétés :

- `articles`: la liste des articles créés (avec leur identifiant interne)
- `errors`: les éventuelles erreurs rencontrées

[opengraph]: http://ogp.me/
