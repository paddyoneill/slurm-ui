package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type NotebookRecord struct {
	ID         uuid.UUID
	SlurmJobID int
	Name       string
	State      string
	Host       *string
	Port       int
	Token      string
	BaseEnv    string
	CreatedAt  time.Time
}

var ErrNotebookNotFound = errors.New("notebook not found")

func CreateNotebook(ctx context.Context, conn *sql.DB, id uuid.UUID, slurmJobId int, name string, state string, host *string, port int, token string, baseEnv string) (NotebookRecord, error) {
	createdAt := time.Now().UTC()

	_, err := conn.ExecContext(ctx, "INSERT INTO notebooks (id, slurm_job_id, slurm_job_name, slurm_job_state, host, token, port, base_env, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)", id.String(), slurmJobId, name, state, host, token, port, baseEnv, createdAt.Format(time.RFC3339))
	if err != nil {
		return NotebookRecord{}, fmt.Errorf("create notebook %s: %w", id, err)
	}

	return NotebookRecord{
		ID:         id,
		SlurmJobID: slurmJobId,
		Name:       name,
		State:      state,
		Host:       host,
		Port:       port,
		Token:      token,
		BaseEnv:    baseEnv,
		CreatedAt:  createdAt,
	}, nil
}

func DeleteNotebookByID(ctx context.Context, conn *sql.DB, id uuid.UUID) error {
	result, err := conn.ExecContext(ctx, "DELETE FROM notebooks WHERE id = ?", id.String())
	if err != nil {
		return fmt.Errorf("delete notebook %s: %w", id, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete notebook %s: get rows affected: %w", id.String(), err)
	}

	switch rowsAffected {
	case 0:
		return ErrNotebookNotFound
	case 1:
		return nil
	default:
		return fmt.Errorf("delete notebook %s: expected to delete 1 row, deleted %d", id, rowsAffected)
	}
}

func GetNotebookByID(ctx context.Context, conn *sql.DB, id uuid.UUID) (NotebookRecord, error) {
	row := conn.QueryRowContext(ctx, "SELECT id, slurm_job_id, slurm_job_name, slurm_job_state, host, port, token, base_env, created_at FROM notebooks WHERE id = ?", id.String())

	notebook, err := scanNotebook(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return NotebookRecord{}, ErrNotebookNotFound
		}
		return NotebookRecord{}, fmt.Errorf("get notebook %s: %w", id, err)
	}

	return notebook, nil
}

func ListNotebooks(ctx context.Context, conn *sql.DB) ([]NotebookRecord, error) {
	rows, err := conn.QueryContext(ctx, "SELECT id, slurm_job_id, slurm_job_name, slurm_job_state, host, port, token, base_env, created_at FROM notebooks ORDER BY created_at DESC")
	if err != nil {
		return nil, fmt.Errorf("query notebooks: %w", err)
	}
	defer rows.Close()

	notebooks := make([]NotebookRecord, 0)
	for rows.Next() {
		notebook, err := scanNotebook(rows)
		if err != nil {
			return nil, err
		}
		notebooks = append(notebooks, notebook)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate notebooks: %w", err)
	}

	return notebooks, nil
}

func UpdateNotebooks(ctx context.Context, conn *sql.DB, id uuid.UUID, name string, state string, host *string) error {
	result, err := conn.ExecContext(ctx, "UPDATE notebooks SET slurm_job_name = ?, slurm_job_state = ?, host = ? WHERE id = ?", name, state, host, id.String())
	if err != nil {
		return fmt.Errorf("update notebook %s: %w", id, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("update notebook %s: get rows affected: %w", id.String(), err)
	}

	switch rowsAffected {
	case 0:
		return ErrNotebookNotFound
	case 1:
		return nil
	default:
		return fmt.Errorf("update notebook %s: expected to delete 1 row, updated %d", id, rowsAffected)
	}
}

func scanNotebook(scanner interface{ Scan(dest ...any) error }) (NotebookRecord, error) {
	var (
		idRaw        string
		slurmJobId   sql.NullInt64
		name         string
		state        string
		host         sql.NullString
		token        string
		port         int
		baseEnv      string
		createdAtRaw string
	)

	if err := scanner.Scan(&idRaw, &slurmJobId, &name, &state, &host, &port, &token, &baseEnv, &createdAtRaw); err != nil {
		return NotebookRecord{}, fmt.Errorf("scan notebook: %w", err)
	}

	id, err := uuid.Parse(idRaw)
	if err != nil {
		return NotebookRecord{}, fmt.Errorf("parse notebook id %s: %w", idRaw, err)
	}

	created_at, err := time.Parse(time.RFC3339, createdAtRaw)
	if err != nil {
		return NotebookRecord{}, fmt.Errorf("parse created_at for notebook %s: %w", idRaw, err)
	}

	return NotebookRecord{
		ID:         id,
		SlurmJobID: int(slurmJobId.Int64),
		Name:       name,
		State:      state,
		Host:       nullStringPtr(host),
		Token:      token,
		Port:       port,
		BaseEnv:    baseEnv,
		CreatedAt:  created_at,
	}, nil
}

func nullStringPtr(value sql.NullString) *string {
	if !value.Valid {
		return nil
	}

	result := value.String
	return &result
}
