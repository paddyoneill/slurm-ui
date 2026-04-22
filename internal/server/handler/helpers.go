package handler

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/google/uuid"

	"github.com/paddyoneill/slurm-ui/internal/db"
	servertypes "github.com/paddyoneill/slurm-ui/internal/server/types"
	"github.com/paddyoneill/slurm-ui/internal/slurm"
)

const (
	notebookBaseURLEnv       = "SLURM_UI_BASE_URL"
	unregisteredNotebookPort = 0
)

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

func buildNotebookEnvironment(notebookID uuid.UUID, token string, registrationToken string, registrationURL string, req servertypes.CreateNotebookRequest) []string {
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
		"NOTEBOOK_TOKEN="+token,
		"NOTEBOOK_BASE_URL="+baseURL,
		"NOTEBOOK_REGISTRATION_TOKEN="+registrationToken,
		"NOTEBOOK_REGISTRATION_URL="+registrationURL,
	)
}

func buildNotebookScript() string {
	return `#!/usr/bin/env bash
	set -eou pipefail

	cleanup() {
		if [[ -n "${NOTEBOOK_PID:-}" ]]; then
			kill "${NOTEBOOK_PID}" 2>/dev/null || true
		fi
		rm -rf "${NOTEBOOK_VENV_DIR}"
	}

	mkdir -p "${NOTEBOOK_DIR}"
	cp -r "${NOTEBOOK_BASE_ENV}" "${NOTEBOOK_VENV_DIR}"
	source "${NOTEBOOK_VENV_DIR}/bin/activate"
	trap cleanup EXIT
	jupyter server --ip=0.0.0.0 --no-browser --ServerApp.token="${NOTEBOOK_TOKEN}" --ServerApp.base_url="${NOTEBOOK_BASE_URL}" --ServerApp.allow_origin='*' &

	export NOTEBOOK_PID="$!"

	NOTEBOOK_PORT="$(python - << 'PYTHON'
import os
import socket
import time

def iter_running_servers():
	providers = []
	try:
		from jupyter_server.serverapp import list_running_servers as list_jp_servers
		providers.append(list_jp_servers)
	except ImportError:
		pass

	try:
		from notebook.notebookapp import list_running_servers as list_nb_servers
		providers.append(list_nb_servers)
	except ImportError:
		pass

	if not providers:
		raise RuntimeError("unable to locate Jupyter runtime server metadata support")

	for provider in providers:
		try:
			yield from provider()
		except Exception:
			continue

def match_server(server, pid, token, base_url):
	if server.get("pid") != pid:
		return False

	server_base_url = server.get("base_url", "")
	if server_base_url and server_base_url != base_url:
		return False

	server_token = server.get("token")
	if server_token is not None and server_token != token:
		return False

	return True

pid = int(os.environ["NOTEBOOK_PID"])
token = os.environ["NOTEBOOK_TOKEN"]
base_url = os.environ["NOTEBOOK_BASE_URL"]
deadline = time.monotonic() + 60

while time.monotonic() < deadline:
	try:
		os.kill(pid, 0)
	except OSError as exc:
		raise RuntimeError("jupyter exited before notebook runtime metadata became available") from exc

	for server in iter_running_servers():
		if not match_server(server, pid, token, base_url):
			continue

		port = int(server["port"])
		try:
			with socket.create_connection(("127.0.0.1", port), timeout=1):
				print(port)
				raise SystemExit(0)
		except OSError:
			break

	time.sleep(1)

raise RuntimeError("timed out waiting for notebook runtime metadata")
PYTHON
)"
	export NOTEBOOK_PORT

	echo "Notebook listening on port ${NOTEBOOK_PORT}"

	python - <<'PYTHON'
import json
import os
import urllib.request

payload = json.dumps({
	"port": int(os.environ["NOTEBOOK_PORT"]),
	"registrationToken": os.environ["NOTEBOOK_REGISTRATION_TOKEN"],
}).encode()

request = urllib.request.Request(
	os.environ["NOTEBOOK_REGISTRATION_URL"],
	data=payload,
	headers={"Content-Type": "application/json"},
	method="POST",
)

with urllib.request.urlopen(request) as response:
	if response.status != 204:
		raise RuntimeError(f"unexpected notebook registration status: {response.status}")
PYTHON

	wait "${NOTEBOOK_PID}"`
}

func notebookEnvironmentLen(values *[]string) int {
	if values == nil {
		return 0
	}

	return len(*values)
}

func buildNotebookRegisterURL(notebookID uuid.UUID) (string, error) {
	baseURL := strings.TrimSpace(os.Getenv(notebookBaseURLEnv))
	if baseURL == "" {
		return "", fmt.Errorf("%s is required", notebookBaseURLEnv)
	}

	parsed, err := url.Parse(baseURL)
	if err != nil {
		return "", fmt.Errorf("parse %s: %w", notebookBaseURLEnv, err)
	}

	if parsed.Scheme == "" || parsed.Host == "" {
		return "", fmt.Errorf("%s must include scheme and host", notebookBaseURLEnv)
	}

	parsed.RawQuery = ""
	parsed.Fragment = ""
	parsed.Path = strings.TrimRight(parsed.Path, "/") + fmt.Sprintf("/api/notebooks/%s/register", notebookID)

	return parsed.String(), nil
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
