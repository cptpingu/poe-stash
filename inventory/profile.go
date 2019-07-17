package inventory

import "time"

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
