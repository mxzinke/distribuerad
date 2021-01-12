package channel

import (
	"fmt"
	"time"
)

type IJob interface {
	ID() string
	SetTrigger(trigger func(data string, ttl time.Duration))
	OnAttach()
	OnDetach()
}

// Returns all available IJob objects.
func (c *Channel) GetJobs() []IJob {
	c.jobsLock.RLock()
	defer c.jobsLock.RUnlock()

	var currentJobs = make([]IJob, len(c.jobs))
	var i = 0
	for _, job := range c.jobs {
		currentJobs[i] = job
		i++
	}

	return currentJobs
}

func (c *Channel) AddJob(job IJob) (IJob, error) {
	c.jobsLock.Lock()
	defer c.jobsLock.Unlock()

	jobID := job.ID()
	if c.jobs[jobID] != nil {
		return c.jobs[jobID], fmt.Errorf("The job with name '%s' already exists! ", jobID)
	}
	c.jobs[jobID] = job
	c.jobs[jobID].OnAttach()

	return c.jobs[jobID], nil
}

func (c *Channel) DeleteJob(jobID string) error {
	c.jobsLock.Lock()
	defer c.jobsLock.Unlock()

	if c.jobs[jobID] == nil {
		return fmt.Errorf("The job '%s' does not exist! ", jobID)
	}

	c.jobs[jobID].OnDetach()
	delete(c.jobs, jobID)

	return nil
}
