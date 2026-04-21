package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"

	"github.com/paddyoneill/slurm-ui/internal/db"
	servertypes "github.com/paddyoneill/slurm-ui/internal/server/types"
	"github.com/paddyoneill/slurm-ui/internal/slurm"
)

func (h *Handler) DeleteNotebooksId(w http.ResponseWriter, r *http.Request, id servertypes.NotebookID) {
	record, err := db.GetNotebookByID(r.Context(), h.db, id)
	if err != nil {
		if errors.Is(err, db.ErrNotebookNotFound) {
			writeError(w, http.StatusNotFound, err)
			return
		}
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	if err := slurm.CancelJob(r.Context(), record.SlurmJobID); err != nil {
		writeError(w, http.StatusBadGateway, err)
		return
	}

	if err := db.DeleteNotebookByID(r.Context(), h.db, id); err != nil {
		if errors.Is(err, db.ErrNotebookNotFound) {
			writeError(w, http.StatusNotFound, err)
			return
		}
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	h.NotifyNotebooksChanged()
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetNotebooks(w http.ResponseWriter, r *http.Request) {
	records, err := db.ListNotebooks(r.Context(), h.db)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	notebooks := make([]servertypes.Notebook, 0)
	for _, record := range records {
		notebooks = append(notebooks, notebookToAPI(record))
	}

	writeJSON(w, http.StatusOK, notebooks)
}

func (h *Handler) GetNotebooksId(w http.ResponseWriter, r *http.Request, id servertypes.JobID) {
	record, err := db.GetNotebookByID(r.Context(), h.db, id)
	if err != nil {
		if errors.Is(err, db.ErrNotebookNotFound) {
			writeError(w, http.StatusNotFound, err)
			return
		}
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	notebook := notebookToAPI(record)
	writeJSON(w, http.StatusOK, notebook)
}

func (h *Handler) PostNotebooks(w http.ResponseWriter, r *http.Request) {
	req, validationErr := decodeCreateNotebookRequest(r)
	if validationErr != nil {
		writeJSON(w, http.StatusBadRequest, validationErr)
		return
	}

	notebookID := uuid.New()
	token := uuid.NewString()
	script := buildNotebookScript()
	environment := buildNotebookEnvironment(notebookID, token, req)

	jobReq := buildSubmitJobRequest(req.LaunchConfig, script, &environment)

	slurmJobID, err := slurm.SubmitJob(r.Context(), jobReq)
	if err != nil {
		writeError(w, http.StatusBadGateway, fmt.Errorf("submit notebook: %w", err))
		return
	}

	jobInfo, err := slurm.GetJob(r.Context(), slurmJobID)
	if err != nil {
		writeError(w, http.StatusBadGateway, fmt.Errorf("get notebook status: %w", err))
		return
	}

	notebookName := req.Name
	if jobInfo.Name != nil && *jobInfo.Name != "" {
		notebookName = *jobInfo.Name
	}

	record, err := db.CreateNotebook(r.Context(), h.db, notebookID, slurmJobID, notebookName, jobInfo.State, jobInfo.Host, notebookPort, token, *req.BaseEnv)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	response := notebookToAPI(record)
	response.Partition = req.Partition
	response.CurrentWorkingDirectory = &req.CurrentWorkingDirectory
	response.Environment = req.Environment
	response.CpusPerTask = req.CpusPerTask
	response.TasksPerNode = req.TasksPerNode
	response.MemoryPerNode = req.MemoryPerNode
	response.TimeLimit = req.TimeLimit

	h.NotifyNotebooksChanged()
	writeJSON(w, http.StatusCreated, response)
}
