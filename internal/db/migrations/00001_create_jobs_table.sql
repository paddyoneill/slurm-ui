-- +goose Up
CREATE TABLE IF NOT EXISTS jobs (
    id TEXT PRIMARY KEY,
    slurm_job_id INTEGER NOT NULL,
    slurm_job_name TEXT NOT NULL,
    slurm_job_state TEXT NOT NULL DEFAULT 'PENDING',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS jobs;
