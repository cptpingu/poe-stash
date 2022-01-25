package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/cptpingu/poe-stash/generate"
	"github.com/cptpingu/poe-stash/misc"
	"github.com/cptpingu/poe-stash/scraper"
	"github.com/pkg/errors"
)

// mandatoryOption ensure an option is not empty.
func mandatoryOption(opt string, name string) bool {
	if opt == "" {
		fmt.Printf("option %#v is mandatory!\n", name)
		return false
	}
	return true
}

// scrapData scraps all data.
func scrapData(account, poeSessID, realm, league string, demo, cache bool, verbosity int) (*scraper.ScrapedData, error) {
	scraper := scraper.NewScraper(account, poeSessID, realm, league)
	scraper.SetDemo(demo)
	if cache {
		scraper.EnableCache()
	}
	scraper.SetVerbosity(verbosity)
	data, err := scraper.ScrapEverything()
	if err != nil {
		return nil, errors.Wrap(err, "ScrapEverything")
	}

	return data, nil
}

// generateData generates html file from the given data.
func generateData(data *scraper.ScrapedData, output string) (resErr error) {
	var file *os.File
	var err error

	if output == "-" {
		file = os.Stdout
	} else {
		file, err = os.Create(output)
		if err != nil {
			return errors.Wrap(err, "os.Create")
		}
		defer func() {
			if err := file.Close(); err != nil {
				if resErr != nil {
					resErr = errors.Wrap(err, "file.Close")
				}
			}
		}()
	}

	w := bufio.NewWriter(file)
	gen := generate.NewGenerator(w)
	if errGen := gen.GenerateHTML(data); errGen != nil {
		return errors.Wrap(err, "GenerateHTML")
	}
	if errFlush := w.Flush(); err != nil {
		return errors.Wrap(errFlush, "Flush")
	}

	return nil
}

// main is the main routine for this CLI.
// This CLI allows to generate an html file which contains all
// inventories, characters and items for a given account.
func main() {
	account := flag.String("account", "", "account name")
	poeSessID := flag.String("poesessid", "", "poesessid got after login on the official website")
	realm := flag.String("realm", "pc", "the realm (pc, sony, xbox)")
	league := flag.String("league", "standard", "league name (anarchy, legion, synthesis, delve...)")
	output := flag.String("output", "", "where to generate html file (put \"-\" for stdin), if empty, a generated name will be created (account-league.html)")
	cache := flag.Bool("cache", false, "do not call distant api, and use local cache if possible, for debug purpose only")
	demo := flag.Bool("demo", false, "use local files to generate example profiles")
	verbosity := flag.Int("verbosity", 0, "set the log verbose level")
	interactive := flag.Bool("interactive", false, "interactive mode")
	version := flag.Bool("version", false, "display the version of this tool")
	flag.Parse()

	if *version {
		fmt.Println(misc.Version)
		return
	}

	if !*interactive {
		valid := true
		valid = mandatoryOption(*account, "account") && valid
		if !*demo {
			valid = mandatoryOption(*poeSessID, "poesessid") && valid
			valid = mandatoryOption(*realm, "realm") && valid
			valid = mandatoryOption(*league, "league") && valid
		}
		if !valid {
			fmt.Println()
			flag.Usage()
			os.Exit(1)
		}
	} else {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Print("Account: ")
		scanner.Scan()
		*account = scanner.Text()
		fmt.Print("PoE Sessid: ")
		scanner.Scan()
		*poeSessID = scanner.Text()
		fmt.Print("League (empty = standard): ")
		scanner.Scan()
		*league = scanner.Text()
		if *league == "" {
			*league = "standard"
		}
		fmt.Print("Realm (empty = pc): ")
		scanner.Scan()
		*realm = scanner.Text()
		if *realm == "" {
			*realm = "pc"
		}
	}

	data, errScrap := scrapData(*account, *poeSessID, *realm, *league, *demo, *cache, *verbosity)
	if errScrap != nil {
		fmt.Println("can't scrap data", errors.WithStack(errScrap))
		os.Exit(2)
	}

	if *output == "" {
		*output = *account + "-" + *league + ".html"
	}

	if errGen := generateData(data, *output); errGen != nil {
		fmt.Println("can't generate data", errGen)
		os.Exit(3)
	}

	fmt.Println("File successfully generated:", *output)
}
