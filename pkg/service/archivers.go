package service

import (
	"context"
	"errors"

	"github.com/ncarlier/reader/pkg/model"
)

// GetArchivers get archivers from current user
func (reg *Registry) GetArchivers(ctx context.Context) (*[]model.Archiver, error) {
	uid := getCurrentUserFromContext(ctx)

	archivers, err := reg.db.GetArchiversByUserID(uid)
	if err != nil {
		reg.logger.Info().Err(err).Uint(
			"uid", uid,
		).Msg("unable to get archivers")
		return nil, err
	}

	return &archivers, err
}

// CreateOrUpdateArchiver create or update a archiver for current user
func (reg *Registry) CreateOrUpdateArchiver(ctx context.Context, form model.ArchiverForm) (*model.Archiver, error) {
	uid := getCurrentUserFromContext(ctx)

	builder := model.NewArchiverBuilder()
	archiver := builder.UserID(uid).Form(&form).Build()
	result, err := reg.db.CreateOrUpdateArchiver(*archiver)
	if err != nil {
		evt := reg.logger.Info().Err(err).Uint(
			"uid", uid,
		).Str("alias", form.Alias)
		if form.ID != nil {
			evt.Uint("id", *form.ID).Msg("unable to update archiver")
		} else {
			evt.Msg("unable to create archiver")
		}
		return nil, err
	}
	return result, err
}

// DeleteArchiver delete a archiver of the current user
func (reg *Registry) DeleteArchiver(ctx context.Context, id uint) (*model.Archiver, error) {
	uid := getCurrentUserFromContext(ctx)

	archiver, err := reg.db.GetArchiverByID(id)
	if err != nil || archiver == nil || *archiver.UserID != uid {
		if err == nil {
			err = errors.New("archiver not found")
		}
		reg.logger.Info().Err(err).Uint(
			"uid", uid,
		).Uint("id", id).Msg("unable to delete archiver")
		return nil, err
	}

	err = reg.db.DeleteArchiver(*archiver)
	if err != nil {
		reg.logger.Info().Err(err).Uint(
			"uid", uid,
		).Uint("id", id).Msg("unable to delete archiver")
		return nil, err
	}
	return archiver, nil
}
