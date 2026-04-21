package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/google/uuid"

	"github.com/paddyoneill/slurm-ui/internal/server/handler"
	servertypes "github.com/paddyoneill/slurm-ui/internal/server/types"
)

const baseUrl = "/api"

func New(db *sql.DB) *http.Server {
	router := http.NewServeMux()
	target, _ := url.Parse("http://localhost:5173")
	viteProxy := httputil.NewSingleHostReverseProxy(target)
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		viteProxy.ServeHTTP(w, r)
	})

	h := handler.New(db)

	router.HandleFunc("GET "+baseUrl+"/jobs", h.GetJobs)
	router.HandleFunc("POST "+baseUrl+"/jobs", h.PostJobs)
	router.HandleFunc("GET "+baseUrl+"/jobs/events", h.GetJobsEvents)
	router.HandleFunc("GET "+baseUrl+"/jobs/{id}", func(w http.ResponseWriter, r *http.Request) {
		jobID, err := parseUUIDPathValue(r, "id")
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, fmt.Errorf("invalid job id: %w", err))
			return
		}

		h.GetJobsId(w, r, jobID)
	})
	router.HandleFunc("DELETE "+baseUrl+"/jobs/{id}", func(w http.ResponseWriter, r *http.Request) {
		jobID, err := parseUUIDPathValue(r, "id")
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, fmt.Errorf("invalid job id: %w", err))
			return
		}

		h.DeleteJobsId(w, r, jobID)
	})

	router.HandleFunc("GET "+baseUrl+"/notebooks", h.GetNotebooks)
	router.HandleFunc("POST "+baseUrl+"/notebooks", h.PostNotebooks)
	router.HandleFunc("GET "+baseUrl+"/notebooks/events", h.GetNotebooksEvents)
	router.HandleFunc("GET "+baseUrl+"/notebooks/{id}", func(w http.ResponseWriter, r *http.Request) {
		jobID, err := parseUUIDPathValue(r, "id")
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, fmt.Errorf("invalid job id: %w", err))
			return
		}

		h.GetNotebooksId(w, r, jobID)
	})
	router.HandleFunc("DELETE "+baseUrl+"/notebooks/{id}", func(w http.ResponseWriter, r *http.Request) {
		jobID, err := parseUUIDPathValue(r, "id")
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, fmt.Errorf("invalid job id: %w", err))
			return
		}

		h.DeleteNotebooksId(w, r, jobID)
	})

	for _, method := range []string{"DELETE", "GET", "PATCH", "POST", "PUT"} {
		router.HandleFunc(method+" "+baseUrl+"/notebooks/{id}/proxy", h.HandleNotebookProxy)
		router.HandleFunc(method+" "+baseUrl+"/notebooks/{id}/proxy/{path...}", h.HandleNotebookProxy)
	}

	go startJobPoller(db, h.NotifyJobsChanged)
	go startNotebookPoller(db, h.NotifyNotebooksChanged)

	return &http.Server{
		Addr:    "localhost:8080",
		Handler: router,
	}
}

func parseUUIDPathValue(r *http.Request, key string) (uuid.UUID, error) {
	return uuid.Parse(r.PathValue(key))
}

func writeJSONError(w http.ResponseWriter, status int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(servertypes.Error{Error: err.Error()})
}
