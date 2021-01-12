package event

import (
	"errors"
	"fmt"
	"github.com/rs/xid"
	"sync"
	"time"
)

type Event struct {
	EventID     string    `json:"eventID"`
	PublishTime time.Time `json:"publishedAt"`
	Data        string    `json:"data"`
	LivesUntil  time.Time `json:"validUntil"`
	IsLocked    bool      `json:"isLocked"`

	// Private for managing IsLocked state
	L *sync.Mutex `json:"-"`
	T *time.Timer `json:"-"`
}

const DefaultTTL = 10 * time.Minute
const DefaultEventLockTTL = 1 * time.Minute

func New(data string, publishAt time.Time, ttl time.Duration) *Event {
	if ttl == 0 {
		ttl = DefaultTTL
	}
	if publishAt.IsZero() {
		publishAt = time.Now()
	}

	return &Event{
		EventID:     xid.New().String(),
		PublishTime: publishAt,
		Data:        data,
		LivesUntil:  publishAt.Add(ttl),
	}
}

func (e *Event) ID() string {
	return e.EventID
}
func (e *Event) ValidUntil() time.Time {
	return e.LivesUntil
}
func (e *Event) PublishedAt() time.Time {
	return e.PublishTime
}

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
