package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	dbpkg "github.com/paddyoneill/slurm-ui/internal/db"
	servertypes "github.com/paddyoneill/slurm-ui/internal/server/types"
)

func TestGetJobsReturnsJobsJSON(t *testing.T) {
	h, conn := newTestHandler(t)

	job, err := dbpkg.CreateJob(context.Background(), conn, uuid.New(), 1234, "example-job", "RUNNING")
	if err != nil {
		t.Fatalf("create job: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/jobs", nil)
	rr := httptest.NewRecorder()

	h.GetJobs(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	var jobs []servertypes.Job
	if err := json.NewDecoder(rr.Body).Decode(&jobs); err != nil {
		t.Fatalf("decode jobs response: %v", err)
	}

	if len(jobs) != 1 {
		t.Fatalf("expected 1 job, got %d", len(jobs))
	}
	if jobs[0].Id != job.ID || jobs[0].SlurmJobId != 1234 || jobs[0].Name != "example-job" || jobs[0].State != "RUNNING" {
		t.Fatalf("unexpected job payload: %+v", jobs[0])
	}
}

func TestGetJobsIdNotFoundReturns404(t *testing.T) {
	h, _ := newTestHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/api/jobs/missing", nil)
	rr := httptest.NewRecorder()

	h.GetJobsId(rr, req, uuid.New())

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", rr.Code)
	}

	var response servertypes.Error
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("decode error response: %v", err)
	}
	if response.Error == "" {
		t.Fatal("expected error message in response")
	}
}

func TestDeleteJobsIdNotFoundReturns404(t *testing.T) {
	h, _ := newTestHandler(t)

	req := httptest.NewRequest(http.MethodDelete, "/api/jobs/missing", nil)
	rr := httptest.NewRecorder()

	h.DeleteJobsId(rr, req, uuid.New())

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", rr.Code)
	}
}

func TestHandleJobsEventsWritesInitialEventAndHeaders(t *testing.T) {
	h, conn := newTestHandler(t)

	_, err := dbpkg.CreateJob(context.Background(), conn, uuid.New(), 4321, "stream-job", "PENDING")
	if err != nil {
		t.Fatalf("create job: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	req := httptest.NewRequest(http.MethodGet, "/api/jobs/events", nil).WithContext(ctx)
	rr := httptest.NewRecorder()

	done := make(chan struct{})
	go func() {
		defer close(done)
		h.GetJobsEvents(rr, req)
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
	if !strings.Contains(body, `"name":"stream-job"`) {
		t.Fatalf("expected SSE payload to include job name, got %q", body)
	}
}

func newTestHandler(t *testing.T) (*Handler, *sql.DB) {
	t.Helper()

	dbPath := filepath.Join(t.TempDir(), "test.db")
	conn, err := dbpkg.Open(dbPath)
	if err != nil {
		t.Fatalf("open test db: %v", err)
	}
	t.Cleanup(func() {
		_ = conn.Close()
		_ = os.Remove(dbPath)
	})

	return New(conn), conn
}
