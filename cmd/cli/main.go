package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"gitlab.perso/poe-stash/generate"
	"gitlab.perso/poe-stash/scraper"
)

// mandatoryOption ensure an option is not empty.
func mandatoryOption(opt string, name string) {
	if opt == "" {
		fmt.Printf("option %#v is mandatory!\n", name)
		os.Exit(1)
	}
}

// main is the main routine for this CLI.
// This CLI allows to generate an html file which contains all
// account, characters and items for a given account.
func main() {
	account := flag.String("account", "", "account name")
	poeSessID := flag.String("poesessid", "", "poesessid got after login on the official website")
	realm := flag.String("realm", "pc", "the realm (pc, ps4, xbox)")
	league := flag.String("league", "Standard", "league name (anarchy, legion, synthesis, delve...)")
	output := flag.String("output", "-", "where to genreate html file (put \"-\" for stdin")
	cache := flag.Bool("cache", false, "do not call distant api, and use local cache if possible")
	flag.Parse()
	mandatoryOption(*account, "account")
	mandatoryOption(*poeSessID, "poesessid")
	mandatoryOption(*realm, "realm")
	mandatoryOption(*league, "league")
	mandatoryOption(*output, "output")

	scraper := scraper.NewScraper(*account, *poeSessID, *realm, *league, *cache)
	data, errScrap := scraper.ScrapEverything()
	if errScrap != nil {
		fmt.Println("can't scrap data", errScrap)
		os.Exit(2)
	}

	var file *os.File
	var err error
	if *output == "-" {
		file = os.Stdout
	} else {
		file, err = os.Create(*output)
		if err != nil {
			panic(err)
		}
		defer func() {
			if err := file.Close(); err != nil {
				panic(err)
			}
		}()
	}

	w := bufio.NewWriter(file)
	gen := generate.NewGenerator(w)
	if errGen := gen.GenerateHTML(data); errGen != nil {
		fmt.Println("can't generate data", errGen)
		os.Exit(3)
	}
	w.Flush()
}
