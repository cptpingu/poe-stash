package main

import (
	"bufio"
	"flag"
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/cptpingu/poe-stash/cmd/server/page"
	"github.com/cptpingu/poe-stash/generate"
	"github.com/cptpingu/poe-stash/misc"
	"github.com/cptpingu/poe-stash/scraper"
)

// EnvMiddleware will add env to query.
func EnvMiddleware(verbosity int) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("verbosity", verbosity)
		c.Next()
	}
}

// setupRouter setups the http server and all its pages.
func setupRouter(passwords map[string]string, verbosity int) *gin.Engine {
	router := gin.Default()
	router.Use(EnvMiddleware(verbosity))

	t := template.Must(generate.LoadAllTemplates())
	router.SetHTMLTemplate(t)
	router.NoRoute(page.CustomErrorHandler)

	router.Static("/data", scraper.DataDir)
	router.GET("/", page.MainPageHandler)
	router.GET("/view/:account", page.ViewAccountHandler)
	router.GET("/download/:account", page.DownloadFileHandler)

	if passwords != nil {
		authorized := router.Group("/", gin.BasicAuth(passwords))
		authorized.GET("/gen/:account", page.GenAccountHandler)
	} else {
		router.GET("/gen/:account", page.GenAccountHandler)
	}

	return router
}

// loadPasswords load passwords from a given file.
// Format is:
//   login:pass
//   login:pass
//   login:pass
//   ...
func loadPasswords(filename string) (r map[string]string, mainErr error) {
	res := make(map[string]string, 2)
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := file.Close()
		if err == nil {
			mainErr = err
		}
	}()

	reader := bufio.NewReader(file)
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		parts := strings.Split(string(line), ":")
		// Invalid line is skipped
		if len(parts) != 2 {
			fmt.Println("Skipped invalid line:", string(line))
		} else {
			res[parts[0]] = parts[1]
		}
	}

	return res, nil
}

// dirExists checks if directory exists.
func dirExists(dir string) bool {
	stat, err := os.Stat(dir)
	return err == nil && stat.IsDir()
}

// main is the main routine which launch the http server.
// This server allows to generate and view account characters,
// stash and items for given users.
func main() {
	port := flag.Int("port", 2121, "port")
	passwordFile := flag.String("passwords", "", "password file (containing login:pass in plain text)")
	version := flag.Bool("version", false, "display the version of this tool")
	verbosity := flag.Int("verbosity", 0, "set the log verbose level")
	flag.Parse()

	if *version {
		fmt.Println(misc.Version)
		return
	}

	// Check data dir exists.
	if !dirExists(scraper.DataDir) {
		// Try to cd on the directory where the binary is launched.
		currentDir := filepath.Dir(os.Args[0])
		fmt.Println("No data directory, trying to fallback on", currentDir)
		err := os.Chdir(currentDir)
		if err != nil {
			panic(err)
		}
		if !dirExists(scraper.DataDir) {
			panic("can't find any data directory!")
		}
	}

	var passwords map[string]string
	if *passwordFile != "" {
		var err error
		_, err = loadPasswords(*passwordFile)
		if err != nil {
			panic(err)
		}
	}

	r := setupRouter(passwords, *verbosity)
	r.Run(fmt.Sprintf(":%d", *port))
}
