+++
title = "Scripting"
description = "Script to personalize the integration of the articles"
weight = 1
+++

An incoming webhook can run a script on each article it receives.

## Features

This script gives control over a set of features:

- Change the title of the article (`setTitle(string)`)
- Change the text of the article (`setText(string)`)
- Categorize the article (`setCategory(string)`)
- Notify connected devices of the article (`sendNotification()`)
- Call an outgoing webhook (`triggerWebhook(string)`)

The script expects a boolean value in return.
This value decides the fate of the article:

- `return true;`: the article will be integrated
- `return false;`: the article will be ignored

> Note that if the article is ignored then the above features have no effect.

## Syntax

You are free to implement the orchestration logic of these actions with a [simple syntax](https://github.com/skx/evalfilter).

You can access the following attributes within your script:

- `Title`: the title of the article
- `Text`: the text of the article
- `HTML`: the HTML content of the article
- `URL`: the URL of the article
- `Tags`: the tags of the article (string array)

## Examples

Here are some examples to illustrate the possibilities:

### Categorize an article according to its tags

```c
if ("news" in Tags) {
    setCategory("News");
}
return true;
```

### Add a prefix to the title

```c
setTitle(sprintf("[From my awesome webhook] %s", Title);
return true;
```

### Send a notification on a particular topic

```c
if ("news" in Tags) {
    setCategory("News");
    if (Title ~= /breaking|important/i ) {
        sendNotification();
    }
}
return true;
```

### Filter a topic without interest

```c
if (Title ~= /boring|stupid/i ) {
    return false;
}
return true;
```
