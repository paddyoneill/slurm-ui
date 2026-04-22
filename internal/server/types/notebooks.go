package types

import "github.com/google/uuid"

// CreateNotebookRequest is the JSON payload accepted when creating a notebook.
type CreateNotebookRequest struct {
	LaunchConfig
	BaseEnv *string `json:"baseEnv,omitempty"`
}

// RegisterNotebookRequest is the JSON payload accepted from notebook port registration requests.
type RegisterNotebookRequest struct {
	Port              int    `json:"port"`
	RegistrationToken string `json:"registrationToken"`
}

// Notebook is the JSON representation returned for notebook resources.
type Notebook struct {
	BaseEnv                 string    `json:"baseEnv"`
	CpusPerTask             *int      `json:"cpusPerTask,omitempty"`
	CreatedAt               Timestamp `json:"createdAt"`
	CurrentWorkingDirectory *string   `json:"currentWorkingDirectory,omitempty"`
	Environment             *[]string `json:"environment,omitempty"`
	Host                    *string   `json:"host,omitempty"`
	Id                      uuid.UUID `json:"id"`
	MemoryPerNode           *int      `json:"memoryPerNode,omitempty"`
	Name                    string    `json:"name"`
	Partition               *string   `json:"partition,omitempty"`
	Port                    int       `json:"port"`
	SlurmJobId              int       `json:"slurmJobId"`
	State                   string    `json:"state"`
	TasksPerNode            *int      `json:"tasksPerNode,omitempty"`
	TimeLimit               *int      `json:"timeLimit,omitempty"`
	Token                   string    `json:"token"`
}

// NotebookID identifies a persisted notebook record.
type NotebookID = uuid.UUID

// ProxyPath is the proxied notebook-relative path segment.
type ProxyPath = string
