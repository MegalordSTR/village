package simulation

import (
	"io"
	"log"
	"testing"

	"github.com/vano44/village/internal/economy"
)

func init() {
	// Suppress log output during tests
	log.SetOutput(io.Discard)
}

func TestStringToResourceType_KnownLegacy(t *testing.T) {
	tests := []struct {
		input    string
		expected economy.ResourceType
	}{
		{"food", economy.ResourceGrain},
		{"ore", economy.ResourceIronOre},
		{"tool", economy.ResourceTools},
		{"wheat", economy.ResourceGrain},
		{"gold", economy.ResourceIronOre},
		{"meat", economy.ResourceGrain},
		{"wood", economy.ResourceWood},
		{"stone", economy.ResourceStone},
		{"iron", economy.ResourceIron},
		{"flour", economy.ResourceFlour},
		{"bread", economy.ResourceBread},
		{"planks", economy.ResourcePlanks},
		{"cloth", economy.ResourceCloth},
		{"wool", economy.ResourceWool},
		{"vegetables", economy.ResourceVegetables},
		{"iron_ore", economy.ResourceIronOre},
		{"tools", economy.ResourceTools},
		{"furniture", economy.ResourceFurniture},
		{"weapons", economy.ResourceWeapons},
		{"clothing", economy.ResourceClothing},
	}
	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			got, err := StringToResourceType(tc.input)
			if err != nil {
				t.Errorf("StringToResourceType(%q) returned error: %v", tc.input, err)
			}
			if got != tc.expected {
				t.Errorf("StringToResourceType(%q) = %v, want %v", tc.input, got, tc.expected)
			}
		})
	}
}

func TestStringToResourceType_Unknown(t *testing.T) {
	// Unknown string should return an error
	got, err := StringToResourceType("unknown_resource")
	if err == nil {
		t.Errorf("StringToResourceType(\"unknown_resource\") should return error, got %v", got)
	}
	// Returned type may be ResourceGrain (default) but error is primary
	if got != economy.ResourceGrain {
		t.Errorf("StringToResourceType(\"unknown_resource\") = %v, want %v on error", got, economy.ResourceGrain)
	}
}

func TestStringToResourceType_ValidType(t *testing.T) {
	// If the string matches a valid economy.ResourceType constant, return it unchanged
	tests := []economy.ResourceType{
		economy.ResourceGrain,
		economy.ResourceWood,
		economy.ResourceStone,
		economy.ResourceIronOre,
		economy.ResourceIron,
		economy.ResourceFlour,
		economy.ResourceBread,
		economy.ResourcePlanks,
		economy.ResourceCloth,
		economy.ResourceWool,
		economy.ResourceVegetables,
		economy.ResourceTools,
		economy.ResourceFurniture,
		economy.ResourceWeapons,
		economy.ResourceClothing,
	}
	for _, rt := range tests {
		t.Run(string(rt), func(t *testing.T) {
			got, err := StringToResourceType(string(rt))
			if err != nil {
				t.Errorf("StringToResourceType(%q) returned error: %v", rt, err)
			}
			if got != rt {
				t.Errorf("StringToResourceType(%q) = %v, want %v", rt, got, rt)
			}
		})
	}
}

func TestLoadInventoryFromGameState_SizeLimit(t *testing.T) {
	inv := economy.NewInventory()
	resources := make([]Resource, 10001) // exceed limit
	for i := range resources {
		resources[i] = Resource{
			Type:     "grain",
			Quantity: 1,
			Quality:  1.0,
		}
	}
	err := LoadInventoryFromGameState(inv, resources, "global", economy.GameDate{})
	if err == nil {
		t.Error("LoadInventoryFromGameState should reject oversized array")
	}
}

func TestLoadInventoryFromGameState_UnknownResource(t *testing.T) {
	inv := economy.NewInventory()
	resources := []Resource{
		{
			Type:     "unknown_resource_type",
			Quantity: 10,
			Quality:  1.0,
		},
	}
	err := LoadInventoryFromGameState(inv, resources, "global", economy.GameDate{})
	if err == nil {
		t.Error("LoadInventoryFromGameState should reject unknown resource type")
	}
}
