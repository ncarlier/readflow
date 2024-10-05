package service

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/ncarlier/readflow/internal/exporter"

	"github.com/ncarlier/readflow/pkg/defaults"
	"github.com/ncarlier/readflow/pkg/downloader"
	"github.com/ncarlier/readflow/pkg/utils"
)

// DownloadArticle get article as a binary file
func (reg *Registry) DownloadArticle(ctx context.Context, idArticle uint, format string) (*downloader.WebAsset, error) {
	uid := getCurrentUserIDFromContext(ctx)
	start := time.Now()

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

	// get timeout from user plan
	var timeout time.Duration
	if plan, _ := reg.GetCurrentUserPlan(ctx); plan != nil {
		timeout = plan.DownloadTimout.Duration
	}
	if timeout == 0 {
		timeout = defaults.Timeout
	}

	// get downloadable article from the cache
	key := utils.Hash(format, article.Hash)
	data, err := reg.downloadCache.Get(key)
	if err != nil {
		logger.Info().Err(err).Msg(ErrArticleDownload.Error())
	}
	if data != nil {
		logger.Debug().Msg("get article downloadable asset from cache")
		return downloader.NewWebAsset(data)
	}

	// export article to the downloadable format
	logger.Debug().Msg("preparing article downloadable asset...")
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
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

	logger.Info().Dur("took", time.Since(start)).Msg("article downloadable asset created")

	return result, nil
}

// Download web asset
func (reg *Registry) Download(ctx context.Context, url string, header *http.Header) (*downloader.WebAsset, *http.Response, error) {
	return reg.dl.Get(ctx, url, header)
}
