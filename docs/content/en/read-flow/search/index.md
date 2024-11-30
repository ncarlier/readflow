+++
title = "Search"
description = "Use the search engine to find your articles"
weight = 2
+++

In readflow you can find your articles using the search engine.

The search engine supports the following features:

- a word sequence is treated as a logical search with an AND operator.
  - *e.g. `dark vador jedi` will search for all articles containing the words `dark`, `vador` and `jedi`.*
- text inside quotation marks is treated as a phrase search.
  - *e.g. `“dark vador”` will find all articles containing the word sequence `dark vador`.*
- use of the keyword `OR` (logical OR) to search for one expression or another.
  - *e.g. `sith or jedi` will search for all articles containing the word `sith` OR the word `jedi`.*
- use of the `-` character to exclude an expression from the search.
  - *e.g. `sith -jedi` will search for all articles containing the word `sith` WITHOUT the word `jedi`.*
- hashtags
  - *e.g. `#droid` will search for all articles with the tag `#droid` in the title or text.*

You can combine these features to refine your search.
