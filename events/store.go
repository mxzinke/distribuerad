package events

import (
	domain "distribuerad/interface"
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

	return store.channels[name]
}

func (store *ChannelStore) DeleteChannel(name string) {
	store.lock.Lock()
	defer store.lock.Unlock()

	delete(store.channels, name)
}
