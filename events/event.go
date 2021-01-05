package events

import (
	"fmt"
	"github.com/google/uuid"
	"time"
)

func (c *Channel) AddDelayedEvent(data string, publishAt time.Time) *Event {
	event := &Event{
		ID:          uuid.New().String(),
		PublishedAt: publishAt,
		Data:        data,
	}

	time.AfterFunc(publishAt.Sub(time.Now()), func() {
		c.lock.Lock()
		defer c.lock.Unlock()
		c.events = append(c.events, event)
	})

	return event
}

func (c *Channel) AddEvent(data string) *Event {
	c.lock.Lock()
	defer c.lock.Unlock()

	event := &Event{
		ID:          uuid.New().String(),
		PublishedAt: time.Now(),
		Data:        data,
	}

	c.events = append(c.events, event)
	return event
}

func (c *Channel) DeleteEvent(eventID string) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	var index = -1
	for i, e := range c.events {
		if e.ID == eventID {
			index = i
			break
		}
	}

	if index == 1 {
		return fmt.Errorf("No event could be found with ID %s. Maybe it has already been picked up ", eventID)
	}

	c.events = append(c.events[:index], c.events[index:]...)

	return nil
}
