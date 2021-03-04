+++
title = "Wallabag"
description = "Send articles to Wallabag"
weight = 3
+++

![](images/wallabag.webp)

[Wallabag](https://www.wallabag.org) is a service that saves articles and allows you to read them later.

To use readflow with Wallabag you must first create [an API client](https://doc.wallabag.org/fr/developer/api/oauth.html):

![](images/api-credentials.png)

Once obtained, you can [configure your outgoing webhook](https://readflow.app/settings/integrations):

![](../../incoming-webhook/integrations.png)

Click on the `Add` button to add an outgoing webhook.
The webhook add page opens.

1. Enter an alias
1. Choose `Wallabag` as provider
1. Configure the URL of the service if necessary
1. Past you client ID and secret
1. Enter your credentials
1. Click on the checkbox if you want to make it your default service

The default archiving service can be invoked via the keyboard shortcut `s` when viewing an article.

Once configured, you will see a new entry in the context menu of the article.

You can now send an article to Wallabag.
