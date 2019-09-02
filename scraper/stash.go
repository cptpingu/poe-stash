package scraper

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"strconv"

	"github.com/cptpingu/poe-stash/models"
)

// ScrapStash scraps a stash from the official website.
func (s *Scraper) ScrapStash(indexID int) (*models.StashTab, error) {
	var body []byte
	var err error
	if s.demo {
		filename := DemoDir + s.accountName + "/stash_" + strconv.Itoa(indexID) + ".json"
		body, err = ioutil.ReadFile(filename)
		if err != nil {
			return nil, err
		}
	} else {
		url := fmt.Sprintf(StashURL, url.QueryEscape(s.accountName), url.QueryEscape(s.realm), url.QueryEscape(s.league), 1, indexID)
		body, err = s.CallAPI(url)
		if err != nil {
			return nil, err
		}
	}

	stash, errStash := models.ParseStashTab(body)
	if errStash != nil {
		return nil, errStash
	}

	return stash, nil
}

// ScrapWholeStash scraps all tabs in a stash from the official website.
func (s *Scraper) ScrapWholeStash() ([]models.Tab, []*models.StashTab, error) {
	var stashTab []*models.StashTab

	// Scrap first stash to get the number of stash.
	firstStash, err := s.ScrapStash(0)
	if err != nil {
		return nil, nil, err
	}
	stashTab = append(stashTab, firstStash)

	// Scrap the rest.
	for i := 1; i < firstStash.NumTabs; i++ {
		stash, err := s.ScrapStash(i)
		if err != nil {
			return nil, nil, err
		}
		stashTab = append(stashTab, stash)
	}

	return firstStash.Tabs, stashTab, nil
}
