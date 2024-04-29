package job

import (
	"github.com/labstack/gommon/log"
	"time"
)

type Job struct {
	name      string
	frequency time.Duration
	run       func() error
}

func NewJob(name string, frequency time.Duration, run func() error) *Job {
	return &Job{
		name:      name,
		frequency: frequency,
		run:       run,
	}
}

type Scheduler struct {
	jobs []*Job
}

func NewScheduler(jobs []*Job) *Scheduler {
	return &Scheduler{jobs: jobs}
}

func (s *Scheduler) Start() {
	for _, job := range s.jobs {
		go s.runJob(job)
	}
}

func (s *Scheduler) runJob(job *Job) {
	ticker := time.NewTicker(job.frequency)
	defer ticker.Stop()

	for range ticker.C {
		if err := job.run(); err != nil {
			log.Errorf("Failed to run job %s: %v", job.name, err)
		}
	}
}
