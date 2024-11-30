package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/ncarlier/readflow/internal/model"
	"github.com/ncarlier/readflow/pkg/event"
	"github.com/ncarlier/readflow/pkg/utils"
	"github.com/ncarlier/readflow/pkg/validator"
	"github.com/rs/zerolog"
)

const unableToUpdateArticleErrorMsg = "unable to update article"

// UpdateArticle update article
func (reg *Registry) UpdateArticle(ctx context.Context, form model.ArticleUpdateForm) (*model.Article, error) {
	start := time.Now()
	uid := getCurrentUserIDFromContext(ctx)

	logger := reg.logger.With().Uint("uid", uid).Uint("id", form.ID).Logger()

	// Validate update form
	if err := validateArticleUpdateForm(&form); err != nil {
		logger.Info().Err(err).Msg(unableToUpdateArticleErrorMsg)
		return nil, err
	}

	// Validate category_id if set
	// TODO manage category deafectation
	if form.CategoryID != nil {
		if _, err := reg.GetCategory(ctx, *form.CategoryID); err != nil {
			logger.Info().Err(err).Msg(unableToUpdateArticleErrorMsg)
			return nil, err
		}
	}

	// Update loggers values
	updateLoggerContextWithUpdateForm(&logger, &form)
	logger.Debug().Msg("updating article...")

	// Refresh HTML content if asked
	if form.RefreshContent != nil && *form.RefreshContent {
		if err := reg.refreshArticleContent(ctx, &form); err != nil {
			return nil, err
		}
	}

	// HTML updated? Then clean it!
	if form.HTML != nil {
		html := reg.sanitizer.Sanitize(*form.HTML)
		form.HTML = &html
	}

	// Update article
	article, err := reg.db.UpdateArticleForUser(uid, form)
	if err != nil {
		logger.Err(err).Msg(unableToUpdateArticleErrorMsg)
		return nil, err
	}
	if article == nil {
		err := errors.New("article not found")
		logger.Err(err).Msg(unableToUpdateArticleErrorMsg)
		return nil, err
	}

	logger.Info().Dur("took", time.Since(start)).Msg("article updated")

	// Emit update event
	reg.events.Publish(event.NewEvent(EventUpdateArticle, *article))

	return article, nil
}

func (reg *Registry) refreshArticleContent(ctx context.Context, form *model.ArticleUpdateForm) error {
	article, err := reg.GetArticle(ctx, form.ID)
	if err != nil {
		return err
	}

	if article.URL == nil || *article.URL == "" {
		// nothing to refresh
		return nil
	}

	logger := reg.logger.With().Uint("uid", article.UserID).Uint("id", form.ID).Logger()
	logger.Debug().Msg("refreshing article content...")
	page, err := reg.webScraper.Scrap(ctx, *article.URL)
	if err != nil {
		return err
	}
	if form.Title == nil || *form.Title == "" {
		form.Title = &page.Title
	}
	if form.Text == nil || *form.Text == "" {
		form.Text = &page.Text
	}
	form.HTML = &page.HTML
	form.Image = &page.Image
	return nil
}

func updateLoggerContextWithUpdateForm(logger *zerolog.Logger, form *model.ArticleUpdateForm) {
	logger.UpdateContext(func(ctx zerolog.Context) zerolog.Context {
		if form.CategoryID != nil {
			ctx = ctx.Uint("category_id", *form.CategoryID)
		}
		if form.Stars != nil {
			ctx = ctx.Int("stars", *form.Stars)
		}
		if form.Status != nil {
			ctx = ctx.Str("status", *form.Status)
		}
		if form.Title != nil {
			ctx = ctx.Str("title", utils.Truncate(*form.Title, 24))
		}
		if form.Text != nil {
			ctx = ctx.Str("text", utils.Truncate(*form.Text, 24))
		}
		if form.RefreshContent != nil {
			ctx = ctx.Bool("refresh", *form.RefreshContent)
		}
		return ctx
	})
}

func validateArticleUpdateForm(form *model.ArticleUpdateForm) error {
	validations := new(validator.FieldsValidator)
	validations.Validate("stars", func() bool {
		return form.Stars == nil || *form.Stars <= 5
	})
	validations.Validate("title", func() bool {
		if form.Title != nil {
			l := len(strings.TrimSpace(*form.Title))
			return l <= 256
		}
		return true
	})
	validations.Validate("text", func() bool {
		if form.Text != nil {
			l := len(strings.TrimSpace(*form.Text))
			return l <= 512
		}
		return true
	})
	validations.Validate("status", func() bool {
		if form.Status != nil {
			status := *form.Status
			return status == "inbox" || status == "read" || status == "to_read"
		}
		return true
	})
	return validations.Error()
}
