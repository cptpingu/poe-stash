package page

import (
	"github.com/gin-gonic/gin"

	"github.com/poe-stash/scraper"
)

// DownloadFileHandler handles force download of a file.
func DownloadFileHandler(c *gin.Context) {
	account := c.Params.ByName("account")
	filename := account + ".html"
	filepath := scraper.DataDir + filename
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/octet-stream")
	c.File(filepath)
}
