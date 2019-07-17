package main

import (
	"fmt"

	"gitlab.perso/poe-stash/scraper"
)

func main() {
	fmt.Println("start")
	scraper := scraper.NewScraper("cptpingu", "", "pc", "Standard")
	stashes, err := scraper.ScrapAllStashes(0)
	if err != nil {
		fmt.Println("can't scrap stashes", err)
	}
	for i, stash := range stashes {
		fmt.Println("======", i, "======\n", stash.String(), "\n")
	}
}
