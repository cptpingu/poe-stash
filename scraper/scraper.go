package scraper

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"gitlab.perso/poe-stash/inventory"
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
)

// Scraper scraps path of exile site using its API.
type Scraper struct {
	accountName  string
	poeSessionID string
	realm        string
	league       string

	client http.Client
}

// ScrapedData holds everything scrapped.
type ScrapedData struct {
	Characters []*inventory.CharacterInventory
	Stash      []*inventory.StashTab
}

// NewScraper returns a configured scraper.
func NewScraper(accountName, poeSessionID, realm, league string) *Scraper {
	return &Scraper{
		accountName:  accountName,
		poeSessionID: poeSessionID,
		realm:        realm,
		league:       league,
		client:       http.Client{},
	}
}

// CallAPI calls a distant API and returns the content.
func (s *Scraper) CallAPI(url string) ([]byte, error) {
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
	return body, nil
}

// ScrapEverything scraps items, characters, profile, inventory and so on...
func (s *Scraper) ScrapEverything() (*ScrapedData, error) {
	data := &ScrapedData{
		Characters: make([]*inventory.CharacterInventory, 0, 10),
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
			inventory, errInventory := s.ScrapCharacterInventory(character.Name)
			if errInventory != nil {
				return nil, errInventory
			}
			data.Characters = append(data.Characters, inventory)
		}
	}

	// Retrieves the stash of an account.
	stash, errStash := s.ScrapWholeStash(11)
	if errStash != nil {
		return nil, errStash
	}
	data.Stash = stash

	return data, nil
}
