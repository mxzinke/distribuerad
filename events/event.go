package events

import (
	"distribuerad/core"
	"fmt"
	"github.com/robfig/cron/v3"
	"github.com/rs/xid"
	"sync"
	"time"
)

type Channel struct {
	events   []*domain.Event
	jobs     map[string]*jobExecution
	lock     *sync.RWMutex
	jobsLock *sync.RWMutex
	cleanup  *cron.Cron
}

const defaultTTL = 10 * time.Minute

func (c *Channel) GetEvents() []*domain.Event {
	c.lock.RLock()
	defer c.lock.RUnlock()

	validEvents := c.getValidEvents()

	// in case, there are invalid elements
	if len(c.events) > len(validEvents) {
		go c.cleanupInvalidEvents()
	}

	return validEvents
}

func (c *Channel) AddEvent(data string, ttl time.Duration) *domain.Event {
	c.lock.Lock()
	defer c.lock.Unlock()

	if ttl == 0 {
		ttl = defaultTTL
	}

	if ttl == 0 {
		ttl = defaultTTL
	}

	event := &domain.Event{
		ID:          xid.New().String(),
		PublishedAt: time.Now(),
		Data:        data,
		LivesUntil:  time.Now().Add(ttl),
	}

	c.events = append(c.events, event)
	return event
}

func (c *Channel) AddDelayedEvent(data string, publishAt time.Time, ttl time.Duration) *domain.Event {
	if ttl == 0 {
		ttl = defaultTTL
	}
	event := &domain.Event{
		ID:          xid.New().String(),
		PublishedAt: publishAt,
		Data:        data,
		LivesUntil:  publishAt.Add(ttl),
	}

	time.AfterFunc(publishAt.Sub(time.Now()), func() {
		c.lock.Lock()
		defer c.lock.Unlock()
		c.events = append(c.events, event)
	})

	return event
}

func (c *Channel) SetEventLock(eventID string, shouldBeLocked bool, ttl time.Duration) error {
	c.lock.RLock()
	defer c.lock.RUnlock()

	index := c.findEventIndex(eventID)
	if index == -1 {
		return fmt.Errorf("No event could be found with ID %s. Maybe it has already been picked up. ", eventID)
	}

	if shouldBeLocked {
		return c.events[index].Lock(ttl)
	}

	c.events[index].Unlock()
	return nil
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

// Please use the lock before calling the function!
func (c *Channel) getValidEvents() []*domain.Event {
	var validEvents []*domain.Event
	now := time.Now()

	for _, event := range c.events {
		if event.LivesUntil.After(now) {
			validEvents = append(validEvents, event)
		}
	}

	return validEvents
}

// Please use the lock before calling the function!
// Will return -1 if it could not find the event
func (c *Channel) findEventIndex(eventID string) int {
	var index = -1
	for i, e := range c.getValidEvents() {
		if e.ID == eventID {
			index = i
			break
		}
	}

	return index
}

func (c *Channel) cleanupInvalidEvents() {
	c.lock.Lock()
	c.events = c.getValidEvents()
	c.lock.Unlock()
}
