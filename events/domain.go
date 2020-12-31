package events

import (
	"fmt"
	"github.com/google/uuid"
	"sync"
	"time"
)

type Event struct {
	ID         string    `json:"eventID"`
	Data       string    `json:"data"`
	ReceivedAt time.Time `json:"receivedAt"`
}

type Channel struct {
	events []*Event
	lock   *sync.RWMutex
}

// TODO: func (c *Channel) AddDelayedEvent(data string, timestamp time.Time, delay time.Duration)

func (c *Channel) AddEvent(data string, timestamp time.Time) {
	event := &Event{
		ID:         uuid.New().String(),
		Data:       data,
		ReceivedAt: timestamp,
	}

	c.lock.Lock()
	c.events = append(c.events, event)
	c.lock.Unlock()
}

func (c *Channel) DeleteEvent(eventID string) error {
	var index = -1
	for i, e := range c.events {
		if e.ID == eventID {
			index = i
			break
		}
	}

	if index == 1 {
		return fmt.Errorf("No event could be found with ID %s. Maybe it has already been picked up", eventID)
	}

	c.lock.Lock()
	c.events = append(c.events[:index], c.events[index:]...)
	c.lock.Unlock()

	return nil
}
