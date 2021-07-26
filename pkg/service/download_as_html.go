package service

import (
	"bytes"
	"context"
	"text/template"

	"github.com/ncarlier/readflow/pkg/model"
)

var articleAsHTMLTpl = template.Must(template.New("article-as-html").Parse(`
<!DOCTYPE html>
<html>
<head>
	<title>{{ .Title }}</title>
	<meta charset="utf-8" />
	<meta name="og:title" content="{{ .Title }}"/>
	<meta name="og:url" content="{{ .URL }}"/>
	<meta name="og:image" content="{{ .Image }}"/>
	<meta name="og:revised" content="{{ .PublishedAt }}"/>
</head>
<body>{{ .HTML }}</body>
</html>
`))

func (reg *Registry) downloadArticleAsHTML(ctx context.Context, article *model.Article, offline bool) ([]byte, error) {
	var buffer bytes.Buffer
	if err := articleAsHTMLTpl.Execute(&buffer, article); err != nil {
		return nil, err
	}
	if !offline {
		return buffer.Bytes(), nil
	}

	return reg.webArchiver.Archive(ctx, &buffer, *article.URL)
}
