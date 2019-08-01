package generate

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"reflect"
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
	return Generator{
		writer:   writer,
		template: template.Must(LoadAllTemplates()),
	}
}

// LoadAllTemplates load all templates.
func LoadAllTemplates() (*template.Template, error) {
	return FindAndParseTemplates(templateDir, ".tmpl", template.FuncMap{
		"DeducePosX":           DeducePosX,
		"DeducePosY":           DeducePosY,
		"ItemRarity":           ItemRarity,
		"ItemRarityType":       ItemRarityType,
		"ItemRarityHeight":     ItemRarityHeight,
		"GenSpecialBackground": GenSpecialBackground,
		"ColorType":            ColorType,
		"AugmentedType":        AugmentedType,
		"WordWrap":             WordWrap,
		"ConvToCssProgress":    ConvToCssProgress,
		"PoEMarkup":            PoEMarkup,
		"PoEMarkupLinesOnly":   PoEMarkupLinesOnly,
		"ColorToSocketClass":   ColorToSocketClass,
		"SocketRight":          SocketRight,
		"SocketedClass":        SocketedClass,
		"SocketedId":           SocketedId,
		"AltWeaponImage":       AltWeaponImage,
		"SellDescription":      SellDescription,
		"XpToNextLevel":        inventory.XpToNextLevel,
		"CurrentXp":            inventory.CurrentXp,
		"XpNeeded":             inventory.XpNeeded,
		"PrettyPrint":          inventory.PrettyPrint,
		"ContainsPattern":      ContainsPattern,
		"GenProperties":        GenProperties,
		"attr": func(s string) template.HTMLAttr {
			return template.HTMLAttr(s)
		},
		"safe": func(s string) template.HTML {
			return template.HTML(s)
		},
		"add": func(a, b int) int {
			return a + b
		},
		"dict": func(values ...interface{}) (map[string]interface{}, error) {
			if len(values) == 0 {
				return nil, errors.New("invalid dict call")
			}
			dict := make(map[string]interface{})
			for i := 0; i < len(values); i++ {
				key, isset := values[i].(string)
				if !isset {
					if reflect.TypeOf(values[i]).Kind() == reflect.Map {
						m := values[i].(map[string]interface{})
						for i, v := range m {
							dict[i] = v
						}
					} else {
						return nil, errors.New("dict values must be maps")
					}
				} else {
					i++
					if i == len(values) {
						return nil, errors.New("specify the key for non array values")
					}
					dict[key] = values[i]
				}
			}
			return dict, nil
		},
		"nl2br": func(line string) string {
			return strings.Replace(line, "\n", "<br />", -1)
		},
	})
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

// FindAndParseTemplates find all templates and initialize a template with it.
func FindAndParseTemplates(rootDir, ext string, funcMap template.FuncMap) (*template.Template, error) {
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
	return g.template.ExecuteTemplate(g.writer, "profile", data)
}

// DeducePosX transforms relative stash position in
// absolute css position using a given layout.
func DeducePosX(layoutType, inventoryId string, layout map[string]inventory.Layout, x, y, idx int) float64 {
	switch inventory.LayoutType(layoutType) {
	case inventory.CurrencyLayout:
		if value, ok := layout[strconv.Itoa(x)]; ok {
			return value.X
		}
	case inventory.FragmentLayout:
		if value, ok := layout[strconv.Itoa(x)]; ok {
			return value.X
		}
	case inventory.MapLayout:
		return 0
	case inventory.JewelLayout:
		return 287 + float64(idx)*47
	case inventory.InventoryLayout:
		key := inventoryId + "X"
		switch inventoryId {
		case "MainInventory":
			if value, ok := inventory.DefaultInventoryLayout[key]; ok {
				return value + float64(x)*cellSize
			}
		case "Flask":
			key = inventoryId + "X" + strconv.Itoa(x)
		}
		if value, ok := inventory.DefaultInventoryLayout[key]; ok {
			return value
		}
	default:
		return float64(x) * cellSize
	}
	return 0
}

// DeducePosY transforms relative stash position in
// absolute css position using a given layout.
func DeducePosY(layoutType, inventoryId string, layout map[string]inventory.Layout, x, y, idx int) float64 {
	switch inventory.LayoutType(layoutType) {
	case inventory.CurrencyLayout:
		if value, ok := layout[strconv.Itoa(x)]; ok {
			return value.Y
		}
	case inventory.FragmentLayout:
		if value, ok := layout[strconv.Itoa(x)]; ok {
			return value.Y
		}
	case inventory.MapLayout:
		return 0
	case inventory.JewelLayout:
		return -47
	case inventory.InventoryLayout:
		key := inventoryId + "Y"
		switch inventoryId {
		case "MainInventory":
			if value, ok := inventory.DefaultInventoryLayout[key]; ok {
				return value + float64(y)*cellSize
			}
		}
		if value, ok := inventory.DefaultInventoryLayout[key]; ok {
			return value
		}
	default:
		return float64(y) * cellSize
	}
	return 0
}

// rarityCharacteritics return the item visual characteristics to apply.
func rarityCharacteritics(frameType inventory.FrameType) (string, string, string) {
	switch frameType {
	case inventory.NormalItemFrameType:
		return "Normal", "normalPopup", ""
	case inventory.MagicItemFrameType:
		return "Magic", "magicPopup", ""
	case inventory.RareItemFrameType:
		return "Rare", "rarePopup", "doubleLine"
	case inventory.UniqueItemFrameType:
		return "Unique", "uniquePopup", "doubleLine"
	case inventory.GemFrameType:
		return "Gem", "gemPopup", ""
	case inventory.CurrencyFrameType:
		return "Currency", "currencyPopup", ""
	case inventory.DivinationCardFrameType:
		return "Divination Card", "divinationCard", "doubleLine"
	case inventory.QuestItemFrameType:
		return "Quest", "questPopup", ""
	case inventory.ProphecyFrameType:
		return "Normal", "prophecyPopup", ""
	case inventory.RelicFrameType:
		return "Relic", "relicPopup", ""
	default:
		return "", "", ""
	}
}

// ItemRarity return the correct class type from a frame type.
func ItemRarity(frameType inventory.FrameType) string {
	rarity, _, _ := rarityCharacteritics(frameType)
	return rarity
}

// ItemRarityType return the correct class type from a frame type.
func ItemRarityType(frameType inventory.FrameType) string {
	_, frameClass, _ := rarityCharacteritics(frameType)
	return frameClass
}

// ItemRarityHeight return the correct height from a frame type.
func ItemRarityHeight(frameType inventory.FrameType) string {
	_, _, heightClass := rarityCharacteritics(frameType)
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

// AugmentedType deduces the css class to colorize a property
// from a raw number.
func AugmentedType(atype float64) string {
	switch atype {
	case 1:
		return " (augmented)"
	default:
		return ""
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

// ColorToSocketClass convert a color into a socket class.
func ColorToSocketClass(color string) string {
	switch color {
	case "R":
		return "socketStr"
	case "G":
		return "socketDex"
	case "B":
		return "socketInt"
	case "W":
		return "socketGen"
	case "A":
		return "socketAbyss"
	default:
		return ""
	}
}

// SocketRight find if a socket has to be aligned right or not.
// Sockets are listed in a "snake" order.
// 0-1
//   |
// 3-2
// |
// 4-5
func SocketRight(idx int) string {
	switch idx {
	case 1, 2, 5:
		return "socketRight"
	}
	return ""
}

// searchSocketId search the right corresponding socket id.
// return -1 if nothing is found.
func searchSocketId(idx int, socketedItems []inventory.Item) int {
	for socketedIndex, v := range socketedItems {
		if idx == v.Socket {
			return socketedIndex
		}
	}
	return -1
}

// SocketedClass computes if a socket contains an item
// and construct everything needed to display it.
func SocketedClass(idx int, socketedItems []inventory.Item) string {
	// Search the coresponding socket id in socketed.
	iSocket := searchSocketId(idx, socketedItems)
	if iSocket < 0 || iSocket >= len(socketedItems) {
		return ""
	}
	item := socketedItems[iSocket]
	if item.IsAbyssJewel {
		return "socketed abyssJewel"
	}
	switch item.SocketColor {
	case "S":
		return "socketed strGem"
	case "D":
		return "socketed dexGem"
	case "I":
		return "socketed intGem"
	case "G":
		return "socketed genGem"
	default:
		return "socketed"
	}
}

// SocketedId computes id to attach mouseover to.
func SocketedId(idx int, socketedItems []inventory.Item) template.HTMLAttr {
	iSocket := searchSocketId(idx, socketedItems)
	if iSocket < 0 || iSocket >= len(socketedItems) {
		return ""
	}
	item := socketedItems[iSocket]
	return template.HTMLAttr(fmt.Sprintf(`id="item-%s"`, item.Id))
}

// AltWeaponImage returns the miniature image for alternative weapons.
func AltWeaponImage(items []*inventory.Item, filter string) template.HTMLAttr {
	for _, item := range items {
		if item.InventoryId == filter {
			top := 0.0
			switch item.Height {
			case 4:
				top = 6
			case 3:
				top = 14.625
			case 2:
				top = 23.25
			}
			left := 0.0
			switch item.Width {
			case 2:
				left = 5.9869
			case 1:
				left = 14.8357
			}
			return template.HTMLAttr(fmt.Sprintf(`src="%s" alt="" style="width: %fpx; height: %fpx; top: %fpx; left: %fpx;"`,
				item.Icon,
				float64(item.Width)*17.3287,
				float64(item.Height)*17.25,
				top,
				left,
			))
		}
	}
	return ""
}

// SellDescription generates the text for the trade forum.
func SellDescription(item inventory.Item, charName string) string {
	desc := ""
	if !strings.HasPrefix(item.InventoryId, "Stash") {
		desc = ` character="` + charName + `"`
	}
	return fmt.Sprintf(`[linkItem location=%s%s league="%s" x="%d" y="%d"]`,
		item.InventoryId, desc, item.League, item.X, item.Y)
}

// ContainsPattern checks if sentence contains any pattern
// like %0, %1, and so on...
func ContainsPattern(s string) bool {
	for i := 0; i < 10; i++ {
		if strings.Contains(s, "%"+strconv.Itoa(i)) {
			return true
		}
	}
	return false
}

// GenProperties generate properties for item with formatted
// strings like flasks.
func GenProperties(property inventory.ItemProperty) template.HTML {
	var args []interface{}
	for _, value := range property.Values {
		v := value.([]interface{})
		desc := v[0].(string)
		mode := ColorType(v[1].(float64))
		args = append(args, mode, desc)
	}
	pattern := property.Name
	for i := 0; i < 10; i++ {
		pattern = strings.ReplaceAll(
			pattern,
			"%"+strconv.Itoa(i),
			`<span class="%s">%s</span>`,
		)
	}
	return template.HTML(fmt.Sprintf(pattern, args...))
}
