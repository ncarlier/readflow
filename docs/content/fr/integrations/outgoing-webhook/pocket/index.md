+++
title = "Pocket"
description = "Envoyer des articles vers Pocket"
weight = 4
+++

![](images/pocket-logo.png)

[Pocket](https://getpocket.com) est un service qui sauvegarde des articles et vous permet de les lire plus tard.

Pour utiliser readflow avec Pocket vous devez [configurer votre webhook sortant](https://readflow.app/settings/integrations):

![](../../incoming-webhook/integrations.png)

Cliquer sur le bouton `Add` pour ajouter un webhook sortant.
La page d'ajout de webhook s'ouvre:

1. Choisissez `Pocket` comme fournisseur de service
1. Cliquez sur le bouton `Link with Pocket`
1. Vous êtes redirigé sur Pocket qui vous demande si vous voulez autorizer readflow à se connecter à votre compte Pocket
1. En acceptant vous êtes redirigé sur readflow avec les champs pré-remplis
1. Saisissez un alias
1. Cochez la case si vous souhaitez en faire votre service par défaut

![](images/pocket-configuration.png)

Le service d'archivage par défaut peut être invoqué via le racourci clavier `shift+s` lors de la visualisation d'un article.

Une fois configuré, vous verrez une nouvelle entrée dans le menu contextuel des articles.

Vous pouvez désormais envoyer un article vers Pocket.
