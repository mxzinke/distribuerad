package store

import (
	"sync"
)

type IChannel interface {
	GetName() string

	// Reactions
	OnAttach()
	OnDetach()
}

type ChannelStore struct {
	channels map[string]IChannel
	lock     *sync.Mutex
}

func New() *ChannelStore {
	return &ChannelStore{
		lock:     &sync.Mutex{},
		channels: make(map[string]IChannel),
	}
}

// Does return an IChannel which was found by name
func (store *ChannelStore) GetChannel(name string) IChannel {
	return store.channels[name]
}

// Adds an IChannel to the internal store (by name). If the name is already reserved, the saved
// IChannel is returned. While action is performed, no other AddChannel or DeleteChannel is possible.
//
// The IChannel OnAttach() function is called, after it was stored.
func (store *ChannelStore) AddChannel(channel IChannel) IChannel {
	store.lock.Lock()
	defer store.lock.Unlock()

	name := channel.GetName()
	if store.channels[name] != nil {
		return store.channels[name]
	}
	store.channels[name] = channel

	store.channels[name].OnAttach()

	return store.channels[name]

}

// Removes an IChannel by his name from the internal store. If no IChannel was found by the given name,
// no action will be performed. While action is performed, no other AddChannel or DeleteChannel is possible.
//
// The IChannel OnDetach() function is called, before it will get deleted.
func (store *ChannelStore) DeleteChannel(name string) {
	if store.channels[name] == nil {
		return
	}

	store.lock.Lock()
	defer store.lock.Unlock()

	store.channels[name].OnDetach()
	delete(store.channels, name)
}
