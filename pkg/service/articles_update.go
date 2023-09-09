package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/ncarlier/readflow/pkg/event"
	"github.com/ncarlier/readflow/pkg/helper"
	"github.com/ncarlier/readflow/pkg/model"
	"github.com/rs/zerolog"
)

const unableToUpdateArticleErrorMsg = "unable to update article"

// UpdateArticle update article
func (reg *Registry) UpdateArticle(ctx context.Context, form model.ArticleUpdateForm) (*model.Article, error) {
	uid := getCurrentUserIDFromContext(ctx)
	start := time.Now()

	logger := reg.logger.With().Uint("uid", uid).Uint("id", form.ID).Logger()

	// Validate update form
	if err := validateArticleUpdateForm(form); err != nil {
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
	logger = addLoggerContextForUpdateArticle(logger, form)

	// Update article
	logger.Debug().Msg("updating article...")
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

func addLoggerContextForUpdateArticle(logger zerolog.Logger, form model.ArticleUpdateForm) zerolog.Logger {
	ctx := logger.With()
	if form.CategoryID != nil {
		ctx = ctx.Uint("category_id", *form.CategoryID)
	}
	if form.Stars != nil {
		ctx = ctx.Uint("stars", *form.Stars)
	}
	if form.Status != nil {
		ctx = ctx.Str("status", *form.Status)
	}
	if form.Title != nil {
		ctx = ctx.Str("title", helper.Truncate(*form.Title, 24))
	}
	if form.Text != nil {
		ctx = ctx.Str("text", helper.Truncate(*form.Text, 24))
	}
	return ctx.Logger()
}

func validateArticleUpdateForm(form model.ArticleUpdateForm) error {
	validations := new(helper.FieldsValidator)
	validations.Validate("stars", func() bool {
		return form.Stars == nil || *form.Stars <= 5
	})
	validations.Validate("title", func() bool {
		if form.Title != nil {
			l := len(strings.TrimSpace(*form.Title))
			return l > 0 && l <= 256
		}
		return true
	})
	validations.Validate("text", func() bool {
		if form.Text != nil {
			l := len(strings.TrimSpace(*form.Text))
			return l >= 0 && l <= 512
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
