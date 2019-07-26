package inventory

import (
	"encoding/json"
)

// Character profile character
type Character struct {
	Expired         bool   `json:"expired"`
	LastActive      bool   `json:"lastActive"`
	Level           int    `json:"level"`
	AscendancyClass int    `json:"ascendancyClass"`
	ClassID         int    `json:"classId"`
	Experience      int64  `json:"experience"`
	Name            string `json:"name"`
	League          string `json:"league"`
	Class           string `json:"class"`
}

// DefaultInventoryLayout is the default mapping for placing
// inventory items correctly.
var DefaultInventoryLayout = map[string]float64{
	"MainInventoryX": 14,
	"MainInventoryY": 531.568,
	"FlaskX0":        186.339,
	"FlaskX1":        234.803,
	"FlaskX2":        281.268,
	"FlaskX3":        328.732,
	"FlaskX4":        376.197,
	"FlaskY":         418.511,
	"BodyArmourX":    252.059,
	"BodyArmourY":    206.138,
	"RingX":          182.688,
	"RingY":          253.602,
	"Ring2X":         368.895,
	"Ring2Y":         253.602,
	"BootsX":         368.895,
	"BootsY":         312.629,
	"GlovesX":        135.223,
	"GlovesY":        312.629,
	"HelmX":          252.059,
	"HelmY":          99.6471,
	"AmuletX":        368.895,
	"AmuletY":        194.576,
	"WeaponX":        65.8519,
	"WeaponY":        111.209,
	"Weapon2X":       65.8519,
	"Weapon2Y":       111.209,
	"BeltX":          252.059,
	"BeltY":          360.093,
	"OffhandX":       438.266,
	"OffhandY":       111.209,
	"Offhand2X":      438.266,
	"Offhand2Y":      111.209,
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
	Items    []*Item   `json:"items"`
}

func (c *CharacterInventory) String() string {
	json, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		return "<marshalling error>"
	}
	return string(json)
}
