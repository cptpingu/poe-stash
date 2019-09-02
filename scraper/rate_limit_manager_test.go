package scraper

import (
	"testing"
	"time"
)

func noError(t *testing.T, err error) {
	if err != nil {
		t.Helper()
		t.Fatalf("expected not error but got %v\n", err)
	}
}

func errorEqual(t *testing.T, expected, got error) {
	if expected != got {
		t.Helper()
		t.Fatalf("expected %v but got %v\n", expected, got)
	}
}

func rulesEqual(t *testing.T, expected, got RateRules) {
	if expected != got {
		t.Helper()
		t.Fatalf("expected %v but got %v\n", expected, got)
	}
}

func intEqual(t *testing.T, expected, got int) {
	if expected != got {
		t.Helper()
		t.Fatalf("expected %v but got %v\n", expected, got)
	}
}

// TestExtractFirstRuleFromString tests basic rate limiting works.
func TestExtractFirstRuleFromString(t *testing.T) {
	var r RateRules
	var err error

	_, err = ExtractFirstRuleFromString("")
	errorEqual(t, errRuleParseFirst, err)
	_, err = ExtractFirstRuleFromString("errRuleParseFirst")
	errorEqual(t, errRuleParseFirst, err)
	_, err = ExtractFirstRuleFromString("34:34:34:34")
	errorEqual(t, errRuleParseFirst, err)
	_, err = ExtractFirstRuleFromString("a:b:c")
	errorEqual(t, errInvalidRule, err)

	r, err = ExtractFirstRuleFromString("1:2:3")
	noError(t, err)
	rulesEqual(t, RateRules{1, 2, 3}, r)

	r, err = ExtractFirstRuleFromString("2:60:60,2:240:900")
	noError(t, err)
	rulesEqual(t, RateRules{2, 60, 60}, r)
}

// TestNotFoundRateLimitManager tests default rate limit manager is used.
func TestNotFoundRateLimitManager(t *testing.T) {
	r := NewRateLimitManager()
	rateLimiter := r.GetRateLimiter("sess", "http://unknown.com")
	if *rateLimiter != *DefaultRateLimiter {
		t.Fatalf("expected %v but got %v\n", DefaultRateLimiter, rateLimiter)
	}
}

// TestSimpleRateLimitManager tests rate limit manager works.
func TestSimpleRateLimitManager(t *testing.T) {
	r := NewRateLimitManager()

	sess := "sessid"
	url := "http://toto.com"
	rules := RateRules{45, 60, 0}
	state := RateRules{4, 60, 0}

	r.UpdateRateLimiter(sess, url, rules, state)
	rateLimiter := r.GetRateLimiter(sess, url)
	intEqual(t, state.NbQuery, rateLimiter.NbQuery)
	intEqual(t, rules.NbQuery, rateLimiter.NbMaxQuery)
	intEqual(t, rules.WindowSeconds, int(rateLimiter.Window)/int(time.Second))
}
