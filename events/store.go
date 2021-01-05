package events

import (
	"sync"
)

type ChannelStore struct {
	channels map[string]*Channel
	lock     *sync.RWMutex
}

func NewChannelStore() *ChannelStore {
	return &ChannelStore{
		lock: &sync.RWMutex{},
	}
}

func (store *ChannelStore) GetChannel(name string) *Channel {
	store.lock.RLock()
	defer store.lock.RUnlock()

	return store.channels[name]
}

func (store *ChannelStore) AddChannel(name string) *Channel {
	store.lock.Lock()
	defer store.lock.Unlock()

	if store.channels[name] == nil {
		return store.channels[name]
	}
	store.channels[name] = &Channel{
		lock: &sync.RWMutex{},
	}

	return store.channels[name]
}

func (store *ChannelStore) DeleteChannel(name string) {
	store.lock.Lock()
	defer store.lock.Unlock()

	delete(store.channels, name)
}
