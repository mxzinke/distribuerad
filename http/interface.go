package events_http

import "time"

type IChannelStore interface {
	GetChannel(name string) IChannel
	AddChannel(name string) IChannel
	DeleteChannel(name string)
}

type IChannel interface {
	GetEvents() []*Event
	AddEvent(data string) *Event
	AddDelayedEvent(data string, publishAt time.Time) *Event
}

type Event struct {
	ID          string    `json:"eventID"`
	PublishedAt time.Time `json:"publishedAt"`
	Data        string    `json:"data"`
}

type eventList struct {
	ChannelName string   `json:"channelName"`
	Events      []*Event `json:"events"`
}
