// Package models defines the internal data structures shared across the
// backend and the JSON shapes returned by the HTTP API. The JSON tags here are
// the contract the Svelte frontend depends on, so change them with care.
package models

import "time"

// StaffMember is a single person on the sales team and how they map across the
// two data sources (Aircall and the sales spreadsheet).
type StaffMember struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Email           string `json:"email"`
	AircallUserID   int64  `json:"aircallUserId"`
	SpreadsheetName string `json:"spreadsheetName"`
}

// CallMetrics holds Aircall-derived telephony numbers for one person.
type CallMetrics struct {
	CallsIn        int `json:"callsIn"`
	CallsOut       int `json:"callsOut"`
	TotalCalls     int `json:"totalCalls"`
	MinutesOnPhone int `json:"minutesOnPhone"`
}

// SalesMetrics holds spreadsheet-derived sales numbers for one person.
// SalesValue is a currency amount; SalesCount is a unit count. Target is the
// optional quota (0 means "no target"). TargetProgress is SalesValue/Target (or
// SalesCount/Target when no value column exists), expressed as a 0..1+ fraction.
type SalesMetrics struct {
	SalesCount     int     `json:"salesCount"`
	SalesValue     float64 `json:"salesValue"`
	Target         float64 `json:"target"`
	TargetProgress float64 `json:"targetProgress"`
}

// DashboardPerson combines all sources for one staff member plus the derived
// competitive score and rank.
type DashboardPerson struct {
	StaffMember    StaffMember  `json:"staffMember"`
	CallMetrics    CallMetrics  `json:"callMetrics"`
	SalesMetrics   SalesMetrics `json:"salesMetrics"`
	Score          float64      `json:"score"`
	Rank           int          `json:"rank"`
	HasAircallData bool         `json:"hasAircallData"`
	HasSalesData   bool         `json:"hasSalesData"`
}

// TeamTotals are the summed metrics across everyone on the board.
type TeamTotals struct {
	CallsIn        int     `json:"callsIn"`
	CallsOut       int     `json:"callsOut"`
	TotalCalls     int     `json:"totalCalls"`
	MinutesOnPhone int     `json:"minutesOnPhone"`
	SalesCount     int     `json:"salesCount"`
	SalesValue     float64 `json:"salesValue"`
	Headcount      int     `json:"headcount"`
}

// Mover describes the person whose rank improved the most since the previous
// successful refresh. Delta is positive when they climbed the leaderboard.
type Mover struct {
	StaffID string `json:"staffId"`
	Name    string `json:"name"`
	Delta   int    `json:"delta"`
}

// SourceStatus reports the health of a single upstream data source.
type SourceStatus struct {
	Name        string     `json:"name"`
	OK          bool       `json:"ok"`
	Stale       bool       `json:"stale"`
	LastSuccess *time.Time `json:"lastSuccess"`
	LastError   string     `json:"lastError"`
	Mode        string     `json:"mode"` // "live" or "demo"
}

// ScoringConfig is echoed back so the frontend can show how the score is built.
type ScoringConfig struct {
	Weights ScoreWeights `json:"weights"`
}

// ScoreWeights are the configurable contributions of each metric to the score.
type ScoreWeights struct {
	Sales          float64 `json:"sales"`
	CallsOut       float64 `json:"callsOut"`
	CallsIn        float64 `json:"callsIn"`
	MinutesOnPhone float64 `json:"minutesOnPhone"`
}

// SourceStatuses groups the per-source health blocks.
type SourceStatuses struct {
	Aircall SourceStatus `json:"aircall"`
	Excel   SourceStatus `json:"excel"`
}

// DashboardResponse is the full payload returned by GET /api/dashboard.
type DashboardResponse struct {
	Title                  string           `json:"title"`
	GeneratedAt            time.Time        `json:"generatedAt"`
	LastUpdated            time.Time        `json:"lastUpdated"`
	Stale                  bool             `json:"stale"`
	RefreshIntervalSeconds int              `json:"refreshIntervalSeconds"`
	People                 []DashboardPerson `json:"people"`
	TeamTotals             TeamTotals       `json:"teamTotals"`
	TopPerformer           *DashboardPerson `json:"topPerformer"`
	BiggestMover           *Mover           `json:"biggestMover"`
	SourceStatus           SourceStatuses   `json:"sourceStatus"`
	Scoring                ScoringConfig    `json:"scoring"`
}
