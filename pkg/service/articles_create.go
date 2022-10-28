package service

import (
	"context"

	"github.com/ncarlier/readflow/pkg/event"
	"github.com/ncarlier/readflow/pkg/model"

	// activate all content providers
	_ "github.com/ncarlier/readflow/pkg/scraper/content-provider/all"
)

const unableToCreateArticleErrorMsg = "unable to create article"

// ArticleCreationOptions article creation options
type ArticleCreationOptions struct {
	IgnoreHydrateError bool
}

// CreateArticle creates new article
func (reg *Registry) CreateArticle(ctx context.Context, form model.ArticleCreateForm, opts ArticleCreationOptions) (*model.Article, error) {
	uid := getCurrentUserIDFromContext(ctx)
	logger := reg.logger.Info().Uint("uid", uid).Str("title", form.TruncatedTitle())
	debug := reg.logger.Debug().Uint("uid", uid).Str("title", form.TruncatedTitle())

	plan, err := reg.GetCurrentUserPlan(ctx)
	if err != nil {
		return nil, err
	}

	if plan != nil {
		// Check user quota
		req := model.ArticlesPageRequest{}
		totalArticles, err := reg.CountCurrentUserArticles(ctx, req)
		if err != nil {
			logger.Err(err).Msg(unableToCreateArticleErrorMsg)
			return nil, err
		}
		if totalArticles >= plan.TotalArticles {
			err = ErrUserQuotaReached
			debug.Err(err).Uint(
				"total", totalArticles,
			).Msg(unableToCreateArticleErrorMsg)
			return nil, err
		}
	}

	// TODO validate article!

	// Get category if specified
	var category *model.Category
	if form.CategoryID != nil {
		cat, err := reg.GetCategory(ctx, *form.CategoryID)
		if err != nil {
			logger.Err(err).Msg(unableToCreateArticleErrorMsg)
			return nil, err
		}
		category = cat
	}

	if category == nil {
		// Process article by the rule engine
		if err := reg.ProcessArticleByRuleEngine(ctx, &form); err != nil {
			logger.Err(err).Msg(unableToCreateArticleErrorMsg)
			return nil, err
		}
	}

	if form.URL != nil && !form.IsComplete() {
		// Fetch original article in order to extract missing attributes
		if err := reg.scrapOriginalArticle(ctx, &form); err != nil {
			logger.Err(err).Msg("unable to fetch original article")
			// TODO excerpt and image should be extracted from HTML content
			if !opts.IgnoreHydrateError {
				return nil, err
			}
		}
		// update logger field
		logger = logger.Str("title", form.TruncatedTitle())
	}

	// Sanitize HTML content
	if form.HTML != nil {
		content := reg.sanitizer.Sanitize(*form.HTML)
		form.HTML = &content
	}

	debug.Msg("creating article...")
	article, err := reg.db.CreateArticleForUser(uid, form)
	if err != nil {
		logger.Err(err).Msg(unableToCreateArticleErrorMsg)
		return nil, err
	}
	logger.Uint("id", article.ID).Msg("article created")
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

// scrapOriginalArticle add missing attributes form original article
func (reg *Registry) scrapOriginalArticle(ctx context.Context, article *model.ArticleCreateForm) error {
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
