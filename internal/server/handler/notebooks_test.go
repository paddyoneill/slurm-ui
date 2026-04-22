package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
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

	notebook, err := dbpkg.CreateNotebook(context.Background(), conn, uuid.New(), 5678, "notebook-one", "RUNNING", &host, 8888, "token-123", "register-123", "/envs/base")
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

	_, err := dbpkg.CreateNotebook(context.Background(), conn, uuid.New(), 6789, "stream-notebook", "RUNNING", &host, 8888, "token-456", "register-456", "/envs/base")
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

func TestPostNotebooksIdRegisterRegistersPort(t *testing.T) {
	h, conn := newTestHandler(t)

	notebook, err := dbpkg.CreateNotebook(context.Background(), conn, uuid.New(), 5678, "register-notebook", "RUNNING", nil, 0, "token-123", "register-123", "/envs/base")
	if err != nil {
		t.Fatalf("create notebook: %v", err)
	}

	body := bytes.NewBufferString(`{"port":43123,"registrationToken":"register-123"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/notebooks/"+notebook.ID.String()+"/register", body)
	rr := httptest.NewRecorder()

	h.PostNotebooksIdRegister(rr, req, notebook.ID)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", rr.Code)
	}

	updated, err := dbpkg.GetNotebookByID(context.Background(), conn, notebook.ID)
	if err != nil {
		t.Fatalf("get notebook: %v", err)
	}
	if updated.Port != 43123 {
		t.Fatalf("expected notebook port 43123, got %d", updated.Port)
	}
}

func TestPostNotebooksIdRegisterRejectsInvalidToken(t *testing.T) {
	h, conn := newTestHandler(t)

	notebook, err := dbpkg.CreateNotebook(context.Background(), conn, uuid.New(), 5678, "register-notebook", "RUNNING", nil, 0, "token-123", "register-123", "/envs/base")
	if err != nil {
		t.Fatalf("create notebook: %v", err)
	}

	body := bytes.NewBufferString(`{"port":43123,"registrationToken":"wrong-token"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/notebooks/"+notebook.ID.String()+"/register", body)
	rr := httptest.NewRecorder()

	h.PostNotebooksIdRegister(rr, req, notebook.ID)

	if rr.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d", rr.Code)
	}

	updated, err := dbpkg.GetNotebookByID(context.Background(), conn, notebook.ID)
	if err != nil {
		t.Fatalf("get notebook: %v", err)
	}
	if updated.Port != 0 {
		t.Fatalf("expected notebook port to remain unregistered, got %d", updated.Port)
	}
}

func TestHandleNotebookProxyRequiresRegisteredPort(t *testing.T) {
	h, conn := newTestHandler(t)
	host := "compute-03"

	notebook, err := dbpkg.CreateNotebook(context.Background(), conn, uuid.New(), 5678, "proxy-notebook", "RUNNING", &host, 0, "token-123", "register-123", "/envs/base")
	if err != nil {
		t.Fatalf("create notebook: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/notebooks/"+notebook.ID.String()+"/proxy/tree", nil)
	rr := httptest.NewRecorder()

	h.proxyNotebookRequest(rr, req, notebook.ID, "tree")

	if rr.Code != http.StatusBadGateway {
		t.Fatalf("expected status 502, got %d", rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "notebook port is not available") {
		t.Fatalf("expected missing port error, got %q", rr.Body.String())
	}
}

func TestBuildNotebookRegisterURLUsesConfiguredBaseURL(t *testing.T) {
	t.Setenv("SLURM_UI_BASE_URL", "https://slurm-ui.example.com/base")

	notebookID := uuid.MustParse("11111111-1111-1111-1111-111111111111")

	registerURL, err := buildNotebookRegisterURL(notebookID)
	if err != nil {
		t.Fatalf("build register url: %v", err)
	}

	expected := "https://slurm-ui.example.com/base/api/notebooks/11111111-1111-1111-1111-111111111111/register"
	if registerURL != expected {
		t.Fatalf("expected register url %q, got %q", expected, registerURL)
	}
}

func TestBuildNotebookRegisterURLRequiresBaseURL(t *testing.T) {
	notebookID := uuid.MustParse("22222222-2222-2222-2222-222222222222")

	_, err := buildNotebookRegisterURL(notebookID)
	if err == nil {
		t.Fatal("expected missing base URL to fail")
	}
	if !strings.Contains(err.Error(), "SLURM_UI_BASE_URL is required") {
		t.Fatalf("expected missing base URL error, got %v", err)
	}
}

func TestBuildNotebookScriptRegistersAfterNotebookStarts(t *testing.T) {
	script := buildNotebookScript()

	startIndex := strings.Index(script, `jupyter server --ip=0.0.0.0 --no-browser`)
	runtimeIndex := strings.Index(script, `from jupyter_server.serverapp import list_running_servers as list_jp_servers`)
	waitIndex := strings.Index(script, `socket.create_connection(("127.0.0.1", port), timeout=1)`)
	registerIndex := strings.Index(script, `urllib.request.Request(`)
	finalWaitIndex := strings.Index(script, `wait "${NOTEBOOK_PID}"`)

	if startIndex == -1 || runtimeIndex == -1 || waitIndex == -1 || registerIndex == -1 || finalWaitIndex == -1 {
		t.Fatalf("expected notebook script to start notebook, read runtime metadata, wait for the bound port, register the port, and wait on the process")
	}
	if strings.Contains(script, `sock.bind(("0.0.0.0", 0))`) {
		t.Fatalf("expected notebook script not to preselect a port; got script:\n%s", script)
	}
	if !(startIndex < runtimeIndex && runtimeIndex < waitIndex && waitIndex < registerIndex && registerIndex < finalWaitIndex) {
		t.Fatalf("expected notebook registration flow to happen after startup; got script:\n%s", script)
	}
}

func TestRewriteNotebookLocationDoesNotDoublePrefixBasePath(t *testing.T) {
	notebookID := uuid.MustParse("33333333-3333-3333-3333-333333333333")
	upstream, err := url.Parse("http://compute-01:43123")
	if err != nil {
		t.Fatalf("parse upstream: %v", err)
	}

	location := "http://compute-01:43123/api/notebooks/33333333-3333-3333-3333-333333333333/proxy/tree?token=abc"
	rewritten := rewriteNotebookLocation(notebookID, location, upstream)
	expected := "/api/notebooks/33333333-3333-3333-3333-333333333333/proxy/tree?token=abc"

	if rewritten != expected {
		t.Fatalf("expected rewritten location %q, got %q", expected, rewritten)
	}
}

func TestNotebookProxyBasePath(t *testing.T) {
	notebookID := uuid.MustParse("44444444-4444-4444-4444-444444444444")
	expected := "/api/notebooks/44444444-4444-4444-4444-444444444444/proxy"

	if got := notebookProxyBasePath(notebookID); got != expected {
		t.Fatalf("expected proxy base path %q, got %q", expected, got)
	}
}
