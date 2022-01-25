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
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/cptpingu/poe-stash/misc"
	"github.com/cptpingu/poe-stash/models"
	"github.com/cptpingu/poe-stash/scraper"
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
		"InfluenceName":        InfluenceName,
		"GenSpecialBackground": GenSpecialBackground,
		"ColorType":            ColorType,
		"AugmentedType":        AugmentedType,
		"WordWrap":             WordWrap,
		"ConvToCssProgress":    ConvToCssProgress,
		"ReplacePoEMarkup":     ReplacePoEMarkup,
		"PoEMarkup":            PoEMarkup,
		"PoEMarkupLinesOnly":   PoEMarkupLinesOnly,
		"ColorToSocketClass":   ColorToSocketClass,
		"SocketRight":          SocketRight,
		"SocketedClass":        SocketedClass,
		"SocketedId":           SocketedId,
		"AltWeaponImage":       AltWeaponImage,
		"SellDescription":      SellDescription,
		"XpToNextLevel":        models.XpToNextLevel,
		"CurrentXp":            models.CurrentXp,
		"XpNeeded":             models.XpNeeded,
		"PrettyPrint":          models.PrettyPrint,
		"ContainsPattern":      ContainsPattern,
		"GenProperties":        GenProperties,
		"SearchItem":           SearchItem,
		"GenNaiveSearchIndex":  GenNaiveSearchIndex,
		"ItemCategory":         ItemCategory,
		"Version": func() string {
			return misc.Version
		},
		"attr": func(s string) template.HTMLAttr {
			return template.HTMLAttr(s)
		},
		"ieq": func(a, b string) bool {
			return strings.EqualFold(a, b)
		},
		"safe": func(s string) template.HTML {
			return template.HTML(s)
		},
		"add": func(a, b int) int {
			return a + b
		},
		"percentage": func(a, b int) float64 {
			if b == 0 {
				return 0
			}
			return (float64(a) / float64(b)) * 100
		},
		"squeeze": func(s string) string {
			return strings.Map(
				func(r rune) rune {
					if unicode.IsLetter(r) {
						return r
					}
					return -1
				},
				s,
			)
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
		"PrettyDate": func() string {
			return time.Now().Format("2006-01-02 15:04:05")
		},
		"DateFormat": func(d time.Time) string {
			return d.Format("2006-01-02")
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
			_, e2 = t.Parse(string(b))
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
func DeducePosX(layoutType, inventoryId string, layout map[string]models.Layout, x, y, idx int) float64 {
	switch models.LayoutType(layoutType) {
	case models.CurrencyLayout,
		models.FragmentLayout,
		models.EssenceLayout:
		if value, ok := layout[strconv.Itoa(x)]; ok {
			return value.X
		}
	case models.MapLayout:
		return 0
	case models.UniqueLayout:
		return 0
	case models.JewelLayout:
		return 287 + float64(idx)*47
	case models.InventoryLayout:
		key := inventoryId + "X"
		switch inventoryId {
		case "MainInventory":
			if value, ok := models.DefaultInventoryLayout[key]; ok {
				return value + float64(x)*cellSize
			}
		case "Flask":
			key = inventoryId + "X" + strconv.Itoa(x)
		}
		if value, ok := models.DefaultInventoryLayout[key]; ok {
			return value
		}
	case models.QuadLayout:
		return float64(x) * cellSize / 2
	default:
		return float64(x) * cellSize
	}
	return 0
}

// DeducePosY transforms relative stash position in
// absolute css position using a given layout.
func DeducePosY(layoutType, inventoryId string, layout map[string]models.Layout, x, y, idx int) float64 {
	switch models.LayoutType(layoutType) {
	case models.CurrencyLayout,
		models.FragmentLayout,
		models.EssenceLayout:
		if value, ok := layout[strconv.Itoa(x)]; ok {
			return value.Y
		}
	case models.MapLayout:
		return 0
	case models.UniqueLayout:
		return 0
	case models.JewelLayout:
		return -47
	case models.InventoryLayout:
		key := inventoryId + "Y"
		switch inventoryId {
		case "MainInventory":
			if value, ok := models.DefaultInventoryLayout[key]; ok {
				return value + float64(y)*cellSize
			}
		}
		if value, ok := models.DefaultInventoryLayout[key]; ok {
			return value
		}
	case models.QuadLayout:
		return float64(y) * cellSize / 2
	default:
		return float64(y) * cellSize
	}
	return 0
}

// rarityCharacteritics return the item visual characteristics to apply.
func rarityCharacteritics(frameType models.FrameType) (string, string, string) {
	switch frameType {
	case models.NormalItemFrameType:
		return "normal", "normalPopup", ""
	case models.MagicItemFrameType:
		return "magic", "magicPopup", ""
	case models.RareItemFrameType:
		return "rare", "rarePopup", "doubleLine"
	case models.UniqueItemFrameType:
		return "unique", "uniquePopup", "doubleLine"
	case models.GemFrameType:
		return "gem", "gemPopup", ""
	case models.CurrencyFrameType:
		return "currency", "currencyPopup", ""
	case models.DivinationCardFrameType:
		return "divinationcard", "divinationCard", "doubleLine"
	case models.QuestItemFrameType:
		return "quest", "questPopup", ""
	case models.ProphecyFrameType:
		return "normal", "prophecyPopup", ""
	case models.RelicFrameType:
		return "relic", "relicPopup", "doubleLine"
	default:
		return "", "", ""
	}
}

// ItemRarity return the correct class type from a frame type.
func ItemRarity(frameType models.FrameType) string {
	rarity, _, _ := rarityCharacteritics(frameType)
	return rarity
}

// ItemRarityType return the correct class type from a frame type.
func ItemRarityType(frameType models.FrameType) string {
	_, frameClass, _ := rarityCharacteritics(frameType)
	return frameClass
}

// ItemRarityHeight return the correct height from a frame type.
func ItemRarityHeight(frameType models.FrameType) string {
	_, _, heightClass := rarityCharacteritics(frameType)
	return heightClass
}

// InfluenceName returns the influence name.
// The left icon follows a predetermined hierarchy with Shaper always having top
// priority. It seems to follow the following order:
// Shaper > Elder > Crusader > Redeemer > Hunter > Warlord.
func InfluenceName(item models.Item, leftInfluence bool) string {
	res := ""

	// Shaper/Elder are top display priority.
	if item.Influences.Shaper && (!leftInfluence || res == "") {
		res = "shaper"
	}
	if item.Influences.Elder && (!leftInfluence || res == "") {
		res = "elder"
	}

	// Influences are display in a defined order.
	if item.Influences.Crusader && (!leftInfluence || res == "") {
		res = "crusader"
	}
	if item.Influences.Redeemer && (!leftInfluence || res == "") {
		res = "redeemer"
	}
	if item.Influences.Hunter && (!leftInfluence || res == "") {
		res = "hunter"
	}
	if item.Influences.Warlord && (!leftInfluence || res == "") {
		res = "warlord"
	}

	// Synthesised can't be influenced.
	if item.IsSynthesised && res == "" {
		res = "synthetic"
	}
	// Influenced veiled items are displayed as influenced items.
	// Very low priority.
	if item.IsVeiled && res == "" {
		res = "veiled"
	}

	return res
}

// GenSpecialBackground generates a special background
// like shaper or elder ones.
func GenSpecialBackground(item models.Item) string {
	pattern := ""
	if item.Influences.Shaper {
		pattern = "style='background-image: url(\"https://www.pathofexile.com/image/inventory/ShaperBackground.png?w=%d&h=%d&x=%d&y=%d\");'"
	}
	if item.Influences.Elder {
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

// ReplacePoEMarkup returns the line interpreted after markup interpretation.
// Grammar examples:
//  <property>{text}
//	<property>{<property>{text}}
func ReplacePoEMarkup(raw string, small bool) string {
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
		fontSize := float64(nb) * 0.5
		if small {
			fontSize = float64(nb) * 0.35
		}
		style = " style=\"font-size:" + strconv.FormatFloat(fontSize, 'f', -1, 64) + "px\""
		property = ""
	} else {
		property = " " + property
	}

	suffix := raw[bracketR+1:]

	return prefix +
		"<span class=\"PoEMarkup" + property + "\"" + style + ">" +
		ReplacePoEMarkup(raw[bracketL+1:bracketR], small) +
		"</span>" +
		ReplacePoEMarkup(suffix, small)
}

// PoEMarkup converts a raw string containing markup into HTML.
func PoEMarkup(raw string, small bool) template.HTML {
	line := ReplacePoEMarkup(raw, small)
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
func PoEMarkupLinesOnly(lines []string, small bool) template.HTML {
	res := ReplacePoEMarkup(strings.Join(lines, "\n"), small)
	res = strings.Replace(res, "\n", "<br />", -1)
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
func searchSocketId(idx int, socketedItems []models.Item) int {
	for socketedIndex, v := range socketedItems {
		if idx == v.Socket {
			return socketedIndex
		}
	}
	return -1
}

// SocketedClass computes if a socket contains an item
// and construct everything needed to display it.
func SocketedClass(idx int, socketedItems []models.Item) string {
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
func SocketedId(idx int, socketedItems []models.Item) template.HTMLAttr {
	iSocket := searchSocketId(idx, socketedItems)
	if iSocket < 0 || iSocket >= len(socketedItems) {
		return ""
	}
	item := socketedItems[iSocket]
	return template.HTMLAttr(fmt.Sprintf(`id="item-%s"`, item.Id))
}

// AltWeaponImage returns the miniature image for alternative weapons.
func AltWeaponImage(items []*models.Item, filter string) template.HTMLAttr {
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
func SellDescription(item models.Item, charName string) string {
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
func GenProperties(property models.ItemProperty) template.HTML {
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

// SearchItem search an item by its name.
func SearchItem(items []models.Item, name string) models.Item {
	for _, item := range items {
		if item.Type == name {
			return item
		}
	}
	return models.Item{}
}

// extractWords extracts relevant words from a sentence.
func extractWords(line string) []string {
	line = strings.ToLower(line)
	line = strings.ReplaceAll(line, "'", "")
	line = strings.ReplaceAll(line, ":", "")

	splitted := strings.Split(line, " ")
	res := make([]string, 0, len(splitted))
	for _, v := range splitted {
		if strings.Trim(v, " ") != "" {
			res = append(res, v)
		}
	}
	return res
}

// GenNaiveSearchIndex generates very naive indexing for an item description.
// It's just a list of selected unique sorted words.
func GenNaiveSearchIndex(item models.Item) string {
	words := make(map[string]struct{})

	// Extract name.
	for _, v := range extractWords(item.Name) {
		words[v] = struct{}{}
	}

	// Extract type of item.
	for _, v := range extractWords(item.Type) {
		words[v] = struct{}{}
	}

	// Extract properties.
	for _, mod := range item.ExplicitMods {
		for _, v := range extractWords(mod) {
			words[v] = struct{}{}
		}
	}
	for _, mod := range item.ImplicitMods {
		for _, v := range extractWords(mod) {
			words[v] = struct{}{}
		}
	}
	for _, mod := range item.UtilityMods {
		for _, v := range extractWords(mod) {
			words[v] = struct{}{}
		}
	}
	for _, mod := range item.EnchantMods {
		for _, v := range extractWords(mod) {
			words[v] = struct{}{}
		}
	}
	for _, mod := range item.CraftedMods {
		for _, v := range extractWords(mod) {
			words[v] = struct{}{}
		}
	}

	// Construct final string with sorted keywords.
	keys := make([]string, 0, len(words))
	for key := range words {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return strings.Join(keys, " ")
}

// itemCategoryAttribute returns attributes of an item category.
func itemCategoryAttribute(item models.Item) []string {
	res := make([]string, 0, 5)
	hasInfluence := false

	if item.Influences.Elder {
		hasInfluence = true
		res = append(res, "elder")
	}
	if item.Influences.Shaper {
		hasInfluence = true
		res = append(res, "shaper")
	}
	if item.Influences.Crusader {
		hasInfluence = true
		res = append(res, "crusader")
	}
	if item.Influences.Hunter {
		hasInfluence = true
		res = append(res, "hunter")
	}
	if item.Influences.Redeemer {
		hasInfluence = true
		res = append(res, "redeemer")
	}
	if item.Influences.Warlord {
		hasInfluence = true
		res = append(res, "warlord")
	}
	if hasInfluence {
		res = append(res, "influences", "influence")
	}

	if item.IsIdentified {
		res = append(res, "identified")
	} else {
		res = append(res, "unidentified")
	}
	if item.IsCorrupted {
		res = append(res, "corrupted", "corrupt")
	}
	if item.IsVeiled {
		res = append(res, "veiled", "veil")
	}
	if item.IsRelic {
		res = append(res, "relic")
	}
	if item.IsVerified {
		res = append(res, "verified")
	}
	if item.IsAbyssJewel {
		res = append(res, "abyss")
	}
	if len(item.ProphecyText) > 0 {
		res = append(res, "prophecy")
	}
	if item.Hybrid.IsVaalGem {
		res = append(res, "vaalgem", "vaal")
	}
	if len(item.ArtFilename) > 0 {
		res = append(res, "divination", "divine", "divcard", "divinationcard")
	}

	return res
}

// itemCategoryType returns types of an item category.
func itemCategoryType(categories models.Category) []string {
	res := make([]string, 0, 5)

	if categories.Armor != nil {
		res = append(res, models.Armors...)
		res = append(res, *categories.Armor...)
	}
	if categories.Accessories != nil {
		res = append(res, models.Accessories...)
		res = append(res, *categories.Accessories...)
	}
	if categories.Currency != nil {
		res = append(res, models.Currencies...)
		res = append(res, *categories.Currency...)
	}
	if categories.Jewels != nil {
		res = append(res, models.Jewels...)
		res = append(res, *categories.Jewels...)
	}
	if categories.Weapons != nil {
		res = append(res, models.Weapons...)
		res = append(res, *categories.Weapons...)
	}
	if categories.Gems != nil {
		res = append(res, models.Gems...)
		res = append(res, *categories.Gems...)
	}
	if categories.Maps != nil {
		res = append(res, models.Maps...)
		res = append(res, *categories.Maps...)
	}

	return res
}

// ItemCategory returns a text item category from categories.
func ItemCategory(item models.Item) string {
	attributes := itemCategoryAttribute(item)
	types := itemCategoryType(item.Category)
	if len(types) == 0 {
		types = tryDeduceTypeOfObject(item)
	}
	res := append(attributes, types...)
	sort.Strings(res)
	return strings.Join(res, " ")
}

// tryDeduceTypeOfObject tries to deduce categories of an object using object
// name when the Categories section is missing.
func tryDeduceTypeOfObject(item models.Item) []string {
	// Try to match known base type.
	if attributes, ok := models.BaseTypes[item.Type]; ok {
		return attributes
	}

	// Try to deduce item type with the frame type.
	switch item.FrameType {
	case models.GemFrameType:
		return models.Gems
	case models.CurrencyFrameType:
		if strings.Contains(item.Type, "Oil") {
			return models.Oils
		}
		if strings.Contains(item.Type, "Fossil") {
			return models.Fossils
		}
		if strings.Contains(item.Type, "Essence") {
			return models.Essences
		}
		return models.Currencies
	case models.DivinationCardFrameType:
		return models.Cards
	case models.QuestItemFrameType:
		return models.Quests
	}

	// Detect jewels type.
	if strings.Contains(item.Type, "Jewel") {
		return models.Jewels
	}

	// Detect map type.
	if strings.Contains(item.Type, "Map") {
		return models.Maps
	}

	return []string{}
}
