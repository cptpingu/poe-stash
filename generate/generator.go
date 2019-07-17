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

// GenerateCharacters generates HTML part for characters.
func (g *Generator) GenerateCharacters(characters []*inventory.CharacterInventory) error {
	for _, character := range characters {
		fmt.Fprint(g.writer, character.CharDesc.Name, "\n")
	}
	return nil
}

// GenerateInventory generates HTML part for inventory of characters.
func (g *Generator) GenerateInventory(characters []*inventory.CharacterInventory) error {
	for i, character := range characters {
		for _, item := range character.Items {
			fmt.Fprint(g.writer, characters[i].CharDesc.Name, " got ", item.Name, "\n")
		}
	}
	return nil
}

// GenerateStash generates the whole stash of an account.
func (g *Generator) GenerateStash(stashTabs []*inventory.StashTab) error {
	for i, tab := range stashTabs {
		fmt.Fprint(g.writer, "    item's tab:", i)
		if len(tab.Items) > 0 {
			fmt.Fprint(g.writer, " Object example:", tab.Items[0].Type)
		}
		fmt.Fprint(g.writer, "\n")
	}
	return nil
}
