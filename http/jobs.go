package events_http

import (
	domain "distribuerad/interface"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

func handleCreateNewJob(store domain.IChannelStore) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		channel := resolveChannelName(store, p)

		var job struct {
			Name    string `json:"name"`
			Data    string `json:"data"`
			CronDef string `json:"cronDef"`
		}
		if err := json.NewDecoder(r.Body).Decode(&job); err != nil {
			log.Printf("Could not read the request data of new-event request: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Validate
		if job.Name == "" || job.CronDef == "" {
			errorResponse(w, fmt.Errorf("Parameters 'name' and 'cronDef' has to be set!"),
				http.StatusBadRequest)
			return
		}

		newJob, err := channel.AddJob(job.Name, job.Data, job.CronDef)
		if err != nil {
			errorResponse(w, err, http.StatusConflict)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(newJob); err != nil {
			log.Printf("Error writing body with JSON data: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func handleGetAllJobs(store domain.IChannelStore) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		channel := resolveChannelName(store, p)

		events := jobsList{
			ChannelName: p.ByName("channelName"),
			Jobs:        channel.GetJobs(),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(events); err != nil {
			log.Printf("Error writing body with JSON data: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func handleDeleteJob(store domain.IChannelStore) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		if err := resolveChannelName(store, p).DeleteJob(p.ByName("job-id")); err != nil {
			errorResponse(w, err, http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusAccepted)
	}
}
