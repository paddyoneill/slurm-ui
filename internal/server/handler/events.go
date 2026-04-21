package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/paddyoneill/slurm-ui/internal/db"
	servertypes "github.com/paddyoneill/slurm-ui/internal/server/types"
)

type listener chan struct{}

func (h *Handler) GetJobsEvents(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	l := h.subsubscribeJobs()
	defer h.unsubscribeJobs(l)

	if err := h.writeJobsEvent(r.Context(), w); err != nil {
		http.Error(w, "failed to write initial event", http.StatusInternalServerError)
		return
	}

	flusher.Flush()

	for {
		select {
		case <-r.Context().Done():
			return
		case <-l:
			if err := h.writeJobsEvent(r.Context(), w); err != nil {
				return
			}
			flusher.Flush()
		}
	}
}

func (h *Handler) subsubscribeJobs() listener {
	l := make(listener, 1)

	h.jobListenersMu.Lock()
	h.jobListeners[l] = struct{}{}
	h.jobListenersMu.Unlock()

	return l
}

func (h *Handler) unsubscribeJobs(l listener) {
	h.jobListenersMu.Lock()
	delete(h.jobListeners, l)
	h.jobListenersMu.Unlock()

	close(l)
}

func (h *Handler) NotifyJobsChanged() {
	h.jobListenersMu.Lock()
	defer h.jobListenersMu.Unlock()

	for l := range h.jobListeners {
		select {
		case l <- struct{}{}:
		default:
		}
	}
}

func (h *Handler) writeJobsEvent(ctx context.Context, w http.ResponseWriter) error {
	records, err := db.ListJobs(ctx, h.db)
	if err != nil {
		return err
	}

	jobs := make([]servertypes.Job, 0, len(records))
	for _, record := range records {
		jobs = append(jobs, jobToAPI(record))
	}

	payload, err := json.Marshal(jobs)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(w, "data: %s\n\n", payload)
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) GetNotebooksEvents(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	l := h.subsubscribeNotebooks()
	defer h.unsubscribeNotebooks(l)

	if err := h.writeNotebooksEvent(r.Context(), w); err != nil {
		http.Error(w, "failed to write initial event", http.StatusInternalServerError)
		return
	}

	flusher.Flush()

	for {
		select {
		case <-r.Context().Done():
			return
		case <-l:
			if err := h.writeNotebooksEvent(r.Context(), w); err != nil {
				return
			}
			flusher.Flush()
		}
	}
}

func (h *Handler) subsubscribeNotebooks() listener {
	l := make(listener, 1)

	h.notebookListenersMu.Lock()
	h.notebookListeners[l] = struct{}{}
	h.notebookListenersMu.Unlock()

	return l
}

func (h *Handler) unsubscribeNotebooks(l listener) {
	h.notebookListenersMu.Lock()
	delete(h.notebookListeners, l)
	h.notebookListenersMu.Unlock()

	close(l)
}

func (h *Handler) NotifyNotebooksChanged() {
	h.notebookListenersMu.Lock()
	defer h.notebookListenersMu.Unlock()

	for l := range h.notebookListeners {
		select {
		case l <- struct{}{}:
		default:
		}
	}
}

func (h *Handler) writeNotebooksEvent(ctx context.Context, w http.ResponseWriter) error {
	records, err := db.ListNotebooks(ctx, h.db)
	if err != nil {
		return err
	}

	notebooks := make([]servertypes.Notebook, 0, len(records))
	for _, record := range records {
		notebooks = append(notebooks, notebookToAPI(record))
	}

	payload, err := json.Marshal(notebooks)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(w, "data: %s\n\n", payload)
	if err != nil {
		return err
	}

	return nil
}
