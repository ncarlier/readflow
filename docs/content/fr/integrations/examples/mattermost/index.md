+++
title = "Mattermost ou Slack"
description = "Envoyer un article vers un canal Mattermost ou Slack"
weight = 2
+++

Envoyer un article vers un canal Mattermost ou Slack consiste à utiliser un webhook sortant avec un format de sortie JSON approprié.

Vous allez devoir créer un document JSON avec une propriété `text` comme ceci:

```json
{
	"text": ":tada: {{.Title}} (<{{.URL}}|voir plus>) cc @all",
}
```
