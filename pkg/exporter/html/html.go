package html

import (
	"bytes"
	"context"
	"text/template"

	"github.com/ncarlier/readflow/pkg/constant"
	"github.com/ncarlier/readflow/pkg/exporter"
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

// HTMLExporter convert an article to HTML format
type HTMLExporter struct{}

func newHTMLExporter(downloader exporter.Downloader) (exporter.ArticleExporter, error) {
	return &HTMLExporter{}, nil
}

// Export an article to HTML format
func (exp *HTMLExporter) Export(ctx context.Context, article *model.Article) (*model.FileAsset, error) {
	var buffer bytes.Buffer
	if err := articleAsHTMLTpl.Execute(&buffer, article); err != nil {
		return nil, err
	}
	return &model.FileAsset{
		Data:        buffer.Bytes(),
		ContentType: constant.ContentTypeHTML,
		Name:        article.Title + ".html",
	}, nil
}

func init() {
	exporter.Register("html", newHTMLExporter)
}
