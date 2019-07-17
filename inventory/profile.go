package inventory

import (
	"time"
)

// Profile website account profile
type Profile struct {
	GuildName   string
	GuildURL    string
	GuildID     int
	ForumPosts  int
	JoinedAt    time.Time
	LastVisited time.Time
	Badges      []*Badge
	Characters  []*Character
}

// Badge user profile badge
type Badge struct {
	Name string
	URL  string
}
