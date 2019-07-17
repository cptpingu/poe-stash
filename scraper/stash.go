package scraper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"gitlab.perso/poe-stash/inventory"
)

// parseStashTab parses a Path of Exile stash tabulation.
func parseStashTab(data []byte) (*inventory.StashTab, error) {
	stash := inventory.StashTab{}
	if err := json.Unmarshal(data, &stash); err != nil {
		return nil, err
	}

	// Clean useless markers in the json.
	for _, item := range stash.Items {
		item.Name = strings.TrimPrefix(item.Name, "<<set:MS>><<set:M>><<set:S>>")
		item.Type = strings.TrimPrefix(item.Type, "<<set:MS>><<set:M>><<set:S>>")
	}

	return &stash, nil
}

// ScrapStash scraps a stash from the official website.
func (s *Scraper) ScrapStash(charID, indexID int) (stash *inventory.StashTab, err error) {
	client := &http.Client{}
	url := fmt.Sprintf(StashURL, s.accountName, s.realm, s.league, charID, indexID)
	req, errRequest := http.NewRequest("GET", url, nil)
	if errRequest != nil {
		return nil, errRequest
	}
	cookie := http.Cookie{
		Name:  "POESESSID",
		Value: s.poeSessionID,
	}
	req.AddCookie(&cookie)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, errResponse := client.Do(req)
	if errResponse != nil {
		return nil, errResponse
	}
	defer func() {
		if err == nil {
			err = resp.Body.Close()
		}
	}()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("error while calling PoE API (code %d), using this url: %s", resp.StatusCode, url)
	}

	body, errRead := ioutil.ReadAll(resp.Body)
	if errRead != nil {
		return nil, errRead
	}

	stash, errStash := parseStashTab(body)
	if errStash != nil {
		return nil, errStash
	}

	return stash, nil
}

// ScrapAllStashes scraps all stashes from the official website.
func (s *Scraper) ScrapAllStashes(charID int) (stash []*inventory.StashTab, err error) {
	var stashes []*inventory.StashTab
	maxIndexID := 11
	for i := 0; i <= maxIndexID; i++ {
		stash, err := s.ScrapStash(0, i)
		if err != nil {
			return nil, err
		}
		stashes = append(stashes, stash)
	}
	return stashes, nil
}
