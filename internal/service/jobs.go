package service

import "github.com/ncarlier/readflow/pkg/job"

// AddJob add job to the service
func (reg *Registry) AddJob(j job.Job) {
	reg.scheduler.Add(j)
}
