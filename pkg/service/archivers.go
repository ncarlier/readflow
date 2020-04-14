package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/ncarlier/readflow/pkg/service/archive"

	"github.com/ncarlier/readflow/pkg/model"
)

// GetArchivers get archivers from current user
func (reg *Registry) GetArchivers(ctx context.Context) (*[]model.Archiver, error) {
	uid := getCurrentUserFromContext(ctx)

	archivers, err := reg.db.GetArchiversByUser(uid)
	if err != nil {
		reg.logger.Info().Err(err).Uint(
			"uid", uid,
		).Msg("unable to get archivers")
		return nil, err
	}

	return &archivers, err
}

// GetArchiver get an archiver of the current user
func (reg *Registry) GetArchiver(ctx context.Context, id uint) (*model.Archiver, error) {
	uid := getCurrentUserFromContext(ctx)

	archiver, err := reg.db.GetArchiverByID(id)
	if err != nil || archiver == nil || *archiver.UserID != uid {
		if err == nil {
			err = ErrArchiverNotFound
		}
		return nil, err
	}
	return archiver, nil
}

// CreateArchiver create an archiver for current user
func (reg *Registry) CreateArchiver(ctx context.Context, form model.ArchiverCreateForm) (*model.Archiver, error) {
	uid := getCurrentUserFromContext(ctx)

	// Validate archiver configuration
	dummy := model.Archiver{
		Provider: form.Provider,
		Config:   form.Config,
	}
	_, err := archive.NewArchiveProvider(dummy)
	if err != nil {
		reg.logger.Info().Err(err).Uint(
			"uid", uid,
		).Msg("unable to configure archiver")
		return nil, err
	}

	result, err := reg.db.CreateArchiverForUser(uid, form)
	if err != nil {
		reg.logger.Info().Err(err).Uint(
			"uid", uid,
		).Str("alias", form.Alias).Msg("unable to create archiver")
		return nil, err
	}
	return result, err
}

// UpdateArchiver update a archiver for current user
func (reg *Registry) UpdateArchiver(ctx context.Context, form model.ArchiverUpdateForm) (*model.Archiver, error) {
	uid := getCurrentUserFromContext(ctx)

	if form.Provider != nil && form.Config != nil {
		// Validate archiver configuration
		dummy := model.Archiver{
			Provider: *form.Provider,
			Config:   *form.Config,
		}
		_, err := archive.NewArchiveProvider(dummy)
		if err != nil {
			reg.logger.Info().Err(err).Uint(
				"uid", uid,
			).Msg("unable to configure archiver")
			return nil, err
		}
	} else {
		// Provider can only be modify with its configuration
		form.Provider = nil
		form.Config = nil
	}

	result, err := reg.db.UpdateArchiverForUser(uid, form)
	if err != nil {
		reg.logger.Info().Err(err).Uint(
			"uid", uid,
		).Uint(
			"id", form.ID,
		).Msg("unable to update archiver")
		return nil, err
	}
	return result, err
}

// DeleteArchiver delete a archiver of the current user
func (reg *Registry) DeleteArchiver(ctx context.Context, id uint) (*model.Archiver, error) {
	uid := getCurrentUserFromContext(ctx)

	archiver, err := reg.GetArchiver(ctx, id)
	if err != nil {
		return nil, err
	}

	err = reg.db.DeleteArchiverByUser(uid, id)
	if err != nil {
		reg.logger.Info().Err(err).Uint(
			"uid", uid,
		).Uint("id", id).Msg("unable to delete archiver")
		return nil, err
	}
	return archiver, nil
}

// DeleteArchivers delete archivers of the current user
func (reg *Registry) DeleteArchivers(ctx context.Context, ids []uint) (int64, error) {
	uid := getCurrentUserFromContext(ctx)
	idsStr := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(ids)), ","), "[]")

	nb, err := reg.db.DeleteArchiversByUser(uid, ids)
	if err != nil {
		reg.logger.Info().Err(err).Uint(
			"uid", uid,
		).Str("ids", idsStr).Msg("unable to delete archivers")
		return 0, err
	}
	reg.logger.Debug().Err(err).Uint(
		"uid", uid,
	).Str("ids", idsStr).Int64("nb", nb).Msg("archivers deleted")
	return nb, nil
}

// ArchiveArticle archive an article using a archive provider
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
		logger.Info().Err(err).Msg("unable to archive article")
		return err
	}

	if alias != nil {
		logger = logger.With().Str("archiver", *alias).Logger()
	}

	archiverConf, err := reg.db.GetArchiverByUserAndAlias(uid, alias)
	if err != nil || archiverConf == nil {
		if err == nil {
			err = errors.New("archiver not found")
		}
		logger.Info().Err(err).Msg("unable to archive article")
		return err
	}

	provider, err := archive.NewArchiveProvider(*archiverConf)
	if err != nil {
		logger.Info().Err(err).Msg("unable to archive article")
		return err
	}

	logger.Debug().Msg("archiving article...")
	err = provider.Archive(ctx, *article)
	if err != nil {
		logger.Info().Err(err).Msg("unable to archive article")
		return err
	}
	logger.Info().Msg("article archived")
	return nil
}
