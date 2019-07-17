package inventory

import (
	"encoding/json"
)

// Character profile character
type Character struct {
	Expired         bool
	LastActive      bool
	Level           int
	AscendancyClass int `json:"ascendancyClass"`
	ClassID         int `json:"classId"`
	Experience      int64
	Name            string
	League          string
	Class           string
}

func (c *Character) String() string {
	json, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		return "<marshalling error>"
	}
	return string(json)
}

// CharacterInventory holds inventory of a character.
type CharacterInventory struct {
	CharDesc Character `json:"character"`
	Items    []*Item
}

func (c *CharacterInventory) String() string {
	json, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		return "<marshalling error>"
	}
	return string(json)
}
