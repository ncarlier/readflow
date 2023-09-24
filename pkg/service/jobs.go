package service

import "github.com/ncarlier/readflow/pkg/job"

// AddJob add job to the service
func (reg *Registry) AddJob(job job.Job) {
	reg.scheduler.Add(job)
}
