package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	dbpkg "github.com/paddyoneill/slurm-ui/internal/db"
	servertypes "github.com/paddyoneill/slurm-ui/internal/server/types"
)

func TestGetNotebooksReturnsNotebooksJSON(t *testing.T) {
	h, conn := newTestHandler(t)
	host := "compute-01"

	notebook, err := dbpkg.CreateNotebook(context.Background(), conn, uuid.New(), 5678, "notebook-one", "RUNNING", &host, 8888, "token-123", "/envs/base")
	if err != nil {
		t.Fatalf("create notebook: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/notebooks", nil)
	rr := httptest.NewRecorder()

	h.GetNotebooks(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	var notebooks []servertypes.Notebook
	if err := json.NewDecoder(rr.Body).Decode(&notebooks); err != nil {
		t.Fatalf("decode notebooks response: %v", err)
	}

	if len(notebooks) != 1 {
		t.Fatalf("expected 1 notebook, got %d", len(notebooks))
	}
	if notebooks[0].Id != notebook.ID || notebooks[0].SlurmJobId != 5678 || notebooks[0].Name != "notebook-one" {
		t.Fatalf("unexpected notebook payload: %+v", notebooks[0])
	}
	if notebooks[0].Host == nil || *notebooks[0].Host != host {
		t.Fatalf("expected host %q, got %#v", host, notebooks[0].Host)
	}
}

func TestGetNotebooksIdNotFoundReturns404(t *testing.T) {
	h, _ := newTestHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/api/notebooks/missing", nil)
	rr := httptest.NewRecorder()

	h.GetNotebooksId(rr, req, uuid.New())

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", rr.Code)
	}
}

func TestDeleteNotebooksIdNotFoundReturns404(t *testing.T) {
	h, _ := newTestHandler(t)

	req := httptest.NewRequest(http.MethodDelete, "/api/notebooks/missing", nil)
	rr := httptest.NewRecorder()

	h.DeleteNotebooksId(rr, req, uuid.New())

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", rr.Code)
	}
}

func TestHandleNotebooksEventsWritesInitialEventAndHeaders(t *testing.T) {
	h, conn := newTestHandler(t)
	host := "compute-02"

	_, err := dbpkg.CreateNotebook(context.Background(), conn, uuid.New(), 6789, "stream-notebook", "RUNNING", &host, 8888, "token-456", "/envs/base")
	if err != nil {
		t.Fatalf("create notebook: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	req := httptest.NewRequest(http.MethodGet, "/api/notebooks/events", nil).WithContext(ctx)
	rr := httptest.NewRecorder()

	done := make(chan struct{})
	go func() {
		defer close(done)
		h.GetNotebooksEvents(rr, req)
	}()

	time.Sleep(10 * time.Millisecond)
	cancel()
	<-done

	if got := rr.Header().Get("Content-Type"); got != "text/event-stream" {
		t.Fatalf("expected text/event-stream content type, got %q", got)
	}
	if got := rr.Header().Get("Cache-Control"); got != "no-cache" {
		t.Fatalf("expected no-cache header, got %q", got)
	}

	body := rr.Body.String()
	if !strings.HasPrefix(body, "data: ") {
		t.Fatalf("expected SSE body to start with data prefix, got %q", body)
	}
	if !strings.Contains(body, `"name":"stream-notebook"`) {
		t.Fatalf("expected SSE payload to include notebook name, got %q", body)
	}
}
