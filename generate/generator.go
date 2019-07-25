package generate

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gitlab.perso/poe-stash/inventory"
	"gitlab.perso/poe-stash/scraper"
)

const (
	templateDir = scraper.DataDir + "template/"
	cellSize    = 47.4645
	tokenErr    = "#Error(PoEMarkup)"
)

// Generator construct html files from a scraped user.
type Generator struct {
	writer   io.Writer
	template *template.Template
}

// NewGenerator constructs a new generator.
func NewGenerator(writer io.Writer) Generator {
	t := template.Must(findAndParseTemplates(templateDir, ".tmpl", template.FuncMap{
		"DeducePosX":           DeducePosX,
		"DeducePosY":           DeducePosY,
		"ItemRarityType":       ItemRarityType,
		"ItemRarityHeight":     ItemRarityHeight,
		"GenSpecialBackground": GenSpecialBackground,
		"ColorType":            ColorType,
		"WordWrap":             WordWrap,
		"ConvToCssProgress":    ConvToCssProgress,
		"PoEMarkup":            PoEMarkup,
		"PoEMarkupLinesOnly":   PoEMarkupLinesOnly,
		"attr": func(s string) template.HTMLAttr {
			return template.HTMLAttr(s)
		},
		"safe": func(s string) template.HTML {
			return template.HTML(s)
		},
	}))
	return Generator{
		writer:   writer,
		template: t,
	}
}

// WordWrap take string and apply an html wordwrap on it.
func WordWrap(s string) template.HTML {
	maxTextSize := 53
	parts := strings.SplitAfter(s, " ")
	res := ""
	nb := 0
	for _, part := range parts {
		nb += len(part)
		res += part
		if nb > maxTextSize {
			nb = 0
			res += "<br />"
		}
	}
	return template.HTML(res)
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

// DeducePosX transforms relative stash position in
// absolute css position using a given layout.
func DeducePosX(layout map[string]inventory.Layout, x, y int) float64 {
	if len(layout) > 0 {
		if value, ok := layout[strconv.Itoa(x)]; ok {
			return value.X
		}
		return 0
	}
	return float64(x) * cellSize
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
	return float64(y) * cellSize
}

// rarityCharacteritics return the item visual characteristics to apply.
func rarityCharacteritics(frameType inventory.FrameType) (string, string) {
	switch frameType {
	case inventory.NormalItemFrameType:
		return "normalPopup", ""
	case inventory.MagicItemFrameType:
		return "magicPopup", ""
	case inventory.RareItemFrameType:
		return "rarePopup", "doubleLine"
	case inventory.UniqueItemFrameType:
		return "uniquePopup", "doubleLine"
	case inventory.GemFrameType:
		return "gemPopup", ""
	case inventory.CurrencyFrameType:
		return "currencyPopup", ""
	case inventory.DivinationCardFrameType:
		return "divinationCard", "doubleLine"
	case inventory.QuestItemFrameType:
		return "questPopup", ""
	case inventory.ProphecyFrameType:
		return "prophecyPopup", ""
	case inventory.RelicFrameType:
		return "relicPopup", ""
	default:
		return "", ""
	}
}

// ItemRarityType return the correct class type from a frame type.
func ItemRarityType(frameType inventory.FrameType) string {
	frameClass, _ := rarityCharacteritics(frameType)
	return frameClass
}

// ItemRarityHeight return the correct height from a frame type.
func ItemRarityHeight(frameType inventory.FrameType) string {
	_, heightClass := rarityCharacteritics(frameType)
	return heightClass
}

// GenSpecialBackground generates a special background
// like shaper or elder ones.
func GenSpecialBackground(item inventory.Item) string {
	pattern := ""
	if item.IsShaper {
		pattern = "style='background-image: url(\"https://www.pathofexile.com/image/inventory/ShaperBackground.png?w=%d&h=%d&x=%d&y=%d\");'"
	}
	if item.IsElder {
		pattern = "style='background-image: url(\"https://www.pathofexile.com/image/inventory/ElderBackground.png?w=%d&h=%d&x=%d&y=%d\");'"
	}
	if pattern == "" {
		return ""
	}
	return fmt.Sprintf(pattern, item.Width, item.Height,
		int(float64(item.X)*cellSize), int(float64(item.Y)*cellSize))
}

// ColorType deduces the css class to colorize a property
// from a raw number.
func ColorType(colorType float64) string {
	switch colorType {
	case 1:
		return "colourAugmented"
	default:
		return "colourDefault"
	}
}

// ConvToCssProgress convert a progress into css percentage.
func ConvToCssProgress(progress float64) string {
	return strconv.Itoa(int(math.Round(progress*100))) + "%"
}

// replacePoEMarkup returns the line interpreted after markup interpretation.
// Grammar examples:
//  <property>{text}
//	<property>{<property>{text}}
func replacePoEMarkup(raw string) string {
	// Just a raw text, return it.
	first := strings.Index(raw, "<")
	if first < 0 {
		return raw
	}
	prefix := raw[:first]

	second := strings.Index(raw, ">")
	if second < 0 {
		return tokenErr
	}
	property := raw[first+1 : second]

	bracketL := strings.Index(raw, "{")
	if bracketL < 0 {
		return tokenErr
	}

	// Search matching "}".
	open := 1
	bracketR := -1
	for i := bracketL + 1; i < len(raw); i++ {
		if raw[i] == '{' {
			open++
		}
		if raw[i] == '}' {
			open--
		}
		if open == 0 {
			bracketR = i
			break
		}
	}
	if bracketR < 0 {
		return tokenErr
	}

	style := ""
	if strings.HasPrefix(property, "size:") {
		rawNb := property[len("size:"):]
		nb, err := strconv.Atoi(rawNb)
		if err != nil {
			return tokenErr
		}
		fontSize := float64(nb) / 2.0
		style = " style=\"font-size:" + strconv.FormatFloat(fontSize, 'f', -1, 64) + "px\""
		property = ""
	} else {
		property = " " + property
	}

	suffix := raw[bracketR+1:]

	return prefix +
		"<span class=\"PoEMarkup" + property + "\"" + style + ">" +
		replacePoEMarkup(raw[bracketL+1:bracketR]) +
		"</span>" +
		replacePoEMarkup(suffix)
}

// PoEMarkup converts a raw string containing markup into HTML.
func PoEMarkup(raw string) template.HTML {
	line := replacePoEMarkup(raw)
	lines := strings.Split(line, "\r\n")
	res := ""
	for _, line := range lines {
		res += "<div class=\"explicitMod\">\n"
		res += "  <span class=\"lc\">\n"
		res += "    " + line + "\n"
		res += "  </span>\n"
		res += "</div>\n"
	}
	return template.HTML(res)
}

// PoEMarkupLinesOnly converts a raw string containing markup into HTML.
// It is expexcted to only have lines separated by end of lines.
func PoEMarkupLinesOnly(lines []string) template.HTML {
	res := replacePoEMarkup(strings.Join(lines, "\n"))
	strings.Replace(res, "\n", "<br />", -1)
	return template.HTML(res)
}
