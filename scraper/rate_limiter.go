package scraper

import (
	"time"
)

var (
	// DefaultRateLimiter is 45 queries in 60 second.
	// So roughly less than 1 query per second (1350ms between queries).
	DefaultRateLimiter = NewRateLimiter(45, 60*time.Second)
)

// RateLimiter holds rule for limiting queries in time.
// Strategy is as follows:
// While we are under the rate limiting rule, just download everything.
// When hitting the limit, wait and query at a slower pace.
// The goal here is to not wait for few queries, but still be gentle
// for the one with a lot of call to make.
type RateLimiter struct {
	// NbMaxQuery is the number of queries allowed.
	NbMaxQuery int
	// Window is the window used to count the number of queries.
	Window time.Duration
	// NbQuery is the number of query currently made.
	NbQuery int

	slowModeEnabled    bool
	slowMode           bool
	timeBetweenQueries time.Duration
	lastQuery          time.Time
}

// NewRateLimiter returns a new rate limiter.
func NewRateLimiter(nbMaxQuery int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		NbMaxQuery: nbMaxQuery,
		Window:     window,
		NbQuery:    0,

		slowModeEnabled:    false,
		slowMode:           false,
		timeBetweenQueries: time.Duration(int64(window) / int64(nbMaxQuery)),
		lastQuery:          time.Time{},
	}
}

// nextQuery is the internal method for NextQuery with a choosable "now"
// mainly used to ease unit testing.
func (r *RateLimiter) nextQuery(now time.Time) (time.Duration, func(time.Time)) {
	queryDone := func(t time.Time) {
		r.NbQuery++
		r.lastQuery = t
	}

	if r.NbQuery < r.NbMaxQuery {
		return time.Duration(0), queryDone
	}

	elapsedTime := now.Sub(r.lastQuery)

	// Check if slow mode can be desactivated.
	if elapsedTime > r.Window {
		if elapsedTime > 2*r.Window {
			r.slowMode = false
		}
		r.NbQuery = 0
		return time.Duration(0), queryDone
	}

	// We hit the rate limit for the first time, so wait a bit
	// then switch into "slow mode".
	if !r.slowMode || !r.slowModeEnabled {
		r.slowMode = true
		return r.Window - elapsedTime, queryDone
	}

	return r.timeBetweenQueries - elapsedTime, queryDone
}

// EnableSlowMode enables slow mode.
// This mode allows to query slowly instead of querying by batch.
func (r *RateLimiter) EnableSlowMode() {
	r.slowModeEnabled = true
}

// NextQuery returns the time to wait before performing the next query.
// A closure function is returns which has to be called after the query has been
// called.
func (r *RateLimiter) NextQuery() (time.Duration, func()) {
	waitTime, queryDone := r.nextQuery(time.Now())
	return waitTime, func() { queryDone(time.Now()) }
}
