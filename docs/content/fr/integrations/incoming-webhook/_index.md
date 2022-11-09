+++
title = "Webhook entrant"
description = "Ajouter du contenu à readflow depuis des services tiers"
weight = 1
+++

Un webhook entrant permet à un service tiers de créer du contenu dans readflow.

Pour ajouter un webhook entrant, allez sur [l'écran de configuration des intégrations](https://readflow.app/settings/integrations).

![](integrations.png)

Cliquez sur le bouton `Add` pour ajouter un webhook entrant:

![](add-incoming-webhook.png)

Saisissez un alias pour votre webhook.

Vous pouvez personnaliser l'intégration des articles par ce webhook via un [script](scripting).

Le service utilisant le webhook doit implémenter l'[API d'intégration](integration-api) de readflow.

> Le nombre de webhooks entrants est limité par un quota dépendant du plan utilisateur.
> L'utilisation d'un webhook entrant est également contrôlée par une limitation dedit.

{{%children style="h4" description="true"%}}
