+++
title = "Le moteur de règle"
description = "Utiliser le moteur de règle pour organiser vos flux de lecture"
weight = 2
+++

Lors de l'ajout d'un article, le moteur de règle va entrer en action et appliquer les règles selon l'ordre de priorité.
A la première règle validée, l'article est placé dans la catégorie cible.
Si aucune règle ne matche alors l'article n'aura pas de catégorie.

Pour gérer les règles, allez sur [l'écran de configuration](https://readflow.app/settings/rules):

![](images/rules.png)

Cliquez sur le bouton `Add rule` pour ajouter une règle:

![](images/add-rule.png)


Une règle est composée:

- D'un alias (son identifiant visuel)
- D'une priorité (son ordre d'exécution)
- D'une définition
- Et d'une catégorie cible

## La définition d'une règle

La définition d'une règle est un pseudo code dont le résultat doit être vrai ou faux.

Au sein de la règle il est possible de faire référence à l'article (`article`) et la clé d'API (`key`).

Un article contient les attributs suivants:

- `article.Title`: son titre
- `article.Text`: son content texte (le résumé)
- `article.HTML`: son contenu HTML
- `article.URL`: son URL d'origine
- `article.Image`: son illustration

### La syntaxe

#### Les opérateurs

- `==` (égal)
- `!=` (non égal)
- `matches` (valide une expression régulière)
- `not ("foo" matches "bar")` (ne valide pas une expression régulière)

#### Les opérateurs logiques

- `not` ou `!` (non)
- `and` ou `&&` (et)
- `or` ou `||` (ou)

#### Les autres opérateurs

- `~` (concatenation)
  *Exemple:* `'Harry' ~ ' ' ~ 'Potter'` donnera `Harry Potter`
- `in` (contient)
- `not in` (ne contient pas)
  *Exemple:* `key in ["test", "bookmarklet"]`

#### Les fonctions

- `len` (longueur de la chaine de caractères)
   *Exemple:* `len(article.Text) == 0`

### Exemples:

Classer les articles dont la clé d'API est "foo":

```js
key == "foo"
```

Classer les articles dont la clé d'API est "foo" ou "bar":

```js
key == "foo" || key == "bar"
// Peut aussi s'écrire:
key in ["foo", "bar"]
```

Classer les articles dont le titre contient "Amazon" et "Alexa":

```js
article.Title matches "Amazon" and article.Title matches "Alexa"
```

Classer les articles qui viennent de CNN:

```js
article.URL matches "^https://edition.cnn.com"
```
