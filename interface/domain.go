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
	AddJob(jobID, data string, every time.Duration) *Job
	AddCronJob(jobID, data, cronDef string) *Job
	DeleteJob(jobID string)
}

type Event struct {
	ID          string    `json:"eventID"`
	PublishedAt time.Time `json:"publishedAt"`
	Data        string    `json:"data"`
}

type Job struct {
	ID    string        `json:"jobID"`
	Cron  string        `json:"cron,omitempty"`
	Every time.Duration `json:"every,omitempty"`
	Data  string        `json:"data"`
}
