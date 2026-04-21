package handler

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/paddyoneill/slurm-ui/internal/db"
	servertypes "github.com/paddyoneill/slurm-ui/internal/server/types"
	"github.com/paddyoneill/slurm-ui/internal/slurm"
)

const notebookPort = 8888

func jobToAPI(record db.JobRecord) servertypes.Job {
	return servertypes.Job{
		Id:         record.ID,
		SlurmJobId: record.SlurmJobId,
		Name:       record.Name,
		State:      record.State,
		CreatedAt:  record.CreatedAt,
	}
}

func notebookToAPI(record db.NotebookRecord) servertypes.Notebook {
	return servertypes.Notebook{
		Id:         record.ID,
		SlurmJobId: record.SlurmJobID,
		Name:       record.Name,
		State:      record.State,
		Host:       record.Host,
		Port:       record.Port,
		Token:      record.Token,
		BaseEnv:    record.BaseEnv,
		CreatedAt:  record.CreatedAt,
	}
}

func buildNotebookEnvironment(notebookID uuid.UUID, token string, req servertypes.CreateNotebookRequest) []string {
	notebookDir := fmt.Sprintf("%s/notebooks/%s", req.CurrentWorkingDirectory, notebookID)
	venvDir := fmt.Sprintf("%s/venv", notebookDir)
	baseURL := fmt.Sprintf("/api/notebooks/%s/proxy/", notebookID)

	env := make([]string, 0, notebookEnvironmentLen(req.Environment)+7)
	if req.Environment != nil {
		env = append(env, (*req.Environment)...)
	}

	return append(env,
		"NOTEBOOK_DIR="+notebookDir,
		"NOTEBOOK_BASE_ENV="+*req.BaseEnv,
		"NOTEBOOK_VENV_DIR="+venvDir,
		fmt.Sprintf("NOTEBOOK_PORT=%d", notebookPort),
		"NOTEBOOK_TOKEN="+token,
		"NOTEBOOK_BASE_URL="+baseURL,
	)
}

func buildNotebookScript() string {
	return `#!/usr/bin/env bash

	mkdir -p "${NOTEBOOK_DIR}"
	cp -r "${NOTEBOOK_BASE_ENV}" "${NOTEBOOK_VENV_DIR}"
	source "${NOTEBOOK_VENV_DIR}/bin/activate"
	trap 'rm -rf "${NOTEBOOK_VENV_DIR}"' EXIT
	jupyter notebook --ip=0.0.0.0 --port "${NOTEBOOK_PORT}" --no-browser --ServerApp.token="${NOTEBOOK_TOKEN}" --ServerApp.base_url="${NOTEBOOK_BASE_URL}" --ServerApp.allow_origin='*'`
}

func notebookEnvironmentLen(values *[]string) int {
	if values == nil {
		return 0
	}

	return len(*values)
}

func buildSubmitJobRequest(config servertypes.LaunchConfig, script string, environment *[]string) slurm.SubmitJobRequest {
	if environment == nil {
		environment = config.Environment
	}

	return slurm.SubmitJobRequest{
		Name:                    config.Name,
		Script:                  script,
		Partition:               config.Partition,
		CurrentWorkingDirectory: config.CurrentWorkingDirectory,
		Environment:             environment,
		CpusPerTask:             config.CpusPerTask,
		TasksPerNode:            config.TasksPerNode,
		MemoryPerNode:           config.MemoryPerNode,
		TimeLimit:               config.TimeLimit,
	}
}
