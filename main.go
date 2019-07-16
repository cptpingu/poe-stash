package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	// StashURL is the official URL for the getting a user account stashes.
	StashURL = "https://www.pathofexile.com/character-window/get-stash-items?accountName=%s&realm=%s&league=%s&tabs=%d&tabIndex=%d"
	// ViewProfileURL is the official URL for the getting a user account main profile information.
	ViewProfileURL = "http://www.pathofexile.com/account/view-profile/%s"
	// ProfileCharactersURL is the official URL for the getting a user account characters.
	ProfileCharactersURL = "https://pathofexile.com/character-window/get-characters?accountName=%s"
	// ProfileCharacterItems is the official URL for the getting a user account inventories.
	ProfileCharacterItems = "https://www.pathofexile.com/character-window/get-items"
)

// Profile website account profile
type Profile struct {
	GuildName   string
	GuildURL    string
	GuildID     int
	JoinedAt    time.Time
	ForumPosts  int
	LastVisited time.Time
	Badges      []*Badge
	Characters  []*Character
}

// Badge user profile badge
type Badge struct {
	Name string
	URL  string
}

// Character profile character
type Character struct {
	Name            string
	Level           int
	League          string
	Class           string
	AscendancyClass int `json:"ascendancyClass"`
	ClassID         int `json:"classId"`
	Items           []*Item
}

// CharacterItems items of the profile character
type CharacterItems struct {
	Items []*Item `json:"items"`
}

// Stash holds everything about a stash.
type Stash struct {
	AccountName       string `json:"accountName"`
	LastCharacterName string `json:"lastCharacterName"`
	Id                string `json:"id"`
	Label             string `json:"stash"`
	Type              string `json:"stashType"`
	Items             []Item `json:"items"`
	IsPublic          bool   `json:"public"`
}

// Socket describes a socket.
type Socket struct {
	GroupId   int    `json:"group"`
	Attribute string `json:"attr"`
}

// ItemProperty holds item properties (name, how to display, ...).
type ItemProperty struct {
	Name        string        `json:"name"`
	Values      []interface{} `json:"values"`
	DisplayMode int           `json:"displayMode"`
}

// FrameType is a type of rarity of an item.
type FrameType int

const (
	NormalItemFrameType FrameType = iota
	MagicItemFrameType
	RareItemFrameType
	UniqueItemFrameType
	GemFrameType
	CurrencyFrameType
	DivinationCardFrameType
	QuestItemFrameType
	ProphecyFrameType
	RelicFrameType
)

// Item is a description of all properties of a single item.
type Item struct {
	// Names for some items may include markup. For example: <<set:MS>><<set:M>><<set:S>>Roth's Reach
	Name string `json:"name"`
	Type string `json:"typeLine"`

	Properties   []ItemProperty `json:"properties"`
	Requirements []ItemProperty `json:"requirements"`

	Sockets []Socket `json:"sockets"`

	ExplicitMods []string `json:"explicitMods"`
	ImplicitMods []string `json:"implicitMods"`
	UtilityMods  []string `json:"utilityMods"`
	EnchantMods  []string `json:"enchantMods"`
	CraftedMods  []string `json:"craftedMods"`
	CosmeticMods []string `json:"cosmeticMods"`

	Note string `json:"note"`

	IsVerified             bool      `json:"verified"`
	Width                  int       `json:"w"`
	Height                 int       `json:"h"`
	ItemLevel              int       `json:"ilvl"`
	Icon                   string    `json:"icon"`
	League                 string    `json:"league"`
	Id                     string    `json:"id"`
	IsIdentified           bool      `json:"identified"`
	IsCorrupted            bool      `json:"corrupted"`
	IsLockedToCharacter    bool      `json:"lockedToCharacter"`
	IsSupport              bool      `json:"support"`
	DescriptionText        string    `json:"descrText"`
	SecondDescriptionText  string    `json:"secDescrText"`
	FlavorText             []string  `json:"flavourText"`
	ArtFilename            string    `json:"artFilename"`
	FrameType              FrameType `json:"frameType"`
	StackSize              int       `json:"stackSize"`
	MaxStackSize           int       `json:"maxStackSize"`
	X                      int       `json:"x"`
	Y                      int       `json:"y"`
	InventoryId            string    `json:"inventoryId"`
	SocketedItems          []Item    `json:"socketedItems"`
	IsRelic                bool      `json:"isRelic"`
	TalismanTier           int       `json:"talismanTier"`
	ProphecyText           string    `json:"prophecyText"`
	ProphecyDifficultyText string    `json:"prophecyDiffText"`
}

// Color holds primary colors.
type Color struct {
	R int
	G int
	B int
}

// Tab describes a tabulation style (background color,
// custom name, style, and so on)
type Tab struct {
	Name            string `json:"n"`
	Index           int    `json:"i"`
	Id              string `json:"id"`
	Type            string `json:"type"`
	Hidden          bool   `json:"hidden"`
	Selected        bool   `json:"selected"`
	BackgroundColor Color  `json:"colour"`
	ImgL            string `json:"srcL"`
	ImgC            string `json:"srcC"`
	ImgR            string `json:"srcR"`
}

// StashTab holds all stash tabulations (thus all items).
type StashTab struct {
	NumTabs int
	Tabs    []Tab
	Items   []Item
	// CurrencyLayout FIXME
}

func (s *StashTab) String() string {
	json, err := json.MarshalIndent(s, "", "\t")
	if err != nil {
		return "<marshalling error>"
	}
	return string(json)
}

// ParseStashTab parses a Path of Exile stash tabulation.
func ParseStashTab(data []byte) (*StashTab, error) {
	stash := StashTab{}
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

// Scraper scraps path of exile site using its API.
type Scraper struct {
	accountName  string
	poeSessionID string
	realm        string
	league       string
}

// NewScraper returns a configured scraper.
func NewScraper(accountName, poeSessionID, realm, league string) *Scraper {
	return &Scraper{
		accountName:  accountName,
		poeSessionID: poeSessionID,
		realm:        realm,
		league:       league,
	}
}

// ScrapStash scraps a stash from the official website.
func (s *Scraper) ScrapStash(charID, indexID int) (stash *StashTab, err error) {
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

	stash, errStash := ParseStashTab(body)
	if errStash != nil {
		return nil, errStash
	}

	return stash, nil
}

func main() {
	fmt.Println("start")
	var stashes []*StashTab
	scraper := NewScraper("cptpingu", "", "pc", "Standard")
	maxIndexID := 0
	for i := 0; i <= maxIndexID; i++ {
		if stash, err := scraper.ScrapStash(0, i); err != nil {
			fmt.Println("Failed to scrap stash:", err, "index:", i)
		} else {
			stashes = append(stashes, stash)
		}
	}
	for _, stash := range stashes {
		fmt.Println(stash.String())
	}
}
