package service

import (
	"context"
	"errors"
	"time"

	"github.com/ncarlier/readflow/internal/model"
)

// CountCurrentUserArticles count current user articles
func (reg *Registry) CountCurrentUserArticles(ctx context.Context, req model.ArticlesPageRequest) (uint, error) {
	uid := getCurrentUserIDFromContext(ctx)
	return reg.CountUserArticles(ctx, uid, req)
}

// CountUserArticles count user articles
func (reg *Registry) CountUserArticles(ctx context.Context, uid uint, req model.ArticlesPageRequest) (uint, error) {
	result, err := reg.db.CountArticlesByUser(uid, req)
	if err != nil {
		reg.logger.Info().Err(err).Uint(
			"uid", uid,
		).Msg("unable to count articles")
		return 0, err
	}
	return result, nil
}

// GetArticles get articles
func (reg *Registry) GetArticles(ctx context.Context, req model.ArticlesPageRequest) (*model.ArticlesPageResponse, error) {
	uid := getCurrentUserIDFromContext(ctx)

	result, err := reg.db.GetPaginatedArticlesByUser(uid, req)
	if err != nil {
		reg.logger.Info().Err(err).Uint(
			"uid", uid,
		).Msg("unable to get articles")
		return nil, err
	}
	return result, nil
}

// GetArticle get article
func (reg *Registry) GetArticle(ctx context.Context, id uint) (*model.Article, error) {
	uid := getCurrentUserIDFromContext(ctx)

	article, err := reg.db.GetArticleByID(id)
	if err != nil || article == nil || article.UserID != uid {
		if err == nil {
			err = errors.New("article not found")
		}
		reg.logger.Debug().Err(err).Uint(
			"uid", uid,
		).Uint("id", id).Msg("unable to get article")
		return nil, err
	}

	return article, nil
}

// MarkAllArticlesAsRead set status to read for all user's articles of a specific status and category
func (reg *Registry) MarkAllArticlesAsRead(ctx context.Context, status string, categoryID *uint) (int64, error) {
	start := time.Now()
	uid := getCurrentUserIDFromContext(ctx)

	nb, err := reg.db.MarkAllArticlesAsReadByUser(uid, status, categoryID)
	if err != nil {
		reg.logger.Info().Err(err).Uint(
			"uid", uid,
		).Msg("unable to mark all articles as read")
		return 0, err
	}
	reg.logger.Debug().Uint(
		"uid", uid,
	).Dur(
		"took", time.Since(start),
	).Msg("all articles marked as read")

	return nb, nil
}

// CleanHistory remove all read articles
func (reg *Registry) CleanHistory(ctx context.Context) (int64, error) {
	start := time.Now()
	uid := getCurrentUserIDFromContext(ctx)

	nb, err := reg.db.DeleteAllReadArticlesByUser(uid)
	if err != nil {
		reg.logger.Info().Err(err).Uint(
			"uid", uid,
		).Msg("unable to purge history")
		return 0, err
	}
	reg.logger.Debug().Uint(
		"uid", uid,
	).Dur(
		"took", time.Since(start),
	).Msg("history purged")

	return nb, nil
}
