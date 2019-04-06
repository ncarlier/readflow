package job

import "github.com/ncarlier/reader/pkg/db"

// Scheduler is a job scheduler
type Scheduler struct {
	cleanDatabaseJob *CleanDatabaseJob
}

// StartNewScheduler create and start new job scheduler
func StartNewScheduler(_db db.DB) *Scheduler {
	return &Scheduler{
		cleanDatabaseJob: NewCleanDatabaseJob(_db),
	}
}

// Shutdown job scheduler
func (s *Scheduler) Shutdown() {
	s.cleanDatabaseJob.Stop()
}
