+++
title = "Scripting"
description = "Script de personalisation de l'intégration des articles"
weight = 1
+++

Un webhook entrant peut exécuter un script sur chaque article qu'il reçoit.

## Fonctionnalités

Ce script donne la main sur un ensemble de fonctionnalités:

- Changer le titre de l'article (`setTitle(string)`)
- Changer le texte de l'article (`setText(string)`)
- Catégoriser l'article (`setCategory(string)`)
- Marquer l'article comme lu (`markAsRead()`)
- Marquer l'article comme à lire (`markAsToRead()`)
- Notifier l'article aux appareils connectés (`sendNotification()`)
- Appeler un webhook sortant (`triggerWebhook(string)`)
- Désactiver la politique de notification globale pour ce webhook (`disableGlobalNotification()`)

Le script attend en retour une valeur booléenne.
Cette valeur décide du sort de l'article:

- `return true;`: l'article va être intégré
- `return false;`: l'article va être ignoré

> Notez que si l'article est ignoré alors les fonctionnalités ci-dessus n'ont aucun effet.

## Syntaxe

Vous êtes libre d'implémenter la logique d'orchestration de ces actions grâce à une [syntaxe simple](https://github.com/skx/evalfilter).

Vous pouvez accéder au sein de votre script aux attributs suivants:

- `Title`: le titre de l'article
- `Text`: le texte de l'article
- `HTML`: le contenu HTML de l'article
- `URL`: l'URL de l'article
- `Origin`: l'origine de l'article
- `Tags`: les *tags* **et** *hashtags* de l'article (tableau de string)

> Notez que les *hashtags* sont extraits du titre et du texte de l'article, tandis que les *tags* sont fournis avec l'article en entrée.
> Seuls les hashtags sont conservés (car ils font partie du titre et du texte).
> Les tags d'entrée ne sont pas conservées et ne sont utilisées que par le script.
> Si vous souhaitez conserver les tags, il est recommandé de les inclure en tant que hashtags dans le titre ou le texte.

## Exemples

Voici quelques exemples qui illustrent les possibilités:

### Catégoriser un article en fonction de ses tags

```c
if ("news" in Tags) {
    setCategory("News");
}
return true;
```

### Ajouter un préfixe au titre

```c
setTitle(sprintf("[From my awesome webhook] %s", Title));
return true;
```

### Envoyer une notification sur un sujet particulier

```c
if ("news" in Tags) {
    setCategory("News");
    if (Title ~= /breaking|important/i) {
        sendNotification();
    }
}
return true;
```

### Filtrer un sujet sans intérêt

```c
if (Title ~= /boring|stupid/i) {
    return false;
}
return true;
```

### Accepter les article provenant d'un autre utilisateur readflow

```c
if (Origin == "johndoe@example.com") {
    setTitle(sprintf("[From %s] %s", Origin, Title));
    return true;
}
return false;
```
