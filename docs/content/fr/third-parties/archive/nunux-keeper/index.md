+++
title = "Nunux Keeper"
description = "Archiver des articles vers Nunux Keeper"
weight = 1
+++

![](images/nunux-keeper.png)

[Nunux Keeper](https://keeper.nunux.org) est un service Open Source d'archivage d'article.

Pour utiliser readflow avec Nunux Keeper vous devez au préalable obetnnir [une clé d'API dans Nunux Keeper](https://app.nunux.org/keeper/settings/api-key):

![](images/nunux-keeper-key.png)

Une fois obtenue, vous pouvez [configurer votre service d'archivage](https://readflow.app/settings/archive-services):

![](images/archive-services.png)

Cliquer sur le bouton `Add archive service`. La page d'ajout de service s'ouvre:

![](images/new-archive-service.png)

1. Saisissez un alias
1. Choisissez `Nunux Keeper` comme fournisseur de service
1. Configurez si nécessaire l'URL du service
1. Coller votre clé d'API
1. Cochez la case si vous souhaitez en faire votre service par défaut

Le service d'archivage par défaut peut être invoqué via le racourci clavier `s` lors de la visualisation d'un article.

Une fois configuré, vous verrez une nouvelle entrée dans le menu contextuel des articles:

![](images/save-to-keeper.png)

Vous pouvez désormais envoyer un article vers Nunux Keeper.
