package economy

import (
	"math/rand"
)

// Scavenge returns resources found in a given location with a given effort.
// Location can be "abandoned_mine", "fallen_forest", etc.
// Effort is a measure of how much work is put into scavenging (0..1).
// Returns a slice of resources (may be empty).
func Scavenge(location string, effort float64) []Resource {
	if effort <= 0 {
		return nil
	}
	switch location {
	case "abandoned_mine":
		return []Resource{
			NewResource(ResourceIronOre, 10*effort),
			NewResource(ResourceStone, 5*effort),
		}
	case "fallen_forest":
		return []Resource{
			NewResource(ResourceWood, 15*effort),
		}
	default:
		return nil
	}
}

// RecommendUse decides whether to use a scarce resource now or save for future.
// currentModifier is the current scarcity multiplier (e.g., 1.5 for rich mine).
// futureNeed is a factor indicating how critical the resource will be in future (>0).
// Returns true if the resource should be used now, false if it should be saved.
func RecommendUse(currentModifier, futureNeed float64) bool {
	// If current modifier is low (resource abundant) or future need is low, use now.
	// Threshold arbitrary: 0.5
	return currentModifier*futureNeed < 0.5
}

// DiscoverResource discovers a new resource source at the given location type.
// Location type can be "mine_vein" or "forest_patch".
// Returns a Mine or Forest with randomized parameters.
func DiscoverResource(locationType string) interface{} {
	switch locationType {
	case "mine_vein":
		totalOre := 100.0 + rand.Float64()*900.0    // 100-1000
		depletionRate := 0.05 + rand.Float64()*0.15 // 0.05-0.2
		return Mine{
			TotalOre:      totalOre,
			RemainingOre:  totalOre,
			DepletionRate: depletionRate,
		}
	case "forest_patch":
		health := 0.2 + rand.Float64()*0.8        // 0.2-1.0
		regrowthRate := 0.05 + rand.Float64()*0.1 // 0.05-0.15
		return Forest{
			Health:       health,
			RegrowthRate: regrowthRate,
		}
	default:
		return nil
	}
}

// ApplyScarcityModifier applies a scarcity multiplier to a production yield.
// Returns the adjusted yield.
func ApplyScarcityModifier(yield, modifier float64) float64 {
	return yield * modifier
}
