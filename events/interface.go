package events

import (
	"sync"
	"time"
)

type Event struct {
	ID          string    `json:"eventID"`
	PublishedAt time.Time `json:"publishedAt"`
	Data        string    `json:"data"`
}

type Channel struct {
	events []*Event
	lock   *sync.RWMutex
}
