package events

import (
	domain "distribuerad/interface"
	"fmt"
	"github.com/robfig/cron/v3"
)

type jobExecution struct {
	definition *domain.Job
	runner     *cron.Cron
}

func (c *Channel) GetJobs() []*domain.Job {
	c.jobsLock.RLock()
	defer c.jobsLock.RUnlock()

	var currentJobs = make([]*domain.Job, len(c.jobs))
	var i = 0
	for _, value := range c.jobs {
		currentJobs[i] = value.definition
		i++
	}

	return currentJobs
}

func (c *Channel) AddJob(jobID, data, cronDef string) (*domain.Job, error) {
	c.jobsLock.Lock()
	defer c.jobsLock.Unlock()

	if c.jobs[jobID] != nil {
		return c.jobs[jobID].definition, fmt.Errorf("The job with name %s already exists! ", jobID)
	}

	// add the execution cronJob / intervalJob
	execution := cron.New(cron.WithChain(
		cron.Recover(cron.DefaultLogger),
	))
	if _, err := execution.AddFunc(cronDef, func() {
		c.AddEvent(data)
	}); err != nil {
		return nil, fmt.Errorf("Error at cronJob: %v", err)
	}
	execution.Start()

	c.jobs[jobID] = &jobExecution{
		definition: &domain.Job{
			ID:      jobID,
			Data:    data,
			CronDef: cronDef,
		},
		runner: execution,
	}
	return c.jobs[jobID].definition, nil
}

func (c *Channel) DeleteJob(jobID string) error {
	c.jobsLock.Lock()
	defer c.jobsLock.Unlock()

	if c.jobs[jobID] == nil {
		return fmt.Errorf("The job %s does not exist! ", jobID)
	}

	c.jobs[jobID].runner.Stop()
	delete(c.jobs, jobID)

	return nil
}
