package events_http

import (
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

func StartHTTP(bindAddr string, store IChannelStore) {
	router := httprouter.New()

	router.POST("/events/:channelName", handleCreateNewEvent(store))
	router.GET("/events/:channelName", handleGetAllEvents(store))

	// TODO: Implementing the two left routes:
	// router.DELETE("/events/:channelName/:eventID", )

	log.Fatalln(http.ListenAndServe(bindAddr, router))
}

func resolveChannelName(store IChannelStore, params httprouter.Params) IChannel {
	channelName := params.ByName("channelName")
	channel := store.GetChannel(channelName)
	if channel == nil {
		channel = store.AddChannel(channelName)
	}

	return channel
}
