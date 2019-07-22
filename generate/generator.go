package generate

import (
	"html/template"
	"io"

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
	return Generator{
		writer:   writer,
		template: template.Must(template.ParseFiles(templateDir + "main.tpl")),
	}
}

// GenerateHTML generates HTML from scraped data.
func (g *Generator) GenerateHTML(data *scraper.ScrapedData) error {
	if err := g.template.Execute(g.writer, data); err != nil {
		return err
	}

	return nil
}
