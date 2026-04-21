package server

import (
	"context"
	"database/sql"
	"time"

	"github.com/paddyoneill/slurm-ui/internal/db"
	"github.com/paddyoneill/slurm-ui/internal/slurm"
)

const (
	jobPollInterval      = 10 * time.Second
	jobPollTimeout       = 30 * time.Second
	notebookPollInterval = 5 * time.Second
	notebookPollTimeout  = 30 * time.Second
)

func startJobPoller(dbConn *sql.DB, onJobsChanged func()) {
	ticker := time.NewTicker(jobPollInterval)
	defer ticker.Stop()

	for range ticker.C {
		ctx, cancel := context.WithTimeout(context.Background(), jobPollTimeout)
		pollJobs(ctx, dbConn, onJobsChanged)
		cancel()
	}
}

func pollJobs(ctx context.Context, dbConn *sql.DB, onJobsChanged func()) {
	jobs, err := db.ListJobs(ctx, dbConn)
	if err != nil {
		return
	}

	updated := false
	for _, job := range jobs {
		jobInfo, err := slurm.GetJob(ctx, job.SlurmJobId)
		if err != nil {
			continue
		}

		jobName := job.Name
		if jobInfo.Name != nil && *jobInfo.Name != "" {
			jobName = *jobInfo.Name
		}

		if jobName == job.Name && jobInfo.State == job.State {
			continue
		}

		if err := db.UpdateJobs(ctx, dbConn, job.ID, jobName, jobInfo.State); err != nil {
			continue
		}

		updated = true
	}

	if updated {
		onJobsChanged()
	}
}

func startNotebookPoller(dbConn *sql.DB, onJobsChanged func()) {
	ticker := time.NewTicker(notebookPollInterval)
	defer ticker.Stop()

	for range ticker.C {
		ctx, cancel := context.WithTimeout(context.Background(), notebookPollTimeout)
		pollNotebooks(ctx, dbConn, onJobsChanged)
		cancel()
	}
}

func pollNotebooks(ctx context.Context, dbConn *sql.DB, onJobsChanged func()) {
	notebooks, err := db.ListNotebooks(ctx, dbConn)
	if err != nil {
		return
	}

	updated := false
	for _, notebook := range notebooks {
		jobInfo, err := slurm.GetJob(ctx, notebook.SlurmJobID)
		if err != nil {
			continue
		}

		notebookName := notebook.Name
		if jobInfo.Name != nil && *jobInfo.Name != "" {
			notebookName = *jobInfo.Name
		}

		if notebookName == notebook.Name && jobInfo.State == notebook.State && equalStringPtrs(jobInfo.Host, notebook.Host) {
			continue
		}

		if err := db.UpdateNotebooks(ctx, dbConn, notebook.ID, notebookName, jobInfo.State, jobInfo.Host); err != nil {
			continue
		}

		updated = true
	}

	if updated {
		onJobsChanged()
	}
}

func equalStringPtrs(a *string, b *string) bool {
	switch {
	case a == nil && b == nil:
		return true
	case a == nil || b == nil:
		return false
	default:
		return *a == *b
	}
}
