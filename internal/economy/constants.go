package economy

// CategoryForResource returns the category of a resource type.
func CategoryForResource(rt ResourceType) ResourceCategory {
	switch rt {
	case ResourceGrain, ResourceVegetables, ResourceWood, ResourceStone, ResourceIronOre, ResourceWool:
		return CategoryRaw
	case ResourceFlour, ResourceBread, ResourcePlanks, ResourceIron, ResourceCloth:
		return CategoryProcessed
	case ResourceTools, ResourceFurniture, ResourceWeapons, ResourceClothing:
		return CategoryAdvanced
	default:
		return CategoryRaw // fallback
	}
}

// SpoilageRate returns the base spoilage rate per week for a resource type.
// Food spoils faster than materials. Non-food resources have zero spoilage.
// Base rates: grain 1% (1x), vegetables 2% (2x), bread 0.5% (0.5x), flour 1% (1x).
func SpoilageRate(rt ResourceType) float64 {
	switch rt {
	case ResourceGrain, ResourceFlour:
		return 0.01 // 1% per week
	case ResourceVegetables:
		return 0.02 // 2% per week
	case ResourceBread:
		return 0.005 // 0.5% per week
	default:
		return 0.0 // all other resources do not spoil
	}
}

// BaseValue returns the base value per unit of a resource type.
func BaseValue(rt ResourceType) float64 {
	switch rt {
	// Raw materials
	case ResourceGrain:
		return 1.0
	case ResourceVegetables:
		return 1.2
	case ResourceWood:
		return 2.0
	case ResourceStone:
		return 3.0
	case ResourceIronOre:
		return 5.0
	case ResourceWool:
		return 4.0
	// Processed goods
	case ResourceFlour:
		return 1.5
	case ResourceBread:
		return 2.5
	case ResourcePlanks:
		return 3.0
	case ResourceIron:
		return 8.0
	case ResourceCloth:
		return 6.0
	// Advanced goods
	case ResourceTools:
		return 15.0
	case ResourceFurniture:
		return 20.0
	case ResourceWeapons:
		return 25.0
	case ResourceClothing:
		return 12.0
	default:
		return 0.0
	}
}
