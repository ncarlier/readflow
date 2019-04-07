package event

import (
	"github.com/ncarlier/readflow/pkg/metric"
	"github.com/ncarlier/readflow/pkg/model"
)

func init() {
	bus.Subscribe(CreateArticle, func(payload ...interface{}) {
		if article, ok := payload[0].(model.Article); ok {
			metric.IncNewArticlesCounter(article)
		}
	})
}
