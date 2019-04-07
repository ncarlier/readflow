package service

import (
	"context"
	"errors"
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

	apiKeys, err := reg.db.GetAPIKeysByUserID(uid)
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
			err = errors.New("API key not found")
		}
		return nil, err
	}
	return apiKey, nil
}

// CreateOrUpdateAPIKey create or update an API key for current user
func (reg *Registry) CreateOrUpdateAPIKey(ctx context.Context, id *uint, alias string) (*model.APIKey, error) {
	uid := getCurrentUserFromContext(ctx)

	builder := model.NewAPIKeyBuilder()
	apiKey := builder.UserID(uid).Alias(alias).Build()
	apiKey.ID = id
	result, err := reg.db.CreateOrUpdateAPIKey(*apiKey)
	if err != nil {
		evt := reg.logger.Info().Err(err).Uint(
			"uid", uid,
		).Str("alias", alias)
		if id != nil {
			evt.Uint("id", *id).Msg("unable to update API key")
		} else {
			evt.Msg("unable to create API key")
		}
		return nil, err
	}
	return result, err
}

// DeleteAPIKey delete an API key of the current user
func (reg *Registry) DeleteAPIKey(ctx context.Context, id uint) (*model.APIKey, error) {
	uid := getCurrentUserFromContext(ctx)

	apiKey, err := reg.db.GetAPIKeyByID(id)
	if err != nil || apiKey == nil || apiKey.UserID != uid {
		if err == nil {
			err = errors.New("API key not found")
		}
		reg.logger.Info().Err(err).Uint(
			"uid", uid,
		).Uint("id", id).Msg("unable to delete API key")
		return nil, err
	}

	err = reg.db.DeleteAPIKey(*apiKey)
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

	nb, err := reg.db.DeleteAPIKeys(uid, ids)
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
