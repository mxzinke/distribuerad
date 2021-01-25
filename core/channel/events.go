package channel

import (
	"fmt"
	"time"
)

type IEvent interface {
	ID() string
	PublishedAt() time.Time
	ValidUntil() time.Time
}

func (c *Channel) GetEvents() []IEvent {
	c.lock.RLock()
	defer c.lock.RUnlock()

	validEvents := c.getValidEvents()

	// in case, there are invalid elements
	if len(c.events) > len(validEvents) {
		go c.cleanupInvalidEvents()
	}

	return validEvents
}

func (c *Channel) GetEvent(eventID string) IEvent {
	c.lock.RLock()
	defer c.lock.RUnlock()

	index := c.findEventIndex(eventID)
	if index == -1 {
		// return fmt.Errorf("No event could be found with ID %s. Maybe it has already been picked up. ", eventID)
		return nil
	}

	return c.events[index]
}

func (c *Channel) AddEvent(event IEvent) IEvent {
	if event.PublishedAt().After(time.Now()) {
		time.AfterFunc(event.PublishedAt().Sub(time.Now()), func() {
			c.events = append(c.events, event)
		})
	} else {
		c.lock.Lock()
		c.events = append(c.events, event)
		c.lock.Unlock()
	}

	return event
}

func (c *Channel) DeleteEvent(eventID string) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	index := c.findEventIndex(eventID)
	if index == -1 {
		return fmt.Errorf("No event could be found with ID %s. Maybe it has already been picked up.", eventID)
	}

	if index == 0 {
		c.events = c.events[1:]
	} else {
		c.events = append(c.events[:index], c.events[index+1:]...)
	}

	return nil
}
