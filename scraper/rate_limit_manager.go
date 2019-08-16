package scraper

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

var (
	errRuleParseLines = errors.New("can't parse any rate limit rules")
	errRuleParseFirst = errors.New("can't parse the first rate limit rules")
	errInvalidRule    = errors.New("invalid rule")
)

// RateRules are the rate limit rules
type RateRules struct {
	NbQuery        int
	WindowSeconds  int // In seconds
	BanTimeSeconds int // In seconds
}

// ExtractFirstRuleFromString extracts rate limit rules from a string.
// Format is: <NbQueries>:<MaxQueries>:<BanTimeSeconds> coma separated.
// Example:
//   2:60:60,2:240:900
// Not that it get only the very first rule.
func ExtractFirstRuleFromString(s string) (RateRules, error) {
	res := RateRules{}
	lines := strings.Split(s, ",")
	if len(lines) <= 0 {
		return res, errRuleParseLines
	}
	rules := strings.Split(lines[0], ":")
	if len(rules) != 3 {
		return res, errRuleParseFirst
	}

	var err error
	if res.NbQuery, err = strconv.Atoi(rules[0]); err != nil {
		return res, errInvalidRule
	}
	if res.WindowSeconds, err = strconv.Atoi(rules[1]); err != nil {
		return res, errInvalidRule
	}
	if res.BanTimeSeconds, err = strconv.Atoi(rules[2]); err != nil {
		return res, errInvalidRule
	}

	return res, nil
}

// keyRules is the key used for associating rate limit rules.
type keyRules struct {
	poesessid string
	url       string
}

// RateLimitManager holds RateLimiter
// Every URL has its own rate limit, so it's maintained separately.
type RateLimitManager struct {
	defaultRateLimiter *RateLimiter
	rules              map[keyRules]*RateLimiter
}

// NewRateLimitManager returns a new rate limit manager.
func NewRateLimitManager() RateLimitManager {
	return RateLimitManager{
		defaultRateLimiter: DefaultRateLimiter,
		rules:              make(map[keyRules]*RateLimiter),
	}
}

// GetRateLimiter returns the best rate limiter given the session and the url
// provided. If none match, a default rate limiter is given.
func (r *RateLimitManager) GetRateLimiter(poesessid, url string) *RateLimiter {
	if rateLimiter, ok := r.rules[keyRules{poesessid, url}]; ok {
		return rateLimiter
	}
	return DefaultRateLimiter
}

// UpdateRateLimiter updates or creates the corresponding rate limiter.
func (r *RateLimitManager) UpdateRateLimiter(poesessid, url string, rules, state RateRules) {
	key := keyRules{poesessid, url}
	rateLimiter, ok := r.rules[key]
	if !ok {
		r.rules[key] = &RateLimiter{}
		rateLimiter = r.rules[key]
	}
	rateLimiter.NbQuery = state.NbQuery
	rateLimiter.NbMaxQuery = rules.NbQuery
	rateLimiter.Window = time.Duration(state.WindowSeconds) * time.Second
	rateLimiter.timeBetweenQueries = time.Duration(int64(rateLimiter.Window) / int64(rateLimiter.NbMaxQuery))
}

// NewPoeRateLimitManager returns a new rate limit manager already fill for path
// of exile API.
//
// GGG rules I got so far:
//  * https://www.pathofexile.com/character-window/get-stash-items
//    x-rate-limit-ip: 45:60:60,200:120:900
//    No more than 45 queries by minute (1 min ban) => 0.75 q/s
//    No more than 200 queries by 2 minutes (15 min ban) => 1.3 q/s
//  * https://www.pathofexile.com/character-window/get-characters
//    x-rate-limit-ip: 60:60:60,200:120:900
//    No more than 60 queries by minute (1 min ban) => 1 q/s
//    No more than 200 queries by 2 minutes (15 min ban) => 1.3 q/s
//  * https://www.pathofexile.com/character-window/get-passive-skills
//    No limit defined, so let's put 120:60:0
func NewPoeRateLimitManager(poesessid string) RateLimitManager {
	r := NewRateLimitManager()
	r.rules[keyRules{poesessid, "pathofexile.com/character-window/get-stash-items"}] = NewRateLimiter(45, 60*time.Second)
	r.rules[keyRules{poesessid, "pathofexile.com/character-window/get-characters"}] = NewRateLimiter(60, 60*time.Second)
	r.rules[keyRules{poesessid, "pathofexile.com/character-window/get-passive-skills"}] = NewRateLimiter(120, 60*time.Second)
	return r
}
