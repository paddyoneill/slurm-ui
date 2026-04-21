package types

import "github.com/google/uuid"

// CreateJobRequest is the JSON payload accepted when creating a Slurm job.
type CreateJobRequest struct {
	LaunchConfig
	Script string `json:"script"`
}

// Job is the JSON representation returned for Slurm jobs.
type Job struct {
	CpusPerTask             *int      `json:"cpusPerTask,omitempty"`
	CreatedAt               Timestamp `json:"createdAt"`
	CurrentWorkingDirectory *string   `json:"currentWorkingDirectory,omitempty"`
	Environment             *[]string `json:"environment,omitempty"`
	Id                      uuid.UUID `json:"id"`
	MemoryPerNode           *int      `json:"memoryPerNode,omitempty"`
	Name                    string    `json:"name"`
	Partition               *string   `json:"partition,omitempty"`
	Script                  *string   `json:"script,omitempty"`
	SlurmJobId              int       `json:"slurmJobId"`
	State                   string    `json:"state"`
	TasksPerNode            *int      `json:"tasksPerNode,omitempty"`
	TimeLimit               *int      `json:"timeLimit,omitempty"`
}

// JobID identifies a persisted job record.
type JobID = uuid.UUID
