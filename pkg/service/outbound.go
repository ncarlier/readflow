package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/ncarlier/readflow/pkg/integration/outbound"
	"github.com/ncarlier/readflow/pkg/model"
)

// GetOutboundServices get outbound services from current user
func (reg *Registry) GetOutboundServices(ctx context.Context) (*[]model.OutboundService, error) {
	uid := getCurrentUserFromContext(ctx)

	outboundServices, err := reg.db.GetOutboundServicesByUser(uid)
	if err != nil {
		reg.logger.Info().Err(err).Uint(
			"uid", uid,
		).Msg("unable to get outbound services")
		return nil, err
	}

	return &outboundServices, err
}

// GetOutboundService get an outbound service of the current user
func (reg *Registry) GetOutboundService(ctx context.Context, id uint) (*model.OutboundService, error) {
	uid := getCurrentUserFromContext(ctx)

	outboundService, err := reg.db.GetOutboundServiceByID(id)
	if err != nil || outboundService == nil || *outboundService.UserID != uid {
		if err == nil {
			err = ErrOutboundServiceNotFound
		}
		return nil, err
	}
	return outboundService, nil
}

// CreateOutboundService create an outbound service for current user
func (reg *Registry) CreateOutboundService(ctx context.Context, form model.OutboundServiceCreateForm) (*model.OutboundService, error) {
	uid := getCurrentUserFromContext(ctx)

	// Validate outbound service configuration
	dummy := model.OutboundService{
		Provider: form.Provider,
		Config:   form.Config,
	}
	_, err := outbound.NewOutboundServiceProvider(dummy)
	if err != nil {
		reg.logger.Info().Err(err).Uint(
			"uid", uid,
		).Msg("unable to configure outbound service")
		return nil, err
	}

	result, err := reg.db.CreateOutboundServiceForUser(uid, form)
	if err != nil {
		reg.logger.Info().Err(err).Uint(
			"uid", uid,
		).Str("alias", form.Alias).Msg("unable to create outbound service")
		return nil, err
	}
	return result, err
}

// UpdateOutboundService update an outbound service for current user
func (reg *Registry) UpdateOutboundService(ctx context.Context, form model.OutboundServiceUpdateForm) (*model.OutboundService, error) {
	uid := getCurrentUserFromContext(ctx)

	if form.Provider != nil && form.Config != nil {
		// Validate outbound service configuration
		dummy := model.OutboundService{
			Provider: *form.Provider,
			Config:   *form.Config,
		}
		_, err := outbound.NewOutboundServiceProvider(dummy)
		if err != nil {
			reg.logger.Info().Err(err).Uint(
				"uid", uid,
			).Msg("unable to configure outbound service")
			return nil, err
		}
	} else {
		// Provider can only be modify with its configuration
		form.Provider = nil
		form.Config = nil
	}

	result, err := reg.db.UpdateOutboundServiceForUser(uid, form)
	if err != nil {
		reg.logger.Info().Err(err).Uint(
			"uid", uid,
		).Uint(
			"id", form.ID,
		).Msg("unable to update outbound service")
		return nil, err
	}
	return result, err
}

// DeleteOutboundService delete an outbound service of the current user
func (reg *Registry) DeleteOutboundService(ctx context.Context, id uint) (*model.OutboundService, error) {
	uid := getCurrentUserFromContext(ctx)

	outboundService, err := reg.GetOutboundService(ctx, id)
	if err != nil {
		return nil, err
	}

	err = reg.db.DeleteOutboundServiceByUser(uid, id)
	if err != nil {
		reg.logger.Info().Err(err).Uint(
			"uid", uid,
		).Uint("id", id).Msg("unable to delete outbound service")
		return nil, err
	}
	return outboundService, nil
}

// DeleteOutboundServices delete outbound services of the current user
func (reg *Registry) DeleteOutboundServices(ctx context.Context, ids []uint) (int64, error) {
	uid := getCurrentUserFromContext(ctx)
	idsStr := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(ids)), ","), "[]")

	nb, err := reg.db.DeleteOutboundServicesByUser(uid, ids)
	if err != nil {
		reg.logger.Info().Err(err).Uint(
			"uid", uid,
		).Str("ids", idsStr).Msg("unable to delete outbound services")
		return 0, err
	}
	reg.logger.Debug().Err(err).Uint(
		"uid", uid,
	).Str("ids", idsStr).Int64("nb", nb).Msg("outbound services deleted")
	return nb, nil
}

// ArchiveArticle archive an article using an outbound service
func (reg *Registry) ArchiveArticle(ctx context.Context, idArticle uint, alias *string) error {
	uid := getCurrentUserFromContext(ctx)

	logger := reg.logger.With().Uint(
		"uid", uid,
	).Uint("article", idArticle).Logger()

	article, err := reg.db.GetArticleByID(idArticle)
	if err != nil || article == nil || article.UserID != uid {
		if err == nil {
			err = errors.New("article not found")
		}
		logger.Info().Err(err).Msg(ErrOutboundServiceSend.Error())
		return err
	}

	if alias != nil {
		logger = logger.With().Str("alias", *alias).Logger()
	}

	outboundServiceConf, err := reg.db.GetOutboundServiceByUserAndAlias(uid, alias)
	if err != nil || outboundServiceConf == nil {
		if err == nil {
			err = errors.New("outbound service not found")
		}
		logger.Info().Err(err).Msg(ErrOutboundServiceSend.Error())
		return err
	}

	provider, err := outbound.NewOutboundServiceProvider(*outboundServiceConf)
	if err != nil {
		logger.Info().Err(err).Msg(ErrOutboundServiceSend.Error())
		return err
	}

	logger.Debug().Msg("sending article...")
	err = provider.Send(ctx, *article)
	if err != nil {
		logger.Info().Err(err).Msg(ErrOutboundServiceSend.Error())
		return err
	}
	logger.Info().Msg("article sent to outbound service")
	return nil
}
