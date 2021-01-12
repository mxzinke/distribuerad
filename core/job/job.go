package job

import (
	"github.com/robfig/cron/v3"
	"log"
	"time"
)

// Job describes a Cron-/Interval Job doing an action each time it triggers
type Job struct {
	JobID     string        `json:"jobID"`
	CronDef   string        `json:"cron"`
	Data      string        `json:"data"`
	TTL       time.Duration `json:"ttl"`
	CreatedAt time.Time     `json:"createdAt"`

	// For internal processes
	cronJob *cron.Cron
	trigger TriggerFunc
}

// TriggerFunc is a shortcut for a function providing event-data, used for the job trigger
type TriggerFunc func(data string, ttl time.Duration)

func New(jobID, data, cronDef string, eventTTL time.Duration) *Job {
	return &Job{
		JobID:     jobID,
		CronDef:   cronDef,
		Data:      data,
		TTL:       eventTTL,
		CreatedAt: time.Now(),
	}
}

func (j *Job) ID() string {
	return j.JobID
}

func (j *Job) SetTrigger(trigger TriggerFunc) {
	j.trigger = trigger
}

func (j *Job) OnAttach() {
	// add the execution cronJob / intervalJob
	j.cronJob = cron.New(cron.WithChain(
		cron.Recover(cron.DefaultLogger),
	))
	if _, err := j.cronJob.AddFunc(j.CronDef, func() {
		if j.trigger == nil {
			// Nothing will happen, because we can not call a trigger anyway
			return
		}
		j.trigger(j.Data, j.TTL)
	}); err != nil {
		log.Printf("Error at starting (cron-)job: %v", err)
	}
	j.cronJob.Start()
}

func (j *Job) OnDetach() {
	j.cronJob.Stop()
}
