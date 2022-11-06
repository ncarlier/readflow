+++
title = "Incoming webhook"
description = "Add content to readflow from third-party services"
weight = 1
+++

An incoming webhook allows a third party service to create content in readflow.

To add an incoming webhook, go to [the integration configuration screen](https://readflow.app/settings/integrations).

![](integrations.png)

Click the `Add` button to add an incoming webhook:

![](add-incoming-webhook.png)

Enter an alias for your webhook.

You can customize the integration of the articles by this webhook via a [script](scripting).

The service using the webhook must implement the [Integration API](integration-api):

{{%children style="h4" description="true"%}}
