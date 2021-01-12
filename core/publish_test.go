package core

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func BenchmarkQueueWith1000Events(b *testing.B) {
	c := Init()
	channelName := "test-channel"

	amountElem := 1000
	for i := 0; i < amountElem; i++ {
		c.AddChannelEvent(channelName, "Test_data", 0, time.Now())
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.ChannelQueue(channelName)
	}
}

func BenchmarkPublishOneChannel(b *testing.B) {
	c := Init()
	channelName := "test-channel"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.AddChannelEvent(channelName, "Test_data", 0, time.Now())
	}
}

func BenchmarkMultiPublishAtOneChannel(b *testing.B) {
	c := Init()
	channelName := "test-channel"

	wg := sync.WaitGroup{}
	wg.Add(b.N)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go func() {
			c.AddChannelEvent(channelName, "Test_data", 0, time.Now())
			wg.Done()
		}()
	}
	wg.Wait()
}

func BenchmarkPublish2Channels(b *testing.B) {
	benchmarkPublishNChannels(2, b)
}

func BenchmarkPublish4Channels(b *testing.B) {
	benchmarkPublishNChannels(4, b)
}

func BenchmarkPublish10Channels(b *testing.B) {
	benchmarkPublishNChannels(10, b)
}

func benchmarkPublishNChannels(n int, b *testing.B) {
	c := Init()
	wg := sync.WaitGroup{}
	wg.Add(n)
	var channels []string
	for i := 0; i < n; i++ {
		channels = append(channels, fmt.Sprintf("test-%d", n))
	}
	everyChannelEvents := b.N / n

	b.ResetTimer()
	for _, channel := range channels {
		go func(channelName string) {
			for i := 0; i < everyChannelEvents; i++ {
				c.AddChannelEvent(channelName, "Test_data", 0, time.Now())
			}
			wg.Done()
		}(channel)
	}
	wg.Wait()
}
