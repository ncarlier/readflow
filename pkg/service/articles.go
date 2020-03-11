package service

import (
	"context"
	"errors"

	"github.com/ncarlier/readflow/pkg/model"
)

// CountCurrentUserArticles count current user articles
func (reg *Registry) CountCurrentUserArticles(ctx context.Context, req model.ArticlesPageRequest) (uint, error) {
	uid := getCurrentUserFromContext(ctx)
	return reg.CountUserArticles(ctx, uid, req)
}

// CountUserArticles count user articles
func (reg *Registry) CountUserArticles(ctx context.Context, uid uint, req model.ArticlesPageRequest) (uint, error) {
	result, err := reg.db.CountArticlesByUserID(uid, req)
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
	uid := getCurrentUserFromContext(ctx)

	result, err := reg.db.GetPaginatedArticlesByUserID(uid, req)
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
	uid := getCurrentUserFromContext(ctx)

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

// UpdateArticleStatus update article status
func (reg *Registry) UpdateArticleStatus(ctx context.Context, id uint, status string) (*model.Article, error) {
	uid := getCurrentUserFromContext(ctx)

	article, err := reg.GetArticle(ctx, id)
	if err != nil {
		return nil, err
	}

	article.Status = status
	article, err = reg.db.CreateOrUpdateArticle(*article)
	if err != nil {
		reg.logger.Info().Err(err).Uint(
			"uid", uid,
		).Uint("id", id).Msg("unable to update article")
		return nil, err
	}
	// TODO maybe too verbose... debug level is maybe an option here
	reg.logger.Info().Uint(
		"uid", uid,
	).Uint("id", *article.ID).Str("status", article.Status).Msg("article status updated")

	return article, nil
}

// MarkAllArticlesAsRead set status to read for all articles (of a category if provided)
func (reg *Registry) MarkAllArticlesAsRead(ctx context.Context, categoryID *uint) (int64, error) {
	uid := getCurrentUserFromContext(ctx)

	nb, err := reg.db.MarkAllArticlesAsRead(uid, categoryID)
	if err != nil {
		reg.logger.Info().Err(err).Uint(
			"uid", uid,
		).Msg("unable to mark all articles as read")
		return 0, err
	}
	reg.logger.Debug().Uint(
		"uid", uid,
	).Msg("all articles marked as read")

	return nb, nil
}

// CleanHistory remove all read articles
func (reg *Registry) CleanHistory(ctx context.Context) (int64, error) {
	uid := getCurrentUserFromContext(ctx)

	nb, err := reg.db.DeleteAllReadArticles(uid)
	if err != nil {
		reg.logger.Info().Err(err).Uint(
			"uid", uid,
		).Msg("unable to purge history")
		return 0, err
	}
	reg.logger.Debug().Uint(
		"uid", uid,
	).Msg("history purged")

	return nb, nil
}
