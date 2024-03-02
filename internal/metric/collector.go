package metric

import (
	"sync"
	"time"

	"github.com/ncarlier/readflow/internal/db"
	"github.com/ncarlier/readflow/pkg/logger"
	"github.com/rs/zerolog"
)

var defaultInterval = time.Duration(30) * time.Second

// Collector is a metric collector
type Collector interface {
	Collect() error
}

// Collectors is a go routine in charge of the metrics collector.
type Collectors struct {
	interval   time.Duration
	collectors []Collector
	logger     zerolog.Logger
	stopChan   chan bool
	stopWait   *sync.WaitGroup
}

var collectors *Collectors

// StartCollectors starts metrics collectors
func StartCollectors(_db db.DB) {
	if collectors != nil {
		return
	}
	collectors = &Collectors{
		interval: defaultInterval,
		logger:   logger.With().Str("component", "metrics-collector").Logger(),
		stopChan: make(chan bool),
		stopWait: &sync.WaitGroup{},
		collectors: []Collector{
			newArticleMetricsCollector(_db),
			newUserMetricsCollector(_db),
		},
	}
	collectors.start()
}

// StopCollectors stop metrics collectors
func StopCollectors() {
	if collectors == nil {
		return
	}
	collectors.stop()
}

func (c *Collectors) collect() {
	c.logger.Debug().Msg("collecting metrics...")
	for _, collector := range c.collectors {
		if err := collector.Collect(); err != nil {
			c.logger.Error().Err(err).Msg("unable to collect metrics")
		}
	}
}

func (c *Collectors) start() {
	c.logger.Debug().Msg("starting...")
	ticker := time.NewTicker(c.interval)
	c.stopWait.Add(1)
	go func() {
		c.collect()
		for range ticker.C {
			c.collect()
		}
	}()

	go func() {
		<-c.stopChan
		ticker.Stop()
		c.logger.Debug().Msg("stopped")
		c.stopWait.Done()
	}()
}

func (c *Collectors) stop() {
	c.logger.Debug().Msg("stopping...")
	c.stopChan <- true
	c.stopWait.Wait()
}
