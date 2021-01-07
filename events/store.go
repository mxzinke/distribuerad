package events

import (
	domain "distribuerad/interface"
	"github.com/robfig/cron/v3"
	"log"
	"sync"
)

type ChannelStore struct {
	channels map[string]*Channel
	lock     *sync.Mutex
}

func NewChannelStore() domain.IChannelStore {
	return &ChannelStore{
		lock:     &sync.Mutex{},
		channels: map[string]*Channel{},
	}
}

// Returns a channel, if the channel does not exist yet, it will create one
func (store *ChannelStore) GetChannel(name string) domain.IChannel {
	if store.channels[name] != nil {
		return store.channels[name]
	}

	return store.AddChannel(name)
}

func (store *ChannelStore) AddChannel(name string) domain.IChannel {
	store.lock.Lock()
	defer store.lock.Unlock()

	if store.channels[name] != nil {
		return store.channels[name]
	}
	store.channels[name] = &Channel{
		lock:     &sync.RWMutex{},
		jobs:     map[string]*jobExecution{},
		jobsLock: &sync.RWMutex{},
	}

	// Startup: Cleanup Background task
	store.channels[name].cleanup = cron.New(cron.WithChain(
		cron.Recover(cron.DefaultLogger),
	))
	if _, err := store.channels[name].cleanup.AddFunc(
		"@every 10m",
		store.channels[name].cleanupInvalidEvents,
	); err != nil {
		log.Printf("Error: The cleanup background task could not be started.")
	}
	store.channels[name].cleanup.Start()

	return store.channels[name]
}

func (store *ChannelStore) DeleteChannel(name string) {
	store.lock.Lock()
	defer store.lock.Unlock()

	store.channels[name].cleanup.Stop()
	delete(store.channels, name)
}
