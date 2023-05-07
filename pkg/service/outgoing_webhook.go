package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/ncarlier/readflow/pkg/constant"
	"github.com/ncarlier/readflow/pkg/integration/webhook"
	webhookProvider "github.com/ncarlier/readflow/pkg/integration/webhook"

	// import all outgoing webhook providers
	_ "github.com/ncarlier/readflow/pkg/integration/webhook/all"
	"github.com/ncarlier/readflow/pkg/model"
)

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

	webhook, err := reg.db.GetOutgoingWebhookByID(id)
	if err != nil || webhook == nil || *webhook.UserID != uid {
		if err == nil {
			err = ErrOutgoingWebhookNotFound
		}
		return nil, err
	}
	return webhook, nil
}

// CreateOutgoingWebhook create an outgoing webhook for current user
func (reg *Registry) CreateOutgoingWebhook(ctx context.Context, form model.OutgoingWebhookCreateForm) (*model.OutgoingWebhook, error) {
	uid := getCurrentUserIDFromContext(ctx)

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
	_, err := webhookProvider.NewOutgoingWebhookProvider(dummy, reg.conf)
	if err != nil {
		reg.logger.Info().Err(err).Uint(
			"uid", uid,
		).Msg("unable to configure outgoing webhook")
		return nil, err
	}

	result, err := reg.db.CreateOutgoingWebhookForUser(uid, form)
	if err != nil {
		reg.logger.Info().Err(err).Uint(
			"uid", uid,
		).Str("alias", form.Alias).Msg("unable to create outgoing webhook")
		return nil, err
	}
	return result, err
}

// UpdateOutgoingWebhook update an outgoing webhook for current user
func (reg *Registry) UpdateOutgoingWebhook(ctx context.Context, form model.OutgoingWebhookUpdateForm) (*model.OutgoingWebhook, error) {
	uid := getCurrentUserIDFromContext(ctx)

	if form.Provider != nil && form.Config != nil && form.Secrets != nil {
		// Validate outgoing webhook configuration
		dummy := model.OutgoingWebhook{
			Provider: *form.Provider,
			Config:   *form.Config,
			Secrets:  *form.Secrets,
		}
		_, err := webhookProvider.NewOutgoingWebhookProvider(dummy, reg.conf)
		if err != nil {
			reg.logger.Info().Err(err).Uint(
				"uid", uid,
			).Msg("unable to configure outgoing webhook")
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

	result, err := reg.db.UpdateOutgoingWebhookForUser(uid, form)
	if err != nil {
		reg.logger.Info().Err(err).Uint(
			"uid", uid,
		).Uint(
			"id", form.ID,
		).Msg("unable to update outgoing webhook")
		return nil, err
	}
	return result, err
}

// DeleteOutgoingWebhook delete an outgoing webhook of the current user
func (reg *Registry) DeleteOutgoingWebhook(ctx context.Context, id uint) (*model.OutgoingWebhook, error) {
	uid := getCurrentUserIDFromContext(ctx)

	webhook, err := reg.GetOutgoingWebhook(ctx, id)
	if err != nil {
		return nil, err
	}

	err = reg.db.DeleteOutgoingWebhookByUser(uid, id)
	if err != nil {
		reg.logger.Info().Err(err).Uint(
			"uid", uid,
		).Uint("id", id).Msg("unable to delete outgoing webhook")
		return nil, err
	}
	return webhook, nil
}

// DeleteOutgoingWebhooks delete outgoing webhooks of the current user
func (reg *Registry) DeleteOutgoingWebhooks(ctx context.Context, ids []uint) (int64, error) {
	uid := getCurrentUserIDFromContext(ctx)
	idsStr := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(ids)), ","), "[]")

	nb, err := reg.db.DeleteOutgoingWebhooksByUser(uid, ids)
	if err != nil {
		reg.logger.Info().Err(err).Uint(
			"uid", uid,
		).Str("ids", idsStr).Msg("unable to delete outgoing webhooks")
		return 0, err
	}
	reg.logger.Debug().Err(err).Uint(
		"uid", uid,
	).Str("ids", idsStr).Int64("nb", nb).Msg("outgoing webhooks deleted")
	return nb, nil
}

// SendArticle send an article to outgoing webhook
func (reg *Registry) SendArticle(ctx context.Context, idArticle uint, alias *string) (*webhook.Result, error) {
	user, err := reg.GetCurrentUser(ctx)
	if err != nil {
		return nil, err
	}

	logger := reg.logger.With().Uint(
		"uid", *user.ID,
	).Uint("article", idArticle).Logger()

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

	provider, err := webhookProvider.NewOutgoingWebhookProvider(*webhookConf, reg.conf)
	if err != nil {
		logger.Info().Err(err).Msg(ErrOutgoingWebhookSend.Error())
		return nil, err
	}

	logger.Debug().Msg("sending article...")
	// HACK: put downloader inside the context
	// This is needed by some providers (S3 for instance)
	webhookContext := context.WithValue(ctx, constant.ContextDownloader, reg.dl)
	// Add user to the context
	if value := ctx.Value(constant.ContextUser); value == nil {
		webhookContext = context.WithValue(webhookContext, constant.ContextUser, user)
	}
	if _, isDeadline := ctx.Deadline(); !isDeadline {
		// no outgoing webhook deadline defined... setting one
		var cancel context.CancelFunc
		webhookContext, cancel = context.WithTimeout(webhookContext, 3*constant.DefaultTimeout)
		defer cancel()
	}

	result, err := provider.Send(webhookContext, *article)
	if err != nil {
		logger.Info().Err(err).Msg(ErrOutgoingWebhookSend.Error())
		return nil, err
	}
	logger.Info().Msg("article sent to outgoing webhook")
	return result, nil
}
