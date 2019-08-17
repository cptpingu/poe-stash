package page

import (
	"errors"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/poe-stash/inventory"
	"github.com/poe-stash/scraper"
)

// listAllAccounts list all fetch accounts.
func listAllAccounts(rootDir, ext string) ([]string, error) {
	var accounts []string
	cleanRoot := filepath.Clean(rootDir)
	err := filepath.Walk(cleanRoot, func(curPath string, info os.FileInfo, e1 error) error {
		if info == nil {
			return errors.New("can't list files")
		}
		if !info.IsDir() && strings.HasSuffix(curPath, ext) {
			if e1 != nil {
				return e1
			}
			filename := path.Base(curPath)
			accounts = append(accounts, filename[:len(filename)-len(ext)])
		}
		return nil
	})

	return accounts, err
}

// MainPageHandler serves the main page of this website.
func MainPageHandler(c *gin.Context) {
	accounts, err := listAllAccounts("data", ".html")
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error", err)
		return
	}

	scraper := scraper.NewScraper("", "", "", "", false)
	leagues, errLeagues := scraper.GetLeagues()
	if errLeagues != nil {
		c.HTML(http.StatusInternalServerError, "error", errLeagues)
	}

	c.HTML(http.StatusOK, "main", struct {
		Accounts []string
		Leagues  []*inventory.League
	}{
		Accounts: accounts,
		Leagues:  leagues,
	})
}
