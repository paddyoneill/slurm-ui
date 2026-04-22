-- +goose Up
ALTER TABLE notebooks ADD COLUMN registration_token TEXT NOT NULL DEFAULT '';

-- +goose Down
CREATE TABLE notebooks_rollback (
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

INSERT INTO notebooks_rollback (id, slurm_job_id, slurm_job_name, slurm_job_state, host, token, port, base_env, created_at)
SELECT id, slurm_job_id, slurm_job_name, slurm_job_state, host, token, port, base_env, created_at
FROM notebooks;

DROP TABLE notebooks;
ALTER TABLE notebooks_rollback RENAME TO notebooks;
