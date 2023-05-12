package service

import (
	"context"
	"fmt"

	"github.com/ncarlier/readflow/pkg/constant"
	"github.com/ncarlier/readflow/pkg/event"
	"github.com/ncarlier/readflow/pkg/model"
	"github.com/ncarlier/readflow/pkg/scripting"

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
	logger := reg.logger.With().Uint("uid", uid).Str("title", form.TruncatedTitle()).Logger()

	plan, err := reg.GetCurrentUserPlan(ctx)
	if err != nil {
		return nil, err
	}

	if plan != nil {
		// check user quota
		req := model.ArticlesPageRequest{}
		totalArticles, err := reg.CountCurrentUserArticles(ctx, req)
		if err != nil {
			logger.Info().Err(err).Msg(unableToCreateArticleErrorMsg)
			return nil, err
		}
		if totalArticles >= plan.ArticlesLimit {
			err = ErrUserQuotaReached
			logger.Debug().Err(err).Uint(
				"total", totalArticles,
			).Msg(unableToCreateArticleErrorMsg)
			return nil, err
		}
	}

	// TODO validate article!
	// validate category
	if form.CategoryID != nil {
		if _, err := reg.GetCategory(ctx, *form.CategoryID); err != nil {
			logger.Info().Err(err).Msg(unableToCreateArticleErrorMsg)
			return nil, err
		}
	}

	if form.URL != nil && !form.IsComplete() {
		// fetch original article in order to extract missing attributes
		if err := reg.scrapOriginalArticle(ctx, &form); err != nil {
			logger.Info().Err(err).Str("URL", *form.URL).Msg("unable to fetch original article")
			// TODO excerpt and image should be extracted from HTML content
			if !opts.IgnoreHydrateError {
				return nil, err
			}
		}
	}

	// sanitize HTML content
	if form.HTML != nil {
		content := reg.sanitizer.Sanitize(*form.HTML)
		form.HTML = &content
	}

	var ops scripting.OperationStack
	if webhook, ok := ctx.Value(constant.ContextIncomingWebhook).(*model.IncomingWebhook); ok {
		// process article by the script engine if comming from webhook
		if ops, err = reg.processArticleByScriptEngine(ctx, webhook, &form); err != nil {
			logger.Debug().Err(err).Msg("unable to process article by script engine")
			text := err.Error()
			if form.Text != nil {
				text = fmt.Sprintf("%s\n%s", text, *form.Text)
			}
			form.Text = &text
		}
	}
	// drop if asked
	if ops.Contains(scripting.OpDrop) {
		return nil, nil
	}
	// exec set operations
	reg.execSetOperations(ctx, ops, &form)

	logger.Debug().Msg("creating article...")
	// persist article
	article, err := reg.db.CreateArticleForUser(uid, form)
	if err != nil {
		logger.Error().Err(err).Msg(unableToCreateArticleErrorMsg)
		return nil, err
	}
	logger.Info().Uint("id", article.ID).Msg("article created")
	// exec asynchronously other operations
	go func() {
		if err := reg.execOtherOperations(ctx, ops, article); err != nil {
			logger.Info().Err(err).Msg(unableToCreateArticleErrorMsg)
		}
	}()
	// emit article creation event
	var evtOpts event.EventOption
	evtOpts.AddIf(NoNotificationEventOption, ops.Contains(scripting.OpDisableGlobalNotification))
	reg.events.Publish(event.NewEventWithOption(EventCreateArticle, *article, evtOpts))
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
	ctx, cancel := context.WithTimeout(ctx, constant.DefaultTimeout)
	defer cancel()
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
