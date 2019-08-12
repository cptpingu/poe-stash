package main

import (
	"bufio"
	"flag"
	"fmt"
	"html/template"
	"io"
	"os"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/poe-stash/cmd/server/page"
	"github.com/poe-stash/generate"
	"github.com/poe-stash/scraper"
)

// setupRouter setups the http server and all its pages.
func setupRouter(passwords map[string]string) *gin.Engine {
	router := gin.Default()

	t := template.Must(generate.LoadAllTemplates())
	router.SetHTMLTemplate(t)
	router.NoRoute(page.CustomErrorHandler)

	router.Static("/data", scraper.DataDir)
	router.GET("/", page.MainPageHandler)
	router.GET("/view/:account", page.ViewAccountHandler)

	authorized := router.Group("/", gin.BasicAuth(passwords))
	authorized.GET("/gen/:account", page.GenAccountHandler)

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

// main is the main routine which launch the http server.
// This server allows to generate and view account characters,
// stash and items for given users.
func main() {
	port := flag.Int("port", 2121, "port")
	passwordFile := flag.String("passwords", "./pass.txt", "password file (containing login:pass in plain text)")
	flag.Parse()

	passwords, err := loadPasswords(*passwordFile)
	if err != nil {
		panic(err)
	}
	r := setupRouter(passwords)
	r.Run(fmt.Sprintf(":%d", *port))
}
