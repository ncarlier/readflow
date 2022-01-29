package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/ncarlier/readflow/pkg/model"
)

// GetIncomingWebhookByToken returns an incoming webhook by its token
func (reg *Registry) GetIncomingWebhookByToken(token string) (*model.IncomingWebhook, error) {
	return reg.db.GetIncomingWebhookByToken(token)
}

// GetIncomingWebhooks get incoming webhook from current user
func (reg *Registry) GetIncomingWebhooks(ctx context.Context) (*[]model.IncomingWebhook, error) {
	uid := getCurrentUserIDFromContext(ctx)

	result, err := reg.db.GetIncomingWebhooksByUser(uid)
	if err != nil {
		reg.logger.Info().Err(err).Uint(
			"uid", uid,
		).Msg("unable to get incoming webhooks")
		return nil, err
	}

	return &result, err
}

// GetIncomingWebhook get an incoming webhook of the current user
func (reg *Registry) GetIncomingWebhook(ctx context.Context, id uint) (*model.IncomingWebhook, error) {
	uid := getCurrentUserIDFromContext(ctx)

	result, err := reg.db.GetIncomingWebhookByID(id)
	if err != nil || result == nil || result.UserID != uid {
		if err == nil {
			err = ErrIncomingWebhookNotFound
		}
		return nil, err
	}
	return result, nil
}

// CreateIncomingWebhook create an incoming webhook for current user
func (reg *Registry) CreateIncomingWebhook(ctx context.Context, form model.IncomingWebhookCreateForm) (*model.IncomingWebhook, error) {
	uid := getCurrentUserIDFromContext(ctx)

	// Validate user quota
	plan, err := reg.GetCurrentUserPlan(ctx)
	if err != nil {
		return nil, err
	}
	if plan != nil {
		totalWebhooks, err := reg.db.CountIncomingWebhooksByUser(uid)
		if err != nil {
			reg.logger.Info().Err(err).Uint(
				"uid", uid,
			).Msg("unable to create incoming webhook")
			return nil, err
		}
		if totalWebhooks >= plan.TotalWebhooks {
			err = ErrUserQuotaReached
			reg.logger.Info().Err(err).Uint(
				"uid", uid,
			).Uint(
				"total", plan.TotalCategories,
			).Msg("unable to create incoming webhook")
			return nil, err
		}
	}

	result, err := reg.db.CreateIncomingWebhookForUser(uid, form)
	if err != nil {
		reg.logger.Info().Err(err).Uint(
			"uid", uid,
		).Str("alias", form.Alias).Msg("unable to create incoming webhook")
		return nil, err
	}
	return result, err
}

// UpdateIncomingWebhook update an incoming webhook for current user
func (reg *Registry) UpdateIncomingWebhook(ctx context.Context, form model.IncomingWebhookUpdateForm) (*model.IncomingWebhook, error) {
	uid := getCurrentUserIDFromContext(ctx)

	result, err := reg.db.UpdateIncomingWebhookForUser(uid, form)
	if err != nil {
		reg.logger.Info().Err(err).Uint(
			"uid", uid,
		).Uint(
			"id", form.ID,
		).Msg("unable to update incoming webhook")
		return nil, err
	}
	return result, err
}

// DeleteIncomingWebhook delete an incoming webhook of the current user
func (reg *Registry) DeleteIncomingWebhook(ctx context.Context, id uint) (*model.IncomingWebhook, error) {
	uid := getCurrentUserIDFromContext(ctx)

	result, err := reg.GetIncomingWebhook(ctx, id)
	if err != nil {
		return nil, err
	}

	err = reg.db.DeleteIncomingWebhookByUser(uid, id)
	if err != nil {
		reg.logger.Info().Err(err).Uint(
			"uid", uid,
		).Uint("id", id).Msg("unable to delete incoming webhook")
		return nil, err
	}
	return result, nil
}

// DeleteIncomingWebhooks delete incoming webhooks of the current user
func (reg *Registry) DeleteIncomingWebhooks(ctx context.Context, ids []uint) (int64, error) {
	uid := getCurrentUserIDFromContext(ctx)
	idsStr := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(ids)), ","), "[]")

	nb, err := reg.db.DeleteIncomingWebhooksByUser(uid, ids)
	if err != nil {
		reg.logger.Info().Err(err).Uint(
			"uid", uid,
		).Str("ids", idsStr).Msg("unable to delete incoming webhooks")
		return 0, err
	}
	reg.logger.Debug().Err(err).Uint(
		"uid", uid,
	).Str("ids", idsStr).Int64("nb", nb).Msg("incoming webhooks deleted")
	return nb, nil
}
