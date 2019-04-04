package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	ruleengine "github.com/ncarlier/reader/pkg/service/rule-engine"

	"github.com/ncarlier/reader/pkg/model"
)

// GetRules get rules from current user
func (reg *Registry) GetRules(ctx context.Context) (*[]model.Rule, error) {
	uid := getCurrentUserFromContext(ctx)

	rules, err := reg.db.GetRulesByUserID(uid)
	if err != nil {
		reg.logger.Info().Err(err).Uint(
			"uid", uid,
		).Msg("unable to get rules")
		return nil, err
	}

	return &rules, err
}

// GetRule get an rule of the current user
func (reg *Registry) GetRule(ctx context.Context, id uint) (*model.Rule, error) {
	uid := getCurrentUserFromContext(ctx)

	rule, err := reg.db.GetRuleByID(id)
	if err != nil || rule == nil || *rule.UserID != uid {
		if err == nil {
			err = errors.New("Rule not found")
		}
		return nil, err
	}
	return rule, nil
}

// CreateOrUpdateRule create or update a rule for current user
func (reg *Registry) CreateOrUpdateRule(ctx context.Context, form model.RuleForm) (*model.Rule, error) {
	uid := getCurrentUserFromContext(ctx)

	builder := model.NewRuleBuilder()
	rule := builder.UserID(uid).Form(&form).Build()

	// Validate rule configuration
	_, err := ruleengine.NewRuleProcessor(rule)
	if err != nil {
		reg.logger.Info().Err(err).Uint(
			"uid", uid,
		).Msg("unable to configure rule")
		return nil, err
	}

	result, err := reg.db.CreateOrUpdateRule(*rule)
	if err != nil {
		evt := reg.logger.Info().Err(err).Uint(
			"uid", uid,
		).Str("alias", form.Alias)
		if form.ID != nil {
			evt.Uint("id", *form.ID).Msg("unable to update rule")
		} else {
			evt.Msg("unable to create rule")
		}
		return nil, err
	}
	return result, err
}

// DeleteRule delete a rule of the current user
func (reg *Registry) DeleteRule(ctx context.Context, id uint) (*model.Rule, error) {
	uid := getCurrentUserFromContext(ctx)

	rule, err := reg.db.GetRuleByID(id)
	if err != nil || rule == nil || *rule.UserID != uid {
		if err == nil {
			err = errors.New("rule not found")
		}
		reg.logger.Info().Err(err).Uint(
			"uid", uid,
		).Uint("id", id).Msg("unable to delete rule")
		return nil, err
	}

	err = reg.db.DeleteRule(*rule)
	if err != nil {
		reg.logger.Info().Err(err).Uint(
			"uid", uid,
		).Uint("id", id).Msg("unable to delete rule")
		return nil, err
	}
	return rule, nil
}

// DeleteRules delete rules of the current user
func (reg *Registry) DeleteRules(ctx context.Context, ids []uint) (int64, error) {
	uid := getCurrentUserFromContext(ctx)
	idsStr := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(ids)), ","), "[]")

	nb, err := reg.db.DeleteRules(uid, ids)
	if err != nil {
		reg.logger.Info().Err(err).Uint(
			"uid", uid,
		).Str("ids", idsStr).Msg("unable to delete rules")
		return 0, err
	}
	reg.logger.Debug().Err(err).Uint(
		"uid", uid,
	).Str("ids", idsStr).Int64("nb", nb).Msg("rules deleted")
	return nb, nil
}
