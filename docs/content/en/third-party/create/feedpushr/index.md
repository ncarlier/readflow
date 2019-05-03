+++
title = "RSS feed with feedpushr"
description = "Use feedpushr to integrate RSS feeds"
weight = 2
+++

![](images/feedpushr.png)

[Feedpushr](https://github.com/ncarlier/feedpushr) is a powerful open source RSS feed aggregator able to process and send articles to third party services.
And among these services, we have readflow.

To use readflow with feedpushr, you must use and configure its plugin:

```bash
$ # Set readflow URL
$ export APP_READFLOW_URL=https://api.readflow.app
$ # Set API key
$ export APP_READFLOW_API_KEY=d70*************456
$ # Start feedpushr
$ feedpushr --log-pretty --plugin ./feedpushr-readflow.so --output readflow://
```

You should see this on the feedpushr UI:

![](images/feedpushr-ui.png)

You can then import your OPML subscriptions into feedpushr and see your articles in readflow.
