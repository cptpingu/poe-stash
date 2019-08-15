package page

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/poe-stash/scraper"
)

// DownloadFileHandler handles force download of a file.
func DownloadFileHandler(c *gin.Context) {
	account := c.Params.ByName("account")
	realm := "pc"
	league := "standard"
	date := "today"
	query := c.Request.URL.Query()
	if param, ok := query["realm"]; ok {
		realm = param[0]
	}
	if param, ok := query["league"]; ok {
		league = param[0]
	}
	if param, ok := query["date"]; ok {
		date = param[0]
	}

	// File name serve to the user (cosmetic).
	filename := fmt.Sprintf("%s-%s-%s-%s.html", account, league, realm, date)
	// Name the file really has.
	realFilename := fmt.Sprintf("%s-%s-%s.html", account, league, realm)
	filepath := scraper.DataDir + realFilename
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/octet-stream")
	c.File(filepath)
}
