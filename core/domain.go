package domain

import (
	"time"
)

type IChannelStore interface {
	GetChannel(name string) IChannel
	AddChannel(name string) IChannel
	DeleteChannel(name string)
}

type IChannel interface {
	GetEvents() []*Event
	AddEvent(data string, duration time.Duration) *Event
	AddDelayedEvent(data string, publishAt time.Time, ttl time.Duration) *Event
	DeleteEvent(eventID string) error

	// Related to jobs:
	GetJobs() []*Job
	AddJob(jobID, data, cronDef string, ttl time.Duration) (*Job, error)
	DeleteJob(jobID string) error
}

type Job struct {
	ID        string        `json:"jobID"`
	CronDef   string        `json:"cron"`
	Data      string        `json:"data"`
	TTL       time.Duration `json:"ttl"`
	CreatedAt time.Time     `json:"createdAt"`
}
