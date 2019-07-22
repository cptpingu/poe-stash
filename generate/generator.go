package generate

import (
	"fmt"
	"io"

	"gitlab.perso/poe-stash/inventory"
	"gitlab.perso/poe-stash/scraper"
)

// Generator construct html files from a scraped user.
type Generator struct {
	writer io.Writer
}

// NewGenerator constructs a new generator.
func NewGenerator(writer io.Writer) Generator {
	return Generator{
		writer: writer,
	}
}

// GenerateHTML generates HTML from scraped data.
func (g *Generator) GenerateHTML(data *scraper.ScrapedData) error {
	if err := g.GenerateWealth(data.Wealth); err != nil {
		return err
	}
	if err := g.GenerateCharacters(data.Characters); err != nil {
		return err
	}
	if err := g.GenerateInventory(data.Characters); err != nil {
		return err
	}
	if err := g.GenerateStash(data.Stash); err != nil {
		return err
	}
	return nil
}

// GenerateWealth generates HTML part for wealth account in chaos orbs.
func (g *Generator) GenerateWealth(wealth int) error {
	_, err := fmt.Fprint(g.writer, "Wealth: ", wealth, "\n")
	return err
}

// GenerateCharacters generates HTML part for characters.
func (g *Generator) GenerateCharacters(characters []*inventory.CharacterInventory) error {
	for _, character := range characters {
		if _, err := fmt.Fprint(g.writer, character.CharDesc.Name, "\n"); err != nil {
			return err
		}
	}
	return nil
}

// GenerateInventory generates HTML part for inventory of characters.
func (g *Generator) GenerateInventory(characters []*inventory.CharacterInventory) error {
	for i, character := range characters {
		for _, item := range character.Items {
			if _, err := fmt.Fprint(g.writer, characters[i].CharDesc.Name, " got ", item.Name, "\n"); err != nil {
				return err
			}
		}
	}
	return nil
}

// GenerateStash generates the whole stash of an account.
func (g *Generator) GenerateStash(stashTabs []*inventory.StashTab) error {
	for i, tab := range stashTabs {
		if _, err := fmt.Fprint(g.writer, "    item's tab:", i); err != nil {
			return err
		}
		if len(tab.Items) > 0 {
			if _, err := fmt.Fprint(g.writer, " Object example:", tab.Items[0].Type); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprint(g.writer, "\n"); err != nil {
			return err
		}
	}
	return nil
}
