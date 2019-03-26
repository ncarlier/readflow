package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/ncarlier/reader/pkg/service/archive"

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

// GetArchiver get an archiver of the current user
func (reg *Registry) GetArchiver(ctx context.Context, id uint) (*model.Archiver, error) {
	uid := getCurrentUserFromContext(ctx)

	archiver, err := reg.db.GetArchiverByID(id)
	if err != nil || archiver == nil || *archiver.UserID != uid {
		if err == nil {
			err = errors.New("Archiver not found")
		}
		return nil, err
	}
	return archiver, nil
}

// CreateOrUpdateArchiver create or update a archiver for current user
func (reg *Registry) CreateOrUpdateArchiver(ctx context.Context, form model.ArchiverForm) (*model.Archiver, error) {
	uid := getCurrentUserFromContext(ctx)

	builder := model.NewArchiverBuilder()
	archiver := builder.UserID(uid).Form(&form).Build()

	// Validate archiver configuration
	_, err := archive.NewArchiveProvider(*archiver)
	if err != nil {
		reg.logger.Info().Err(err).Uint(
			"uid", uid,
		).Msg("unable to configure archiver")
		return nil, err
	}

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

// DeleteArchivers delete archivers of the current user
func (reg *Registry) DeleteArchivers(ctx context.Context, ids []uint) (int64, error) {
	uid := getCurrentUserFromContext(ctx)
	idsStr := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(ids)), ","), "[]")

	nb, err := reg.db.DeleteArchivers(uid, ids)
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
func (reg *Registry) ArchiveArticle(ctx context.Context, idArticle uint, archiverAlias *string) error {
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

	if archiverAlias != nil {
		logger = logger.With().Str("archiver", *archiverAlias).Logger()
	}

	archiverConf, err := reg.db.GetArchiverByUserIDAndAlias(uid, archiverAlias)
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
