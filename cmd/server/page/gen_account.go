package page

import (
	"bufio"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"gitlab.perso/poe-stash/generate"
	"gitlab.perso/poe-stash/scraper"
)

// GenAccountHandler handles refresh of an account.
func GenAccountHandler(c *gin.Context) {
	user := c.MustGet(gin.AuthUserKey).(string)
	_ = user

	account := c.Params.ByName("account")
	poeSessID := c.Params.ByName("poesessid")
	realm := "pc"
	league := "Standard"

	query := c.Request.URL.Query()
	if param, ok := query["realm"]; ok {
		realm = param[0]
	}
	if param, ok := query["league"]; ok {
		league = param[0]
	}

	scrap := scraper.NewScraper(account, poeSessID, realm, league, false)
	data, errScrap := scrap.ScrapEverything()
	if errScrap != nil {
		c.HTML(http.StatusOK, "error", errScrap)
		return
	}

	output := scraper.DataDir + account + ".html"
	file, err := os.Create(output)
	if err != nil {
		c.HTML(http.StatusOK, "error", err)
		return
	}
	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	w := bufio.NewWriter(file)
	gen := generate.NewGenerator(w)
	if errGen := gen.GenerateHTML(data); errGen != nil {
		c.HTML(http.StatusOK, "error", errGen)
		return
	}
	w.Flush()
	c.HTML(http.StatusOK, "redirect", "/view/"+account)
}
