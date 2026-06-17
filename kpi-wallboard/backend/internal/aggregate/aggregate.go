// Package aggregate combines Aircall call metrics and spreadsheet sales metrics
// per staff member, scores and ranks everyone, and derives team totals plus the
// top performer and biggest mover.
package aggregate

import (
	"sort"

	"wallboard/internal/excel"
	"wallboard/internal/models"
	"wallboard/internal/scoring"
)

// Inputs are everything needed to build the per-person board.
type Inputs struct {
	Staff   []models.StaffMember
	Weights models.ScoreWeights

	// Calls is keyed by Aircall user ID. Nil/empty when Aircall data is
	// unavailable this cycle.
	Calls map[int64]models.CallMetrics
	// HasAircall indicates whether Calls reflects a successful fetch.
	HasAircall bool

	// Sales is keyed by the normalised spreadsheet name. Nil/empty when Excel
	// data is unavailable this cycle.
	Sales map[string]models.SalesMetrics
	// HasExcel indicates whether Sales reflects a successful parse.
	HasExcel bool

	// PrevRanks maps staff ID to their rank in the previous build, used to find
	// the biggest mover.
	PrevRanks map[string]int
}

// Result is the assembled board.
type Result struct {
	People         []models.DashboardPerson
	Totals         models.TeamTotals
	Top            *models.DashboardPerson
	Mover          *models.Mover
	UnmatchedSales []string // spreadsheet names not found in the staff mapping
}

// Build assembles the board from the given inputs.
func Build(in Inputs) Result {
	people := make([]models.DashboardPerson, 0, len(in.Staff))
	matchedSales := make(map[string]bool)

	for _, s := range in.Staff {
		p := models.DashboardPerson{StaffMember: s}

		if in.HasAircall {
			if cm, ok := in.Calls[s.AircallUserID]; ok && s.AircallUserID != 0 {
				p.CallMetrics = cm
				p.HasAircallData = true
			}
		}

		if in.HasExcel {
			key := excel.NormName(s.SpreadsheetName)
			if sm, ok := in.Sales[key]; ok {
				p.SalesMetrics = sm
				p.HasSalesData = true
				matchedSales[key] = true
			}
		}

		people = append(people, p)
	}

	// Score, then rank.
	scoring.Apply(people, in.Weights)
	sortByScore(people)
	for i := range people {
		people[i].Rank = i + 1
	}

	res := Result{
		People:  people,
		Totals:  teamTotals(people),
		Mover:   biggestMover(people, in.PrevRanks),
		Top:     topPerformer(people),
		UnmatchedSales: unmatched(in.Sales, matchedSales, in.HasExcel),
	}
	return res
}

// sortByScore orders people by score desc with stable, sensible tie-breakers.
func sortByScore(people []models.DashboardPerson) {
	sort.SliceStable(people, func(i, j int) bool {
		a, b := people[i], people[j]
		if a.Score != b.Score {
			return a.Score > b.Score
		}
		if a.SalesMetrics.SalesValue != b.SalesMetrics.SalesValue {
			return a.SalesMetrics.SalesValue > b.SalesMetrics.SalesValue
		}
		if a.CallMetrics.TotalCalls != b.CallMetrics.TotalCalls {
			return a.CallMetrics.TotalCalls > b.CallMetrics.TotalCalls
		}
		return a.StaffMember.Name < b.StaffMember.Name
	})
}

func teamTotals(people []models.DashboardPerson) models.TeamTotals {
	var t models.TeamTotals
	t.Headcount = len(people)
	for _, p := range people {
		t.CallsIn += p.CallMetrics.CallsIn
		t.CallsOut += p.CallMetrics.CallsOut
		t.TotalCalls += p.CallMetrics.TotalCalls
		t.MinutesOnPhone += p.CallMetrics.MinutesOnPhone
		t.SalesCount += p.SalesMetrics.SalesCount
		t.SalesValue += p.SalesMetrics.SalesValue
	}
	return t
}

func topPerformer(people []models.DashboardPerson) *models.DashboardPerson {
	if len(people) == 0 {
		return nil
	}
	top := people[0]
	return &top
}

// biggestMover finds the person whose rank improved most since PrevRanks. Ties
// are broken by the better (lower) current rank.
func biggestMover(people []models.DashboardPerson, prev map[string]int) *models.Mover {
	if len(prev) == 0 {
		return nil
	}
	var best *models.Mover
	for _, p := range people {
		old, ok := prev[p.StaffMember.ID]
		if !ok {
			continue
		}
		delta := old - p.Rank // positive => climbed
		if delta <= 0 {
			continue
		}
		if best == nil || delta > best.Delta {
			best = &models.Mover{
				StaffID: p.StaffMember.ID,
				Name:    p.StaffMember.Name,
				Delta:   delta,
			}
		}
	}
	return best
}

func unmatched(sales map[string]models.SalesMetrics, matched map[string]bool, hasExcel bool) []string {
	if !hasExcel {
		return nil
	}
	var out []string
	for name := range sales {
		if !matched[name] {
			out = append(out, name)
		}
	}
	sort.Strings(out)
	return out
}

// RanksByID returns a staffID->rank snapshot, used as PrevRanks for the next
// build so the biggest mover can be computed.
func RanksByID(people []models.DashboardPerson) map[string]int {
	m := make(map[string]int, len(people))
	for _, p := range people {
		m[p.StaffMember.ID] = p.Rank
	}
	return m
}
