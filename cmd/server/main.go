package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"gitlab.perso/poe-stash/generate"
	"gitlab.perso/poe-stash/scraper"
)

// viewAccountHandler handles viewing an account
// profile (characters, stash, and so on...).
func viewAccountHandler(c *gin.Context) {
	account := c.Params.ByName("account")

	content, err := ioutil.ReadFile(scraper.DataDir + account + ".html")
	if err != nil {
		c.Data(http.StatusNotFound, "text/html; charset=utf-8", []byte("Account "+account+" not found!"))
		return
	}

	c.Data(http.StatusOK, "text/html; charset=utf-8", content)
}

// genAccountHandler handles refresh of an account.
func genAccountHandler(c *gin.Context) {
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

	scrap := scraper.NewScraper(account, poeSessID, realm, league)
	data, errScrap := scrap.ScrapEverything()
	if errScrap != nil {
		fmt.Println("can't scrap data", errScrap)
		os.Exit(2)
	}

	output := scraper.DataDir + account + ".html"
	file, err := os.Create(output)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	w := bufio.NewWriter(file)
	gen := generate.NewGenerator(w)
	if errGen := gen.GenerateHTML(data); errGen != nil {
		fmt.Println("can't generate data", errGen)
		os.Exit(3)
	}
	w.Flush()
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(
		"Generated<br />",
	))
}

// setupRouter setups the http server and all its pages.
func setupRouter() *gin.Engine {
	router := gin.Default()

	authorized := router.Group("/", gin.BasicAuth(gin.Accounts{
		"***":  "***",
		"****": "****",
	}))

	router.Static("/data", scraper.DataDir)
	router.GET("/view/:account", viewAccountHandler)
	authorized.GET("/gen/:account/:poesessid", genAccountHandler)

	return router
}

// main is the main routine which launch the http server.
// This server allows to generate and view account characters,
// stash and items for given users.
func main() {
	r := setupRouter()
	r.Run(":2121")
}
