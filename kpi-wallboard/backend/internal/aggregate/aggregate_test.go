package aggregate

import (
	"testing"

	"wallboard/internal/excel"
	"wallboard/internal/models"
)

func testStaff() []models.StaffMember {
	return []models.StaffMember{
		{ID: "a", Name: "Alice", AircallUserID: 1, SpreadsheetName: "Alice"},
		{ID: "b", Name: "Bob", AircallUserID: 2, SpreadsheetName: "Bob"},
		{ID: "c", Name: "Carol", AircallUserID: 3, SpreadsheetName: "Carol"},
	}
}

func defaultWeights() models.ScoreWeights {
	return models.ScoreWeights{Sales: 0.5, CallsOut: 0.25, CallsIn: 0.1, MinutesOnPhone: 0.15}
}

func TestBuild_RanksScoresAndTotals(t *testing.T) {
	in := Inputs{
		Staff:      testStaff(),
		Weights:    defaultWeights(),
		HasAircall: true,
		HasExcel:   true,
		Calls: map[int64]models.CallMetrics{
			1: {CallsIn: 10, CallsOut: 40, TotalCalls: 50, MinutesOnPhone: 100},
			2: {CallsIn: 5, CallsOut: 20, TotalCalls: 25, MinutesOnPhone: 50},
			3: {CallsIn: 2, CallsOut: 10, TotalCalls: 12, MinutesOnPhone: 25},
		},
		Sales: map[string]models.SalesMetrics{
			excel.NormName("Alice"): {SalesCount: 10, SalesValue: 30000, Target: 20000},
			excel.NormName("Bob"):   {SalesCount: 6, SalesValue: 18000, Target: 20000},
			excel.NormName("Carol"): {SalesCount: 3, SalesValue: 9000, Target: 20000},
		},
	}
	res := Build(in)

	if len(res.People) != 3 {
		t.Fatalf("people = %d, want 3", len(res.People))
	}
	// Alice leads every metric, so she must be rank 1 with score 100.
	if res.People[0].StaffMember.ID != "a" {
		t.Errorf("rank 1 = %s, want a (Alice)", res.People[0].StaffMember.ID)
	}
	if res.People[0].Score != 100 {
		t.Errorf("top score = %v, want 100 (leads all metrics)", res.People[0].Score)
	}
	if res.People[0].Rank != 1 || res.People[2].Rank != 3 {
		t.Errorf("ranks not assigned 1..n: %+v", []int{res.People[0].Rank, res.People[1].Rank, res.People[2].Rank})
	}
	// Scores must be strictly descending given the inputs.
	if !(res.People[0].Score > res.People[1].Score && res.People[1].Score > res.People[2].Score) {
		t.Errorf("scores not descending: %v %v %v", res.People[0].Score, res.People[1].Score, res.People[2].Score)
	}

	// Team totals.
	if res.Totals.CallsOut != 70 {
		t.Errorf("team callsOut = %d, want 70", res.Totals.CallsOut)
	}
	if res.Totals.SalesValue != 57000 {
		t.Errorf("team salesValue = %v, want 57000", res.Totals.SalesValue)
	}
	if res.Totals.Headcount != 3 {
		t.Errorf("headcount = %d, want 3", res.Totals.Headcount)
	}

	// Top performer.
	if res.Top == nil || res.Top.StaffMember.ID != "a" {
		t.Errorf("top performer = %+v, want Alice", res.Top)
	}
}

func TestBuild_MissingDataFlags(t *testing.T) {
	in := Inputs{
		Staff:      testStaff(),
		Weights:    defaultWeights(),
		HasAircall: true,
		HasExcel:   true,
		Calls: map[int64]models.CallMetrics{
			1: {CallsOut: 10, TotalCalls: 10},
			// Bob (2) absent from Aircall.
			3: {CallsOut: 5, TotalCalls: 5},
		},
		Sales: map[string]models.SalesMetrics{
			excel.NormName("Alice"): {SalesValue: 1000},
			// Carol absent from sales.
			excel.NormName("Bob"): {SalesValue: 500},
		},
	}
	res := Build(in)

	byID := map[string]models.DashboardPerson{}
	for _, p := range res.People {
		byID[p.StaffMember.ID] = p
	}
	if !byID["a"].HasAircallData || !byID["a"].HasSalesData {
		t.Errorf("Alice should have both sources")
	}
	if byID["b"].HasAircallData {
		t.Errorf("Bob should be missing Aircall data")
	}
	if !byID["b"].HasSalesData {
		t.Errorf("Bob should have sales data")
	}
	if byID["c"].HasSalesData {
		t.Errorf("Carol should be missing sales data")
	}
	// Everyone still appears even with partial data.
	if len(res.People) != 3 {
		t.Errorf("all staff should appear; got %d", len(res.People))
	}
}

func TestBuild_AircallUnavailableKeepsPeople(t *testing.T) {
	in := Inputs{
		Staff:      testStaff(),
		Weights:    defaultWeights(),
		HasAircall: false, // source down
		HasExcel:   true,
		Sales: map[string]models.SalesMetrics{
			excel.NormName("Alice"): {SalesValue: 1000},
			excel.NormName("Bob"):   {SalesValue: 2000},
			excel.NormName("Carol"): {SalesValue: 1500},
		},
	}
	res := Build(in)
	if len(res.People) != 3 {
		t.Fatalf("people = %d, want 3 even with Aircall down", len(res.People))
	}
	for _, p := range res.People {
		if p.HasAircallData {
			t.Errorf("%s should have no Aircall data when source is down", p.StaffMember.ID)
		}
	}
	// Bob has the most sales, so should rank first.
	if res.People[0].StaffMember.ID != "b" {
		t.Errorf("rank 1 = %s, want b (highest sales)", res.People[0].StaffMember.ID)
	}
}

func TestBuild_BiggestMover(t *testing.T) {
	in := Inputs{
		Staff:      testStaff(),
		Weights:    defaultWeights(),
		HasAircall: false,
		HasExcel:   true,
		Sales: map[string]models.SalesMetrics{
			excel.NormName("Alice"): {SalesValue: 1000},
			excel.NormName("Bob"):   {SalesValue: 3000},
			excel.NormName("Carol"): {SalesValue: 2000},
		},
		// Previously Carol was 1st, Bob 3rd. Now Bob leads => Bob climbed +2.
		PrevRanks: map[string]int{"c": 1, "a": 2, "b": 3},
	}
	res := Build(in)
	if res.Mover == nil {
		t.Fatal("expected a biggest mover")
	}
	if res.Mover.StaffID != "b" {
		t.Errorf("mover = %s, want b (Bob)", res.Mover.StaffID)
	}
	if res.Mover.Delta != 2 {
		t.Errorf("mover delta = %d, want 2", res.Mover.Delta)
	}
}

func TestBuild_NoPrevRanksNoMover(t *testing.T) {
	in := Inputs{
		Staff:      testStaff(),
		Weights:    defaultWeights(),
		HasAircall: false,
		HasExcel:   true,
		Sales: map[string]models.SalesMetrics{
			excel.NormName("Alice"): {SalesValue: 1000},
		},
	}
	if res := Build(in); res.Mover != nil {
		t.Errorf("expected no mover on first build, got %+v", res.Mover)
	}
}

func TestBuild_UnmatchedSales(t *testing.T) {
	in := Inputs{
		Staff:      testStaff(),
		Weights:    defaultWeights(),
		HasAircall: false,
		HasExcel:   true,
		Sales: map[string]models.SalesMetrics{
			excel.NormName("Alice"):   {SalesValue: 1000},
			excel.NormName("Ghost"):   {SalesValue: 500}, // not in staff
		},
	}
	res := Build(in)
	if len(res.UnmatchedSales) != 1 || res.UnmatchedSales[0] != excel.NormName("Ghost") {
		t.Errorf("unmatched = %v, want [ghost]", res.UnmatchedSales)
	}
}
