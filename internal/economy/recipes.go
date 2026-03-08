package economy

// ResourceRequirement defines a required input for a recipe.
type ResourceRequirement struct {
	Type       ResourceType `json:"type"`
	Quantity   float64      `json:"quantity"`
	MinQuality QualityTier  `json:"min_quality"` // minimum quality needed
}

// ResourceOutput defines a produced output from a recipe.
type ResourceOutput struct {
	Type               ResourceType `json:"type"`
	Quantity           float64      `json:"quantity"`
	QualityInheritance float64      `json:"quality_inheritance"` // factor 0..1, 0 = ignore input quality, 1 = full inheritance
}

// Recipe defines a production recipe.
type Recipe struct {
	ID         string                `json:"id"`
	Name       string                `json:"name"`
	Building   string                `json:"building"` // building type required
	Skill      string                `json:"skill"`    // skill ID required
	Inputs     []ResourceRequirement `json:"inputs"`
	Outputs    []ResourceOutput      `json:"outputs"`
	Time       float64               `json:"time"`        // hours per batch
	BaseYield  float64               `json:"base_yield"`  // multiplier for output quantity
	MaxWorkers int                   `json:"max_workers"` // maximum workers that can be assigned
}

// AllRecipes returns the list of all defined recipes.
func AllRecipes() []Recipe {
	return []Recipe{
		// Food chain
		{
			ID:       "mill_flour",
			Name:     "Mill Flour",
			Building: "mill",
			Skill:    "milling",
			Inputs: []ResourceRequirement{
				{Type: ResourceGrain, Quantity: 10, MinQuality: QualityPoor},
			},
			Outputs: []ResourceOutput{
				{Type: ResourceFlour, Quantity: 8, QualityInheritance: 0.8},
			},
			Time:       4.0,
			BaseYield:  1.0,
			MaxWorkers: 3,
		},
		{
			ID:       "bakery_bread",
			Name:     "Bake Bread",
			Building: "bakery",
			Skill:    "baking",
			Inputs: []ResourceRequirement{
				{Type: ResourceFlour, Quantity: 5, MinQuality: QualityPoor},
			},
			Outputs: []ResourceOutput{
				{Type: ResourceBread, Quantity: 10, QualityInheritance: 0.9},
			},
			Time:       3.0,
			BaseYield:  1.0,
			MaxWorkers: 2,
		},
		// Wood chain
		{
			ID:       "sawmill_planks",
			Name:     "Saw Planks",
			Building: "sawmill",
			Skill:    "carpentry",
			Inputs: []ResourceRequirement{
				{Type: ResourceWood, Quantity: 5, MinQuality: QualityPoor},
			},
			Outputs: []ResourceOutput{
				{Type: ResourcePlanks, Quantity: 4, QualityInheritance: 0.7},
			},
			Time:       6.0,
			BaseYield:  1.0,
			MaxWorkers: 4,
		},
		{
			ID:       "workshop_furniture",
			Name:     "Craft Furniture",
			Building: "workshop",
			Skill:    "carpentry",
			Inputs: []ResourceRequirement{
				{Type: ResourcePlanks, Quantity: 20, MinQuality: QualityNormal},
			},
			Outputs: []ResourceOutput{
				{Type: ResourceFurniture, Quantity: 1, QualityInheritance: 0.8},
			},
			Time:       24.0,
			BaseYield:  1.0,
			MaxWorkers: 5,
		},
		// Metal chain
		{
			ID:       "smelter_iron",
			Name:     "Smelt Iron",
			Building: "smelter",
			Skill:    "smithing",
			Inputs: []ResourceRequirement{
				{Type: ResourceIronOre, Quantity: 10, MinQuality: QualityPoor},
			},
			Outputs: []ResourceOutput{
				{Type: ResourceIron, Quantity: 5, QualityInheritance: 0.6},
			},
			Time:       8.0,
			BaseYield:  1.0,
			MaxWorkers: 6,
		},
		{
			ID:       "forge_tools",
			Name:     "Forge Tools",
			Building: "forge",
			Skill:    "smithing",
			Inputs: []ResourceRequirement{
				{Type: ResourceIron, Quantity: 2, MinQuality: QualityNormal},
			},
			Outputs: []ResourceOutput{
				{Type: ResourceTools, Quantity: 1, QualityInheritance: 0.9},
			},
			Time:       12.0,
			BaseYield:  1.0,
			MaxWorkers: 3,
		},
		{
			ID:       "forge_weapons",
			Name:     "Forge Weapons",
			Building: "forge",
			Skill:    "smithing",
			Inputs: []ResourceRequirement{
				{Type: ResourceIron, Quantity: 5, MinQuality: QualityGood},
			},
			Outputs: []ResourceOutput{
				{Type: ResourceWeapons, Quantity: 1, QualityInheritance: 0.95},
			},
			Time:       20.0,
			BaseYield:  1.0,
			MaxWorkers: 4,
		},
		// Textile chain
		{
			ID:       "spinner_yarn",
			Name:     "Spin Yarn",
			Building: "spinner",
			Skill:    "weaving",
			Inputs: []ResourceRequirement{
				{Type: ResourceWool, Quantity: 5, MinQuality: QualityPoor},
			},
			Outputs: []ResourceOutput{
				{Type: ResourceCloth, Quantity: 3, QualityInheritance: 0.8},
			},
			Time:       5.0,
			BaseYield:  1.0,
			MaxWorkers: 2,
		},
		{
			ID:       "tailor_clothing",
			Name:     "Sew Clothing",
			Building: "tailor",
			Skill:    "tailoring",
			Inputs: []ResourceRequirement{
				{Type: ResourceCloth, Quantity: 4, MinQuality: QualityNormal},
			},
			Outputs: []ResourceOutput{
				{Type: ResourceClothing, Quantity: 1, QualityInheritance: 0.85},
			},
			Time:       10.0,
			BaseYield:  1.0,
			MaxWorkers: 3,
		},
		// Additional food chain: vegetables -> cooked meal? but we don't have a resource. Let's add another grain to bread variant.
		{
			ID:       "bakery_bread2",
			Name:     "Bake Bread (vegetables)",
			Building: "bakery",
			Skill:    "baking",
			Inputs: []ResourceRequirement{
				{Type: ResourceVegetables, Quantity: 8, MinQuality: QualityPoor},
			},
			Outputs: []ResourceOutput{
				{Type: ResourceBread, Quantity: 12, QualityInheritance: 0.7},
			},
			Time:       4.0,
			BaseYield:  1.0,
			MaxWorkers: 2,
		},
	}
}

// FindRecipeByID returns the recipe with the given ID, or nil if not found.
func FindRecipeByID(id string) *Recipe {
	for _, r := range AllRecipes() {
		if r.ID == id {
			return &r
		}
	}
	return nil
}
