package service

import (
	"context"
	"strings"

	"github.com/ncarlier/readflow/pkg/event"
	"github.com/ncarlier/readflow/pkg/helper"
	"github.com/ncarlier/readflow/pkg/model"
	"github.com/rs/zerolog"
)

const unableToUpdateArticleErrorMsg = "unable to update article"

// UpdateArticle update article
func (reg *Registry) UpdateArticle(ctx context.Context, form model.ArticleUpdateForm) (*model.Article, error) {
	uid := getCurrentUserIDFromContext(ctx)

	logger := reg.logger.Info().Uint("uid", uid).Uint("id", form.ID)
	debug := reg.logger.Debug().Uint("uid", uid).Uint("id", form.ID)

	// Validate update form
	if err := validateArticleUpdateForm(form); err != nil {
		logger.Err(err).Msg(unableToUpdateArticleErrorMsg)
		return nil, err
	}

	// Validate category_id if set
	// TODO manage category deafectation
	if form.CategoryID != nil {
		if _, err := reg.GetCategory(ctx, *form.CategoryID); err != nil {
			logger.Err(err).Msg(unableToUpdateArticleErrorMsg)
			return nil, err
		}
	}

	// Update loggers values
	addLoggerContextForUpdateArticle(logger, form)
	addLoggerContextForUpdateArticle(debug, form)

	// Update article
	debug.Msg("updating article...")
	article, err := reg.db.UpdateArticleForUser(uid, form)
	if err != nil {
		logger.Err(err).Msg(unableToUpdateArticleErrorMsg)
		return nil, err
	}

	logger.Msg("article updated")

	// Emit update event
	event.Emit(event.UpdateArticle, *article)

	return article, nil
}

func addLoggerContextForUpdateArticle(logger *zerolog.Event, form model.ArticleUpdateForm) {
	if form.CategoryID != nil {
		logger.Uint("category_id", *form.CategoryID)
	}
	if form.Stars != nil {
		logger.Uint("stars", *form.Stars)
	}
	if form.Status != nil {
		logger.Str("status", *form.Status)
	}
	if form.Title != nil {
		logger.Str("title", helper.Truncate(*form.Title, 24))
	}
	if form.Text != nil {
		logger.Str("text", helper.Truncate(*form.Text, 24))
	}
}

func validateArticleUpdateForm(form model.ArticleUpdateForm) error {
	validations := new(helper.FieldsValidator)
	validations.Validate("stars", func() bool {
		return form.Stars == nil || (*form.Stars >= 0 && *form.Stars <= 5)
	})
	validations.Validate("title", func() bool {
		if form.Title != nil {
			l := len(strings.TrimSpace(*form.Title))
			return l >= 0 && l <= 128
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
