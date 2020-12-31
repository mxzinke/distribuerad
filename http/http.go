package events_http

import (
	"../events"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"time"
)

func StartHTTP(bindAddr string, store *events.ChannelStore) {
	router := httprouter.New()

	router.POST("/events/:channelName", CreateNewEvent(store))

	// TODO: Implementing the two left routes:
	// router.GET("/events/:channelName", )
	// router.DELETE("/events/:channelName/:eventID", )

	log.Fatalln(http.ListenAndServe(bindAddr, router))
}

func CreateNewEvent(store *events.ChannelStore) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		channelName := p.ByName("channelName")
		channel := store.FindChannel(channelName)
		if channel == nil {
			channel = store.AddChannel(channelName)
		}

		var reqData []byte
		if _, err := r.Body.Read(reqData); err != nil {
			log.Println("Could not read the request data of new-event request.")
			w.WriteHeader(500)
			return
		}

		var event struct {
			Data      string    `json:"data"`
			Timestamp time.Time `json:"timestamp"`
		}
		if err := json.Unmarshal(reqData, &event); err != nil {
			log.Printf("Error parsing JSON from request data: %v", err)
			w.WriteHeader(500)
			return
		}

		if event.Timestamp.IsZero() {
			channel.AddEvent(event.Data, time.Now())
		} else {
			channel.AddEvent(event.Data, event.Timestamp)
		}
	}
}
