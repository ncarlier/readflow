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
	uid := getCurrentUserFromContext(ctx)

	categories, err := reg.db.GetCategoriesByUserID(uid)
	if err != nil {
		reg.logger.Info().Err(err).Uint(
			"uid", uid,
		).Msg("unable to get categories")
	}
	return categories, err
}

// CountCurrentUserCategories get total categories of current user
func (reg *Registry) CountCurrentUserCategories(ctx context.Context) (uint, error) {
	uid := getCurrentUserFromContext(ctx)
	return reg.db.CountCategoriesByUserID(uid)
}

// GetCategory get a category of the current user
func (reg *Registry) GetCategory(ctx context.Context, id uint) (*model.Category, error) {
	uid := getCurrentUserFromContext(ctx)

	category, err := reg.db.GetCategoryByID(id)
	if err != nil || category == nil || *category.UserID != uid {
		if err == nil {
			err = ErrCategoryNotFound
		}
		return nil, err
	}
	return category, nil
}

// CreateOrUpdateCategory create or update a category for current user
func (reg *Registry) CreateOrUpdateCategory(ctx context.Context, form model.CategoryForm) (*model.Category, error) {
	uid := getCurrentUserFromContext(ctx)

	plan, err := reg.GetCurrentUserPlan(ctx)
	if err != nil {
		return nil, err
	}

	if plan != nil && form.ID == nil {
		// Check user quota
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

	// Create category
	builder := model.NewCategoryBuilder()
	category := builder.Form(&form).UserID(uid).Build()

	// Validate category's rule
	_, err = ruleengine.NewRuleProcessor(*category)
	if err != nil {
		reg.logger.Info().Err(err).Uint(
			"uid", uid,
		).Msg("invalid category rule")
		return nil, err
	}

	result, err := reg.db.CreateOrUpdateCategory(*category)
	if err != nil {
		evt := reg.logger.Info().Err(err).Uint(
			"uid", uid,
		).Str("title", *form.Title)
		if form.ID != nil {
			evt.Uint("id", *form.ID).Msg("unable to update category")
		} else {
			evt.Msg("unable to create category")
		}
		return nil, err
	}

	// Force to refresh the rule engine cache
	reg.ruleEngineCache.Evict(uid)

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

	// Force to refresh the rule engine cache
	reg.ruleEngineCache.Evict(uid)

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

	// Force to refresh the rule engine cache
	reg.ruleEngineCache.Evict(uid)

	return nb, nil
}
