+++
title = "Webhook"
description = "Archiver des articles vers un webhook"
weight = 1
+++

![](images/webhook.png)

Un webhook est un simple point d'accès HTTP accessible en POST.
L'article est envoyé en JSON selon le format suivant:

```json
[
    {
        "title": "Le titre",
        "text": "Le text de l'article (le résumé)",
        "html": "Le contenu HTML de l'article",
        "url": "L'URL d'origine de larticle",
        "image": "L'URL de l'illustration de l'image",
        "published_at": "La date de publication de l'article"
    }
]
```

Pour ajouter un webhook, vous devez [configurer votre service d'archivage](https://readflow.app/settings/archive-services):

![](images/archive-services.png)

Cliquez sur le bouton `Add archive service`. La page d'ajout de service s'ouvre:

![](images/new-archive-service.png)

1. Saisissez un alias
1. Choisissez `webhook` comme fournisseur de service
1. Configurez l'URL du webhook
1. Cochez la case si vous souhaitez en faire votre service par défaut

Une fois configuré, vous verrez une nouvelle entrée dans le menu contextuel des articles:

![](images/save-to-webhook.png)

Vous pouvez désormais envoyer un article vers un point d'accès HTTP.

## Partager un article avec un autre readflow

Il est à noter que le format JSON utilisé est compatible avec l'API d'intégration de readflow.
Il est donc possible de configurer un webhook pour envoyer l'article vers un autre compte ou une autre instance readflow.

Pour ce faire vous devez configurer l'URL comme ceci: `https://api:<API_KEY>@api.readflow.app/articles`

En remplacent `<API_KEY>` par une API key du compte cible.
