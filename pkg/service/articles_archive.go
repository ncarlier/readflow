package service

import (
	"bytes"
	"context"
	"errors"
	"text/template"

	"github.com/go-shiori/obelisk"
	"github.com/ncarlier/readflow/pkg/constant"
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

// ArchiveArticle save artice as a single HTML page
func (reg *Registry) ArchiveArticle(ctx context.Context, idArticle uint) ([]byte, error) {
	uid := getCurrentUserIDFromContext(ctx)

	logger := reg.logger.With().Uint(
		"uid", uid,
	).Uint("article", idArticle).Logger()

	article, err := reg.db.GetArticleByID(idArticle)
	if err != nil || article == nil || article.UserID != uid {
		if err == nil {
			err = errors.New("article not found")
		}
		logger.Info().Err(err).Msg(ErrArticleArchiving.Error())
		return nil, err
	}

	if article.HTML == nil || article.URL == nil {
		err := errors.New("missing require attributes")
		logger.Info().Err(err).Msg(ErrArticleArchiving.Error())
		return nil, err
	}

	data, err := reg.downloadCache.Get(article.Hash)
	if err != nil {
		logger.Info().Err(err).Msg(ErrArticleArchiving.Error())
	}
	if data != nil {
		reg.logger.Debug().Uint("uid", uid).Uint("id", idArticle).Msg("article archive get from cache")
		return data, nil
	}

	var buffer bytes.Buffer
	err = articleAsHTMLTpl.Execute(&buffer, article)
	if err != nil {
		logger.Info().Err(err).Msg(ErrArticleArchiving.Error())
	}

	req := obelisk.Request{
		Input: &buffer,
		URL:   *article.URL,
	}

	arc := obelisk.Archiver{
		UserAgent:      constant.UserAgent,
		RequestTimeout: constant.DefaultTimeout,
		EnableLog:      reg.logger.Debug().Enabled(),
	}
	arc.Validate()

	result, _, err := arc.Archive(ctx, req)
	if err != nil {
		logger.Info().Err(err).Msg(ErrArticleArchiving.Error())
		return nil, err
	}

	if err := reg.downloadCache.Put(article.Hash, result); err != nil {
		logger.Info().Err(err).Msg(ErrArticleArchiving.Error())
	}

	reg.logger.Info().Uint("uid", uid).Uint("id", idArticle).Msg("article archive created")

	return result, nil
}
