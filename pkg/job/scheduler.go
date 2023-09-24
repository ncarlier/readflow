package job

// Scheduler is a job scheduler
type Scheduler struct {
	jobs []Job
}

// NewScheduler create new job scheduler
func NewScheduler(jobs ...Job) *Scheduler {
	result := &Scheduler{
		jobs: jobs,
	}
	for _, job := range result.jobs {
		go job.Start()
	}

	return result
}

// Add and start job
func (s *Scheduler) Add(job Job) {
	job.Start()
	s.jobs = append(s.jobs, job)
}

// Shutdown job scheduler
func (s *Scheduler) Shutdown() {
	for _, job := range s.jobs {
		job.Stop()
	}
}
