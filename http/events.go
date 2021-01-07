package events_http

import (
	"distribuerad/interface"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"time"
)

func handleCreateNewEvent(store domain.IChannelStore) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		channel := resolveChannelName(store, p)

		var event struct {
			Data      string    `json:"data"`
			PublishAt time.Time `json:"publishAt,omitempty"`
			TTL       string    `json:"ttl,omitempty"`
		}
		if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
			log.Printf("Could not read the request data of new-event request: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		ttl, err := time.ParseDuration(event.TTL)
		if err != nil {
			errorResponse(w, fmt.Errorf("Parameter 'ttl' should be in a duration format (e.g. '1h30m10s')!"),
				http.StatusBadRequest)
		}

		var newEvent *domain.Event
		if event.PublishAt.IsZero() {
			newEvent = channel.AddEvent(event.Data, ttl)
		} else {
			newEvent = channel.AddDelayedEvent(event.Data, event.PublishAt, ttl)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(newEvent); err != nil {
			log.Printf("Error writing header with JSON data: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func handleGetAllEvents(store domain.IChannelStore) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		channel := resolveChannelName(store, p)

		eventList := channel.GetEvents()
		events := queue{
			ChannelName: p.ByName("channel-name"),
			Count:       len(eventList),
			Events:      eventList,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(events); err != nil {
			log.Printf("Error writing header with JSON data: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func handleDeleteEvent(store domain.IChannelStore) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		if err := resolveChannelName(store, p).DeleteEvent(p.ByName("event-id")); err != nil {
			errorResponse(w, err, http.StatusGone)
			return
		}

		w.WriteHeader(http.StatusAccepted)
	}
}
