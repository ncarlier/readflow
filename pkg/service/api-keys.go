package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/ncarlier/readflow/pkg/model"
)

// GetAPIKeyByToken returns an API key by its token
func (reg *Registry) GetAPIKeyByToken(token string) (*model.APIKey, error) {
	return reg.db.GetAPIKeyByToken(token)
}

// GetAPIKeys get API keys from current user
func (reg *Registry) GetAPIKeys(ctx context.Context) (*[]model.APIKey, error) {
	uid := getCurrentUserFromContext(ctx)

	apiKeys, err := reg.db.GetAPIKeysByUser(uid)
	if err != nil {
		reg.logger.Info().Err(err).Uint(
			"uid", uid,
		).Msg("unable to get API keys")
		return nil, err
	}

	return &apiKeys, err
}

// GetAPIKey get an API key of the current user
func (reg *Registry) GetAPIKey(ctx context.Context, id uint) (*model.APIKey, error) {
	uid := getCurrentUserFromContext(ctx)

	apiKey, err := reg.db.GetAPIKeyByID(id)
	if err != nil || apiKey == nil || apiKey.UserID != uid {
		if err == nil {
			err = ErrAPIKeyNotFound
		}
		return nil, err
	}
	return apiKey, nil
}

// CreateAPIKey create an API key for current user
func (reg *Registry) CreateAPIKey(ctx context.Context, form model.APIKeyCreateForm) (*model.APIKey, error) {
	uid := getCurrentUserFromContext(ctx)

	result, err := reg.db.CreateAPIKeyForUser(uid, form)
	if err != nil {
		reg.logger.Info().Err(err).Uint(
			"uid", uid,
		).Str("alias", form.Alias).Msg("unable to create API key")
		return nil, err
	}
	return result, err
}

// UpdateAPIKey update an API key for current user
func (reg *Registry) UpdateAPIKey(ctx context.Context, form model.APIKeyUpdateForm) (*model.APIKey, error) {
	uid := getCurrentUserFromContext(ctx)

	result, err := reg.db.UpdateAPIKeyForUser(uid, form)
	if err != nil {
		reg.logger.Info().Err(err).Uint(
			"uid", uid,
		).Str("alias", form.Alias).Uint(
			"id", form.ID,
		).Msg("unable to update API key")
		return nil, err
	}
	return result, err
}

// DeleteAPIKey delete an API key of the current user
func (reg *Registry) DeleteAPIKey(ctx context.Context, id uint) (*model.APIKey, error) {
	uid := getCurrentUserFromContext(ctx)

	apiKey, err := reg.GetAPIKey(ctx, id)
	if err != nil {
		return nil, err
	}

	err = reg.db.DeleteAPIKeyByUser(uid, id)
	if err != nil {
		reg.logger.Info().Err(err).Uint(
			"uid", uid,
		).Uint("id", id).Msg("unable to delete API key")
		return nil, err
	}
	return apiKey, nil
}

// DeleteAPIKeys delete API keys of the current user
func (reg *Registry) DeleteAPIKeys(ctx context.Context, ids []uint) (int64, error) {
	uid := getCurrentUserFromContext(ctx)
	idsStr := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(ids)), ","), "[]")

	nb, err := reg.db.DeleteAPIKeysByUser(uid, ids)
	if err != nil {
		reg.logger.Info().Err(err).Uint(
			"uid", uid,
		).Str("ids", idsStr).Msg("unable to delete API keys")
		return 0, err
	}
	reg.logger.Debug().Err(err).Uint(
		"uid", uid,
	).Str("ids", idsStr).Int64("nb", nb).Msg("API keys deleted")
	return nb, nil
}
