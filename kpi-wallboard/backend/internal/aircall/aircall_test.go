package aircall

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestFetchCallMetrics_LivePaginationAndAggregation(t *testing.T) {
	var srv *httptest.Server
	srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify basic auth is sent.
		if id, tok, ok := r.BasicAuth(); !ok || id != "id" || tok != "tok" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Query().Get("page") == "2" {
			fmt.Fprintf(w, `{
				"calls": [
					{"direction":"inbound","duration":30,"user":{"id":2}},
					{"direction":"outbound","duration":0,"user":{"id":1}},
					{"direction":"inbound","duration":90,"user":{"id":1}}
				],
				"meta": {"next_page_link": "", "current_page": 2}
			}`)
			return
		}
		// Page 1, with a link to page 2.
		fmt.Fprintf(w, `{
			"calls": [
				{"direction":"inbound","duration":120,"user":{"id":1}},
				{"direction":"outbound","duration":200,"user":{"id":1}},
				{"direction":"outbound","duration":60,"user":{"id":2}},
				{"direction":"inbound","duration":45,"user":null}
			],
			"meta": {"next_page_link": "%s/calls?page=2", "current_page": 1}
		}`, srv.URL)
	}))
	defer srv.Close()

	c := New(Options{APIID: "id", APIToken: "tok", BaseURL: srv.URL})
	if c.Mode() != "live" {
		t.Fatalf("Mode = %q, want live", c.Mode())
	}

	now := time.Now()
	got, err := c.FetchCallMetrics(context.Background(), DayRange(now, time.UTC), []int64{1, 2})
	if err != nil {
		t.Fatalf("FetchCallMetrics: %v", err)
	}

	// User 1: in=2 (120,90), out=2 (200,0), seconds=410 -> 7 min.
	u1 := got[1]
	if u1.CallsIn != 2 || u1.CallsOut != 2 || u1.TotalCalls != 4 {
		t.Errorf("user1 counts = %+v, want in2 out2 total4", u1)
	}
	if u1.MinutesOnPhone != 7 {
		t.Errorf("user1 minutes = %d, want 7", u1.MinutesOnPhone)
	}

	// User 2: in=1 (30), out=1 (60), seconds=90 -> 2 min (rounded).
	u2 := got[2]
	if u2.CallsIn != 1 || u2.CallsOut != 1 || u2.TotalCalls != 2 {
		t.Errorf("user2 counts = %+v, want in1 out1 total2", u2)
	}
	if u2.MinutesOnPhone != 2 {
		t.Errorf("user2 minutes = %d, want 2", u2.MinutesOnPhone)
	}

	// Call with null user must be ignored entirely.
	if len(got) != 2 {
		t.Errorf("got %d users, want 2 (null-user call ignored)", len(got))
	}
}

func TestFetchCallMetrics_AuthErrorNoRetryLoop(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "nope", http.StatusUnauthorized)
	}))
	defer srv.Close()

	c := New(Options{APIID: "id", APIToken: "bad", BaseURL: srv.URL})
	_, err := c.FetchCallMetrics(context.Background(), DayRange(time.Now(), time.UTC), []int64{1})
	if err == nil {
		t.Fatal("expected auth error")
	}
}

func TestDemoMode_DeterministicAndPopulated(t *testing.T) {
	c := New(Options{}) // no creds -> demo
	if c.Mode() != "demo" {
		t.Fatalf("Mode = %q, want demo", c.Mode())
	}

	// Fix the time so the day fraction is stable.
	now := time.Date(2026, 6, 17, 17, 0, 0, 0, time.UTC)
	r := DayRange(now, time.UTC)
	ids := []int64{1001, 1002, 1003}

	a, _ := c.FetchCallMetrics(context.Background(), r, ids)
	b, _ := c.FetchCallMetrics(context.Background(), r, ids)

	if len(a) != 3 {
		t.Fatalf("demo returned %d users, want 3", len(a))
	}
	for _, id := range ids {
		if a[id] != b[id] {
			t.Errorf("demo not deterministic for %d: %+v vs %+v", id, a[id], b[id])
		}
		if a[id].TotalCalls <= 0 {
			t.Errorf("demo user %d has no calls: %+v", id, a[id])
		}
		if a[id].CallsIn+a[id].CallsOut != a[id].TotalCalls {
			t.Errorf("demo user %d totals inconsistent: %+v", id, a[id])
		}
	}
}

func TestDayRange(t *testing.T) {
	now := time.Date(2026, 6, 17, 14, 30, 0, 0, time.UTC)
	r := DayRange(now, time.UTC)
	if r.From.Hour() != 0 || r.From.Minute() != 0 {
		t.Errorf("DayRange.From = %v, want midnight", r.From)
	}
	if !r.To.Equal(now) {
		t.Errorf("DayRange.To = %v, want %v", r.To, now)
	}
}
