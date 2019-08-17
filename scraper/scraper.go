package scraper

import (
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/poe-stash/inventory"
)

const (
	// StashURL is the official URL for the getting of a user account stash.
	StashURL = "https://www.pathofexile.com/character-window/get-stash-items?accountName=%s&realm=%s&league=%s&tabs=%d&tabIndex=%d"
	// ViewProfileURL is the official URL for the getting a user account main profile information.
	ViewProfileURL = "http://www.pathofexile.com/account/view-profile/%s"
	// ProfileCharactersURL is the official URL for the getting a user account characters.
	ProfileCharactersURL = "https://pathofexile.com/character-window/get-characters?accountName=%s"
	// ProfileCharacterItemsURL is the official URL for the getting a user account inventories.
	ProfileCharacterItemsURL = "https://www.pathofexile.com/character-window/get-items?accountName=%s&realm=%s&character=%s"
	// ProfileCharacterSkillsURL is the official URL for getting a user skills and jewels/abyss put in it.
	ProfileCharacterSkillsURL = "https://www.pathofexile.com/character-window/get-passive-skills?character=%s&accountName=%s"
	// LeaguesURL is the official URL for getting the list of all leagues.
	LeaguesURL = "http://api.pathofexile.com/leagues?type=main&compact=1"

	// DataDir is where all data are.
	DataDir = "data/"
	// DataCacheDir is the cache directory.
	DataCacheDir = DataDir + "cache/"
)

// Scraper scraps path of exile site using its API.
type Scraper struct {
	verbosity int
	cache     bool
	cacheDir  string

	accountName  string
	poeSessionID string
	realm        string
	league       string

	client           http.Client
	rateLimitManager RateLimitManager
}

// ScrapedData holds everything scrapped.
type ScrapedData struct {
	Account    string
	League     string
	Realm      string
	Date       time.Time
	Characters []*inventory.CharacterInventory
	Skills     []*inventory.CharacterSkills
	Stash      []*inventory.StashTab
	Wealth     inventory.WealthBreakdown
}

// NewScraper returns a configured scraper.
func NewScraper(accountName, poeSessionID, realm, league string) *Scraper {
	return &Scraper{
		cacheDir:         DataCacheDir,
		accountName:      accountName,
		poeSessionID:     poeSessionID,
		realm:            realm,
		league:           league,
		client:           http.Client{},
		rateLimitManager: NewPoeRateLimitManager(poeSessionID),
	}
}

// EnableCache enable caching of queries.
// Useful for debug, do not enable it in production.
func (s *Scraper) EnableCache() {
	s.cache = true
}

// SetVerbosity set verbosity of logs.
func (s *Scraper) SetVerbosity(v int) {
	s.verbosity = v
}

// hash url into a number.
func hash(s string) string {
	h := fnv.New32a()
	h.Write([]byte(s))
	return strconv.Itoa(int(h.Sum32()))
}

// CallAPI calls a distant API and returns the content.
func (s *Scraper) CallAPI(apiURL string) ([]byte, error) {
	var fileCache string
	if s.cache {
		fileCache = s.cacheDir + hash(apiURL)
		if b, err := ioutil.ReadFile(fileCache); err != nil {
			fmt.Println("can't read cache", err)
		} else {
			return b, nil
		}
	}

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}
	cookie := http.Cookie{
		Name:  "POESESSID",
		Value: s.poeSessionID,
	}
	req.AddCookie(&cookie)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Handle rate limiting.
	baseURL := req.URL.Hostname() + req.URL.EscapedPath()
	rateLimiter := s.rateLimitManager.GetRateLimiter(s.poeSessionID, baseURL)

	waitTime, queryDone := rateLimiter.NextQuery()
	if s.verbosity > 0 {
		fmt.Println("wait:", waitTime, "query:", apiURL)
		if s.verbosity > 1 {
			fmt.Println("request:", req)
		}
	}
	time.Sleep(waitTime)

	// Query the server.
	resp, errResponse := s.client.Do(req)
	queryDone()
	if errResponse != nil {
		return nil, errResponse
	}

	// Let check if there are some rate limiting rules
	rateLimitRules := resp.Header.Get("X-Rate-Limit-Account")
	if rateLimitRules == "" {
		rateLimitRules = resp.Header.Get("X-Rate-Limit-Ip")
	}
	rateLimitState := resp.Header.Get("X-Rate-Limit-Account-State")
	if rateLimitState == "" {
		rateLimitState = resp.Header.Get("X-Rate-Limit-Ip-State")
	}
	rules, errRule := ExtractFirstRuleFromString(rateLimitRules)
	state, errState := ExtractFirstRuleFromString(rateLimitState)
	// If so, then update our current rate limit counters with the ones
	// the server see from its side (for better accuracy).
	if errRule == nil && errState == nil {
		s.rateLimitManager.UpdateRateLimiter(s.poeSessionID, baseURL, rules, state)
	}
	if s.verbosity > 0 {
		r := s.rateLimitManager.GetRateLimiter(s.poeSessionID, baseURL)
		fmt.Println("Status:", resp.StatusCode, "Rate:", r.NbQuery, "/", r.NbMaxQuery, "ServerRate:", rateLimitState, rateLimitRules)
		if s.verbosity > 1 {
			fmt.Println("Response:", resp)
		}
	}

	defer func() {
		localErr := resp.Body.Close()
		if err == nil {
			err = localErr
		}
	}()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("error while calling PoE API (code %d), using this url: %s", resp.StatusCode, apiURL)
	}

	body, errRead := ioutil.ReadAll(resp.Body)
	if errRead != nil {
		return nil, errRead
	}

	if s.cache {
		if err := ioutil.WriteFile(fileCache, body, 0644); err != nil {
			fmt.Println("can't write to cache", err)
		}
	}

	return body, nil
}

// ScrapEverything scraps items, characters, profile, inventory and so on...
func (s *Scraper) ScrapEverything() (*ScrapedData, error) {
	data := &ScrapedData{
		Account:    s.accountName,
		League:     s.league,
		Realm:      s.realm,
		Date:       time.Now(),
		Characters: make([]*inventory.CharacterInventory, 0, 10),
		Skills:     make([]*inventory.CharacterSkills, 0, 10),
		Stash:      nil,
	}

	// Get the list of all characters of a user.
	characters, errChar := s.ScrapCharacters()
	if errChar != nil {
		return nil, errChar
	}

	// Get inventory of every characters found.
	for _, character := range characters {
		if !character.Expired {
			charInventory, errInventory := s.ScrapCharacterInventory(character.Name)
			if errInventory != nil {
				return nil, errInventory
			}
			data.Characters = append(data.Characters, charInventory)
			charSkills, errSkills := s.ScrapCharacterSkills(character.Name)
			if errSkills != nil {
				return nil, errSkills
			}
			data.Skills = append(data.Skills, charSkills)
		}
	}

	// Retrieves the stash of an account.
	stash, errStash := s.ScrapWholeStash()
	if errStash != nil {
		return nil, errStash
	}
	data.Stash = stash
	data.Wealth = inventory.ComputeWealth(data.Stash, data.Characters)

	return data, nil
}

// GetLeagues retrieves all available league names.
func (s *Scraper) GetLeagues() ([]*inventory.League, error) {
	body, errRequest := s.CallAPI(LeaguesURL)
	if errRequest != nil {
		return nil, errRequest
	}
	leagues, errLeagues := inventory.ParseLeagues(body)
	if errLeagues != nil {
		return nil, errLeagues
	}

	return leagues, nil
}
