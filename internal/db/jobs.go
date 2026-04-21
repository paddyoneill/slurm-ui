package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type JobRecord struct {
	ID         uuid.UUID
	SlurmJobId int
	Name       string
	State      string
	CreatedAt  time.Time
}

var ErrJobNotFound = errors.New("job not found")

func CreateJob(ctx context.Context, conn *sql.DB, id uuid.UUID, slurmJobId int, name string, state string) (JobRecord, error) {
	createdAt := time.Now().UTC()

	_, err := conn.ExecContext(ctx, "INSERT INTO jobs (id, slurm_job_id, slurm_job_name, slurm_job_state, created_at) VALUES (?, ?, ?, ?, ?)", id.String(), slurmJobId, name, state, createdAt.Format(time.RFC3339))
	if err != nil {
		return JobRecord{}, fmt.Errorf("create job %s: %w", id, err)
	}

	return JobRecord{
		ID:         id,
		SlurmJobId: slurmJobId,
		Name:       name,
		State:      state,
		CreatedAt:  createdAt,
	}, nil
}

func DeleteJobByID(ctx context.Context, conn *sql.DB, id uuid.UUID) error {
	result, err := conn.ExecContext(ctx, "DELETE FROM jobs WHERE id = ?", id.String())
	if err != nil {
		return fmt.Errorf("delete job %s: %w", id, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete job %s: get rows affected: %w", id, err)
	}

	switch rowsAffected {
	case 0:
		return ErrJobNotFound
	case 1:
		return nil
	default:
		return fmt.Errorf("delete job %s: expected to delete 1 row, deleted %d", id, rowsAffected)
	}
}

func GetJobByID(ctx context.Context, conn *sql.DB, id uuid.UUID) (JobRecord, error) {
	row := conn.QueryRowContext(ctx, "SELECT id, slurm_job_id, slurm_job_name, slurm_job_state, created_at FROM jobs WHERE ID = ?", id.String())

	job, err := scanJob(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return JobRecord{}, ErrJobNotFound
		}
		return JobRecord{}, fmt.Errorf("get job %s: %w", id, err)
	}

	return job, nil
}

func ListJobs(ctx context.Context, conn *sql.DB) ([]JobRecord, error) {
	rows, err := conn.QueryContext(ctx, "SELECT id, slurm_job_id, slurm_job_name, slurm_job_state, created_at FROM jobs ORDER BY created_at DESC")
	if err != nil {
		return nil, fmt.Errorf("query jobs: %w", err)
	}
	defer rows.Close()

	jobs := make([]JobRecord, 0)
	for rows.Next() {
		job, err := scanJob(rows)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, job)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate jobs: %w", err)
	}

	return jobs, nil
}

func UpdateJobs(ctx context.Context, conn *sql.DB, id uuid.UUID, name string, state string) error {
	result, err := conn.ExecContext(ctx, "UPDATE jobs SET slurm_job_name = ?, slurm_job_state = ? WHERE id = ?", name, state, id.String())
	if err != nil {
		return fmt.Errorf("update job %s: %w", id, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("update job %s: got rows affected: %w", id, err)
	}

	switch rowsAffected {
	case 0:
		return ErrJobNotFound
	case 1:
		return nil
	default:
		return fmt.Errorf("update job %s: expected to update 1 row, updated %d", id, rowsAffected)
	}
}

func scanJob(scanner interface{ Scan(dest ...any) error }) (JobRecord, error) {
	var (
		idRaw       string
		slurmJobId  int
		name        string
		state       string
		createAtRaw string
	)

	if err := scanner.Scan(
		&idRaw,
		&slurmJobId,
		&name,
		&state,
		&createAtRaw,
	); err != nil {
		return JobRecord{}, fmt.Errorf("scan job: %w", err)
	}

	id, err := uuid.Parse(idRaw)
	if err != nil {
		return JobRecord{}, fmt.Errorf("parse job id %s: %w", idRaw, err)
	}

	created_at, err := time.Parse(time.RFC3339, createAtRaw)
	if err != nil {
		return JobRecord{}, fmt.Errorf("parse created_at for job %s: %w", idRaw, err)
	}

	return JobRecord{
		ID:         id,
		SlurmJobId: slurmJobId,
		Name:       name,
		State:      state,
		CreatedAt:  created_at.UTC(),
	}, nil
}
