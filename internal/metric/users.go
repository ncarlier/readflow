package metric

import (
	"github.com/ncarlier/readflow/internal/db"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	totalUsers = promauto.NewGauge(prometheus.GaugeOpts{
		Name: metricName("users_total"),
		Help: "Total number of users",
	})
)

type userMetricsCollector struct {
	db db.DB
}

func newUserMetricsCollector(_db db.DB) Collector {
	return &userMetricsCollector{
		db: _db,
	}
}

// Collect user's metrics
func (c *userMetricsCollector) Collect() error {
	nb, err := c.db.CountUsers()
	if err != nil {
		return err
	}
	totalUsers.Set(float64(nb))
	return nil
}
