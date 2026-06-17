// Package refresh owns the periodic data pipeline: fetch Aircall, parse the
// sales spreadsheet, aggregate, and publish a new dashboard into the store. It
// is the only place that mutates the store's data, and it never blanks the
// board on partial failure — each source falls back to its last-good values.
package refresh

import (
	"context"
	"log"
	"sync"
	"time"

	"wallboard/internal/aggregate"
	"wallboard/internal/aircall"
	"wallboard/internal/config"
	"wallboard/internal/excel"
	"wallboard/internal/models"
	"wallboard/internal/store"
)

// Service runs the refresh pipeline.
type Service struct {
	cfg    *config.Config
	air    *aircall.Client
	store  *store.Store
	logger *log.Logger

	mu        sync.Mutex
	salesPath string
}

// New constructs a refresh Service.
func New(cfg *config.Config, air *aircall.Client, st *store.Store, logger *log.Logger) *Service {
	return &Service{
		cfg:       cfg,
		air:       air,
		store:     st,
		logger:    logger,
		salesPath: cfg.SalesXLSXPath,
	}
}

// SetSalesPath swaps the spreadsheet read on subsequent refreshes (used by the
// upload endpoint).
func (s *Service) SetSalesPath(path string) {
	s.mu.Lock()
	s.salesPath = path
	s.mu.Unlock()
}

func (s *Service) currentSalesPath() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.salesPath
}

// Run performs an immediate refresh, then refreshes on the configured interval
// until the context is cancelled.
func (s *Service) Run(ctx context.Context) {
	s.RefreshOnce(ctx)

	ticker := time.NewTicker(s.cfg.RefreshInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			s.logger.Printf("refresh: stopping (%v)", ctx.Err())
			return
		case <-ticker.C:
			s.RefreshOnce(ctx)
		}
	}
}

// RefreshOnce executes a single refresh cycle.
func (s *Service) RefreshOnce(ctx context.Context) {
	now := time.Now()

	s.refreshAircall(ctx, now)
	s.refreshExcel(now)
	s.rebuild(now)
}

func (s *Service) refreshAircall(ctx context.Context, now time.Time) {
	userIDs := make([]int64, 0, len(s.cfg.Staff))
	for _, st := range s.cfg.Staff {
		if st.AircallUserID != 0 {
			userIDs = append(userIDs, st.AircallUserID)
		}
	}

	r := aircall.DayRange(now, s.cfg.Location)
	calls, err := s.air.FetchCallMetrics(ctx, r, userIDs)
	if err != nil {
		s.logger.Printf("refresh: aircall fetch failed: %v (keeping last-good data)", err)
	} else {
		s.logger.Printf("refresh: aircall ok (%s) — metrics for %d users", s.air.Mode(), len(calls))
	}
	s.store.RecordAircall(calls, s.air.Mode(), now, err)
}

func (s *Service) refreshExcel(now time.Time) {
	path := s.currentSalesPath()
	res, err := excel.ParseFile(path, s.cfg.Excel)
	if err != nil {
		s.logger.Printf("refresh: excel parse failed for %s: %v (keeping last-good data)", path, err)
		s.store.RecordExcel(nil, "live", now, err)
		return
	}
	s.logger.Printf("refresh: excel ok — %d rows from %s (%d skipped)", res.Rows, path, res.SkippedRows)
	s.store.RecordExcel(res.BySpreadsheetName, "live", now, nil)
}

// rebuild assembles the dashboard from whatever last-good data exists.
func (s *Service) rebuild(now time.Time) {
	calls, haveCalls, sales, haveSales := s.store.Sources()

	res := aggregate.Build(aggregate.Inputs{
		Staff:      s.cfg.Staff,
		Weights:    s.cfg.Weights,
		Calls:      calls,
		HasAircall: haveCalls,
		Sales:      sales,
		HasExcel:   haveSales,
		PrevRanks:  s.store.PrevRanks(),
	})

	if len(res.UnmatchedSales) > 0 {
		s.logger.Printf("refresh: %d spreadsheet names had no staff mapping: %v",
			len(res.UnmatchedSales), res.UnmatchedSales)
	}

	status := s.store.Status(now)
	resp := &models.DashboardResponse{
		Title:                  s.cfg.DashboardTitle,
		GeneratedAt:            now,
		LastUpdated:            now,
		RefreshIntervalSeconds: int(s.cfg.RefreshInterval.Seconds()),
		People:                 res.People,
		TeamTotals:             res.Totals,
		TopPerformer:           res.Top,
		BiggestMover:           res.Mover,
		SourceStatus:           status,
		Scoring:                models.ScoringConfig{Weights: s.cfg.Weights},
	}

	s.store.SetDashboard(resp, aggregate.RanksByID(res.People))
}
