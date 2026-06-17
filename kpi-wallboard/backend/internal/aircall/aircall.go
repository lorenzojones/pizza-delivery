// Package aircall fetches and aggregates per-user call metrics from the Aircall
// REST API. When no credentials are configured it produces deterministic demo
// data so the wallboard runs out of the box.
package aircall

import (
	"context"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io"
	"log"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"wallboard/internal/models"
)

// Range is a time window to fetch calls for. Designed so the backend can later
// support weekly/monthly views by passing a different range.
type Range struct {
	From time.Time
	To   time.Time
}

// DayRange returns the window covering "today" in the given location, from
// midnight to now.
func DayRange(now time.Time, loc *time.Location) Range {
	local := now.In(loc)
	start := time.Date(local.Year(), local.Month(), local.Day(), 0, 0, 0, 0, loc)
	return Range{From: start, To: local}
}

// WeekRange returns the window from the start of the current ISO week (Monday)
// to now. Provided for future weekly views.
func WeekRange(now time.Time, loc *time.Location) Range {
	local := now.In(loc)
	weekday := (int(local.Weekday()) + 6) % 7 // Monday=0
	start := time.Date(local.Year(), local.Month(), local.Day(), 0, 0, 0, 0, loc).AddDate(0, 0, -weekday)
	return Range{From: start, To: local}
}

// MonthRange returns the window from the first of the current month to now.
func MonthRange(now time.Time, loc *time.Location) Range {
	local := now.In(loc)
	start := time.Date(local.Year(), local.Month(), 1, 0, 0, 0, 0, loc)
	return Range{From: start, To: local}
}

// Client talks to the Aircall API (or produces demo data).
type Client struct {
	apiID    string
	apiToken string
	baseURL  string
	http     *http.Client
	demo     bool
	logger   *log.Logger
}

// Options configures a Client.
type Options struct {
	APIID    string
	APIToken string
	BaseURL  string
	Logger   *log.Logger
}

// New constructs a Client. If APIID or APIToken is empty the client runs in
// demo mode.
func New(opt Options) *Client {
	if opt.BaseURL == "" {
		opt.BaseURL = "https://api.aircall.io/v1"
	}
	if opt.Logger == nil {
		opt.Logger = log.Default()
	}
	return &Client{
		apiID:    opt.APIID,
		apiToken: opt.APIToken,
		baseURL:  opt.BaseURL,
		http:     &http.Client{Timeout: 30 * time.Second},
		demo:     opt.APIID == "" || opt.APIToken == "",
		logger:   opt.Logger,
	}
}

// Mode reports "live" or "demo".
func (c *Client) Mode() string {
	if c.demo {
		return "demo"
	}
	return "live"
}

// FetchCallMetrics returns call metrics keyed by Aircall user ID for the given
// time range. Only users that appear in wantUserIDs are guaranteed in the demo
// output; live mode returns metrics for every user that had calls.
func (c *Client) FetchCallMetrics(ctx context.Context, r Range, wantUserIDs []int64) (map[int64]models.CallMetrics, error) {
	if c.demo {
		return c.demoMetrics(r, wantUserIDs), nil
	}
	return c.liveMetrics(ctx, r)
}

// --- live implementation ---

// apiCall is the subset of the Aircall call object we use.
type apiCall struct {
	Direction string `json:"direction"` // "inbound" | "outbound"
	Duration  int    `json:"duration"`  // seconds
	User      *struct {
		ID int64 `json:"id"`
	} `json:"user"`
}

type callsPage struct {
	Calls []apiCall `json:"calls"`
	Meta  struct {
		NextPageLink string `json:"next_page_link"`
		CurrentPage  int    `json:"current_page"`
	} `json:"meta"`
}

const maxPages = 200 // hard safety cap (200 * 50 = 10k calls/day)

func (c *Client) liveMetrics(ctx context.Context, r Range) (map[int64]models.CallMetrics, error) {
	out := make(map[int64]models.CallMetrics)

	next := fmt.Sprintf("%s/calls?from=%d&to=%d&per_page=50&order=asc",
		c.baseURL, r.From.Unix(), r.To.Unix())

	for page := 0; next != "" && page < maxPages; page++ {
		body, err := c.getWithRetry(ctx, next)
		if err != nil {
			return nil, err
		}
		var p callsPage
		if err := json.Unmarshal(body, &p); err != nil {
			return nil, fmt.Errorf("decoding calls page: %w", err)
		}
		for _, call := range p.Calls {
			if call.User == nil {
				continue
			}
			m := out[call.User.ID]
			switch call.Direction {
			case "inbound":
				m.CallsIn++
			case "outbound":
				m.CallsOut++
			default:
				// Unknown direction still counts toward totals below.
			}
			m.TotalCalls = m.CallsIn + m.CallsOut
			// Accumulate seconds in MinutesOnPhone temporarily; convert below.
			m.MinutesOnPhone += call.Duration
			out[call.User.ID] = m
		}
		next = p.Meta.NextPageLink
	}

	// Convert accumulated seconds to whole minutes.
	for id, m := range out {
		m.MinutesOnPhone = int(math.Round(float64(m.MinutesOnPhone) / 60.0))
		out[id] = m
	}
	return out, nil
}

// getWithRetry performs a GET with basic auth, retrying on 429/5xx with respect
// for Retry-After, up to a small number of attempts.
func (c *Client) getWithRetry(ctx context.Context, rawURL string) ([]byte, error) {
	const attempts = 4
	var lastErr error
	for attempt := 0; attempt < attempts; attempt++ {
		if attempt > 0 {
			backoff := time.Duration(math.Pow(2, float64(attempt))) * time.Second
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(backoff):
			}
		}

		// Ensure the URL is valid/absolute (next_page_link is absolute).
		if _, err := url.Parse(rawURL); err != nil {
			return nil, fmt.Errorf("bad url %q: %w", rawURL, err)
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
		if err != nil {
			return nil, err
		}
		req.SetBasicAuth(c.apiID, c.apiToken)
		req.Header.Set("Accept", "application/json")

		resp, err := c.http.Do(req)
		if err != nil {
			lastErr = err
			continue
		}
		body, readErr := io.ReadAll(resp.Body)
		resp.Body.Close()
		if readErr != nil {
			lastErr = readErr
			continue
		}

		switch {
		case resp.StatusCode == http.StatusOK:
			return body, nil
		case resp.StatusCode == http.StatusTooManyRequests:
			wait := parseRetryAfter(resp.Header.Get("Retry-After"))
			c.logger.Printf("aircall: rate limited, waiting %s", wait)
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(wait):
			}
			lastErr = fmt.Errorf("rate limited (429)")
		case resp.StatusCode >= 500:
			lastErr = fmt.Errorf("aircall server error %d: %s", resp.StatusCode, truncate(body))
		case resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden:
			// Auth errors will not recover by retrying.
			return nil, fmt.Errorf("aircall auth error %d: %s", resp.StatusCode, truncate(body))
		default:
			return nil, fmt.Errorf("aircall error %d: %s", resp.StatusCode, truncate(body))
		}
	}
	return nil, fmt.Errorf("aircall request failed after %d attempts: %w", attempts, lastErr)
}

func parseRetryAfter(h string) time.Duration {
	if h == "" {
		return 5 * time.Second
	}
	if secs, err := strconv.Atoi(h); err == nil {
		return time.Duration(secs) * time.Second
	}
	return 5 * time.Second
}

func truncate(b []byte) string {
	const max = 200
	if len(b) > max {
		return string(b[:max]) + "…"
	}
	return string(b)
}

// --- demo implementation ---

// demoMetrics produces deterministic, plausible call numbers per user that are
// stable within a given day but grow as the day progresses, so the wallboard
// looks alive during demos.
func (c *Client) demoMetrics(r Range, userIDs []int64) map[int64]models.CallMetrics {
	out := make(map[int64]models.CallMetrics, len(userIDs))
	day := r.To.Format("2006-01-02")
	// Fraction of the working day elapsed (0..1), so numbers ramp up.
	frac := dayFraction(r.To)

	for _, id := range userIDs {
		seed := hashSeed(fmt.Sprintf("%s|%d", day, id))
		// Base daily targets vary per person/day.
		baseOut := 25 + int(seed%40)   // 25..64 outbound by end of day
		baseIn := 8 + int((seed>>8)%22) // 8..29 inbound by end of day
		avgMin := 2 + int((seed>>16)%5) // 2..6 avg minutes per call

		callsOut := int(math.Round(float64(baseOut) * frac))
		callsIn := int(math.Round(float64(baseIn) * frac))
		total := callsOut + callsIn
		out[id] = models.CallMetrics{
			CallsIn:        callsIn,
			CallsOut:       callsOut,
			TotalCalls:     total,
			MinutesOnPhone: total * avgMin,
		}
	}
	return out
}

// dayFraction returns how far through a 9am–6pm working day the given time is,
// clamped to [0.05, 1].
func dayFraction(t time.Time) float64 {
	startHour, endHour := 9.0, 18.0
	h := float64(t.Hour()) + float64(t.Minute())/60.0
	frac := (h - startHour) / (endHour - startHour)
	if frac < 0.05 {
		frac = 0.05
	}
	if frac > 1 {
		frac = 1
	}
	return frac
}

func hashSeed(s string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	return h.Sum64()
}
