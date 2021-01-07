package events

import (
	"distribuerad/interface"
	"fmt"
	"github.com/rs/xid"
	"sync"
	"time"
)

type Channel struct {
	events   []*domain.Event
	jobs     map[string]*jobExecution
	lock     *sync.RWMutex
	jobsLock *sync.RWMutex
}

const defaultTTL = 10 * time.Minute

func (c *Channel) GetEvents() []*domain.Event {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.events
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
