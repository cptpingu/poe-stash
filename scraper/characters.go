package scraper

import (
	"fmt"
	"io/ioutil"
	"net/url"

	"github.com/cptpingu/poe-stash/models"
)

// ScrapCharacters scraps all characters owned by a user.
func (s *Scraper) ScrapCharacters() ([]*models.Character, error) {
	var body []byte
	var err error
	if s.demo {
		filename := DemoDir + s.accountName + "/characters.json"
		body, err = ioutil.ReadFile(filename)
		if err != nil {
			return nil, err
		}
	} else {
		url := fmt.Sprintf(ProfileCharactersURL, url.QueryEscape(s.accountName))
		body, err = s.CallAPI(url)
		if err != nil {
			return nil, err
		}
	}

	characters, errCharacters := models.ParseCharacters(body)
	if errCharacters != nil {
		return nil, errCharacters
	}

	return characters, nil
}

// ScrapCharacterInventory scraps the inventory of a given character.
func (s *Scraper) ScrapCharacterInventory(charName string) (*models.CharacterInventory, error) {
	var body []byte
	var err error
	if s.demo {
		filename := DemoDir + s.accountName + "/" + charName + "_inventory.json"
		body, err = ioutil.ReadFile(filename)
		if err != nil {
			return nil, err
		}
	} else {
		url := fmt.Sprintf(ProfileCharacterItemsURL, url.QueryEscape(s.accountName), url.QueryEscape(s.realm), url.QueryEscape(charName))
		body, err = s.CallAPI(url)
		if err != nil {
			return nil, err
		}
	}

	inventory, errInventory := models.ParseInventory(body)
	if errInventory != nil {
		return nil, errInventory
	}

	return inventory, nil
}

// ScrapCharacterSkills scraps the inventory of a given character.
func (s *Scraper) ScrapCharacterSkills(charName string) (*models.CharacterSkills, error) {
	var body []byte
	var err error
	if s.demo {
		filename := DemoDir + s.accountName + "/" + charName + "_skills.json"
		body, err = ioutil.ReadFile(filename)
		if err != nil {
			return nil, err
		}
	} else {
		url := fmt.Sprintf(ProfileCharacterSkillsURL, url.QueryEscape(charName), url.QueryEscape(s.accountName))
		body, err = s.CallAPI(url)
		if err != nil {
			return nil, err
		}
	}

	inventory, errInventory := models.ParseSkills(body)
	if errInventory != nil {
		return nil, errInventory
	}

	return inventory, nil
}
