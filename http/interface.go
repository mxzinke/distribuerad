package events_http

import (
	"distribuerad/interface"
)

type eventList struct {
	ChannelName string          `json:"channelName"`
	Events      []*domain.Event `json:"events"`
}
