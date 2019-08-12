package scraper

import (
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"net/http"
	"strconv"

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

	// DataDir is where all data are.
	DataDir = "data/"
	// DataCacheDir is the cache directory.
	DataCacheDir = DataDir + "cache/"
)

// Scraper scraps path of exile site using its API.
type Scraper struct {
	cache    bool
	cacheDir string

	accountName  string
	poeSessionID string
	realm        string
	league       string

	client http.Client
}

// ScrapedData holds everything scrapped.
type ScrapedData struct {
	Version string

	Characters []*inventory.CharacterInventory
	Skills     []*inventory.CharacterSkills
	Stash      []*inventory.StashTab
	Wealth     inventory.WealthBreakdown
}

// NewScraper returns a configured scraper.
func NewScraper(accountName, poeSessionID, realm, league string, cache bool) *Scraper {
	return &Scraper{
		cache:        cache,
		cacheDir:     DataCacheDir,
		accountName:  accountName,
		poeSessionID: poeSessionID,
		realm:        realm,
		league:       league,
		client:       http.Client{},
	}
}

// hash url into a number.
func hash(s string) string {
	h := fnv.New32a()
	h.Write([]byte(s))
	return strconv.Itoa(int(h.Sum32()))
}

// CallAPI calls a distant API and returns the content.
func (s *Scraper) CallAPI(url string) ([]byte, error) {
	var fileCache string
	if s.cache {
		fileCache = s.cacheDir + hash(url)
		if b, err := ioutil.ReadFile(fileCache); err != nil {
			fmt.Println("can't read cache", err)
		} else {
			return b, nil
		}
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	cookie := http.Cookie{
		Name:  "POESESSID",
		Value: s.poeSessionID,
	}
	req.AddCookie(&cookie)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, errResponse := s.client.Do(req)
	if errResponse != nil {
		return nil, errResponse
	}
	defer func() {
		localErr := resp.Body.Close()
		if err == nil {
			err = localErr
		}
	}()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("error while calling PoE API (code %d), using this url: %s", resp.StatusCode, url)
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
		Version:    "v0.1",
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
