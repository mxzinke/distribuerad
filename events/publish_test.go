package events

import (
	"testing"
)

func BenchmarkPublishOneChannel(b *testing.B) {
	channelStore := NewChannelStore()
	channel := channelStore.AddChannel("test")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		channel.AddEvent("HELLO1234-Event")
	}
}
