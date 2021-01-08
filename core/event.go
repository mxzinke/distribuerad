package domain

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type Event struct {
	ID          string    `json:"eventID"`
	PublishedAt time.Time `json:"publishedAt"`
	Data        string    `json:"data"`
	LivesUntil  time.Time `json:"livesUntil"`
	IsLocked    bool      `json:"isLocked"`

	// Private for managing IsLocked state
	L *sync.Mutex `json:"-"`
	T *time.Timer `json:"-"`
}

const DefaultEventLockTTL = 1 * time.Minute

func (e *Event) Lock(ttl time.Duration) error {
	L := e.createL()
	L.Lock()
	defer L.Unlock()

	if e.IsLocked {
		return errors.New(fmt.Sprintln("The lock is already set."))
	}
	e.IsLocked = true

	// setup ttl timer
	if e.T != nil {
		e.T.Stop()
	}
	if ttl == 0 {
		ttl = DefaultEventLockTTL
	}
	e.T = time.AfterFunc(ttl, e.Unlock)

	return nil
}

func (e *Event) Unlock() {
	L := e.createL()
	L.Lock()
	defer L.Unlock()

	e.IsLocked = false
	e.L = nil

	if e.T != nil {
		e.T.Stop()
	}
	e.T = nil
}

func (e *Event) createL() *sync.Mutex {
	if e.L == nil {
		e.L = &sync.Mutex{}
	}
	return e.L
}
