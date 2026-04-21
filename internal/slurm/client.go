package slurm

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	slurmgen "github.com/SlinkyProject/slurm-client/api/v0043"
)

func newClient() (*slurmgen.ClientWithResponses, error) {
	baseURL := os.Getenv("SLURM_RESTAPI_URL")
	if baseURL == "" {
		return nil, fmt.Errorf("SLURM_RESTAPI_URL is unset")
	}

	token := os.Getenv("SLURM_JWT")
	if token == "" {
		return nil, fmt.Errorf("SLURM_JWT is unset")
	}

	return slurmgen.NewClientWithResponses(baseURL, slurmgen.WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
		req.Header.Set("X-SLURM-USER-TOKEN", token)
		return nil
	}))
}

func unexpectedStatus(operation string, statusCode int, body []byte) error {
	message := strings.TrimSpace(string(body))
	if message == "" {
		return fmt.Errorf("%s: unexpected status %d", operation, statusCode)
	}
	return fmt.Errorf("%s: unexpected status %d: %s", operation, statusCode, message)
}
