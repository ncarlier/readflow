package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/ncarlier/readflow/pkg/model"
	ruleengine "github.com/ncarlier/readflow/pkg/rule-engine"
)

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

	// Validate user quota
	plan, err := reg.GetCurrentUserPlan(ctx)
	if err != nil {
		return nil, err
	}
	if plan != nil {
		totalCategories, err := reg.CountCurrentUserCategories(ctx)
		if err != nil {
			reg.logger.Info().Err(err).Uint(
				"uid", uid,
			).Msg("unable to create category")
			return nil, err
		}
		if totalCategories >= plan.TotalCategories {
			err = ErrUserQuotaReached
			reg.logger.Info().Err(err).Uint(
				"uid", uid,
			).Uint(
				"total", plan.TotalCategories,
			).Msg("unable to create category")
			return nil, err
		}
	}

	// Validate category's rule
	if err := validateCategoryRule(form.Rule); err != nil {
		reg.logger.Info().Err(err).Uint(
			"uid", uid,
		).Msg("invalid category rule")
		return nil, err
	}

	// Create category
	result, err := reg.db.CreateCategoryForUser(uid, form)
	if err != nil {
		reg.logger.Info().Err(err).Uint(
			"uid", uid,
		).Str("title", form.Title).Msg("unable to create category")
		return nil, err
	}

	// Force to refresh the rule engine cache
	reg.ruleEngineCache.Evict(uid)

	return result, err
}

// UpdateCategory update a category for current user
func (reg *Registry) UpdateCategory(ctx context.Context, form model.CategoryUpdateForm) (*model.Category, error) {
	uid := getCurrentUserIDFromContext(ctx)

	// Validate category's rule
	if err := validateCategoryRule(form.Rule); err != nil {
		reg.logger.Info().Err(err).Uint(
			"uid", uid,
		).Msg("invalid category rule")
		return nil, err
	}

	// Update category
	result, err := reg.db.UpdateCategoryForUser(uid, form)
	if err != nil {
		reg.logger.Info().Err(err).Uint(
			"uid", uid,
		).Str("title", *form.Title).Uint(
			"id", form.ID,
		).Msg("unable to update category")
		return nil, err
	}

	// Force to refresh the rule engine cache
	reg.ruleEngineCache.Evict(uid)

	return result, err
}

// DeleteCategory delete a category of the current user
func (reg *Registry) DeleteCategory(ctx context.Context, id uint) (*model.Category, error) {
	uid := getCurrentUserIDFromContext(ctx)

	category, err := reg.GetCategory(ctx, id)
	if err != nil {
		return nil, err
	}

	err = reg.db.DeleteCategoryByUser(uid, id)
	if err != nil {
		reg.logger.Info().Err(err).Uint(
			"uid", uid,
		).Uint("id", id).Msg("unable to delete category")
		return nil, err
	}

	// Force to refresh the rule engine cache
	reg.ruleEngineCache.Evict(uid)

	return category, nil
}

// DeleteCategories delete categories of the current user
func (reg *Registry) DeleteCategories(ctx context.Context, ids []uint) (int64, error) {
	uid := getCurrentUserIDFromContext(ctx)
	idsStr := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(ids)), ","), "[]")

	nb, err := reg.db.DeleteCategoriesByUser(uid, ids)
	if err != nil {
		reg.logger.Info().Err(err).Uint(
			"uid", uid,
		).Str("ids", idsStr).Msg("unable to delete categories")
		return 0, err
	}
	reg.logger.Debug().Err(err).Uint(
		"uid", uid,
	).Str("ids", idsStr).Int64("nb", nb).Msg("categories deleted")

	// Force to refresh the rule engine cache
	reg.ruleEngineCache.Evict(uid)

	return nb, nil
}

func validateCategoryRule(rule *string) error {
	if rule == nil {
		return nil
	}
	// Create dummy category in order to validate rule
	category := model.Category{
		Rule: rule,
	}
	// Validate category's rule
	_, err := ruleengine.NewRuleProcessor(category)
	return err
}
