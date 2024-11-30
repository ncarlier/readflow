+++
title = "Chercher"
description = "Utiliser le moteur de recherche pour retrouver vos articles"
weight = 2
+++

Dans readflow, vous pouvez retrouver vos articles en utilisant le moteur de recherche.

Le moteur de recherche prend en charge les fonctionnalités suivantes :

- une séquence de mots est traitée comme une recherche logique avec un opérateur AND.
  - *Par exemple, la requête `dark vador jedi` recherchera tous les articles contenant les mots `dark`, `vador` et `jedi`.*
- Le texte entre guillemets est traité comme une recherche de phrase.
  - *Par exemple, la requête `« dark vador »` recherchera tous les articles contenant la séquence de mots `dark vador`.*
- L'utilisation du mot-clé `OR` (logical OR) pour rechercher une expression ou une autre.
  - *Par exemple, la requête `sith ou jedi` recherchera tous les articles contenant le mot `sith` OU le mot `jedi`.*
- L'utilisation du caractère `-` pour exclure une expression de la recherche.
  - *Par exemple, la requête `sith -jedi` recherchera tous les articles contenant le mot `sith` SANS le mot `jedi`.*
- Les hashtags
  - *Par exemple, la requête `#droid` recherchera tous les articles contenant le tag `#droid` dans le titre ou le texte*.

Vous pouvez combiner ces fonctionnalités pour affiner votre recherche.
