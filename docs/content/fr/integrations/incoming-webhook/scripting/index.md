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
- Notifier l'article aux appareils connectés (`sendNotification()`)
- Appeler un webhook sortant (`triggerWebhook(string)`)

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
- `Tags`: les tags de l'article (tableau de string)

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
setTitle(sprintf("[From my awesome webhook] %s", Title);
return true;
```

### Envoyer une notification sur un sujet particulier

```c
if ("news" in Tags) {
    setCategory("News");
    if (Title ~= /breaking|important/i ) {
        sendNotification();
    }
}
return true;
```

### Filtrer un sujet sans intérêt

```c
if (Title ~= /boring|stupid/i ) {
    return false;
}
return true;
```
