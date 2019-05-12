package service

import (
	"context"
	"errors"

	"github.com/ncarlier/readflow/pkg/event"
	"github.com/ncarlier/readflow/pkg/model"
	"github.com/ncarlier/readflow/pkg/readability"
)

// ArticleCreationOptions article creation options
type ArticleCreationOptions struct {
	IgnoreHydrateError bool
}

// CreateArticle creates new article
func (reg *Registry) CreateArticle(ctx context.Context, data model.ArticleForm, opts ArticleCreationOptions) (*model.Article, error) {
	uid := getCurrentUserFromContext(ctx)

	// TODO validate article!

	var category *model.Category
	if data.CategoryID != nil {
		cat, err := reg.GetCategory(ctx, *data.CategoryID)
		if err != nil {
			reg.logger.Info().Err(err).Uint(
				"uid", uid,
			).Str("title", data.Title).Msg("unable to create article")
			return nil, err
		}
		category = cat
	}

	if category == nil {
		// Process article by the rule engine
		if err := reg.ProcessArticleByRuleEngine(ctx, &data); err != nil {
			reg.logger.Info().Err(err).Uint(
				"uid", uid,
			).Str("title", data.Title).Msg("unable to create article")
			return nil, err
		}
	}

	builder := model.NewArticleBuilder()
	article := builder.UserID(
		uid,
	).Form(&data).Build()

	if article.URL != nil && (article.Image == nil || article.Text == nil || article.HTML == nil) {
		// Fetch original article to extract missing attributes
		if err := reg.HydrateArticle(ctx, article); err != nil {
			reg.logger.Info().Err(err).Uint(
				"uid", uid,
			).Str("title", article.Title).Msg("unable to fetch original article")
			// TODO excerpt and image should be extracted from HTML content
			if !opts.IgnoreHydrateError {
				return nil, err
			}
		}
	}

	reg.logger.Debug().Uint(
		"uid", uid,
	).Str("title", article.Title).Msg("creating article...")
	newArticle, err := reg.db.CreateOrUpdateArticle(*article)
	if err != nil {
		reg.logger.Info().Err(err).Uint(
			"uid", uid,
		).Str("title", article.Title).Msg("unable to create article")
		return nil, err
	}
	reg.logger.Info().Uint(
		"uid", uid,
	).Str("title", newArticle.Title).Uint("id", *newArticle.ID).Msg("article created")
	event.Emit(event.CreateArticle, *newArticle)
	return newArticle, nil
}

// CreateArticles creates new articles
func (reg *Registry) CreateArticles(ctx context.Context, data []model.ArticleForm) *model.Articles {
	result := model.Articles{}
	for _, art := range data {
		article, err := reg.CreateArticle(ctx, art, ArticleCreationOptions{
			IgnoreHydrateError: true,
		})
		if err != nil {
			result.Errors = append(result.Errors, err)
		}
		if article != nil {
			result.Articles = append(result.Articles, article)
		}
	}
	return &result
}

// CountArticles count articles
func (reg *Registry) CountArticles(ctx context.Context, req model.ArticlesPageRequest) (uint, error) {
	uid := getCurrentUserFromContext(ctx)

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

// HydrateArticle add missimg attributes form original article
func (reg *Registry) HydrateArticle(ctx context.Context, article *model.Article) error {
	art, err := readability.FetchArticle(ctx, *article.URL)
	if art == nil {
		return err
	}
	if article.Title == "" {
		article.Title = art.Title
	}
	if article.HTML == nil {
		article.HTML = art.HTML
	}
	if article.Text == nil {
		article.Text = art.Text
	}
	if article.Image == nil {
		article.Image = art.Image
	}

	return err
}
