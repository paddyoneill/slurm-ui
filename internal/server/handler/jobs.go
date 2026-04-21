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

func (h *Handler) DeleteJobsId(w http.ResponseWriter, r *http.Request, id servertypes.JobID) {
	job, err := db.GetJobByID(r.Context(), h.db, id)
	if err != nil {
		if errors.Is(err, db.ErrJobNotFound) {
			writeError(w, http.StatusNotFound, err)
			return
		}
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	if err := slurm.CancelJob(r.Context(), job.SlurmJobId); err != nil {
		writeError(w, http.StatusBadGateway, err)
		return
	}

	if err := db.DeleteJobByID(r.Context(), h.db, id); err != nil {
		if errors.Is(err, db.ErrJobNotFound) {
			writeError(w, http.StatusNotFound, err)
			return
		}
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	h.NotifyJobsChanged()
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetJobs(w http.ResponseWriter, r *http.Request) {
	records, err := db.ListJobs(r.Context(), h.db)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	jobs := make([]servertypes.Job, 0, len(records))
	for _, record := range records {
		jobs = append(jobs, jobToAPI(record))
	}

	writeJSON(w, http.StatusOK, jobs)
}

func (h *Handler) GetJobsId(w http.ResponseWriter, r *http.Request, id servertypes.JobID) {
	record, err := db.GetJobByID(r.Context(), h.db, id)
	if err != nil {
		if errors.Is(err, db.ErrJobNotFound) {
			writeError(w, http.StatusNotFound, err)
			return
		}
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	job := jobToAPI(record)
	writeJSON(w, http.StatusOK, job)
}

func (h *Handler) PostJobs(w http.ResponseWriter, r *http.Request) {
	req, validationErr := decodeCreateJobRequest(r)
	if validationErr != nil {
		writeJSON(w, http.StatusBadRequest, validationErr)
		return
	}

	slurmJobID, err := slurm.SubmitJob(r.Context(), buildSubmitJobRequest(req.LaunchConfig, req.Script, nil))
	if err != nil {
		writeError(w, http.StatusBadGateway, fmt.Errorf("submit job: %w", err))
		return
	}

	jobInfo, err := slurm.GetJob(r.Context(), slurmJobID)
	if err != nil {
		writeError(w, http.StatusBadGateway, fmt.Errorf("submit job: %w", err))
		return
	}

	jobName := req.Name
	if jobInfo.Name != nil && *jobInfo.Name != "" {
		jobName = *jobInfo.Name
	}

	record, err := db.CreateJob(r.Context(), h.db, uuid.New(), slurmJobID, jobName, jobInfo.State)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	response := jobToAPI(record)
	response.Script = &req.Script
	response.Partition = req.Partition
	response.CurrentWorkingDirectory = &req.CurrentWorkingDirectory
	response.Environment = req.Environment
	response.CpusPerTask = req.CpusPerTask
	response.TasksPerNode = req.TasksPerNode
	response.MemoryPerNode = req.MemoryPerNode
	response.TimeLimit = req.TimeLimit

	h.NotifyJobsChanged()
	writeJSON(w, http.StatusCreated, response)
}
