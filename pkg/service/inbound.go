package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/ncarlier/readflow/pkg/integration/inbound"
	"github.com/ncarlier/readflow/pkg/model"
)

// GetInboundServiceByToken returns an inbound service by its token
func (reg *Registry) GetInboundServiceByToken(token string) (*model.InboundService, error) {
	return reg.db.GetInboundServiceByToken(token)
}

// GetInboundServices get inbound service from current user
func (reg *Registry) GetInboundServices(ctx context.Context) (*[]model.InboundService, error) {
	uid := getCurrentUserFromContext(ctx)

	result, err := reg.db.GetInboundServicesByUser(uid)
	if err != nil {
		reg.logger.Info().Err(err).Uint(
			"uid", uid,
		).Msg("unable to get inbound services")
		return nil, err
	}

	return &result, err
}

// GetInboundService get an inbound service of the current user
func (reg *Registry) GetInboundService(ctx context.Context, id uint) (*model.InboundService, error) {
	uid := getCurrentUserFromContext(ctx)

	result, err := reg.db.GetInboundServiceByID(id)
	if err != nil || result == nil || result.UserID != uid {
		if err == nil {
			err = ErrInboundServiceNotFound
		}
		return nil, err
	}
	return result, nil
}

// CreateInboundService create an inbound service for current user
func (reg *Registry) CreateInboundService(ctx context.Context, form model.InboundServiceCreateForm) (*model.InboundService, error) {
	uid := getCurrentUserFromContext(ctx)

	// Validate inbound service configuration
	dummy := model.InboundService{
		Provider: form.Provider,
		Config:   form.Config,
	}
	_, err := inbound.NewInboundServiceProvider(dummy)
	if err != nil {
		reg.logger.Info().Err(err).Uint(
			"uid", uid,
		).Msg("unable to configure inbound service")
		return nil, err
	}

	// Disable access token for indbound service with import function
	if inbound.Services[form.Provider].Type == inbound.Pull {
		form.Token = ""
	}

	result, err := reg.db.CreateInboundServiceForUser(uid, form)
	if err != nil {
		reg.logger.Info().Err(err).Uint(
			"uid", uid,
		).Str("alias", form.Alias).Msg("unable to create inbound service")
		return nil, err
	}
	return result, err
}

// UpdateInboundService update an inbound service for current user
func (reg *Registry) UpdateInboundService(ctx context.Context, form model.InboundServiceUpdateForm) (*model.InboundService, error) {
	uid := getCurrentUserFromContext(ctx)

	if form.Provider != nil && form.Config != nil {
		// Validate inbound service configuration
		dummy := model.InboundService{
			Provider: *form.Provider,
			Config:   *form.Config,
		}
		_, err := inbound.NewInboundServiceProvider(dummy)
		if err != nil {
			reg.logger.Info().Err(err).Uint(
				"uid", uid,
			).Msg("unable to configure inbound service")
			return nil, err
		}
	} else {
		// Provider can only be modify with its configuration
		form.Provider = nil
		form.Config = nil
	}

	result, err := reg.db.UpdateInboundServiceForUser(uid, form)
	if err != nil {
		reg.logger.Info().Err(err).Uint(
			"uid", uid,
		).Uint(
			"id", form.ID,
		).Msg("unable to update inbound service")
		return nil, err
	}
	return result, err
}

// DeleteInboundService delete an inbound service of the current user
func (reg *Registry) DeleteInboundService(ctx context.Context, id uint) (*model.InboundService, error) {
	uid := getCurrentUserFromContext(ctx)

	result, err := reg.GetInboundService(ctx, id)
	if err != nil {
		return nil, err
	}

	err = reg.db.DeleteInboundServiceByUser(uid, id)
	if err != nil {
		reg.logger.Info().Err(err).Uint(
			"uid", uid,
		).Uint("id", id).Msg("unable to delete API key")
		return nil, err
	}
	return result, nil
}

// DeleteInboundServices delete inbound services of the current user
func (reg *Registry) DeleteInboundServices(ctx context.Context, ids []uint) (int64, error) {
	uid := getCurrentUserFromContext(ctx)
	idsStr := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(ids)), ","), "[]")

	nb, err := reg.db.DeleteInboundServicesByUser(uid, ids)
	if err != nil {
		reg.logger.Info().Err(err).Uint(
			"uid", uid,
		).Str("ids", idsStr).Msg("unable to delete inbound services")
		return 0, err
	}
	reg.logger.Debug().Err(err).Uint(
		"uid", uid,
	).Str("ids", idsStr).Int64("nb", nb).Msg("inbound services deleted")
	return nb, nil
}
