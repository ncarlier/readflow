package db

import (
	"time"

	"github.com/ncarlier/readflow/pkg/job"
	"github.com/ncarlier/readflow/pkg/logger"
	"github.com/rs/zerolog"
)

const (
	maximumArticleRetentionDuration = 48 * time.Hour
	maximumDeviceInactivityDuration = 30 * 24 * time.Hour
)

// CleanupDatabaseJob is a job to clean the database
type CleanupDatabaseJob struct {
	db     DB
	ticker *time.Ticker
	logger zerolog.Logger
}

// NewCleanupDatabaseJob create and start new job to clean the database
func NewCleanupDatabaseJob(db DB) job.Job {
	return &CleanupDatabaseJob{
		db:     db,
		ticker: time.NewTicker(time.Hour),
		logger: logger.With().Str("component", "scheduler").Str("job", "clean-db").Logger(),
	}
}

// Start the cleanup job
func (cdj *CleanupDatabaseJob) Start() {
	cdj.logger.Debug().Msg("job started")
	for range cdj.ticker.C {
		start := time.Now()
		cdj.logger.Debug().Msg("running job...")
		nb, err := cdj.db.DeleteReadArticlesOlderThan(maximumArticleRetentionDuration)
		if err != nil {
			cdj.logger.Error().Err(err).Msg("unable to clean old articles from the database")
			break
		}
		// Using info level only for effective cleanup
		evt := cdj.logger.Debug()
		if nb > 0 {
			evt = cdj.logger.Info()
		}
		evt.Int64("removed_articles", nb).Dur("took", time.Since(start)).Msg("cleanup done")
		start = time.Now()
		nb, err = cdj.db.DeleteInactiveDevicesOlderThan(maximumDeviceInactivityDuration)
		if err != nil {
			cdj.logger.Error().Err(err).Msg("unable to clean old devices from the database")
			break
		}
		// Using info level only for effective cleanup
		evt = cdj.logger.Debug()
		if nb > 0 {
			evt = cdj.logger.Info()
		}
		evt.Int64("removed_devices", nb).Dur("took", time.Since(start)).Msg("cleanup done")
	}
}

// Stop the cleanup job
func (cdj *CleanupDatabaseJob) Stop() {
	cdj.ticker.Stop()
	cdj.logger.Debug().Msg("job stopped")
}
