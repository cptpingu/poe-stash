package scraper

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
