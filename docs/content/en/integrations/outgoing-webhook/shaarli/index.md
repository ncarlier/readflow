+++
title = "Shaarli"
description = "Send articles to Shaarli"
weight = 3
+++

![](images/logo.svg)

[Shaarli](https://github.com/shaarli/Shaarli) is a a personal, minimalist, super-fast, database free, bookmarking service.

To use readflow with Shaarli you must first set the [API secret](https://shaarli.readthedocs.io/en/master/Shaarli-configuration/#shaarli-configurationhttps://doc.wallabag.org/fr/developer/api/oauth.html):

![](images/api-secret.png)

Once done, you can [configure your outgoing webhook](https://readflow.app/settings/integrations):

![](../../incoming-webhook/integrations.png)

Click on the `Add` button to add an outgoing webhook.
The webhook add page opens:

![](images/config.png)

1. Enter an alias
1. Choose `Shaarli` as provider
1. Configure the URL of the service
1. Set your API secret
1. Click on the checkbox if you want to make it your default service

The default archiving service can be invoked via the keyboard shortcut `shift+s` when viewing an article.

Once configured, you will see a new entry in the context menu of the article.

You can now send an article to Shaarli.
