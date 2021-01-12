package core

/*
func BenchmarkPublishOneChannel(b *testing.B) {
	channelStore := NewChannelStore()
	channel := channelStore.AddChannel("test")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		channel.AddEvent("HELLO1234-Event", 0)
	}
}

func BenchmarkMultiPublishAtOneChannel(b *testing.B) {
	channelStore := NewChannelStore()
	channel := channelStore.AddChannel("test")

	wg := sync.WaitGroup{}
	wg.Add(b.N)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go func() {
			channel.AddEvent("HELLO1234-Event", 0)
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
	channelStore := NewChannelStore()
	wg := sync.WaitGroup{}
	wg.Add(n)
	var channels []domain.IChannel
	for i := 0; i < n; i++ {
		channels = append(channels, channelStore.AddChannel(fmt.Sprintf("test-%d", n)))
	}

	everyChannelEvents := b.N / n

	b.ResetTimer()
	for _, channel := range channels {
		go benchmarkAtChannel(channel, &wg, everyChannelEvents)
	}
	wg.Wait()
}

func benchmarkAtChannel(channel domain.IChannel, wg *sync.WaitGroup, nEvents int) {
	for i := 0; i < nEvents; i++ {
		channel.AddEvent("HELLO1234-Event", 0)
	}
	wg.Done()
}
*/
