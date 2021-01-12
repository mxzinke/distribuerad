package job

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestDoesCallTrigger(t *testing.T) {
	testData := "TEST_DATA"
	testTTL := time.Minute
	var hasBeenCalled bool
	trigger := func(data string, ttl time.Duration) {
		hasBeenCalled = true
		assert.Equal(t, data, testData, "should have correct data")
		assert.Equal(t, ttl, testTTL, "should have correct ttl")
	}

	job := &Job{
		JobID:   "test-job",
		CronDef: "@every 1s", // milliseconds are not possible with the underlying framework
		Data:    testData,
		TTL:     testTTL,
	}
	job.SetTrigger(trigger)
	// Simulating that the job is attached to a channel
	job.OnAttach()

	assert.False(t, hasBeenCalled, "should not have called trigger")
	time.Sleep(1 * time.Second)
	assert.True(t, hasBeenCalled, "should have called trigger")
}
