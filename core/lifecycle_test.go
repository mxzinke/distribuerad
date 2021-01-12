package core

/*
func TestIntervalJobLifecycle(t *testing.T) {
	store := NewChannelStore()
	channel := store.AddChannel("test-channel")
	testJobName := "test-job123"

	// Starting state
	assert.Empty(t, channel.GetJobs(), "should not have any job in channel yet")

	// Adding a interval-job (every 1 second)
	job, err := channel.AddJob(testJobName, "EVENT_DATA", "@every 1s", 0)
	assert.Nil(t, err, "should not have an error")
	assert.NotNil(t, job, "should return pointer of job")

	assert.Len(t, channel.GetJobs(), 1, "should have one job in channel")

	// Check if the event has published
	assert.Empty(t, channel.GetEvents(), "should not have an event yet")
	time.Sleep(1*time.Second + 200*time.Millisecond)
	assert.Len(t, channel.GetEvents(), 1, "should have one published event (after 1 second)")

	// Removing a job
	err = channel.DeleteJob(testJobName)
	assert.Nil(t, err, "should not have an error")
	assert.Empty(t, channel.GetJobs(), "should not have any job in channel (after delete)")

	// Check that no more events have been published
	assert.Len(t, channel.GetEvents(), 1, "should still have just one event")
	time.Sleep(1*time.Second + 200*time.Millisecond)
	assert.Len(t, channel.GetEvents(), 1, "should still have just one event (after 1 second)")
}
*/
