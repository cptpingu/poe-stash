package main

import (
	"flag"
	"fmt"
	"html/template"

	"github.com/gin-gonic/gin"

	"gitlab.perso/poe-stash/cmd/server/page"
	"gitlab.perso/poe-stash/generate"
	"gitlab.perso/poe-stash/scraper"
)

// setupRouter setups the http server and all its pages.
func setupRouter() *gin.Engine {
	router := gin.Default()

	t := template.Must(generate.LoadAllTemplates())
	router.SetHTMLTemplate(t)
	// router.Use(errorHandler)
	router.NoRoute(page.CustomErrorHandler)

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
	port := flag.Int("port", 2121, "port")
	flag.Parse()
	r := setupRouter()
	r.Run(fmt.Sprintf(":%d", *port))
}
