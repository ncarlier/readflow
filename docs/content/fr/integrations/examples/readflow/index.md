+++
title = "Readflow"
description = "Partager un article avec un autre readflow"
weight = 3
+++

Envoyer un article vers un autre readflow consiste à utiliser un webhook sortant vers un webhook entrant d'un autre readflow.

Configurer le webhook sortant consiste simplement à formater l'URL comme ceci: `https://api:<TOKEN>@api.readflow.app/articles`

En remplacent `<TOKEN>` par le token du webhook entrant.
