package scraper

import (
	"testing"
	"time"
)

var now = time.Now()

// nextQuery is an helper to ease test writing.
func nextQuery(r *RateLimiter, now time.Time) time.Duration {
	d, queryDone := r.nextQuery(now)
	queryDone(now)
	return d
}

func equalDuration(t *testing.T, expected, got time.Duration) {
	if expected != got {
		t.Helper()
		t.Fatalf("expected %v but got %v\n", expected, got)
	}
}

// TestRateLimit tests basic rate limiting works.
func TestSimpleRateLimit(t *testing.T) {
	r := NewRateLimiter(2, 5*time.Second)
	r.EnableSlowMode()
	equalDuration(t, 2500*time.Millisecond, r.timeBetweenQueries)
	equalDuration(t, time.Duration(0), nextQuery(r, now))
	equalDuration(t, time.Duration(0), nextQuery(r, now))
	equalDuration(t, 5*time.Second, nextQuery(r, now))
	equalDuration(t, 2500*time.Millisecond, nextQuery(r, now))
	equalDuration(t, 2500*time.Millisecond, nextQuery(r, now))
}

// TestLongRateLimit tests rate limiting with a log of queries works.
func TestLongRateLimit(t *testing.T) {
	r := NewRateLimiter(20, 5*time.Second)
	r.EnableSlowMode()
	equalDuration(t, 250*time.Millisecond, r.timeBetweenQueries)
	for i := 0; i < 20; i++ {
		equalDuration(t, time.Duration(0), nextQuery(r, now))
	}
	equalDuration(t, 5*time.Second, nextQuery(r, now))
	equalDuration(t, 250*time.Millisecond, nextQuery(r, now))
	equalDuration(t, 250*time.Millisecond, nextQuery(r, now))
}

// TestSparseRateLimit tests rate limiting is correctly disabled if not querying
// a lot anymore.
func TestSparseRateLimit(t *testing.T) {
	r := NewRateLimiter(20, 5*time.Second)
	r.EnableSlowMode()
	equalDuration(t, 250*time.Millisecond, r.timeBetweenQueries)
	for i := 0; i < 20; i++ {
		equalDuration(t, time.Duration(0), nextQuery(r, now))
	}
	equalDuration(t, 5*time.Second, nextQuery(r, now))
	equalDuration(t, 250*time.Millisecond, nextQuery(r, now))
	equalDuration(t, 250*time.Millisecond, nextQuery(r, now))

	// One hour later, should have been reset.
	equalDuration(t, time.Duration(0), nextQuery(r, now.Add(time.Hour)))
}

// TestContinuousRateLimit tests rate limiting is correctly applied when
// requesting a lot.
func TestContinuousRateLimit(t *testing.T) {
	elapsed := now
	r := NewRateLimiter(20, 5*time.Second)
	r.EnableSlowMode()
	equalDuration(t, 250*time.Millisecond, r.timeBetweenQueries)
	for i := 0; i < 20; i++ {
		elapsed = elapsed.Add(100 * time.Millisecond)
		equalDuration(t, time.Duration(0), nextQuery(r, elapsed))
	}
	// Now, have to wait the size of the window.
	elapsed = elapsed.Add(100 * time.Millisecond)
	equalDuration(t, 4900*time.Millisecond, nextQuery(r, elapsed))
	// Start to the slow pace.
	elapsed = elapsed.Add(100 * time.Millisecond)
	equalDuration(t, 150*time.Millisecond, nextQuery(r, elapsed))

	// Should stay at slow pace.
	elapsed = elapsed.Add(100 * time.Millisecond)
	equalDuration(t, 150*time.Millisecond, nextQuery(r, elapsed))
	elapsed = elapsed.Add(100 * time.Millisecond)
	equalDuration(t, 150*time.Millisecond, nextQuery(r, elapsed))

	// Should have been reset.
	elapsed = elapsed.Add(1 * time.Hour)
	equalDuration(t, time.Duration(0), nextQuery(r, elapsed))
}

// TestManyCyclesCallingAPI tests rate limiting is correctly applied when
// requesting a lot.
func TestManyCyclesCallingAPI(t *testing.T) {
	elapsed := now
	r := NewRateLimiter(45, 60*time.Second)
	// No slow mode, so it will wait for the entire window to be done.
	for j := 0; j < 10; j++ {
		for i := 0; i < 45; i++ {
			equalDuration(t, time.Duration(0), nextQuery(r, elapsed))
		}
		// Now, have to wait the size of the window.
		equalDuration(t, 60*time.Second, nextQuery(r, elapsed))
		elapsed = elapsed.Add(61 * time.Second)
	}
}
