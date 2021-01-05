package events_http

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"time"
)

func handleCreateNewEvent(store IChannelStore) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		channel := resolveChannelName(store, p)

		var reqData []byte
		if _, err := r.Body.Read(reqData); err != nil {
			log.Println("Could not read the request data of new-event request.")
			w.WriteHeader(500)
			return
		}

		var event struct {
			Data      string    `json:"data"`
			PublishAt time.Time `json:"timestamp"`
		}
		if err := json.Unmarshal(reqData, &event); err != nil {
			log.Printf("Error parsing JSON from request data: %v", err)
			w.WriteHeader(500)
			return
		}

		var newEvent *Event
		if event.PublishAt.IsZero() {
			newEvent = channel.AddEvent(event.Data)
		} else {
			newEvent = channel.AddDelayedEvent(event.Data, event.PublishAt)
		}

		marshaled, err := json.Marshal(newEvent)
		if err != nil {
			log.Printf("Error stringify JSON from response data: %v", err)
			w.WriteHeader(500)
			return
		}

		w.WriteHeader(201)
		if _, err := w.Write(marshaled); err != nil {
			log.Printf("Error writing header: %v", err)
		}
	}
}

func handleGetAllEvents(store IChannelStore) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		channel := resolveChannelName(store, p)

		events := eventList{
			ChannelName: p.ByName("channelName"),
			Events:      channel.GetEvents(),
		}

		marshaled, err := json.Marshal(events)
		if err != nil {
			log.Printf("Error stringify JSON from response data: %v", err)
			w.WriteHeader(500)
			return
		}

		w.WriteHeader(200)
		if _, err := w.Write(marshaled); err != nil {
			log.Printf("Error writing header: %v", err)
		}
	}
}
