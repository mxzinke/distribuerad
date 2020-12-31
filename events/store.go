package events

import (
	"sync"
)

type ChannelStore struct {
	channels map[string]*Channel
	lock     *sync.Mutex
}

func NewChannelStore() *ChannelStore {
	return &ChannelStore{
		lock: &sync.Mutex{},
	}
}

func (store *ChannelStore) FindChannel(name string) *Channel {
	return store.channels[name]
}

func (store *ChannelStore) AddChannel(name string) *Channel {
	channel := &Channel{
		lock: &sync.RWMutex{},
	}

	store.lock.Lock()
	store.channels[name] = channel
	store.lock.Unlock()

	return channel
}

func (store *ChannelStore) DeleteChannel(name string) {
	store.lock.Lock()
	delete(store.channels, name)
	store.lock.Unlock()
}
