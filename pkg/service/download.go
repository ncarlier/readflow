package service

import (
	"context"
	"errors"

	"github.com/ncarlier/readflow/pkg/exporter"
	"github.com/ncarlier/readflow/pkg/helper"
	"github.com/ncarlier/readflow/pkg/model"
)

// DownloadArticle get artice as a binary file
func (reg *Registry) DownloadArticle(ctx context.Context, idArticle uint, format string) (*model.FileAsset, error) {
	uid := getCurrentUserIDFromContext(ctx)

	logger := reg.logger.With().Uint(
		"uid", uid,
	).Uint(
		"article", idArticle,
	).Str(
		"format", format,
	).Logger()

	conv, err := exporter.NewArticleExporter(format, reg.downloader)
	if err != nil {
		logger.Info().Err(err).Msg(ErrArticleArchiving.Error())
		return nil, err
	}

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
	logger.Debug().Msg("preparing article download artifact")

	key := helper.Hash(format, article.Hash)
	result, err := reg.downloadCache.Get(key)
	if err != nil {
		logger.Info().Err(err).Msg(ErrArticleArchiving.Error())
	}
	if result != nil {
		reg.logger.Debug().Uint("uid", uid).Uint("id", idArticle).Msg("returns article download artefact from cache")
		return result, nil
	}

	result, err = conv.Export(ctx, article)
	if err != nil {
		logger.Info().Err(err).Msg(ErrArticleArchiving.Error())
		return result, err
	}

	// TODO compute user quota

	if err := reg.downloadCache.Put(key, result); err != nil {
		logger.Info().Err(err).Msg(ErrArticleArchiving.Error())
	}

	reg.logger.Info().Uint("uid", uid).Uint("id", idArticle).Msg("article download artifact created")

	return result, nil
}
