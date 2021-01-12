package core

import (
	"distribuerad/core/channel"
	"distribuerad/core/event"
	"distribuerad/core/job"
	"distribuerad/core/store"
	"time"
)

type Core struct {
	store *store.ChannelStore
}

func Init() *Core {
	return &Core{
		store: store.New(),
	}
}

// --- Event related ---

func (c *Core) ChannelQueue(channelName string) []*event.Event {
	elms := c.findChannel(channelName).GetEvents()

	var castedElms = make([]*event.Event, len(elms))
	for i := range elms {
		castedElms[i] = elms[i].(*event.Event)
	}

	return castedElms
}

func (c *Core) AddChannelEvent(channelName, data string, ttl time.Duration, publishedAt time.Time) *event.Event {
	return c.findChannel(channelName).AddEvent(event.New(data, publishedAt, ttl)).(*event.Event)
}

func (c *Core) DeleteChannelEvent(channelName, eventID string) error {
	return c.findChannel(channelName).DeleteEvent(eventID)
}

// Acquires a lock at the given eventID, no other client can lock it too. The lock will be living only
// for the given Time-to-Live!
//
// Casting to *event.Event is needed, as *channel.Channel implementation uses the channel.IEvent interface
func (c *Core) LockChannelEvent(channelName, eventID string, ttl time.Duration) error {
	return (c.findChannel(channelName).GetEvent(eventID).(*event.Event)).Lock(ttl)
}

// Undoes the lock of LockChannelEvent function. Any client could now re-acquire a new lock.
func (c *Core) UnlockChannelEvent(channelName, eventID string) {
	(c.findChannel(channelName).GetEvent(eventID).(*event.Event)).Unlock()
}

// --- Jobs related ---

func (c *Core) ChannelJobs(channelName string) []*job.Job {
	elms := c.findChannel(channelName).GetJobs()

	var castedElms = make([]*job.Job, len(elms))
	for i := range elms {
		castedElms[i] = elms[i].(*job.Job)
	}

	return castedElms
}

func (c *Core) AddChannelJob(channelName, jobID, data, cronDef string, eventTTL time.Duration) (*job.Job, error) {
	j, err := c.findChannel(channelName).AddJob(job.New(jobID, data, cronDef, eventTTL))
	return j.(*job.Job), err
}

func (c *Core) DeleteChannelJob(channelName, jobID string) error {
	return c.findChannel(channelName).DeleteJob(jobID)
}

// --- Private Functions ---

// Does search for channel in the local store. Creates a new, if non was found.
// Needs to re-cast the store.IChannel interface to *channel.Channel, because it does not know fully about it.
func (c *Core) findChannel(channelName string) *channel.Channel {
	chen := c.store.GetChannel(channelName)
	if chen != nil {
		return chen.(*channel.Channel)
	}

	return c.store.AddChannel(channel.New(channelName)).(*channel.Channel)
}
