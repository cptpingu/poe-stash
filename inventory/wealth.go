package inventory

import (
	"strconv"
	"strings"
)

var (
	// CurrencyExchangeRate holds all the currency converted in chaos.
	CurrencyExchangeRate = map[string]float64{
		"Ancient Orb":                       27,
		"Ancient Shard":                     2,
		"Annulment Shard":                   3,
		"Apprentice Cartographer's Sextant": 1.1,
		"Armourer's Scrap":                  1 / 44,
		"Binding Shard":                     1 / 2,
		"Blacksmith's Whetstone":            1 / 32.3,
		"Blessed Orb":                       3.9,
		"Blessing of Chayula":               162.4,
		"Blessing of Esh":                   3.3,
		"Blessing of Tul":                   4,
		"Blessing of Uul-Netol":             7.4,
		"Blessing of Xoph":                  3,
		"Cartographer's Chisel":             1 / 2.3,
		"Chaos Orb":                         1,
		"Chaos Shard":                       1 / 20,
		"Chromatic Orb":                     1 / 4.3,
		"Divine Orb":                        40,
		"Engineer's Orb":                    1 / 1.2,
		"Engineer's Shard":                  60,
		"Eternal Orb":                       36500,
		"Exalted Orb":                       163.5,
		"Exalted Shard":                     7,
		"Gemcutter's Prism":                 1 / 1.9,
		"Glassblower's Bauble":              1 / 9.4,
		"Harbinger's Orb":                   28,
		"Harbinger's Shard":                 2,
		"Horizon Shard":                     1 / 20,
		"Jeweller's Orb":                    1 / 6.2,
		"Journeyman Cartographer's Sextant": 2.8,
		"Master Cartographer's Sextant":     3.7,
		"Mirror of Kalandra":                38700,
		"Mirror Shard":                      38700 / 20,
		"Orb of Alchemy":                    1 / 2.5,
		"Orb of Alteration":                 1 / 3.6,
		"Orb of Annulment":                  50,
		"Orb of Augmentation":               1 / 37.1,
		"Orb of Binding":                    1 / 3,
		"Orb of Chance":                     1 / 6.5,
		"Orb of Fusing":                     1 / 1.5,
		"Orb of Horizons":                   1,
		"Orb of Regret":                     1.6,
		"Orb of Scouring":                   1 / 1.2,
		"Orb of Transmutation":              1 / 60,
		"Perandus Coin":                     1 / 70,
		"Portal Scroll":                     1 / 54.9,
		"Regal Orb":                         1 / 1.4,
		"Regal Shard":                       (1 / 1.4) / 20,
		"Scroll of Wisdom":                  1 / 111.5,
		"Silver Coin":                       1 / 4.1,
		"Splinter of Chayula":               1.9,
		"Splinter of Esh":                   1 / 8,
		"Splinter of Tul":                   1 / 7.1,
		"Splinter of Uul-Netol":             1 / 3.1,
		"Splinter of Xoph":                  1 / 4.9,
		"Stacked Deck":                      1 / 7.4,
		"Timeless Eternal Empire Splinter":  1 / 7,
		"Timeless Karui Splinter":           1 / 11,
		"Timeless Maraketh Splinter":        1 / 2,
		"Timeless Templar Splinter":         1 / 2,
		"Timeless Vaal Splinter":            1 / 6,
		"Vaal Orb":                          1.6,
	}
)

// getCount retrieves currency count.
func getCount(properties []ItemProperty) int {
	for _, property := range properties {
		if property.Name == "Stack Size" {
			// In: properties/0/values/0/0.
			for _, row := range property.Values {
				rawValues, ok := row.([]interface{})
				if ok {
					for _, subRow := range rawValues {
						switch rawStack := subRow.(type) {
						case string:
							// Format is "1034/20" for example.
							pos := strings.Index(rawStack, "/")
							rawNb := rawStack[:pos]
							nb, err := strconv.Atoi(rawNb)
							if err != nil {
								// Should never happen!
								panic(err)
							}
							return nb
						}
					}
				}
			}
		}
	}
	return 0
}

// WealthBreakdown holds total wealth and its details.
type WealthBreakdown struct {
	EstimatedChaos int
	NbAlch         int
	NbChaos        int
	NbExa          int
}

// ComputeWealth computes the wealth in chaos orbs contained in a stash.
func ComputeWealth(stashTabs []*StashTab, characters []*CharacterInventory) WealthBreakdown {
	var wealth WealthBreakdown
	var estimate float64

	compute := func(item *Item) {
		nb := getCount(item.Properties)
		switch item.Type {
		case "Orb of Alchemy":
			wealth.NbAlch += nb
		case "Chaos Orb":
			wealth.NbChaos += nb
		case "Exalted Orb":
			wealth.NbExa += nb
		}
		if value, ok := CurrencyExchangeRate[item.Type]; ok {
			estimate += float64(nb) * value
		}
	}

	// Get currencies in stash.
	for _, tab := range stashTabs {
		for _, item := range tab.Items {
			compute(&item)
		}
	}
	// Get currencies in inventories.
	for _, character := range characters {
		for _, item := range character.Items {
			compute(item)
		}
	}

	wealth.EstimatedChaos = int(estimate)
	return wealth
}
