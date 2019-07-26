package main

import (
	"html/template"

	"github.com/gin-gonic/gin"

	"gitlab.perso/poe-stash/cmd/server/page"
	"gitlab.perso/poe-stash/scraper"
)

// setupRouter setups the http server and all its pages.
func setupRouter() *gin.Engine {
	router := gin.Default()

	t, err := template.ParseFiles(
		"data/template/main.tmpl",
		"data/template/redirect.tmpl",
		"data/template/error.tmpl",
	)
	if err != nil {
		panic(err)
	}
	router.SetHTMLTemplate(t)

	router.Static("/data", scraper.DataDir)
	router.GET("/", page.MainPageHandler)
	router.GET("/view/:account", page.ViewAccountHandler)

	authorized := router.Group("/", gin.BasicAuth(gin.Accounts{
		"***":  "***",
		"****": "****",
	}))
	authorized.GET("/gen/:account/:poesessid", page.GenAccountHandler)

	return router
}

// main is the main routine which launch the http server.
// This server allows to generate and view account characters,
// stash and items for given users.
func main() {
	r := setupRouter()
	r.Run(":2121")
}
