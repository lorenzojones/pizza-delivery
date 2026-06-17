// Package store holds the most recent successful data from each source plus the
// last fully-built dashboard, in a concurrency-safe way. Keeping last-good data
// per source is what lets the wallboard keep showing numbers when a refresh
// fails.
package store

import (
	"sync"
	"time"

	"wallboard/internal/models"
)

// Store is safe for concurrent use.
type Store struct {
	mu sync.RWMutex

	startedAt  time.Time
	staleAfter time.Duration

	// Last-good per-source data, reused when a refresh of that source fails.
	lastCalls   map[int64]models.CallMetrics
	lastSales   map[string]models.SalesMetrics
	haveCalls   bool
	haveSales   bool

	aircall models.SourceStatus
	excel   models.SourceStatus

	dashboard *models.DashboardResponse
	prevRanks map[string]int
}

// New creates a Store. staleAfter is how old a source's last success may be
// before it is flagged stale.
func New(staleAfter time.Duration) *Store {
	return &Store{
		startedAt:  time.Now(),
		staleAfter: staleAfter,
		aircall:    models.SourceStatus{Name: "Aircall"},
		excel:      models.SourceStatus{Name: "Excel"},
	}
}

// RecordAircall stores the outcome of an Aircall fetch. On success the metrics
// are cached as last-good; on failure the previous last-good is retained.
func (s *Store) RecordAircall(calls map[int64]models.CallMetrics, mode string, now time.Time, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.aircall.Mode = mode
	if err != nil {
		s.aircall.OK = false
		s.aircall.LastError = err.Error()
		return
	}
	s.aircall.OK = true
	s.aircall.LastError = ""
	t := now
	s.aircall.LastSuccess = &t
	s.lastCalls = calls
	s.haveCalls = true
}

// RecordExcel stores the outcome of an Excel parse, mirroring RecordAircall.
func (s *Store) RecordExcel(sales map[string]models.SalesMetrics, mode string, now time.Time, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.excel.Mode = mode
	if err != nil {
		s.excel.OK = false
		s.excel.LastError = err.Error()
		return
	}
	s.excel.OK = true
	s.excel.LastError = ""
	t := now
	s.excel.LastSuccess = &t
	s.lastSales = sales
	s.haveSales = true
}

// Sources returns the cached per-source data and whether each is present.
func (s *Store) Sources() (calls map[int64]models.CallMetrics, haveCalls bool, sales map[string]models.SalesMetrics, haveSales bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.lastCalls, s.haveCalls, s.lastSales, s.haveSales
}

// PrevRanks returns the rank snapshot from the previous build.
func (s *Store) PrevRanks() map[string]int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.prevRanks
}

// SetDashboard records a freshly built dashboard and the rank snapshot to use
// as PrevRanks next cycle.
func (s *Store) SetDashboard(d *models.DashboardResponse, ranks map[string]int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.dashboard = d
	s.prevRanks = ranks
}

// Dashboard returns the latest dashboard with staleness recomputed as of now.
// Returns nil if no successful build has happened yet.
func (s *Store) Dashboard(now time.Time) *models.DashboardResponse {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.dashboard == nil {
		return nil
	}
	// Copy so we can patch time-dependent staleness without locking readers out
	// of a consistent view.
	d := *s.dashboard
	d.SourceStatus.Aircall = s.withStaleness(s.aircall, now)
	d.SourceStatus.Excel = s.withStaleness(s.excel, now)
	d.Stale = d.SourceStatus.Aircall.Stale || d.SourceStatus.Excel.Stale ||
		!d.SourceStatus.Aircall.OK || !d.SourceStatus.Excel.OK
	return &d
}

// Status returns the current per-source statuses with staleness as of now.
func (s *Store) Status(now time.Time) models.SourceStatuses {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return models.SourceStatuses{
		Aircall: s.withStaleness(s.aircall, now),
		Excel:   s.withStaleness(s.excel, now),
	}
}

func (s *Store) withStaleness(st models.SourceStatus, now time.Time) models.SourceStatus {
	if st.LastSuccess == nil {
		st.Stale = true
		return st
	}
	st.Stale = now.Sub(*st.LastSuccess) > s.staleAfter
	return st
}

// StartedAt returns the process start time.
func (s *Store) StartedAt() time.Time { return s.startedAt }
