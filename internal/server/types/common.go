package types

import "time"

// Error is a generic error response payload.
type Error struct {
	Error string `json:"error"`
}

// ValidationIssue describes a signle field-level validation issue.
type ValidationIssue struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationError represents a validation failure response.
type ValidationError struct {
	Error  string             `json:"error"`
	Fields *[]ValidationIssue `json:"fields,omitempty"`
}

// LaunchConfig captures the shared launch-time settings used by jobs and notebooks.
type LaunchConfig struct {
	CpusPerTask             *int      `json:"cpusPerTask,omitempty"`
	CurrentWorkingDirectory string    `json:"currentWorkingDirectory"`
	Environment             *[]string `json:"environment,omitempty"`
	MemoryPerNode           *int      `json:"memoryPerNode,omitempty"`
	Name                    string    `json:"name"`
	Partition               *string   `json:"partition,omitempty"`
	TasksPerNode            *int      `json:"tasksPerNode,omitempty"`
	TimeLimit               *int      `json:"timeLimit,omitempty"`
}

// Timestamp exists so request/response payloads can share the same time type.
type Timestamp = time.Time
