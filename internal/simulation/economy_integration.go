package simulation

import (
	"fmt"
	"github.com/vano44/village/internal/economy"
	"log"
	"math"
)

const maxResourcesPerLoad = 10000

// StringToResourceType converts a legacy string to a valid economy.ResourceType.
// It maps known legacy names (e.g., "food", "ore") to the appropriate economy constants.
// Returns an error for unknown resource types.
func StringToResourceType(s string) (economy.ResourceType, error) {
	rt := economy.ResourceType(s)
	if economy.IsValidType(rt) {
		return rt, nil
	}
	// Map legacy names
	switch s {
	case "food":
		return economy.ResourceGrain, nil
	case "ore":
		return economy.ResourceIronOre, nil
	case "tool":
		return economy.ResourceTools, nil
	case "wheat":
		return economy.ResourceGrain, nil
	case "gold":
		return economy.ResourceIronOre, nil
	case "meat":
		return economy.ResourceGrain, nil
	case "wood":
		return economy.ResourceWood, nil
	case "stone":
		return economy.ResourceStone, nil
	case "iron":
		return economy.ResourceIron, nil
	case "flour":
		return economy.ResourceFlour, nil
	case "bread":
		return economy.ResourceBread, nil
	case "planks":
		return economy.ResourcePlanks, nil
	case "cloth":
		return economy.ResourceCloth, nil
	case "wool":
		return economy.ResourceWool, nil
	case "vegetables":
		return economy.ResourceVegetables, nil
	case "iron_ore":
		return economy.ResourceIronOre, nil
	case "tools":
		return economy.ResourceTools, nil
	case "furniture":
		return economy.ResourceFurniture, nil
	case "weapons":
		return economy.ResourceWeapons, nil
	case "clothing":
		return economy.ResourceClothing, nil
	default:
		return economy.ResourceGrain, fmt.Errorf("unknown resource type %q", s)
	}
}

// IsKnownType returns true if the string is either a valid economy.ResourceType
// or a known legacy string that can be mapped.
func IsKnownType(s string) bool {
	rt := economy.ResourceType(s)
	if economy.IsValidType(rt) {
		return true
	}
	// Check legacy mapping
	switch s {
	case "food", "ore", "tool", "wheat", "gold", "meat", "wood", "stone", "iron", "flour", "bread", "planks", "cloth", "wool", "vegetables", "iron_ore", "tools", "furniture", "weapons", "clothing":
		return true
	}
	return false
}

// CalendarToGameDate converts simulation Calendar to economy GameDate.
func CalendarToGameDate(cal Calendar) economy.GameDate {
	return economy.GameDate{
		Year: cal.Year,
		Week: cal.Week,
	}
}

// ToEconomyResource converts a simulation Resource to an economy Resource.
// Uses default values for missing fields (Location empty, Produced zero, Value base value).
// Returns an error if the resource type is unknown.
func ToEconomyResource(sr Resource) (economy.Resource, error) {
	rt, err := StringToResourceType(string(sr.Type))
	if err != nil {
		return economy.Resource{}, err
	}
	return economy.Resource{
		Type:     rt,
		Quantity: float64(sr.Quantity),
		Quality:  economy.FloatToQuality(sr.Quality),
		Location: "",
		Produced: economy.GameDate{},
		Value:    economy.BaseValue(rt),
	}, nil
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
// Legacy resources are assigned the provided production date (typically current game date).
func LoadInventoryFromGameState(inv *economy.Inventory, resources []Resource, defaultLocation string, produced economy.GameDate) error {
	if len(resources) > maxResourcesPerLoad {
		return fmt.Errorf("resource array too large: %d > %d", len(resources), maxResourcesPerLoad)
	}
	for _, sr := range resources {
		r, err := ToEconomyResource(sr)
		if err != nil {
			return fmt.Errorf("invalid resource %v: %w", sr, err)
		}
		r.Location = defaultLocation
		r.Produced = produced
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

// AddProducedResource adds a newly produced resource to the game state.
// Adds with proper location and production date to Inventory (always present).
func AddProducedResource(state *GameState, building *Building, rt economy.ResourceType, quantity int, qualityFloat float64) error {
	qualityTier := economy.FloatToQuality(qualityFloat)
	er := economy.Resource{
		Type:     rt,
		Quantity: float64(quantity),
		Quality:  qualityTier,
		Location: building.Location,
		Produced: economy.GameDate{Year: state.Calendar.Year, Week: state.Calendar.Week},
		Value:    economy.BaseValue(rt),
	}
	err := state.Inventory.AddResource(building.Location, er)
	if err != nil {
		// Inventory.AddResource already logs the error
		return err
	}
	log.Printf("INFO: operation=AddProducedResource building=%s resource=%s quantity=%d quality=%f", building.Location, rt, quantity, qualityFloat)
	return nil
}

// ConsumeResourceFromState consumes up to the requested amount of a resource type.
// Uses Inventory (always present). Returns the amount actually consumed.
func ConsumeResourceFromState(state *GameState, rt economy.ResourceType, amount int) (int, error) {
	consumed, err := state.Inventory.RemoveResource("global", rt, float64(amount))
	if err != nil {
		// Inventory.RemoveResource already logs the error
		return 0, err
	}
	log.Printf("INFO: operation=ConsumeResourceFromState resource=%s requested=%d consumed=%d", rt, amount, int(consumed))
	return int(consumed), nil
}

// GetAvailableResourceFromState returns the total quantity of a resource type available.
// Uses Inventory (always present).
func GetAvailableResourceFromState(state *GameState, rt economy.ResourceType) int {
	return int(state.Inventory.GetAvailable("global", rt))
}
