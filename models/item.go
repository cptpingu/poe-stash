package models

// Item is a description of all properties of a single item.
type Item struct {
	IsVerified          bool `json:"verified"`
	IsIdentified        bool `json:"identified"`
	IsCorrupted         bool `json:"corrupted"`
	IsSynthesised       bool `json:"synthesised"`
	IsLockedToCharacter bool `json:"lockedToCharacter"`
	IsSupport           bool `json:"support"`
	IsRelic             bool `json:"isRelic"`
	IsAbyssJewel        bool `json:"abyssJewel"`
	IsVeiled            bool `json:"veiled"`
	IsDuplicated        bool `json:"duplicated"`

	Width        int `json:"w"`
	Height       int `json:"h"`
	ItemLevel    int `json:"ilvl"`
	StackSize    int `json:"stackSize"`
	MaxStackSize int `json:"maxStackSize"`
	X            int `json:"x"`
	Y            int `json:"y"`
	TalismanTier int `json:"talismanTier"`
	Socket       int `json:"socket"`

	FrameType FrameType `json:"frameType"`

	// Names for some items may include markup. For example: <<set:MS>><<set:M>><<set:S>>Roth's Reach
	Name                   string `json:"name"`
	Type                   string `json:"typeLine"`
	Icon                   string `json:"icon"`
	League                 string `json:"league"`
	Id                     string `json:"id"`
	DescriptionText        string `json:"descrText"`
	SecondDescriptionText  string `json:"secDescrText"`
	ArtFilename            string `json:"artFilename"`
	InventoryId            string `json:"inventoryId"`
	ProphecyText           string `json:"prophecyText"`
	ProphecyDifficultyText string `json:"prophecyDiffText"`
	Note                   string `json:"note"`
	SocketColor            string `json:"colour"`

	Influences           InfluenceProperty `json:"influences"`
	Properties           []ItemProperty    `json:"properties"`
	AdditionalProperties []ItemProperty    `json:"additionalProperties"`
	Requirements         []ItemProperty    `json:"requirements"`

	Sockets       []Socket `json:"sockets"`
	SocketedItems []Item   `json:"socketedItems"`

	ExplicitMods []string `json:"explicitMods"`
	ImplicitMods []string `json:"implicitMods"`
	UtilityMods  []string `json:"utilityMods"`
	EnchantMods  []string `json:"enchantMods"`
	CraftedMods  []string `json:"craftedMods"`
	CosmeticMods []string `json:"cosmeticMods"`
	VeiledMods   []string `json:"veiledMods"`
	FlavorText   []string `json:"flavourText"`

	Category      Category          `json:"category"`
	Hybrid        HybridType        `json:"hybrid"`
	IncubatedItem IncubatedItemType `json:"incubatedItem"`
}

// Socket describes a socket.
type Socket struct {
	GroupId   int    `json:"group"`
	Attribute string `json:"attr"`
	Color     string `json:"sColour"`
}

// ItemProperty holds item properties (name, how to display, ...).
type ItemProperty struct {
	Name        string        `json:"name"`
	Values      []interface{} `json:"values"`
	DisplayMode int           `json:"displayMode"`
	Progress    float64       `json:"progress"`
}

// InfluenceProperty holds which type of influences affect the item.
type InfluenceProperty struct {
	Elder    bool `json:"elder"`
	Shaper   bool `json:"shaper"`
	Crusader bool `json:"crusader"`
	Hunter   bool `json:"hunter"`
	Redeemer bool `json:"redeemer"`
	Warlord  bool `json:"warlord"`
}

// FrameType is a type of rarity of an item.
type FrameType int

// Frame type represents the type of frame to draw for an item.
const (
	NormalItemFrameType FrameType = iota
	MagicItemFrameType
	RareItemFrameType
	UniqueItemFrameType
	GemFrameType
	CurrencyFrameType
	DivinationCardFrameType
	QuestItemFrameType
	ProphecyFrameType
	RelicFrameType
)

// Category is the type of category of an item.
type Category struct {
	Armor       *[]string `json:"armour"`
	Accessories *[]string `json:"accessories"`
	Currency    *[]string `json:"currency"`
	Jewels      *[]string `json:"jewels"`
	Weapons     *[]string `json:"weapons"`
	Gems        *[]string `json:"gems"`
	Maps        *[]string `json:"maps"`
}

// HybridType represent vaal gems additional properties.
type HybridType struct {
	IsVaalGem             bool           `json:"isVaalGem"`
	BaseTypeName          string         `json:"baseTypeName"`
	SecondDescriptionText string         `json:"secDescrText"`
	Properties            []ItemProperty `json:"properties"`
	ExplicitMods          []string       `json:"explicitMods"`
}

// IncubatedItemType holds information about incubated item attached on an item.
type IncubatedItemType struct {
	Name     string `json:"name"`
	Level    int    `json:"level"`
	Progress int    `json:"progress"`
	Total    int    `json:"total"`
}

// LayoutType is the type of layout.
type LayoutType string

// Layout is the type of layout to use (type of grid to place items).
const (
	DefaultLayout    LayoutType = ""
	CurrencyLayout              = "currency"
	InventoryLayout             = "inventory"
	JewelLayout                 = "jewel"
	FragmentLayout              = "fragment"
	MapLayout                   = "map"
	QuadLayout                  = "quad"
	EssenceLayout               = "essence"
	DivinationLayout            = "divination"
)
