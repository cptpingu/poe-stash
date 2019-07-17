package main

import (
	"fmt"

	"gitlab.perso/poe-stash/scraper"
)

func main() {
	fmt.Println("start")
	scraper := scraper.NewScraper("cptpingu", "", "pc", "Standard")
	data, err := scraper.ScrapEverything()
	if err != nil {
		fmt.Println("can't scrap data", err)
	}
	fmt.Println(data)
}
