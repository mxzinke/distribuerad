package events_http

import (
	domain "distribuerad/core/event"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"time"
)

func StartHTTP(bindAddr string, store domain.IChannelStore) {
	router := httprouter.New()

	// Related to events
	router.GET("/:channel-name/events/", handleGetAllEvents(store))
	router.POST("/:channel-name/events/", handleCreateNewEvent(store))
	router.PATCH("/:channel-name/events/:event-id", handleLockEvent(store))
	router.DELETE("/:channel-name/events/:event-id", handleDeleteEvent(store))

	router.GET("/:channel-name/jobs", handleGetAllJobs(store))
	router.POST("/:channel-name/jobs", handleCreateNewJob(store))
	router.DELETE("/:channel-name/jobs/:job-id", handleDeleteJob(store))

	log.Printf("Make Pub/Sub-Server available at %s", bindAddr)
	log.Fatalln(http.ListenAndServe(bindAddr, router))
}

func resolveChannelName(store domain.IChannelStore, params httprouter.Params) domain.IChannel {
	return store.GetChannel(params.ByName("channel-name"))
}

func parseOptionalDuration(v string) (time.Duration, error) {
	var ttl time.Duration
	if v != "" {
		var err error
		ttl, err = time.ParseDuration(v)
		if err != nil {
			return ttl, errors.New(fmt.Sprintln("Parameter 'ttl' should be in a duration format (e.g. '1h30m10s')!"))
		}
	}

	return ttl, nil
}

func errorResponse(w http.ResponseWriter, message string, status int) {
	w.WriteHeader(status)
	if _, err := w.Write([]byte(message)); err != nil {
		log.Printf("Error at writing error message (status: %d): %v", status, err)
	}
}
