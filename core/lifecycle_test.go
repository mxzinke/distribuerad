package core

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestIntervalJobLifecycle(t *testing.T) {
	c := Init()
	channelName := "test-channel"
	testJobName := "test-job123"

	// Starting state
	assert.Empty(t, c.ChannelJobs(channelName), "should not have any job in channel yet")

	// Adding a interval-job (every 1 second)
	job, err := c.AddChannelJob(channelName, testJobName, "EVENT_DATA", "@every 1s", 0)
	assert.Nil(t, err, "should not have an error")
	assert.NotNil(t, job, "should return pointer of job")

	assert.Len(t, c.ChannelJobs(channelName), 1, "should have one job in channel")

	// Check if the event has published
	assert.Empty(t, c.ChannelQueue(channelName), "should not have an event yet")
	time.Sleep(1*time.Second + 200*time.Millisecond)
	assert.Len(t, c.ChannelQueue(channelName), 1, "should have one published event (after 1 second)")

	// Removing a job
	err = c.DeleteChannelJob(channelName, testJobName)
	assert.Nil(t, err, "should not have an error")
	assert.Empty(t, c.ChannelJobs(channelName), "should not have any job in channel (after delete)")

	// Check that no more events have been published
	assert.Len(t, c.ChannelQueue(channelName), 1, "should still have just one event")
	time.Sleep(1*time.Second + 200*time.Millisecond)
	assert.Len(t, c.ChannelQueue(channelName), 1, "should still have just one event (after 1 second)")
}
