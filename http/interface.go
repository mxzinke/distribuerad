package events_http

import (
	"distribuerad/core/event"
)

type queue struct {
	ChannelName string          `json:"channelName"`
	Count       int             `json:"count"`
	Events      []*domain.Event `json:"events"`
}

type jobsList struct {
	ChannelName string          `json:"channelName"`
	Jobs        []*event.domain `json:"jobs"`
}
