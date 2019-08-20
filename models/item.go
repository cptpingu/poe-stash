package models

// Item is a description of all properties of a single item.
type Item struct {
	IsVerified          bool `json:"verified"`
	IsIdentified        bool `json:"identified"`
	IsCorrupted         bool `json:"corrupted"`
	IsLockedToCharacter bool `json:"lockedToCharacter"`
	IsSupport           bool `json:"support"`
	IsRelic             bool `json:"isRelic"`
	IsElder             bool `json:"elder"`
	IsShaper            bool `json:"shaper"`
	IsAbyssJewel        bool `json:"abyssJewel"`
	IsVeiled            bool `json:"veiled"`

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

	Properties           []ItemProperty `json:"properties"`
	AdditionalProperties []ItemProperty `json:"additionalProperties"`
	Requirements         []ItemProperty `json:"requirements"`

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

	Category Category   `json:"category"`
	Hybrid   HybridType `json:"hybrid"`
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

// FrameType is a type of rarity of an item.
type FrameType int

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
	Armor       []string `json:"armor"`
	Accessories []string `json:"accessories"`
	Currency    []string `json:"currency"`
	Jewels      []string `json:"jewels"`
	Weapons     []string `json:"weapons"`
}

// HybridType represent vaal gems additional properties.
type HybridType struct {
	IsVaalGem             bool           `json:"isVaalGem"`
	BaseTypeName          string         `json:"baseTypeName"`
	SecondDescriptionText string         `json:"secDescrText"`
	Properties            []ItemProperty `json:"properties"`
	ExplicitMods          []string       `json:"explicitMods"`
}

// LayoutType is the type of layout.
type LayoutType string

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
