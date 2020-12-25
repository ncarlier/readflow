+++
title = "Mattermost or Slack"
description = "Send an article to a Mattermost or Slack channel"
weight = 2
+++

Sending an article to a Mattermost or Slack channel involves using an outgoing webhook with an appropriate JSON output format.

You will need to create a JSON document with a `text` property like this:

```json
{
	"text": ":tada: {{.Title}} (<{{.URL}}|voir plus>) cc @all",
}
```
