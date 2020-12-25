+++
title = "Readflow"
description = "Share an article with another readflow"
weight = 3
+++

Sending an article to another readflow consists in using an outgoing webhook with an incoming webhook from another readflow.

Setting up the outgoing webhook is simply a matter of formatting the URL like this: `https://api:<TOKEN>@api.readflow.app/articles`

By replacing `<TOKEN>` with the token of the incoming webhook.
