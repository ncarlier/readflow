package service

import (
	"context"
	"errors"

	"github.com/ncarlier/readflow/pkg/downloader"
	"github.com/ncarlier/readflow/pkg/exporter"
	"github.com/ncarlier/readflow/pkg/helper"
)

// DownloadArticle get article as a binary file
func (reg *Registry) DownloadArticle(ctx context.Context, idArticle uint, format string) (*downloader.WebAsset, error) {
	uid := getCurrentUserIDFromContext(ctx)

	logger := reg.logger.With().Uint(
		"uid", uid,
	).Uint(
		"article", idArticle,
	).Str(
		"format", format,
	).Logger()

	exp, err := exporter.NewArticleExporter(format, reg.dl)
	if err != nil {
		logger.Info().Err(err).Msg(ErrArticleDownload.Error())
		return nil, err
	}

	article, err := reg.db.GetArticleByID(idArticle)
	if err != nil || article == nil || article.UserID != uid {
		if err == nil {
			err = errors.New("article not found")
		}
		logger.Info().Err(err).Msg(ErrArticleDownload.Error())
		return nil, err
	}

	if article.HTML == nil || article.URL == nil {
		err := errors.New("missing require attributes")
		logger.Info().Err(err).Msg(ErrArticleDownload.Error())
		return nil, err
	}
	logger.Debug().Msg("preparing article download artifact")

	// get downloadable article from the cache
	key := helper.Hash(format, article.Hash)
	data, err := reg.downloadCache.Get(key)
	if err != nil {
		logger.Info().Err(err).Msg(ErrArticleDownload.Error())
	}
	if data != nil {
		reg.logger.Debug().Uint("uid", uid).Uint("id", idArticle).Msg("returns article download artefact from cache")
		return downloader.NewWebAsset(data)
	}

	// export article to the downloadable format
	result, err := exp.Export(ctx, article)
	if err != nil {
		logger.Info().Err(err).Msg(ErrArticleDownload.Error())
		return result, err
	}

	// TODO compute user quota

	// put downloadable article into the cache
	value, err := result.Encode()
	if err != nil {
		return nil, err
	}
	if err := reg.downloadCache.Put(key, value); err != nil {
		logger.Info().Err(err).Msg(ErrArticleDownload.Error())
	}

	reg.logger.Info().Uint("uid", uid).Uint("id", idArticle).Msg("article download artifact created")

	return result, nil
}
