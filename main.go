package main

import (
	"fmt"

	"gitlab.perso/poe-stash/scraper"
)

func main() {
	fmt.Println("start")
	scraper := scraper.NewScraper("cptpingu", "", "pc", "Standard")
	characters, errChar := scraper.ScrapCharacters()
	if errChar != nil {
		fmt.Println("can't scrap characters", errChar)
	}
	for i, character := range characters {
		fmt.Println("======", i, "======\n", character.String(), "\n")
		inventory, errInventory := scraper.ScrapCharacterInventory(character.Name)
		if errInventory != nil {
			fmt.Println("can't scrap characters", errInventory)
		}
		fmt.Println("Inventory for ", character.Name, ":\n", inventory.String(), "\n")
	}

	stashes, err := scraper.ScrapAllStashes(0)
	if err != nil {
		fmt.Println("can't scrap stashes", err)
	}
	for i, stash := range stashes {
		fmt.Println("======", i, "======\n", stash.String(), "\n")
	}
}
