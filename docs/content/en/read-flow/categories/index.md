+++
title = "Categories"
description = "Use categories to organize your reading flows"
weight = 2
+++

In readflow you can split your reading flows into categories.
Adding an article to a category will be based on the rules you have configured.

## Categories

To manage categories, go to [the configuration screen](https://readflow.app/settings/categories):

![](images/categories.png)

You can add a category by clicking on `Add category` button.

![](images/add.png)

A category is defined by :

- a title
- an assignment rule
- a notification strategy

## Notification strategy

It is possible to specify the notification strategy for a given category. It is possible to:

- Not to send a notification.
- Send a notification as soon as an article is received (but with a frequency limit).
- Send a notification following the global rules.

The global rules depend on the readflow configuration but follow the following principles:

A notification is sent if:

- An article group has been created;
- The user has not used the application for a while;
- And within a given frequency limit.

## Rule

When adding an article, the rule engine will be activated and apply the rules ordered by categories.
At the first validated rule, the article is placed in the target category.
If no rules match then the article will not have a category.

The definition of a rule is a pseudo code whose result must be true or false.

Within the rule it is possible to refer to some attributes:

- `title`: article title
- `text`: article text content
- `url`: article source URL
- `tags`: article input tags
- `webhook`: incoming webhook alias

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
  *Example:* `webhook in ["test", "bookmarklet"]`

#### Functions

- `len` (length of the character string)
   *Example:* `len(text) >= 100`

### Examples:

Classify articles whose incoming webhook alias is "foo":

```js
webhook == "foo"
```

Classify articles with incoming webhook alias is "foo" or "bar":

```js
webhook == "foo" || webhook == "bar"
// Can also be written:
webhook in ["foo", "bar"]
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
