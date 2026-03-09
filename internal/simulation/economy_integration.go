package simulation

import (
	"github.com/vano44/village/internal/economy"
	"math"
)

// StringToResourceType converts a legacy string to a valid economy.ResourceType.
// It maps known legacy names (e.g., "food", "ore") to the appropriate economy constants.
func StringToResourceType(s string) economy.ResourceType {
	rt := economy.ResourceType(s)
	if economy.IsValidType(rt) {
		return rt
	}
	// Map legacy names
	switch s {
	case "food":
		return economy.ResourceGrain
	case "ore":
		return economy.ResourceIronOre
	case "tool":
		return economy.ResourceTools
	case "wheat":
		return economy.ResourceGrain
	case "gold":
		return economy.ResourceIronOre
	default:
		return economy.ResourceGrain
	}
}

// ToEconomyResource converts a simulation Resource to an economy Resource.
// Uses default values for missing fields (Location empty, Produced zero, Value base value).
func ToEconomyResource(sr Resource) economy.Resource {
	rt := StringToResourceType(string(sr.Type))
	return economy.Resource{
		Type:     rt,
		Quantity: float64(sr.Quantity),
		Quality:  economy.FloatToQuality(sr.Quality),
		Location: "",
		Produced: economy.GameDate{},
		Value:    economy.BaseValue(rt),
	}
}

// FromEconomyResource converts an economy Resource to a simulation Resource.
// Note: simulation.Resource does not have Location, Produced, Value fields.
// Quality is mapped to float64 range [0,1].
func FromEconomyResource(er economy.Resource) Resource {
	return Resource{
		Type:     er.Type,
		Quantity: int(math.Round(er.Quantity)),
		Quality:  economy.QualityToFloat(er.Quality),
	}
}

// LoadInventoryFromGameState imports simulation resources into inventory, assigning them to defaultLocation.
func LoadInventoryFromGameState(inv *economy.Inventory, resources []Resource, defaultLocation string) error {
	for _, sr := range resources {
		r := ToEconomyResource(sr)
		r.Location = defaultLocation
		if err := inv.AddResource(defaultLocation, r); err != nil {
			return err
		}
	}
	return nil
}

// ExportInventoryToGameState exports all inventory resources as simulation resources (location information lost).
func ExportInventoryToGameState(inv *economy.Inventory) []Resource {
	var out []Resource
	resourcesMap := inv.ResourcesMap()
	for _, list := range resourcesMap {
		for _, r := range list {
			out = append(out, FromEconomyResource(r))
		}
	}
	return out
}
