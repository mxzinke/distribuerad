package events_http

import (
	"distribuerad/core"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"time"
)

func handleCreateNewEvent(store domain.IChannelStore) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		channel := resolveChannelName(store, p)

		var payload struct {
			Data      string    `json:"data"`
			PublishAt time.Time `json:"publishAt,omitempty"`
			TTL       string    `json:"ttl,omitempty"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			log.Printf("Could not read the request data of new-event request: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		ttl, err := parseOptionalDuration(payload.TTL)
		if err != nil {
			errorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}

		var newEvent *domain.Event
		if payload.PublishAt.IsZero() {
			newEvent = channel.AddEvent(payload.Data, ttl)
		} else {
			newEvent = channel.AddDelayedEvent(payload.Data, payload.PublishAt, ttl)
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

const (
	lockAction   = "LOCK"
	unlockAction = "UNLOCK"
)

func handleLockEvent(store domain.IChannelStore) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		channel := resolveChannelName(store, p)

		var payload struct {
			Action string `json:"action"`
			TTL    string `json:"ttl,omitempty"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			log.Printf("Could not read the request data of new-event request: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if payload.Action != lockAction && payload.Action != unlockAction {
			errorResponse(w, "Parameter 'action' can only have values 'LOCK' or 'UNLOCK'!",
				http.StatusBadRequest)
			return
		}

		ttl, err := parseOptionalDuration(payload.TTL)
		if err != nil {
			errorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := channel.SetEventLock(p.ByName("event-id"), payload.Action == lockAction, ttl); err != nil {
			errorResponse(w, err.Error(), http.StatusGone)
			return
		}

		w.WriteHeader(http.StatusAccepted)
	}
}

func handleDeleteEvent(store domain.IChannelStore) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		if err := resolveChannelName(store, p).DeleteEvent(p.ByName("event-id")); err != nil {
			errorResponse(w, err.Error(), http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusAccepted)
	}
}
