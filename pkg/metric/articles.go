package metric

import (
	"fmt"

	"github.com/ncarlier/readflow/pkg/db"
	"github.com/ncarlier/readflow/pkg/model"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	createdArticles = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: metricName("articles_created_total"),
		Help: "Total number of created articles by user",
	}, []string{"uid"})
	totalArticles = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: metricName("articles_total"),
		Help: "Total number of articles",
	}, []string{"status"})
)

// IncNewArticlesCounter increments the counter of new articles
func IncNewArticlesCounter(article model.Article) {
	createdArticles.With(prometheus.Labels{"uid": fmt.Sprint(article.UserID)}).Inc()
}

// ArticleMetricsCollector article's metrics collector
type articleMetricsCollector struct {
	db db.DB
}

func newArticleMetricsCollector(_db db.DB) Collector {
	return &articleMetricsCollector{
		db: _db,
	}
}

// Collect article's metrics
func (c *articleMetricsCollector) Collect() error {
	nb, err := c.db.CountArticles("read")
	if err != nil {
		return err
	}
	totalArticles.With(
		prometheus.Labels{"status": "read"},
	).Set(float64(nb))
	nb, err = c.db.CountArticles("unread")
	if err != nil {
		return err
	}
	totalArticles.With(
		prometheus.Labels{"status": "unread"},
	).Set(float64(nb))
	return nil
}
