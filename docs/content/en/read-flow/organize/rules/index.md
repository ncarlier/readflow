+++
title = "The rule engine"
description = "Use the rule engine to organize your reading flows"
weight = 2
+++

When adding an article, the rule engine will be activated and apply the rules ordered by priority.
At the first validated rule, the article is placed in the target category.
If no rules match then the article will not have a category.

To manage the rules, go to[the configuration screen](https://readflow.app/settings/rules):

![](images/rules.png)

Click on the `Add rule` button to add a rule:

![](images/add-rule.png)

A rule is:

- An alias (its visual identifier)
- A priority (its order of execution)
- A definition
- And a target category

## The definition of a rule

The definition of a rule is a pseudo code whose result must be true or false.

Within the rule it is possible to refer to some attributes:

- `title`: article title
- `text`: article text content
- `url`: article source URL
- `tags`: article input tags
- `key`: alias of the used API key

### Syntax

#### Operators

- `==` (equal)
- `!=` (non equal)
- `matches` (validate a regular expression)
- `not ("foo" matches "bar")` (does not validate a regular expression)

#### Logical operators

- `not` or `!`
- `and` or `&&`
- `or` or `||`

#### Other operators

- `~` (concatenation)
  *Example:* `'Harry' ~ ' ' ~ 'Potter'` will be `Harry Potter`
- `in`
- `not in`
  *Example:* `key in ["test", "bookmarklet"]`

#### Functions

- `len` (length of the character string)
   *Example:* `len(text) >= 100`

### Examples:

Classify articles whose API key is "foo":

```js
key == "foo"
```

Classify articles with API key "foo" or "bar":

```js
key == "foo" || key == "bar"
// Peut aussi s'écrire:
key in ["foo", "bar"]
```

Classify articles with "foo" as tag:

```js
"foo" in tags
```

Classify articles with titles containing "Amazon" and "Alexa":

```js
title matches "Amazon" and title matches "Alexa"
```

Classify items that come from CNN:

```js
url matches "^https://edition.cnn.com"
```
