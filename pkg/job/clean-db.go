package job

import (
	"time"

	"github.com/ncarlier/readflow/pkg/db"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	maximumArticleRetentionDuration = 48 * time.Hour
	maximumDeviceInactivityDuration = 30 * 24 * time.Hour
)

// CleanDatabaseJob is a job to clean the database
type CleanDatabaseJob struct {
	db     db.DB
	ticker *time.Ticker
	logger zerolog.Logger
}

// NewCleanDatabaseJob create and start new job to clean the database
func NewCleanDatabaseJob(_db db.DB) *CleanDatabaseJob {
	job := &CleanDatabaseJob{
		db:     _db,
		ticker: time.NewTicker(time.Hour),
		logger: log.With().Str("component", "scheduler").Str("job", "clean-db").Logger(),
	}
	go job.start()
	return job
}

func (cdj *CleanDatabaseJob) start() {
	cdj.logger.Debug().Msg("job started")
	for range cdj.ticker.C {
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
		evt.Int64("removed_articles", nb).Msg("cleanup done")
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
		evt.Int64("removed_devices", nb).Msg("cleanup done")
	}
}

// Stop job
func (cdj *CleanDatabaseJob) Stop() {
	cdj.ticker.Stop()
	cdj.logger.Debug().Msg("job stopped")
}
