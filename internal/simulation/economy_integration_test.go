package simulation

import (
	"testing"

	"github.com/vano44/village/internal/economy"
)

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
			got := StringToResourceType(tc.input)
			if got != tc.expected {
				t.Errorf("StringToResourceType(%q) = %v, want %v", tc.input, got, tc.expected)
			}
		})
	}
}

func TestStringToResourceType_Unknown(t *testing.T) {
	// Unknown string should map to ResourceGrain
	got := StringToResourceType("unknown_resource")
	if got != economy.ResourceGrain {
		t.Errorf("StringToResourceType(\"unknown_resource\") = %v, want %v", got, economy.ResourceGrain)
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
			got := StringToResourceType(string(rt))
			if got != rt {
				t.Errorf("StringToResourceType(%q) = %v, want %v", rt, got, rt)
			}
		})
	}
}
