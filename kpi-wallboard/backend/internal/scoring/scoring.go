// Package scoring computes the configurable competitive score used to rank
// staff. Each metric is normalised against the team's best performer in that
// metric, then combined using the configured weights and scaled to 0..100.
package scoring

import "wallboard/internal/models"

// Apply computes and sets the Score field for every person in place. Scores are
// normalised against the strongest performer per metric so the leaderboard
// always spans a meaningful range. The returned bool reports whether the team
// uses sales value (true) or sales count (false) as the sales metric.
func Apply(people []models.DashboardPerson, w models.ScoreWeights) bool {
	if len(people) == 0 {
		return false
	}

	// Decide which sales metric to use: value if anyone has a value, else count.
	useValue := false
	for i := range people {
		if people[i].SalesMetrics.SalesValue > 0 {
			useValue = true
			break
		}
	}

	salesOf := func(p models.DashboardPerson) float64 {
		if useValue {
			return p.SalesMetrics.SalesValue
		}
		return float64(p.SalesMetrics.SalesCount)
	}

	// Find the max of each metric across the team.
	var maxSales, maxOut, maxIn, maxMin float64
	for i := range people {
		maxSales = max(maxSales, salesOf(people[i]))
		maxOut = max(maxOut, float64(people[i].CallMetrics.CallsOut))
		maxIn = max(maxIn, float64(people[i].CallMetrics.CallsIn))
		maxMin = max(maxMin, float64(people[i].CallMetrics.MinutesOnPhone))
	}

	weightSum := w.Sales + w.CallsOut + w.CallsIn + w.MinutesOnPhone
	if weightSum <= 0 {
		weightSum = 1
	}

	for i := range people {
		p := &people[i]
		combined := w.Sales*ratio(salesOf(*p), maxSales) +
			w.CallsOut*ratio(float64(p.CallMetrics.CallsOut), maxOut) +
			w.CallsIn*ratio(float64(p.CallMetrics.CallsIn), maxIn) +
			w.MinutesOnPhone*ratio(float64(p.CallMetrics.MinutesOnPhone), maxMin)
		p.Score = round1(100 * combined / weightSum)
	}
	return useValue
}

// ratio returns v/maxV, or 0 when maxV is 0.
func ratio(v, maxV float64) float64 {
	if maxV <= 0 {
		return 0
	}
	return v / maxV
}

func round1(f float64) float64 {
	return float64(int(f*10+0.5)) / 10
}
