package events

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTTLOfEvent(t *testing.T) {
	store := NewChannelStore()
	channel := store.AddChannel("test")
	ttl := 200 * time.Microsecond

	e := channel.AddEvent("SOME_DATA", ttl)

	normalExecTime := 10 * time.Microsecond
	assert.WithinDuration(t, e.LivesUntil, time.Now().Add(ttl), normalExecTime,
		"should have valid LivesUntil value")
	assert.True(t, e.LivesUntil.After(e.PublishedAt), "should have relational correct data")
	assert.Len(t, channel.GetEvents(), 1, "should not be empty")

	time.Sleep(ttl)
	assert.Len(t, channel.GetEvents(), 0, "should be gone")
	assert.Error(t, channel.DeleteEvent(e.ID), "should throw error, when try to delete timed out event")
}
