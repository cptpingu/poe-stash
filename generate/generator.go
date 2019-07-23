package generate

import (
	"html/template"
	"io"
	"strconv"

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
	t := template.Must(template.New("").Funcs(template.FuncMap{
		"DeducePosX": DeducePosX,
		"DeducePosY": DeducePosY,
	}).ParseGlob(templateDir + "*.tmpl"))
	return Generator{
		writer:   writer,
		template: t,
	}
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
