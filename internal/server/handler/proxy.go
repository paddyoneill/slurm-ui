package handler

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/google/uuid"

	"github.com/paddyoneill/slurm-ui/internal/db"
	servertypes "github.com/paddyoneill/slurm-ui/internal/server/types"
)

func (h *Handler) HandleNotebookProxy(w http.ResponseWriter, r *http.Request) {
	notebookID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, fmt.Errorf("invalid notebook id: %w", err))
		return
	}

	h.proxyNotebookRequest(w, r, notebookID, r.PathValue("path"))

}

func (h *Handler) proxyNotebookRequest(w http.ResponseWriter, r *http.Request, id servertypes.NotebookID, path string) {
	record, err := db.GetNotebookByID(r.Context(), h.db, id)
	if err != nil {
		if errors.Is(err, db.ErrNotebookNotFound) {
			writeError(w, http.StatusNotFound, err)
			return
		}

		writeError(w, http.StatusInternalServerError, err)
		return
	}

	if record.Host == nil || *record.Host == "" {
		writeError(w, http.StatusBadGateway, fmt.Errorf("notebook host is not available"))
		return
	}

	if record.Port < 1 {
		writeError(w, http.StatusBadGateway, fmt.Errorf("notebook port is not available"))
		return
	}

	if path == "" {
		location := fmt.Sprintf("/api/notebooks/%s/proxy/tree", id)
		if r.URL.RawQuery != "" {
			location += "?" + r.URL.RawQuery
		}
		http.Redirect(w, r, location, http.StatusFound)
		return
	}

	upstream, err := url.Parse(fmt.Sprintf("http://%s:%d", *record.Host, record.Port))
	if err != nil {
		writeError(w, http.StatusBadGateway, fmt.Errorf("invalid notebook upstream: %w", err))
		return
	}

	proxy := &httputil.ReverseProxy{
		Rewrite: func(req *httputil.ProxyRequest) {
			req.SetURL(upstream)
			req.Out.URL.Path = notebookProxyBasePath(id) + "/" + strings.TrimPrefix(path, "/")
			req.Out.URL.RawPath = req.Out.URL.Path
			req.Out.URL.RawQuery = req.In.URL.RawQuery
			req.SetXForwarded()
		},
	}

	proxy.ModifyResponse = func(res *http.Response) error {
		location := res.Header.Get("Location")
		if location == "" {
			return nil
		}

		rewritten := rewriteNotebookLocation(id, location, upstream)
		if rewritten != "" {
			res.Header.Set("Location", rewritten)
		}

		return nil
	}

	proxy.ErrorHandler = func(rw http.ResponseWriter, req *http.Request, err error) {
		http.Error(rw, "Failed to proxy notebook", http.StatusBadGateway)
	}

	proxy.ServeHTTP(w, r)
}

func rewriteNotebookLocation(id servertypes.NotebookID, location string, upstream *url.URL) string {
	parsed, err := url.Parse(location)
	if err != nil {
		return ""
	}

	if parsed.IsAbs() {
		if !sameUpstream(parsed, upstream) {
			return ""
		}
		location = parsed.RequestURI()
		parsed, err = url.Parse(location)
		if err != nil {
			return ""
		}
	}

	if !strings.HasPrefix(parsed.Path, "/") {
		return ""
	}

	basePath := notebookProxyBasePath(id)
	if strings.HasPrefix(parsed.Path, basePath) {
		return notebookLocationWithSuffix(parsed.Path, parsed)
	}

	return notebookLocationWithSuffix(basePath+parsed.Path, parsed)
}

func notebookProxyBasePath(id servertypes.NotebookID) string {
	return fmt.Sprintf("/api/notebooks/%s/proxy", id)
}

func notebookLocationWithSuffix(path string, parsed *url.URL) string {
	rewritten := path
	if parsed.RawQuery != "" {
		rewritten += "?" + parsed.RawQuery
	}
	if parsed.Fragment != "" {
		rewritten += "#" + parsed.Fragment
	}
	return rewritten
}

func sameUpstream(location *url.URL, upstream *url.URL) bool {
	return strings.EqualFold(location.Scheme, upstream.Scheme) && strings.EqualFold(location.Host, upstream.Host)
}
