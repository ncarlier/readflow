package service

import (
	"context"
	"errors"
	"strings"

	"github.com/go-shiori/obelisk"
	"github.com/ncarlier/readflow/pkg/constant"
)

// ArchiveArticle save artice as a single HTML page
func (reg *Registry) ArchiveArticle(ctx context.Context, idArticle uint) ([]byte, error) {
	uid := getCurrentUserFromContext(ctx)

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

	input := strings.NewReader(*article.HTML)

	req := obelisk.Request{
		Input: input,
		URL:   *article.URL,
	}

	arc := obelisk.Archiver{
		UserAgent:      constant.UserAgent,
		RequestTimeout: constant.DefaultTimeout,
		EnableLog:      reg.logger.Debug().Enabled(),
	}
	arc.Validate()

	result, _, err := arc.Archive(context.Background(), req)
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
