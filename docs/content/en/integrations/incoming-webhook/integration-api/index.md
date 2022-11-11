+++
title = "Integration API"
description = "Using the content creation API"
weight = 2
+++

The integration API is a simple URL accessible in POST with a JSON payload.

The payload can be an article or an array of articles.
An article can have the following structure:

- `url`: The original URL of the article.
- `title`: The article title.
- `html`: HTML content of the article.
- `text`: Text content of the article (lmost often a summary).
- `image`: The illustration dof the article.
- `published_at`: Publication date of the article.
- `category`: The target category in readflow.
- `origin`: The origin of the article creation.
- `tags`: Comma separated list of tags.

All properties are optional except the title.

*Example:*

```json
[
    {
        "title": "",
        "url": "http://example.org/foo.html"
    },
    {
        "title": "Lorem ipsum 1",
        "html": "<p>Lorem ipsum <strong>dolor</strong> sit amet, consectetur adipiscing elit, ...</p>"
    },
    {
        "title": "Lorem ipsum 2",
        "html": "<p>Lorem ipsum <strong>dolor</strong> sit amet, consectetur adipiscing elit, ...</p>",
        "text": "Lorem ipsum dolor sit amet, consectetur adipiscing elit, ...",
        "image": "http://example.org/foo.png",
        "published_at": "2019-04-30T14:00:20.786Z",
        "category": "test",
    }
]
```

The API uses the incoming webhook token as basic authentication:

```bash
$ cat payload.json | http \
  -a api:89b5700d-e4da-407e-94a0-7303417189c5 \
  https://api.readflow.app/articles
```

If the URL is present and the article is incomplete (no HTML content, or text, etc...) then the original article is downloaded and processed before integration.

The process consists of:

- clean up HTML content
- extract the title of the page or [OpenGraph][opengraph] tags.
- extract text from HTML content or OpenGraph tags
- extract the illustration from OpenGraph tags

In return the API produces a JSON document with the following properties:

- `articles`: the list of the created articles (with their internal identifier)
- `errors`: any errors found

[opengraph]: http://ogp.me/
