package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ncarlier/readflow/internal/global"
	"github.com/ncarlier/readflow/internal/integration/webhook"
	"github.com/ncarlier/readflow/internal/model"

	"github.com/ncarlier/readflow/pkg/defaults"

	// import all outgoing webhook providers
	_ "github.com/ncarlier/readflow/internal/integration/webhook/all"
)

const unableToCreateOutgoingWebhookErrorMsg = "unable to create outgoing webhook"

// GetOutgoingWebhooks get outgoing webhooks from current user
func (reg *Registry) GetOutgoingWebhooks(ctx context.Context) (*[]model.OutgoingWebhook, error) {
	uid := getCurrentUserIDFromContext(ctx)

	webhooks, err := reg.db.GetOutgoingWebhooksByUser(uid)
	if err != nil {
		reg.logger.Info().Err(err).Uint(
			"uid", uid,
		).Msg("unable to get outgoing webhooks")
		return nil, err
	}

	return &webhooks, err
}

// GetOutgoingWebhook get an outgoing webhook of the current user
func (reg *Registry) GetOutgoingWebhook(ctx context.Context, id uint) (*model.OutgoingWebhook, error) {
	uid := getCurrentUserIDFromContext(ctx)

	whk, err := reg.db.GetOutgoingWebhookByID(id)
	if err != nil || whk == nil || *whk.UserID != uid {
		if err == nil {
			err = ErrOutgoingWebhookNotFound
		}
		return nil, err
	}
	return whk, nil
}

// CreateOutgoingWebhook create an outgoing webhook for current user
func (reg *Registry) CreateOutgoingWebhook(ctx context.Context, form model.OutgoingWebhookCreateForm) (*model.OutgoingWebhook, error) {
	uid := getCurrentUserIDFromContext(ctx)

	logger := reg.logger.With().Uint("uid", uid).Str("alias", form.Alias).Logger()

	// Validate user quota
	plan, err := reg.GetCurrentUserPlan(ctx)
	if err != nil {
		return nil, err
	}
	if plan != nil {
		totalWebhooks, err := reg.db.CountOutgoingWebhooksByUser(uid)
		if err != nil {
			logger.Info().Err(err).Msg(unableToCreateOutgoingWebhookErrorMsg)
			return nil, err
		}
		if totalWebhooks >= plan.OutgoingWebhooksLimit {
			err = ErrUserQuotaReached
			logger.Info().Err(err).Uint("total", totalWebhooks).Msg(unableToCreateOutgoingWebhookErrorMsg)
			return nil, err
		}
	}

	// Seal secrets
	if reg.secretsEngineProvider != nil {
		err := reg.secretsEngineProvider.Seal(&form.Secrets)
		if err != nil {
			return nil, err
		}
	}

	// Validate outgoing webhook configuration
	dummy := model.OutgoingWebhook{
		Provider: form.Provider,
		Config:   form.Config,
		Secrets:  form.Secrets,
	}
	_, err = webhook.NewOutgoingWebhookProvider(dummy, reg.conf)
	if err != nil {
		logger.Info().Err(err).Msg(unableToCreateOutgoingWebhookErrorMsg)
		return nil, err
	}

	logger.Debug().Msg("creating outgoing webhook...")
	result, err := reg.db.CreateOutgoingWebhookForUser(uid, form)
	if err != nil {
		logger.Info().Err(err).Msg(unableToCreateOutgoingWebhookErrorMsg)
		return nil, err
	}
	logger.Info().Uint("id", *result.ID).Msg("outgoing webhook created")
	return result, err
}

// UpdateOutgoingWebhook update an outgoing webhook for current user
func (reg *Registry) UpdateOutgoingWebhook(ctx context.Context, form model.OutgoingWebhookUpdateForm) (*model.OutgoingWebhook, error) {
	uid := getCurrentUserIDFromContext(ctx)

	logger := reg.logger.With().Uint("uid", uid).Uint("id", form.ID).Logger()

	if form.Provider != nil && form.Config != nil && form.Secrets != nil {
		// Validate outgoing webhook configuration
		dummy := model.OutgoingWebhook{
			Provider: *form.Provider,
			Config:   *form.Config,
			Secrets:  *form.Secrets,
		}
		_, err := webhook.NewOutgoingWebhookProvider(dummy, reg.conf)
		if err != nil {
			logger.Info().Err(err).Msg("unable to configure outgoing webhook")
			return nil, err
		}
	} else {
		// Provider can only be modify with its configuration
		form.Provider = nil
		form.Config = nil
	}

	if form.Secrets != nil && reg.secretsEngineProvider != nil {
		err := reg.secretsEngineProvider.Seal(form.Secrets)
		if err != nil {
			return nil, err
		}
	}

	logger.Debug().Msg("updating outgoing webhook...")
	result, err := reg.db.UpdateOutgoingWebhookForUser(uid, form)
	if err != nil {
		logger.Info().Err(err).Msg("unable to update outgoing webhook")
		return nil, err
	}
	logger.Info().Msg("outgoing webhook updated")
	return result, err
}

// DeleteOutgoingWebhook delete an outgoing webhook of the current user
func (reg *Registry) DeleteOutgoingWebhook(ctx context.Context, id uint) (*model.OutgoingWebhook, error) {
	uid := getCurrentUserIDFromContext(ctx)

	logger := reg.logger.With().Uint("uid", uid).Uint("id", id).Logger()

	whk, err := reg.GetOutgoingWebhook(ctx, id)
	if err != nil {
		return nil, err
	}

	logger.Debug().Msg("deleting outgoing webhook...")
	err = reg.db.DeleteOutgoingWebhookByUser(uid, id)
	if err != nil {
		logger.Info().Err(err).Msg("unable to delete outgoing webhook")
		return nil, err
	}
	logger.Info().Msg("outgoing webhook deleted")
	return whk, nil
}

// DeleteOutgoingWebhooks delete outgoing webhooks of the current user
func (reg *Registry) DeleteOutgoingWebhooks(ctx context.Context, ids []uint) (int64, error) {
	uid := getCurrentUserIDFromContext(ctx)
	idsStr := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(ids)), ","), "[]")

	logger := reg.logger.With().Uint("uid", uid).Str("ids", idsStr).Logger()

	logger.Debug().Msg("deleting outgoing webhooks...")
	nb, err := reg.db.DeleteOutgoingWebhooksByUser(uid, ids)
	if err != nil {
		logger.Info().Err(err).Msg("unable to delete outgoing webhooks")
		return 0, err
	}
	logger.Info().Err(err).Int64("nb", nb).Msg("outgoing webhooks deleted")
	return nb, nil
}

// SendArticle send an article to outgoing webhook
func (reg *Registry) SendArticle(ctx context.Context, idArticle uint, alias *string) (*webhook.Result, error) {
	start := time.Now()
	user, err := reg.GetCurrentUser(ctx)
	if err != nil {
		return nil, err
	}

	logger := reg.logger.With().Uint("uid", *user.ID).Uint("article", idArticle).Logger()

	article, err := reg.db.GetArticleByID(idArticle)
	if err != nil || article == nil || article.UserID != *user.ID {
		if err == nil {
			err = errors.New("article not found")
		}
		logger.Info().Err(err).Msg(ErrOutgoingWebhookSend.Error())
		return nil, err
	}

	if alias != nil {
		logger = logger.With().Str("alias", *alias).Logger()
	}

	webhookConf, err := reg.db.GetOutgoingWebhookByUserAndAlias(*user.ID, alias)
	if err != nil || webhookConf == nil {
		if err == nil {
			err = errors.New("outgoing webhook not found")
		}
		logger.Info().Err(err).Msg(ErrOutgoingWebhookSend.Error())
		return nil, err
	}

	// UnSeal secrets
	if reg.secretsEngineProvider != nil {
		err := reg.secretsEngineProvider.UnSeal(&webhookConf.Secrets)
		if err != nil {
			logger.Info().Err(err).Msg(ErrOutgoingWebhookSend.Error())
			return nil, err
		}
	}

	provider, err := webhook.NewOutgoingWebhookProvider(*webhookConf, reg.conf)
	if err != nil {
		logger.Info().Err(err).Msg(ErrOutgoingWebhookSend.Error())
		return nil, err
	}

	// HACK: put downloader inside the context
	// This is needed by some providers (S3 for instance)
	webhookContext := context.WithValue(ctx, global.ContextDownloader, reg.dl)
	// Add user to the context
	if value := ctx.Value(global.ContextUser); value == nil {
		webhookContext = context.WithValue(webhookContext, global.ContextUser, user)
	}
	if _, isDeadline := ctx.Deadline(); !isDeadline {
		// no outgoing webhook deadline defined... setting one
		timeout := defaults.Timeout
		if plan, err := reg.GetCurrentUserPlan(ctx); err == nil && plan != nil {
			timeout = plan.OutgoingWebhooksTimeout.Duration
		}
		var cancel context.CancelFunc
		webhookContext, cancel = context.WithTimeout(webhookContext, timeout)
		defer cancel()
	}

	logger.Debug().Msg("sending article to outgoing webhook...")
	result, err := provider.Send(webhookContext, *article)
	if err != nil {
		logger.Info().Err(err).Msg(ErrOutgoingWebhookSend.Error())
		return nil, err
	}
	logger.Info().Dur("took", time.Since(start)).Msg("article sent to outgoing webhook")
	return result, nil
}
