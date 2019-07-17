package inventory

// Item is a description of all properties of a single item.
type Item struct {
	IsVerified          bool `json:"verified"`
	IsIdentified        bool `json:"identified"`
	IsCorrupted         bool `json:"corrupted"`
	IsLockedToCharacter bool `json:"lockedToCharacter"`
	IsSupport           bool `json:"support"`
	IsRelic             bool `json:"isRelic"`

	Width        int `json:"w"`
	Height       int `json:"h"`
	ItemLevel    int `json:"ilvl"`
	StackSize    int `json:"stackSize"`
	MaxStackSize int `json:"maxStackSize"`
	X            int `json:"x"`
	Y            int `json:"y"`
	TalismanTier int `json:"talismanTier"`

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

	Properties   []ItemProperty `json:"properties"`
	Requirements []ItemProperty `json:"requirements"`

	Sockets       []Socket `json:"sockets"`
	SocketedItems []Item   `json:"socketedItems"`

	ExplicitMods []string `json:"explicitMods"`
	ImplicitMods []string `json:"implicitMods"`
	UtilityMods  []string `json:"utilityMods"`
	EnchantMods  []string `json:"enchantMods"`
	CraftedMods  []string `json:"craftedMods"`
	CosmeticMods []string `json:"cosmeticMods"`

	FlavorText []string `json:"flavourText"`
}

// Socket describes a socket.
type Socket struct {
	GroupId   int    `json:"group"`
	Attribute string `json:"attr"`
}

// ItemProperty holds item properties (name, how to display, ...).
type ItemProperty struct {
	Name        string        `json:"name"`
	Values      []interface{} `json:"values"`
	DisplayMode int           `json:"displayMode"`
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
