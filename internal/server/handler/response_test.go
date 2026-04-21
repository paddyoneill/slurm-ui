package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	servertypes "github.com/paddyoneill/slurm-ui/internal/server/types"
)

func TestWriteJSONWritesStatusContentTypeAndBody(t *testing.T) {
	rr := httptest.NewRecorder()

	writeJSON(rr, http.StatusCreated, servertypes.Error{Error: "created"})

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rr.Code)
	}
	if got := rr.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("expected application/json content type, got %q", got)
	}
	if body := strings.TrimSpace(rr.Body.String()); body != `{"error":"created"}` {
		t.Fatalf("unexpected response body %q", body)
	}
}

func TestWriteErrorWrapsErrorInJSONPayload(t *testing.T) {
	rr := httptest.NewRecorder()

	writeError(rr, http.StatusBadGateway, errors.New("upstream failed"))

	if rr.Code != http.StatusBadGateway {
		t.Fatalf("expected status %d, got %d", http.StatusBadGateway, rr.Code)
	}
	if got := rr.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("expected application/json content type, got %q", got)
	}
	if body := strings.TrimSpace(rr.Body.String()); body != `{"error":"upstream failed"}` {
		t.Fatalf("unexpected error response body %q", body)
	}
}

func TestWriteJSONFallsBackWhenEncodingFails(t *testing.T) {
	rr := httptest.NewRecorder()

	writeJSON(rr, http.StatusAccepted, map[string]any{"bad": make(chan int)})

	if rr.Code != http.StatusAccepted {
		t.Fatalf("expected original status %d, got %d", http.StatusAccepted, rr.Code)
	}
	if got := rr.Header().Get("Content-Type"); got != "text/plain; charset=utf-8" {
		t.Fatalf("expected fallback content type from http.Error, got %q", got)
	}
	if body := strings.TrimSpace(rr.Body.String()); body != `{"error": "failed to encode response"}` {
		t.Fatalf("unexpected fallback body %q", body)
	}
}
