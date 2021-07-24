package service

import (
	"context"

	"github.com/ncarlier/readflow/pkg/event"
	"github.com/ncarlier/readflow/pkg/html"
	"github.com/ncarlier/readflow/pkg/model"

	// activate all content providers
	_ "github.com/ncarlier/readflow/pkg/scraper/content-provider/all"
)

// ArticleCreationOptions article creation options
type ArticleCreationOptions struct {
	IgnoreHydrateError bool
}

// CreateArticle creates new article
func (reg *Registry) CreateArticle(ctx context.Context, form model.ArticleCreateForm, opts ArticleCreationOptions) (*model.Article, error) {
	uid := getCurrentUserIDFromContext(ctx)

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
			).Str("title", form.TruncatedTitle()).Msg("unable to create article")
			return nil, err
		}
		if totalArticles >= plan.TotalArticles {
			err = ErrUserQuotaReached
			reg.logger.Debug().Err(err).Uint(
				"uid", uid,
			).Str("title", form.TruncatedTitle()).Uint(
				"total", totalArticles,
			).Msg("unable to create article")
			return nil, err
		}
	}

	// TODO validate article!

	// Get category if specified
	var category *model.Category
	if form.CategoryID != nil {
		cat, err := reg.GetCategory(ctx, *form.CategoryID)
		if err != nil {
			reg.logger.Info().Err(err).Uint(
				"uid", uid,
			).Str("title", form.TruncatedTitle()).Msg("unable to create article")
			return nil, err
		}
		category = cat
	}

	if category == nil {
		// Process article by the rule engine
		if err := reg.ProcessArticleByRuleEngine(ctx, &form); err != nil {
			reg.logger.Info().Err(err).Uint(
				"uid", uid,
			).Str("title", form.TruncatedTitle()).Msg("unable to create article")
			return nil, err
		}
	}

	if form.URL != nil && !form.IsComplete() {
		// Fetch original article in order to extract missing attributes
		if err := reg.hydrateArticle(ctx, &form); err != nil {
			reg.logger.Info().Err(err).Uint(
				"uid", uid,
			).Str("title", form.TruncatedTitle()).Msg("unable to fetch original article")
			// TODO excerpt and image should be extracted from HTML content
			if !opts.IgnoreHydrateError {
				return nil, err
			}
		}
	}

	// Sanitize HTML content
	if form.HTML != nil {
		content := html.Sanitize(*form.HTML)
		form.HTML = &content
	}

	reg.logger.Debug().Uint(
		"uid", uid,
	).Str("title", form.TruncatedTitle()).Msg("creating article...")
	article, err := reg.db.CreateArticleForUser(uid, form)
	if err != nil {
		reg.logger.Info().Err(err).Uint(
			"uid", uid,
		).Str("title", form.TruncatedTitle()).Msg("unable to create article")
		return nil, err
	}
	reg.logger.Info().Uint(
		"uid", uid,
	).Str("title", form.TruncatedTitle()).Uint("id", article.ID).Msg("article created")
	event.Emit(event.CreateArticle, *article)
	return article, nil
}

// CreateArticles creates new articles
func (reg *Registry) CreateArticles(ctx context.Context, data []model.ArticleCreateForm) *model.CreatedArticlesResponse {
	result := model.CreatedArticlesResponse{}
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
func (reg *Registry) hydrateArticle(ctx context.Context, article *model.ArticleCreateForm) error {
	page, err := reg.webScraper.Scrap(ctx, *article.URL)
	if page == nil {
		return err
	}
	article.URL = &page.URL
	if article.Title == "" {
		article.Title = page.Title
	}
	if article.HTML == nil {
		article.HTML = &page.HTML
	}
	if article.Text == nil {
		article.Text = &page.Text
	}
	if article.Image == nil {
		article.Image = &page.Image
	}

	return err
}
