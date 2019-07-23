package inventory

import "encoding/json"

// Color holds primary colors.
type Color struct {
	R int
	G int
	B int
}

// Tab describes a tabulation style (background color,
// custom name, style, and so on)
type Tab struct {
	Hidden          bool   `json:"hidden"`
	Selected        bool   `json:"selected"`
	Index           int    `json:"i"`
	BackgroundColor Color  `json:"colour"`
	Name            string `json:"n"`
	Id              string `json:"id"`
	Type            string `json:"type"`
	ImgL            string `json:"srcL"`
	ImgC            string `json:"srcC"`
	ImgR            string `json:"srcR"`
}

// Layout is used for custom layout like currency.
type Layout struct {
	X float64
	Y float64
	W int
	H int
}

// StashTab holds all stash tabulations (thus all items).
type StashTab struct {
	NumTabs        int               `json:"numTabs"`
	Items          []Item            `json:"items"`
	CurrencyLayout map[string]Layout `json:"currencyLayout"`
}

func (s *StashTab) String() string {
	json, err := json.MarshalIndent(s, "", "\t")
	if err != nil {
		return "<marshalling error>"
	}
	return string(json)
}
