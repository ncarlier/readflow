package metric

import (
	"fmt"

	"github.com/ncarlier/readflow/pkg/model"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	createdArticles = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "articles_created_total",
		Help: "The total number of created articles by user",
	}, []string{"uid"})
)

// IncNewArticlesCounter increments the counter of new articles
func IncNewArticlesCounter(article model.Article) {
	createdArticles.With(prometheus.Labels{"uid": fmt.Sprint(article.UserID)}).Inc()
}
