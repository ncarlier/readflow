package service

import (
	"context"
	"errors"

	"github.com/ncarlier/reader/pkg/model"
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

	category, err := reg.db.GetCategoryByID(id)
	if err != nil || category == nil || *category.UserID != uid {
		if err == nil {
			err = errors.New("category not found")
		}
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
