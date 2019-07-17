package scraper

import (
	"encoding/json"
	"fmt"
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
func (s *Scraper) ScrapStash(charID, indexID int) (*inventory.StashTab, error) {
	url := fmt.Sprintf(StashURL, s.accountName, s.realm, s.league, charID, indexID)
	body, errRequest := s.CallAPI(url)
	if errRequest != nil {
		return nil, errRequest
	}
	stash, errStash := parseStashTab(body)
	if errStash != nil {
		return nil, errStash
	}

	return stash, nil
}

// ScrapWholeStash scraps all tabs in a stash from the official website.
func (s *Scraper) ScrapWholeStash(charID int) ([]*inventory.StashTab, error) {
	var stashTab []*inventory.StashTab
	maxIndexID := 11
	for i := 0; i <= maxIndexID; i++ {
		stash, err := s.ScrapStash(0, i)
		if err != nil {
			return nil, err
		}
		stashTab = append(stashTab, stash)
	}
	return stashTab, nil
}
