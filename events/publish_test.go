package events

import (
	domain "distribuerad/interface"
	"fmt"
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

func BenchmarkPublish4Channels(b *testing.B) {
	benchmarkPublishNChannels(4, b)
}

func BenchmarkPublish10Channels(b *testing.B) {
	benchmarkPublishNChannels(10, b)
}

func benchmarkPublishNChannels(n int, b *testing.B) {
	channelStore := NewChannelStore()
	onFinish := make(chan int, n+1)
	var channels []domain.IChannel
	for i := 0; i < n; i++ {
		channels = append(channels, channelStore.AddChannel(fmt.Sprintf("test-%d", n)))
	}

	everyChannelEvents := b.N / n

	b.ResetTimer()
	for _, channel := range channels {
		go benchmarkAtChannel(channel, onFinish, everyChannelEvents)
	}
	for i := 0; i < n; i++ {
		<-onFinish
	}
}

func benchmarkAtChannel(channel domain.IChannel, onFinish chan int, nEvents int) {
	for i := 0; i < nEvents; i++ {
		channel.AddEvent("HELLO1234-Event")
	}
	onFinish <- 1
}
