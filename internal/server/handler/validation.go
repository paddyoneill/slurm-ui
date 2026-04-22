package handler

import (
	"encoding/json"
	"net/http"
	"path"
	"regexp"
	"strings"

	servertypes "github.com/paddyoneill/slurm-ui/internal/server/types"
)

var envPattern = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*=.*$`)

func validationError(message string, issues ...servertypes.ValidationIssue) *servertypes.ValidationError {
	if len(issues) == 0 {
		return &servertypes.ValidationError{Error: message}
	}

	fields := make([]servertypes.ValidationIssue, 0, len(issues))

	for _, issue := range issues {
		fields = append(fields, issue)
	}

	return &servertypes.ValidationError{
		Error:  message,
		Fields: &fields,
	}
}

func validationField(field string, message string) servertypes.ValidationIssue {
	return servertypes.ValidationIssue{Field: field, Message: message}
}

func decodeCreateJobRequest(r *http.Request) (servertypes.CreateJobRequest, *servertypes.ValidationError) {
	var req servertypes.CreateJobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return servertypes.CreateJobRequest{}, validationError("Request failed validation", validationField("", "Invalid JSON request body"))
	}

	issues := validateCreateJobRequest(&req)
	if len(issues) > 0 {
		return servertypes.CreateJobRequest{}, validationError("Request failed validation", issues...)
	}

	return req, nil
}

func validateCreateJobRequest(req *servertypes.CreateJobRequest) []servertypes.ValidationIssue {
	fields := validateLaunchConfig(&req.LaunchConfig)

	req.Script = strings.TrimSpace(req.Script)
	if req.Script == "" {
		fields = append(fields, validationField("script", "Script is required"))
	}

	return fields
}

func decodeCreateNotebookRequest(r *http.Request) (servertypes.CreateNotebookRequest, *servertypes.ValidationError) {
	var req servertypes.CreateNotebookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return servertypes.CreateNotebookRequest{}, validationError("Request failed validation", validationField("", "Invalid JSON request body"))
	}

	issues := validateCreateNotebookRequest(&req)
	if len(issues) > 0 {
		return servertypes.CreateNotebookRequest{}, validationError("Request failed validation", issues...)
	}

	return req, nil
}

func validateCreateNotebookRequest(req *servertypes.CreateNotebookRequest) []servertypes.ValidationIssue {
	fields := validateLaunchConfig(&req.LaunchConfig)

	if req.BaseEnv == nil {
		fields = append(fields, validationField("baseEnv", "Base environment is required"))
	} else {
		trimmed := strings.TrimSpace(*req.BaseEnv)
		if trimmed == "" {
			fields = append(fields, validationField("baseEnv", "Base environment cannot be empty"))
		} else {
			req.BaseEnv = &trimmed
		}
	}

	return fields
}

func decodeRegisterNotebookRequest(r *http.Request) (servertypes.RegisterNotebookRequest, *servertypes.ValidationError) {
	var req servertypes.RegisterNotebookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return servertypes.RegisterNotebookRequest{}, validationError("Request faiiled validation", validationField("", "Invalid JSON request body"))
	}

	issues := validateRegisterNotebookRequest(&req)
	if len(issues) > 0 {
		return servertypes.RegisterNotebookRequest{}, validationError("Request failed validation", issues...)
	}

	return req, nil
}

func validateRegisterNotebookRequest(req *servertypes.RegisterNotebookRequest) []servertypes.ValidationIssue {
	fields := make([]servertypes.ValidationIssue, 0)

	req.RegistrationToken = strings.TrimSpace(req.RegistrationToken)
	if req.RegistrationToken == "" {
		fields = append(fields, validationField("registrationToken", "Registration token is required"))
	}

	if req.Port < 1 || req.Port > 65535 {
		fields = append(fields, validationField("port", "Port must be between 1 and 65535"))
	}

	return fields
}

func validateLaunchConfig(req *servertypes.LaunchConfig) []servertypes.ValidationIssue {
	fields := make([]servertypes.ValidationIssue, 0)

	req.Name = strings.TrimSpace(req.Name)
	req.CurrentWorkingDirectory = strings.TrimSpace(req.CurrentWorkingDirectory)

	if req.Name == "" {
		fields = append(fields, validationField("name", "Job name is required"))
	}

	if req.CurrentWorkingDirectory == "" {
		fields = append(fields, validationField("currentWorkingDirectory", "Working directory is required"))
	} else if !path.IsAbs(req.CurrentWorkingDirectory) {
		fields = append(fields, validationField("currentWorkingDirectory", "Must be an absolute path"))
	}

	normaliseOptionalString(&fields, "partition", "Partition cannot be empty", &req.Partition)
	validatePositiveInt(&fields, "cpusPerTask", "Must request at least 1 CPU per task", req.CpusPerTask)
	validatePositiveInt(&fields, "tasksPerNode", "Must request at least 1 task per node", req.TasksPerNode)
	validatePositiveInt(&fields, "memoryPerNode", "Memory must be at least 1 MB", req.MemoryPerNode)
	validatePositiveInt(&fields, "timeLimit", "Time limit must be at least 1 minute", req.TimeLimit)

	if req.Environment != nil {
		normalised := make([]string, 0, len(*req.Environment))
		for _, val := range *req.Environment {
			trimmed := strings.TrimSpace(val)
			if !envPattern.MatchString(trimmed) {
				fields = append(fields, validationField("environment", "Environment variables must look like KEY=value"))
				continue
			}
			normalised = append(normalised, trimmed)
		}
		req.Environment = &normalised
	}

	return fields
}

func normaliseOptionalString(fields *[]servertypes.ValidationIssue, field string, emptyMessage string, value **string) {
	if *value == nil {
		return
	}

	trimmed := strings.TrimSpace(**value)
	if trimmed == "" {
		*fields = append(*fields, validationField(field, emptyMessage))
	}

	*value = &trimmed
}

func validatePositiveInt(fields *[]servertypes.ValidationIssue, field string, message string, value *int) {
	if value != nil && *value < 1 {
		*fields = append(*fields, validationField(field, message))
	}
}
