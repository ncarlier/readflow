package service

import (
	"context"

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

	plan, err := reg.GetCurrentUserPlan(ctx)
	if err != nil {
		return nil, err
	}

	if plan != nil {
		// Check user quota
		req := model.ArticlesPageRequest{}
		totalArticles, err := reg.CountCurrentUserArticles(ctx, req)
		if err != nil {
			reg.logger.Info().Err(err).Uint(
				"uid", uid,
			).Str("title", data.Title).Msg("unable to create article")
			return nil, err
		}
		if totalArticles >= plan.TotalArticles {
			err = ErrUserQuotaReached
			reg.logger.Info().Err(err).Uint(
				"uid", uid,
			).Str("title", data.Title).Uint(
				"total", totalArticles,
			).Msg("unable to create article")
			return nil, err
		}
	}

	// TODO validate article!

	// Get category if specified
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
		if err := reg.hydrateArticle(ctx, article); err != nil {
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

// hydrateArticle add missing attributes form original article
func (reg *Registry) hydrateArticle(ctx context.Context, article *model.Article) error {
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
