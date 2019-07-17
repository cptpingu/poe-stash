package main

import (
	"bufio"
	"fmt"
	"os"

	"gitlab.perso/poe-stash/generate"
	"gitlab.perso/poe-stash/scraper"
)

func main() {
	fmt.Println("start")
	scraper := scraper.NewScraper("cptpingu", "", "pc", "Standard")
	data, errScrap := scraper.ScrapEverything()
	if errScrap != nil {
		fmt.Println("can't scrap data", errScrap)
	}
	w := bufio.NewWriter(os.Stdout)
	gen := generate.NewGenerator(w)
	if errGen := gen.GenerateHTML(data); errGen != nil {
		fmt.Println("can't generate data", errGen)
	}
	w.Flush()
}
