package events_http

import (
	domain "distribuerad/interface"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

func StartHTTP(bindAddr string, store domain.IChannelStore) {
	router := httprouter.New()

	router.POST("/events/:channelName", handleCreateNewEvent(store))
	router.GET("/events/:channelName", handleGetAllEvents(store))
	router.DELETE("/events/:channelName/:event-id", handleDeleteEvent(store))

	log.Printf("Make Pub/Sub-Server available at %s", bindAddr)
	log.Fatalln(http.ListenAndServe(bindAddr, router))
}

func resolveChannelName(store domain.IChannelStore, params httprouter.Params) domain.IChannel {
	channelName := params.ByName("channelName")
	return store.GetChannel(channelName)
}
