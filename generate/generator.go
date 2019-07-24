package generate

import (
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gitlab.perso/poe-stash/inventory"
	"gitlab.perso/poe-stash/scraper"
)

const (
	templateDir = scraper.DataDir + "template/"
)

// Generator construct html files from a scraped user.
type Generator struct {
	writer   io.Writer
	template *template.Template
}

// NewGenerator constructs a new generator.
func NewGenerator(writer io.Writer) Generator {
	t := template.Must(findAndParseTemplates(templateDir, ".tmpl", template.FuncMap{
		"DeducePosX": DeducePosX,
		"DeducePosY": DeducePosY,
	}))
	return Generator{
		writer:   writer,
		template: t,
	}
}

// findAndParseTemplates find all templates and initialize a template with it.
func findAndParseTemplates(rootDir, ext string, funcMap template.FuncMap) (*template.Template, error) {
	cleanRoot := filepath.Clean(rootDir)
	pfx := len(cleanRoot) + 1
	root := template.New("")

	err := filepath.Walk(cleanRoot, func(path string, info os.FileInfo, e1 error) error {
		if !info.IsDir() && strings.HasSuffix(path, ext) {
			if e1 != nil {
				return e1
			}

			b, e2 := ioutil.ReadFile(path)
			if e2 != nil {
				return e2
			}

			name := path[pfx:]
			t := root.New(name).Funcs(funcMap)
			t, e2 = t.Parse(string(b))
			if e2 != nil {
				return e2
			}
		}

		return nil
	})

	return root, err
}

// GenerateHTML generates HTML from scraped data.
func (g *Generator) GenerateHTML(data *scraper.ScrapedData) error {
	return g.template.ExecuteTemplate(g.writer, "layout", data)
}

// AdjustItemPos transforms relative stash position in
// absolute css position.
func AdjustItemPos(pos int) float64 {
	return float64(pos) * 47.4645
}

// DeducePosX transforms relative stash position in
// absolute css position using a given layout.
func DeducePosX(layout map[string]inventory.Layout, x, y int) float64 {
	if len(layout) > 0 {
		if value, ok := layout[strconv.Itoa(x)]; ok {
			return value.X
		}
		return 0
	}
	return float64(x) * 47.4645
}

// DeducePosY transforms relative stash position in
// absolute css position using a given layout.
func DeducePosY(layout map[string]inventory.Layout, x, y int) float64 {
	if len(layout) > 0 {
		if value, ok := layout[strconv.Itoa(x)]; ok {
			return value.Y
		}
		return 0
	}
	return float64(y) * 47.4645
}
