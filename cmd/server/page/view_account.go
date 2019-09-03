package page

import (
	"io/ioutil"
	"net/http"

	"github.com/cptpingu/poe-stash/scraper"
	"github.com/gin-gonic/gin"
)

// ViewAccountHandler handles viewing an account
// profile (characters, stash, and so on...).
func ViewAccountHandler(c *gin.Context) {
	account := c.Params.ByName("account")

	content, err := ioutil.ReadFile(scraper.DataDir + account + ".html")
	if err != nil {
		c.HTML(http.StatusNotFound, "error", "Account "+account+" not found! ("+err.Error()+")")
		return
	}

	c.Data(http.StatusOK, "text/html; charset=utf-8", content)
}
