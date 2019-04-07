package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/ncarlier/readflow/pkg/model"
)

// GetCategories get categories from current user
func (reg *Registry) GetCategories(ctx context.Context) (*[]model.Category, error) {
	uid := getCurrentUserFromContext(ctx)

	categories, err := reg.db.GetCategoriesByUserID(uid)
	if err != nil {
		reg.logger.Info().Err(err).Uint(
			"uid", uid,
		).Msg("unable to get categories")
		return nil, err
	}

	return &categories, err
}

// GetCategory get a category of the current user
func (reg *Registry) GetCategory(ctx context.Context, id uint) (*model.Category, error) {
	uid := getCurrentUserFromContext(ctx)

	category, err := reg.db.GetCategoryByID(id)
	if err != nil || category == nil || *category.UserID != uid {
		if err == nil {
			err = errors.New("category not found")
		}
		return nil, err
	}
	return category, nil
}

// CreateOrUpdateCategory create or update a category for current user
func (reg *Registry) CreateOrUpdateCategory(ctx context.Context, id *uint, title string) (*model.Category, error) {
	uid := getCurrentUserFromContext(ctx)

	category := model.Category{
		ID:     id,
		UserID: &uid,
		Title:  title,
	}
	result, err := reg.db.CreateOrUpdateCategory(category)
	if err != nil {
		evt := reg.logger.Info().Err(err).Uint(
			"uid", uid,
		).Str("title", title)
		if id != nil {
			evt.Uint("id", *id).Msg("unable to update category")
		} else {
			evt.Msg("unable to create category")
		}
		return nil, err
	}
	return result, err
}

// DeleteCategory delete a category of the current user
func (reg *Registry) DeleteCategory(ctx context.Context, id uint) (*model.Category, error) {
	uid := getCurrentUserFromContext(ctx)

	category, err := reg.GetCategory(ctx, id)
	if err != nil {
		reg.logger.Info().Err(err).Uint(
			"uid", uid,
		).Uint("id", id).Msg("unable to delete category")
		return nil, err
	}

	err = reg.db.DeleteCategory(*category)
	if err != nil {
		reg.logger.Info().Err(err).Uint(
			"uid", uid,
		).Uint("id", id).Msg("unable to delete category")
		return nil, err
	}
	return category, nil
}

// DeleteCategories delete categories of the current user
func (reg *Registry) DeleteCategories(ctx context.Context, ids []uint) (int64, error) {
	uid := getCurrentUserFromContext(ctx)
	idsStr := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(ids)), ","), "[]")

	nb, err := reg.db.DeleteCategories(uid, ids)
	if err != nil {
		reg.logger.Info().Err(err).Uint(
			"uid", uid,
		).Str("ids", idsStr).Msg("unable to delete categories")
		return 0, err
	}
	reg.logger.Debug().Err(err).Uint(
		"uid", uid,
	).Str("ids", idsStr).Int64("nb", nb).Msg("categories deleted")
	return nb, nil
}
