package listener

import (
	"github.com/ncarlier/readflow/pkg/event"
	"github.com/ncarlier/readflow/pkg/metric"
	"github.com/ncarlier/readflow/pkg/model"
)

func init() {
	event.Subscribe(event.CreateArticle, func(payload ...interface{}) {
		if article, ok := payload[0].(model.Article); ok {
			metric.IncNewArticlesCounter(article)
		}
	})
}
