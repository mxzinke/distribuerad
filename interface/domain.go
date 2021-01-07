package domain

import "time"

type IChannelStore interface {
	GetChannel(name string) IChannel
	AddChannel(name string) IChannel
	DeleteChannel(name string)
}

type IChannel interface {
	GetEvents() []*Event
	AddEvent(data string) *Event
	AddDelayedEvent(data string, publishAt time.Time) *Event
	DeleteEvent(eventID string) error

	// Related to jobs:
	GetJobs() []*Job
	AddJob(jobID, data, cronDef string) (*Job, error)
	DeleteJob(jobID string) error
}

type Event struct {
	ID          string    `json:"eventID"`
	PublishedAt time.Time `json:"publishedAt"`
	Data        string    `json:"data"`
}

type Job struct {
	ID        string    `json:"jobID"`
	CronDef   string    `json:"cron,omitempty"`
	Data      string    `json:"data"`
	CreatedAt time.Time `json:"createdAt"`
}
