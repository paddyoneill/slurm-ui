# API

This document describes the current HTTP API exposed by the Go backend.

Base URL:

```text
/api
```

All JSON endpoints use `application/json`.

## Error Responses

Generic error:

```json
{
  "error": "something went wrong"
}
```

Validation error:

```json
{
  "error": "Request failed validation",
  "fields": [
    {
      "field": "name",
      "message": "Job name is required"
    }
  ]
}
```

Common status codes:

- `400 Bad Request`: invalid JSON, invalid UUID path parameter, or validation failure
- `404 Not Found`: resource does not exist in the local database
- `500 Internal Server Error`: local database or internal server failure
- `502 Bad Gateway`: upstream Slurm or notebook proxy failure

## Jobs

### `GET /api/jobs`

Returns all jobs currently stored in the local database.

Response:

```json
[
  {
    "id": "d6ff18b2-4159-4b64-9b4d-98d986c01cab",
    "slurmJobId": 1234,
    "name": "example-job",
    "state": "RUNNING",
    "createdAt": "2026-04-21T10:00:00Z"
  }
]
```

### `POST /api/jobs`

Creates and submits a Slurm job, stores the local record, and returns the created job.

Request body:

```json
{
  "name": "example-job",
  "script": "#!/usr/bin/env bash\necho hello",
  "currentWorkingDirectory": "/home/user",
  "partition": "debug",
  "environment": ["FOO=bar"],
  "cpusPerTask": 1,
  "tasksPerNode": 1,
  "memoryPerNode": 512,
  "timeLimit": 30
}
```

Notes:

- `name`, `script`, and `currentWorkingDirectory` are required
- `currentWorkingDirectory` must be an absolute path
- environment entries must look like `KEY=value`
- `cpusPerTask`, `tasksPerNode`, `memoryPerNode`, and `timeLimit` must be positive when supplied

Response:

```json
{
  "id": "d6ff18b2-4159-4b64-9b4d-98d986c01cab",
  "slurmJobId": 1234,
  "name": "example-job",
  "state": "PENDING",
  "createdAt": "2026-04-21T10:00:00Z",
  "script": "#!/usr/bin/env bash\necho hello",
  "currentWorkingDirectory": "/home/user",
  "partition": "debug",
  "environment": ["FOO=bar"],
  "cpusPerTask": 1,
  "tasksPerNode": 1,
  "memoryPerNode": 512,
  "timeLimit": 30
}
```

### `GET /api/jobs/{id}`

Returns a single job by internal UUID.

### `DELETE /api/jobs/{id}`

Cancels the Slurm job first, then removes the local job record.

Response:

- `204 No Content` on success

## Job Events

### `GET /api/jobs/events`

Server-Sent Events stream for jobs.

Headers:

- `Content-Type: text/event-stream`
- `Cache-Control: no-cache`

Behavior:

- sends an initial snapshot immediately on connect
- sends a full jobs array again whenever the local jobs state changes

Event payload format:

```text
data: [{"id":"...","slurmJobId":1234,"name":"example-job","state":"RUNNING","createdAt":"2026-04-21T10:00:00Z"}]

```

## Notebooks

### `GET /api/notebooks`

Returns all notebooks currently stored in the local database.

Response:

```json
[
  {
    "id": "d6ff18b2-4159-4b64-9b4d-98d986c01cab",
    "slurmJobId": 5678,
    "name": "example-notebook",
    "state": "RUNNING",
    "host": "compute-01",
    "port": 8888,
    "token": "secret-token",
    "baseEnv": "/shared/envs/python",
    "createdAt": "2026-04-21T10:00:00Z"
  }
]
```

### `POST /api/notebooks`

Creates a notebook launch script, submits it to Slurm, stores the local notebook record, and returns the created notebook.

Request body:

```json
{
  "name": "example-notebook",
  "baseEnv": "/shared/envs/python",
  "currentWorkingDirectory": "/home/user",
  "partition": "debug",
  "environment": ["FOO=bar"],
  "cpusPerTask": 1,
  "tasksPerNode": 1,
  "memoryPerNode": 512,
  "timeLimit": 30
}
```

Notes:

- `name`, `baseEnv`, and `currentWorkingDirectory` are required
- `baseEnv` must be non-empty
- launch-field validation matches `POST /api/jobs`

Response:

```json
{
  "id": "d6ff18b2-4159-4b64-9b4d-98d986c01cab",
  "slurmJobId": 5678,
  "name": "example-notebook",
  "state": "PENDING",
  "host": null,
  "port": 8888,
  "token": "secret-token",
  "baseEnv": "/shared/envs/python",
  "createdAt": "2026-04-21T10:00:00Z",
  "currentWorkingDirectory": "/home/user",
  "partition": "debug",
  "environment": ["FOO=bar"],
  "cpusPerTask": 1,
  "tasksPerNode": 1,
  "memoryPerNode": 512,
  "timeLimit": 30
}
```

### `GET /api/notebooks/{id}`

Returns a single notebook by internal UUID.

### `DELETE /api/notebooks/{id}`

Cancels the Slurm job first, then removes the local notebook record.

Response:

- `204 No Content` on success

## Notebook Events

### `GET /api/notebooks/events`

Server-Sent Events stream for notebooks.

Headers:

- `Content-Type: text/event-stream`
- `Cache-Control: no-cache`

Behavior:

- sends an initial snapshot immediately on connect
- sends a full notebooks array again whenever the local notebooks state changes

## Notebook Proxy

Notebook HTTP traffic is proxied through the backend to the running Jupyter server.

Routes:

- `GET /api/notebooks/{id}/proxy`
- `GET /api/notebooks/{id}/proxy/{path...}`
- `POST /api/notebooks/{id}/proxy/{path...}`
- `PUT /api/notebooks/{id}/proxy/{path...}`
- `PATCH /api/notebooks/{id}/proxy/{path...}`
- `DELETE /api/notebooks/{id}/proxy/{path...}`

Behavior:

- notebook lookup is done using the local notebook `id`
- if no proxy path is supplied, the backend redirects to:

```text
/api/notebooks/{id}/proxy/tree
```

- the target upstream is:

```text
http://{host}:{port}
```

- upstream `Location` headers are rewritten back under `/api/notebooks/{id}/proxy/...`

Proxy errors:

- `404` if the notebook is not found locally
- `502` if the notebook host is unavailable or the upstream proxy request fails
