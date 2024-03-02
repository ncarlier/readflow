package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/ncarlier/readflow/internal/model"
)

const unableToCreateCategoryErrorMsg = "unable to create category"

// GetCategories get categories from current user
func (reg *Registry) GetCategories(ctx context.Context) ([]model.Category, error) {
	uid := getCurrentUserIDFromContext(ctx)

	categories, err := reg.db.GetCategoriesByUser(uid)
	if err != nil {
		reg.logger.Info().Err(err).Uint(
			"uid", uid,
		).Msg("unable to get categories")
	}
	return categories, err
}

// CountCurrentUserCategories get total categories of current user
func (reg *Registry) CountCurrentUserCategories(ctx context.Context) (uint, error) {
	uid := getCurrentUserIDFromContext(ctx)
	return reg.db.CountCategoriesByUser(uid)
}

// GetCategory get a category of the current user
func (reg *Registry) GetCategory(ctx context.Context, id uint) (*model.Category, error) {
	uid := getCurrentUserIDFromContext(ctx)

	category, err := reg.db.GetCategoryByID(id)
	if err != nil || category == nil || *category.UserID != uid {
		if err == nil {
			err = ErrCategoryNotFound
		}
		return nil, err
	}
	return category, nil
}

// CreateCategory create a category for current user
func (reg *Registry) CreateCategory(ctx context.Context, form model.CategoryCreateForm) (*model.Category, error) {
	uid := getCurrentUserIDFromContext(ctx)

	logger := reg.logger.With().Uint("uid", uid).Logger()

	// Validate user quota
	plan, err := reg.GetCurrentUserPlan(ctx)
	if err != nil {
		return nil, err
	}
	if plan != nil {
		totalCategories, err := reg.CountCurrentUserCategories(ctx)
		if err != nil {
			logger.Info().Err(err).Msg(unableToCreateCategoryErrorMsg)
			return nil, err
		}
		if totalCategories >= plan.CategoriesLimit {
			err = ErrUserQuotaReached
			logger.Info().Err(err).Uint("total", totalCategories).Msg(unableToCreateCategoryErrorMsg)
			return nil, err
		}
	}
	logger = logger.With().Str("title", form.Title).Logger()

	// Create category
	logger.Debug().Msg("creating category...")
	result, err := reg.db.CreateCategoryForUser(uid, form)
	if err != nil {
		logger.Info().Err(err).Msg(unableToCreateCategoryErrorMsg)
		return nil, err
	}
	logger.Info().Uint("id", *result.ID).Msg("category created")

	return result, err
}

// UpdateCategory update a category for current user
func (reg *Registry) UpdateCategory(ctx context.Context, form model.CategoryUpdateForm) (*model.Category, error) {
	uid := getCurrentUserIDFromContext(ctx)

	logger := reg.logger.With().Uint("uid", uid).Uint("id", form.ID).Logger()

	// Update category
	logger.Debug().Msg("updating category...")
	result, err := reg.db.UpdateCategoryForUser(uid, form)
	if err != nil {
		logger.Info().Err(err).Msg("unable to update category")
		return nil, err
	}
	logger.Info().Msg("category updated")

	return result, err
}

// DeleteCategory delete a category of the current user
func (reg *Registry) DeleteCategory(ctx context.Context, id uint) (*model.Category, error) {
	uid := getCurrentUserIDFromContext(ctx)

	logger := reg.logger.With().Uint("uid", uid).Uint("id", id).Logger()

	category, err := reg.GetCategory(ctx, id)
	if err != nil {
		return nil, err
	}

	logger.Debug().Msg("deleting category...")
	err = reg.db.DeleteCategoryByUser(uid, id)
	if err != nil {
		logger.Info().Err(err).Msg("unable to delete category")
		return nil, err
	}
	logger.Info().Msg("category deleted")

	return category, nil
}

// DeleteCategories delete categories of the current user
func (reg *Registry) DeleteCategories(ctx context.Context, ids []uint) (int64, error) {
	uid := getCurrentUserIDFromContext(ctx)
	idsStr := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(ids)), ","), "[]")

	logger := reg.logger.With().Uint("uid", uid).Str("ids", idsStr).Logger()

	logger.Debug().Msg("deleting categories...")
	nb, err := reg.db.DeleteCategoriesByUser(uid, ids)
	if err != nil {
		logger.Info().Err(err).Msg("unable to delete categories")
		return 0, err
	}
	logger.Info().Int64("nb", nb).Msg("categories deleted")

	return nb, nil
}
