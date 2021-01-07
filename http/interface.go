package events_http

import (
	"distribuerad/interface"
)

type queue struct {
	ChannelName string          `json:"channelName"`
	Events      []*domain.Event `json:"events"`
}

type jobsList struct {
	ChannelName string        `json:"channelName"`
	Jobs        []*domain.Job `json:"jobs"`
}
