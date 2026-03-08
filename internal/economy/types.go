package economy

// QualityTier represents the quality level of a resource.
type QualityTier int

const (
	QualityPoor QualityTier = iota
	QualityNormal
	QualityGood
	QualityExcellent
	QualityMasterwork
)

// GameDate represents a date in the game calendar.
type GameDate struct {
	Year int
	Week int // 1-52
}

// ResourceType represents the type of resource.
type ResourceType string

const (
	// Raw materials (6)
	ResourceGrain      ResourceType = "grain"
	ResourceVegetables ResourceType = "vegetables"
	ResourceWood       ResourceType = "wood"
	ResourceStone      ResourceType = "stone"
	ResourceIronOre    ResourceType = "iron_ore"
	ResourceWool       ResourceType = "wool"

	// Processed goods (5)
	ResourceFlour  ResourceType = "flour"
	ResourceBread  ResourceType = "bread"
	ResourcePlanks ResourceType = "planks"
	ResourceIron   ResourceType = "iron"
	ResourceCloth  ResourceType = "cloth"

	// Advanced goods (4)
	ResourceTools     ResourceType = "tools"
	ResourceFurniture ResourceType = "furniture"
	ResourceWeapons   ResourceType = "weapons"
	ResourceClothing  ResourceType = "clothing"
)

// ResourceCategory represents the category of a resource.
type ResourceCategory string

const (
	CategoryRaw       ResourceCategory = "raw"
	CategoryProcessed ResourceCategory = "processed"
	CategoryAdvanced  ResourceCategory = "advanced"
)

// Resource represents a material resource in the economy.
type Resource struct {
	Type     ResourceType `json:"type"`
	Quantity float64      `json:"quantity"` // Units (kg, pieces, etc.)
	Quality  QualityTier  `json:"quality"`
	Location string       `json:"location"`          // Building ID where stored
	Produced GameDate     `json:"produced"`          // Production date (for spoilage)
	Value    float64      `json:"value"`             // Base value in abstract units
	Spoiled  float64      `json:"spoiled,omitempty"` // Amount already spoiled (cumulative)
}
