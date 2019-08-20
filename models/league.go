package models

import (
	"encoding/json"
	"time"
)

// League represents a league and its characteristics.
type League struct {
	Id         string    `json:"id"`
	Realm      string    `json:"realm"`
	Url        string    `json:"url"`
	StartAt    time.Time `json:"startAt"`
	EndAt      time.Time `json:"endAt"`
	DelveEvent bool      `json:"delveEven"`
}

// ParseLeagues parses a Path of Exile leagues list.
func ParseLeagues(data []byte) ([]*League, error) {
	leagues := []*League{}
	if err := json.Unmarshal(data, &leagues); err != nil {
		return nil, err
	}
	return leagues, nil
}
