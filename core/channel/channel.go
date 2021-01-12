package channel

import (
	"github.com/robfig/cron/v3"
	"log"
	"sync"
	"time"
)

type Channel struct {
	name     string
	events   []IEvent
	jobs     map[string]IJob
	lock     *sync.RWMutex
	jobsLock *sync.RWMutex
	cleanup  *cron.Cron
}

func New(name string) *Channel {
	return &Channel{
		name:     name,
		jobs:     map[string]IJob{},
		lock:     &sync.RWMutex{},
		jobsLock: &sync.RWMutex{},
	}
}

func (c *Channel) Name() string {
	return c.name
}

// An action which is performed after adding to a store
func (c *Channel) OnAttach() {
	// Startup: Cleanup Background task
	c.cleanup = cron.New(cron.WithChain(
		cron.Recover(cron.DefaultLogger),
	))
	if _, err := c.cleanup.AddFunc(
		"@every 10m",
		c.cleanupInvalidEvents,
	); err != nil {
		log.Printf("Error: The cleanup background task could not be started.")
	}
	c.cleanup.Start()
}

// An action which is performed before removing from a store
func (c *Channel) OnDetach() {
	if c.cleanup != nil {
		c.cleanup.Stop()
	}
}

// --- Private ---

// Please use the lock before calling the function!
func (c *Channel) getValidEvents() []IEvent {
	var validEvents []IEvent
	now := time.Now()

	for _, event := range c.events {
		if event.ValidUntil().After(now) {
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
		if e.ID() == eventID {
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
