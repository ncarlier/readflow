package service

import (
	"context"
	"errors"

	"github.com/ncarlier/readflow/pkg/constant"
	"github.com/ncarlier/readflow/pkg/helper"
)

// DownloadArticle get artice as a binary file
func (reg *Registry) DownloadArticle(ctx context.Context, idArticle uint, format string) ([]byte, string, error) {
	uid := getCurrentUserIDFromContext(ctx)
	contentType := getFormatContentType(format)

	logger := reg.logger.With().Uint(
		"uid", uid,
	).Uint(
		"article", idArticle,
	).Str(
		"format", format,
	).Logger()

	article, err := reg.db.GetArticleByID(idArticle)
	if err != nil || article == nil || article.UserID != uid {
		if err == nil {
			err = errors.New("article not found")
		}
		logger.Info().Err(err).Msg(ErrArticleArchiving.Error())
		return nil, contentType, err
	}

	if article.HTML == nil || article.URL == nil {
		err := errors.New("missing require attributes")
		logger.Info().Err(err).Msg(ErrArticleArchiving.Error())
		return nil, contentType, err
	}
	logger.Debug().Msg("preparing article download artifact")

	key := helper.Hash(format, article.Hash)
	data, err := reg.downloadCache.Get(key)
	if err != nil {
		logger.Info().Err(err).Msg(ErrArticleArchiving.Error())
	}
	if data != nil {
		reg.logger.Debug().Uint("uid", uid).Uint("id", idArticle).Msg("returns article download artefact from cache")
		return data, contentType, nil
	}

	var result []byte
	switch format {
	case "offline":
		result, err = reg.downloadArticleAsHTML(ctx, article, true)
	default:
		result, err = reg.downloadArticleAsHTML(ctx, article, false)
	}

	if err != nil {
		logger.Info().Err(err).Msg(ErrArticleArchiving.Error())
		return nil, contentType, err
	}

	if err := reg.downloadCache.Put(key, result); err != nil {
		logger.Info().Err(err).Msg(ErrArticleArchiving.Error())
	}

	reg.logger.Info().Uint("uid", uid).Uint("id", idArticle).Msg("article download artifact created")

	return result, contentType, nil
}

func getFormatContentType(format string) string {
	switch format {
	case "makdown":
		return constant.ContentTypeText
	}
	return constant.ContentTypeHTML
}
