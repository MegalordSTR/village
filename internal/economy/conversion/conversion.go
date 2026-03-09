package conversion

import (
	"github.com/vano44/village/internal/economy"
	"github.com/vano44/village/internal/simulation"
	"math"
)

// ToSimulationResource converts an economy Resource to a simulation Resource.
// Note: simulation.Resource does not have Location, Produced, Value fields.
// Quality is mapped to float64 range [0,1].
func ToSimulationResource(er economy.Resource) simulation.Resource {
	return simulation.Resource{
		Type:     string(er.Type),
		Quantity: int(math.Round(er.Quantity)),
		Quality:  economy.QualityToFloat(er.Quality),
	}
}

// FromSimulationResource converts a simulation Resource to an economy Resource.
// Uses default values for missing fields (Location empty, Produced zero, Value base value).
func FromSimulationResource(sr simulation.Resource) economy.Resource {
	rt := economy.ResourceType(sr.Type)
	if !economy.IsValidType(rt) {
		// Fallback to grain if unknown
		rt = economy.ResourceGrain
	}
	return economy.Resource{
		Type:     rt,
		Quantity: float64(sr.Quantity),
		Quality:  economy.FloatToQuality(sr.Quality),
		Location: "",
		Produced: economy.GameDate{},
		Value:    economy.BaseValue(rt),
	}
}

// LoadInventoryFromGameState imports simulation resources into inventory, assigning them to defaultLocation.
func LoadInventoryFromGameState(inv *economy.Inventory, resources []simulation.Resource, defaultLocation string) error {
	for _, sr := range resources {
		r := FromSimulationResource(sr)
		r.Location = defaultLocation
		if err := inv.AddResource(defaultLocation, r); err != nil {
			return err
		}
	}
	return nil
}

// ExportInventoryToGameState exports all inventory resources as simulation resources (location information lost).
func ExportInventoryToGameState(inv *economy.Inventory) []simulation.Resource {
	var out []simulation.Resource
	resourcesMap := inv.ResourcesMap()
	for _, list := range resourcesMap {
		for _, r := range list {
			out = append(out, ToSimulationResource(r))
		}
	}
	return out
}
