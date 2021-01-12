package channel

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type fakeEvent struct {
	EventID     string    `json:"eventID"`
	PublishTime time.Time `json:"publishedAt"`
	Data        string    `json:"data"`
	LivesUntil  time.Time `json:"validUntil"`
}

func (f *fakeEvent) ID() string {
	return f.EventID
}
func (f *fakeEvent) ValidUntil() time.Time {
	return f.LivesUntil
}
func (f *fakeEvent) PublishedAt() time.Time {
	return f.PublishTime
}

func TestTTLOfEvent(t *testing.T) {
	channel := New("test-channel")
	ttl := 200 * time.Microsecond

	e := channel.AddEvent(&fakeEvent{
		EventID:     "21325235",
		PublishTime: time.Now(),
		LivesUntil:  time.Now().Add(ttl),
	})
	assert.Len(t, channel.GetEvents(), 1, "should not be empty")

	time.Sleep(ttl)
	assert.Len(t, channel.GetEvents(), 0, "should be gone")
	assert.Error(t, channel.DeleteEvent(e.ID()), "should throw error, when try to delete timed out event")
}
