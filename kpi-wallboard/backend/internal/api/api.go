// Package api exposes the HTTP endpoints consumed by the Svelte frontend.
package api

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"wallboard/internal/config"
	"wallboard/internal/excel"
	"wallboard/internal/refresh"
	"wallboard/internal/store"
)

// Server wires the HTTP handlers to the store and refresh service.
type Server struct {
	cfg     *config.Config
	store   *store.Store
	refresh *refresh.Service
	logger  *log.Logger
	version string
}

// New constructs a Server.
func New(cfg *config.Config, st *store.Store, rs *refresh.Service, logger *log.Logger, version string) *Server {
	return &Server{cfg: cfg, store: st, refresh: rs, logger: logger, version: version}
}

// Handler returns the root http.Handler with all routes and middleware.
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/dashboard", s.handleDashboard)
	mux.HandleFunc("GET /api/health", s.handleHealth)
	mux.HandleFunc("GET /api/status", s.handleStatus)
	mux.HandleFunc("POST /api/sales/upload", s.handleUpload)
	return withCORS(withLogging(s.logger, mux))
}

func (s *Server) handleDashboard(w http.ResponseWriter, r *http.Request) {
	d := s.store.Dashboard(time.Now())
	if d == nil {
		// No successful build yet — tell the frontend to keep waiting rather
		// than serving a confusing empty board.
		writeJSON(w, http.StatusServiceUnavailable, map[string]any{
			"error":   "dashboard not ready yet",
			"message": "initial data refresh has not completed",
		})
		return
	}
	writeJSON(w, http.StatusOK, d)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	d := s.store.Dashboard(now)
	status := "ok"
	if d == nil {
		status = "starting"
	} else if d.Stale {
		status = "degraded"
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"status":        status,
		"version":       s.version,
		"uptimeSeconds": int(now.Sub(s.store.StartedAt()).Seconds()),
		"time":          now.Format(time.RFC3339),
	})
}

func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	st := s.store.Status(now)
	d := s.store.Dashboard(now)
	resp := map[string]any{
		"sourceStatus":           st,
		"refreshIntervalSeconds": int(s.cfg.RefreshInterval.Seconds()),
		"timezone":               s.cfg.Timezone,
	}
	if d != nil {
		resp["lastUpdated"] = d.LastUpdated.Format(time.RFC3339)
		resp["stale"] = d.Stale
		resp["headcount"] = d.TeamTotals.Headcount
	}
	writeJSON(w, http.StatusOK, resp)
}

// handleUpload accepts a multipart .xlsx upload, validates it, persists it, and
// triggers an immediate refresh.
func (s *Server) handleUpload(w http.ResponseWriter, r *http.Request) {
	const maxUpload = 10 << 20 // 10 MiB
	r.Body = http.MaxBytesReader(w, r.Body, maxUpload)
	if err := r.ParseMultipartForm(maxUpload); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid multipart form: " + err.Error()})
		return
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "missing 'file' field"})
		return
	}
	defer file.Close()

	b, err := io.ReadAll(file)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "reading upload: " + err.Error()})
		return
	}

	// Validate by parsing before we overwrite anything.
	if _, err := excel.ParseBytes(b, s.cfg.Excel); err != nil {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]any{
			"error": "uploaded spreadsheet could not be parsed: " + err.Error(),
		})
		return
	}

	dest := s.cfg.SalesXLSXPath
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "preparing destination: " + err.Error()})
		return
	}
	if err := os.WriteFile(dest, b, 0o644); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "saving upload: " + err.Error()})
		return
	}

	s.refresh.SetSalesPath(dest)
	// Refresh immediately so the upload is reflected without waiting for the
	// next tick. Bounded so a slow Aircall call can't hang the request.
	ctx, cancel := context.WithTimeout(r.Context(), 45*time.Second)
	defer cancel()
	s.refresh.RefreshOnce(ctx)

	s.logger.Printf("api: sales spreadsheet replaced via upload (%s, %d bytes)", header.Filename, len(b))
	writeJSON(w, http.StatusOK, map[string]any{
		"status":   "ok",
		"filename": header.Filename,
		"bytes":    len(b),
		"savedTo":  dest,
	})
}

// --- helpers ---

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	_ = enc.Encode(v)
}

// withCORS allows the frontend dev server (different origin) to call the API.
func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// withLogging logs each request with its status and duration.
func withLogging(logger *log.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		sw := &statusWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(sw, r)
		logger.Printf("%s %s -> %d (%s)", r.Method, r.URL.Path, sw.status, time.Since(start).Round(time.Millisecond))
	})
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (s *statusWriter) WriteHeader(code int) {
	s.status = code
	s.ResponseWriter.WriteHeader(code)
}
