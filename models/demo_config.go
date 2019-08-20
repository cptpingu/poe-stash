package models

import "encoding/json"

// DemoConfig hold pathname of all mocked info.
type DemoConfig struct {
	// Dummy info.
	Account string
	League  string
	Realm   string

	// Pathfiles.
	Characters string
	Skills     []string
	Stash      []string
}

// ParseDemoConfig parses a demo config file.
func ParseDemoConfig(data []byte) (DemoConfig, error) {
	conf := DemoConfig{}
	if err := json.Unmarshal(data, &conf); err != nil {
		return conf, err
	}
	return conf, nil
}
