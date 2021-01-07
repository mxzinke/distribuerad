package events_http

import (
	domain "distribuerad/interface"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

func StartHTTP(bindAddr string, store domain.IChannelStore) {
	router := httprouter.New()

	// Related to events
	router.POST("/:channel-name/events/", handleCreateNewEvent(store))
	router.GET("/:channel-name/events/", handleGetAllEvents(store))
	router.DELETE("/:channel-name/events/:event-id", handleDeleteEvent(store))

	router.POST("/:channel-name/jobs", handleCreateNewJob(store))
	router.GET("/:channel-name/jobs", handleGetAllJobs(store))
	router.DELETE("/:channel-name/jobs/:job-id", handleDeleteJob(store))

	log.Printf("Make Pub/Sub-Server available at %s", bindAddr)
	log.Fatalln(http.ListenAndServe(bindAddr, router))
}

func resolveChannelName(store domain.IChannelStore, params httprouter.Params) domain.IChannel {
	return store.GetChannel(params.ByName("channel-name"))
}

func errorResponse(w http.ResponseWriter, err error, status int) {
	w.WriteHeader(status)
	if _, err := w.Write([]byte(err.Error())); err != nil {
		log.Printf("Error at writing error message (status: %d): %v", status, err)
	}
}
