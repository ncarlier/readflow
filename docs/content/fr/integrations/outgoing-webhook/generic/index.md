+++
title = "Webhook générique"
description = "Envoyer des articles vers un webhook générique"
weight = 1
+++

![](images/webhook.png)

Un webhook générique est un simple point d'accès HTTP accessible en POST.

## Format

Par défaut, l'article est envoyé en JSON selon le format suivant:

```json
{
  "title": "Le titre",
  "text": "Le text de l'article (le résumé)",
  "html": "Le contenu HTML de l'article",
  "url": "L'URL d'origine de larticle",
  "image": "L'URL de l'illustration de l'image",
  "published_at": "La date de publication de l'article"
}
```

Il est possible de personnaliser le contenu envoyé.

Vous pouvez spécifier son type:

- JSON
- HTML
- ou Texte

Et vous pouvez spécifier son format en utilisant la syntaxe des [templates de Golang](https://golang.org/pkg/text/template/).

Pour faire simple, vous avez accès aux propriétés du document JSON ci-dessus mais en les préfixant avec un point et une majuscule le tout entre deux accolades.
Par exemple, la propriété `title` est accessible avec la syntaxe `{{.Title}}`.

## Cinématique

Pour ajouter un webhook générique, vous devez [configurer votre webhook sortant](https://readflow.app/settings/integrations):

![](../../incoming-webhook/images/integrations.png)

Cliquer sur le bouton `Add` pour ajouter un webhook sortant.
La page d'ajout de webhook s'ouvre:

![](images/add-generic-webhook.png)

1. Saisissez un alias
1. Choisissez `generic` comme fournisseur de service
1. Configurez l'URL du webhook
1. Cochez la case si vous souhaitez en faire votre service par défaut

Une fois configuré, vous verrez une nouvelle entrée dans le menu contextuel des articles:

![](images/send-to-webhook.png)

Vous pouvez désormais envoyer un article vers un point d'accès HTTP.
