package service

import (
	"context"

	"github.com/ncarlier/readflow/pkg/model"
	ruleengine "github.com/ncarlier/readflow/pkg/rule-engine"
)

// ProcessArticleByRuleEngine apply user's rules on the article
func (reg *Registry) ProcessArticleByRuleEngine(ctx context.Context, article *model.ArticleCreateForm) error {
	uid := getCurrentUserIDFromContext(ctx)
	// Retrieve pipeline from cache
	pipeline := reg.ruleEngineCache.Get(uid)
	if pipeline == nil {
		reg.logger.Debug().Uint(
			"uid", uid,
		).Msg("loading rules into the cache")
		// Init pipeline if not in cache
		categories, err := reg.GetCategories(ctx)
		if err != nil {
			return err
		}
		pipeline, err = ruleengine.NewProcessorsPipeline(categories)
		if err != nil {
			return err
		}
		reg.ruleEngineCache.Set(uid, pipeline)
	}
	applied, err := pipeline.Apply(ctx, article)
	if err != nil {
		reg.logger.Info().Err(err).Uint(
			"uid", uid,
		).Str("title", article.Title).Msg("unable to apply rules on the article")
		return err
	}
	if applied {
		reg.logger.Debug().Uint(
			"uid", uid,
		).Str("title", article.Title).Uint(
			"category", *article.CategoryID,
		).Msg("rule applied on the article")
	}
	return nil
}
