-- +goose Up
CREATE TABLE IF NOT EXISTS notebooks (
    id TEXT PRIMARY KEY,
    slurm_job_id INTEGER NOT NULL,
    slurm_job_name TEXT NOT NULL,
    slurm_job_state TEXT NOT NULL DEFAULT 'PENDING',
    host TEXT,
    token TEXT NOT NULL,
    port INTEGER NOT NULL DEFAULT 8888,
    base_env TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS notebooks;
