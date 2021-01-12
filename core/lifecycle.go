package core

import "distribuerad/core/store"

type Core struct {
	store *store.ChannelStore
}

func Init() *Core {
	return &Core{
		store: store.New(),
	}
}

// Eventing:
// func ChannelQueue(channelName string) []event.Event
// TODO: returning error! (implement at event package)
// func AddChannelEvent(channelName, data string, ttl time.Duration, publishedAt time.Time) error
// func DeleteChannelEvent(channelName, eventID string) error

// Event Locking:
// func LockChannelEvent(channelName, eventID string, ttl time.Duration) error
// func UnlockChannelEvent(channelName, eventID string)

// Jobs:
// func ChannelJobs(channelName string) []job.Job
// TODO: returning error! (implement at job package)
// func AddChannelJob(jobName, data, cronDef string, ttl time.Duration) error
// func DeleteChannelJob(jobName string) error
