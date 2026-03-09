package economy

import (
	"testing"
)

func TestResourceTypeConstants(t *testing.T) {
	tests := []struct {
		rt   ResourceType
		want string
	}{
		{ResourceGrain, "grain"},
		{ResourceVegetables, "vegetables"},
		{ResourceWood, "wood"},
		{ResourceStone, "stone"},
		{ResourceIronOre, "iron_ore"},
		{ResourceWool, "wool"},
		{ResourceFlour, "flour"},
		{ResourceBread, "bread"},
		{ResourcePlanks, "planks"},
		{ResourceIron, "iron"},
		{ResourceCloth, "cloth"},
		{ResourceTools, "tools"},
		{ResourceFurniture, "furniture"},
		{ResourceWeapons, "weapons"},
		{ResourceClothing, "clothing"},
	}
	for _, tt := range tests {
		if string(tt.rt) != tt.want {
			t.Errorf("ResourceType %v got %q want %q", tt.rt, tt.rt, tt.want)
		}
	}
}

func TestQualityTierValues(t *testing.T) {
	if QualityPoor != 0 {
		t.Errorf("QualityPoor should be 0, got %d", QualityPoor)
	}
	if QualityNormal != 1 {
		t.Errorf("QualityNormal should be 1, got %d", QualityNormal)
	}
	if QualityGood != 2 {
		t.Errorf("QualityGood should be 2, got %d", QualityGood)
	}
	if QualityExcellent != 3 {
		t.Errorf("QualityExcellent should be 3, got %d", QualityExcellent)
	}
	if QualityMasterwork != 4 {
		t.Errorf("QualityMasterwork should be 4, got %d", QualityMasterwork)
	}
}

func TestCategoryForResource(t *testing.T) {
	tests := []struct {
		rt   ResourceType
		want ResourceCategory
	}{
		{ResourceGrain, CategoryRaw},
		{ResourceVegetables, CategoryRaw},
		{ResourceWood, CategoryRaw},
		{ResourceStone, CategoryRaw},
		{ResourceIronOre, CategoryRaw},
		{ResourceWool, CategoryRaw},
		{ResourceFlour, CategoryProcessed},
		{ResourceBread, CategoryProcessed},
		{ResourcePlanks, CategoryProcessed},
		{ResourceIron, CategoryProcessed},
		{ResourceCloth, CategoryProcessed},
		{ResourceTools, CategoryAdvanced},
		{ResourceFurniture, CategoryAdvanced},
		{ResourceWeapons, CategoryAdvanced},
		{ResourceClothing, CategoryAdvanced},
		{ResourceType("unknown"), CategoryRaw}, // default fallback
	}
	for _, tt := range tests {
		got := CategoryForResource(tt.rt)
		if got != tt.want {
			t.Errorf("CategoryForResource(%v) = %v, want %v", tt.rt, got, tt.want)
		}
	}
}

func TestGetCategory(t *testing.T) {
	// Ensure GetCategory matches CategoryForResource
	for _, rt := range []ResourceType{
		ResourceGrain, ResourceVegetables, ResourceWood, ResourceStone,
		ResourceIronOre, ResourceWool, ResourceFlour, ResourceBread,
		ResourcePlanks, ResourceIron, ResourceCloth, ResourceTools,
		ResourceFurniture, ResourceWeapons, ResourceClothing,
	} {
		if got := GetCategory(rt); got != CategoryForResource(rt) {
			t.Errorf("GetCategory(%v) = %v, CategoryForResource = %v", rt, got, CategoryForResource(rt))
		}
	}
}

func TestSpoilageRate(t *testing.T) {
	// Food resources should have positive spoilage
	if SpoilageRate(ResourceGrain) <= 0 {
		t.Errorf("Grain should have positive spoilage")
	}
	if SpoilageRate(ResourceBread) <= 0 {
		t.Errorf("Bread should have positive spoilage")
	}
	// Non-food resources should have zero spoilage
	if SpoilageRate(ResourceWood) != 0 {
		t.Errorf("Wood should have zero spoilage")
	}
	if SpoilageRate(ResourceStone) != 0 {
		t.Errorf("Stone should have zero spoilage")
	}
	if SpoilageRate(ResourceTools) != 0 {
		t.Errorf("Tools should have zero spoilage")
	}
}

func TestBaseValue(t *testing.T) {
	// Advanced goods should be more valuable than raw materials
	if BaseValue(ResourceTools) <= BaseValue(ResourceIron) {
		t.Errorf("Tools should be more valuable than iron")
	}
	if BaseValue(ResourceIron) <= BaseValue(ResourceIronOre) {
		t.Errorf("Processed iron should be more valuable than ore")
	}
	// All defined resource types should have positive base value
	resourceTypes := []ResourceType{
		ResourceGrain, ResourceVegetables, ResourceWood, ResourceStone,
		ResourceIronOre, ResourceWool, ResourceFlour, ResourceBread,
		ResourcePlanks, ResourceIron, ResourceCloth, ResourceTools,
		ResourceFurniture, ResourceWeapons, ResourceClothing,
	}
	for _, rt := range resourceTypes {
		if BaseValue(rt) <= 0 {
			t.Errorf("BaseValue(%v) should be positive, got %f", rt, BaseValue(rt))
		}
	}
	// Unknown resource type should return 0
	if BaseValue(ResourceType("unknown")) != 0 {
		t.Errorf("Unknown resource type should return 0")
	}
}

func TestIsValidType(t *testing.T) {
	// All defined resource types should be valid
	resourceTypes := []ResourceType{
		ResourceGrain, ResourceVegetables, ResourceWood, ResourceStone,
		ResourceIronOre, ResourceWool, ResourceFlour, ResourceBread,
		ResourcePlanks, ResourceIron, ResourceCloth, ResourceTools,
		ResourceFurniture, ResourceWeapons, ResourceClothing,
	}
	for _, rt := range resourceTypes {
		if !IsValidType(rt) {
			t.Errorf("IsValidType(%v) should be true", rt)
		}
	}
	// Unknown resource type should be invalid
	if IsValidType(ResourceType("unknown")) {
		t.Errorf("Unknown resource type should be invalid")
	}
}

func TestIsPerishable(t *testing.T) {
	// Perishable resources
	perishable := []ResourceType{
		ResourceGrain, ResourceVegetables, ResourceFlour, ResourceBread,
	}
	for _, rt := range perishable {
		if !IsPerishable(rt) {
			t.Errorf("%v should be perishable", rt)
		}
	}
	// Non-perishable resources (sample)
	nonPerishable := []ResourceType{
		ResourceWood, ResourceStone, ResourceIronOre, ResourceWool,
		ResourcePlanks, ResourceIron, ResourceCloth,
		ResourceTools, ResourceFurniture, ResourceWeapons, ResourceClothing,
	}
	for _, rt := range nonPerishable {
		if IsPerishable(rt) {
			t.Errorf("%v should not be perishable", rt)
		}
	}
}

func TestNewResource(t *testing.T) {
	r := NewResource(ResourceWood, 10.5)
	if r.Type != ResourceWood {
		t.Errorf("NewResource type mismatch: got %v, want %v", r.Type, ResourceWood)
	}
	if r.Quantity != 10.5 {
		t.Errorf("NewResource quantity mismatch: got %v, want 10.5", r.Quantity)
	}
	if r.Quality != QualityNormal {
		t.Errorf("NewResource quality mismatch: got %v, want %v", r.Quality, QualityNormal)
	}
	if r.Value != BaseValue(ResourceWood) {
		t.Errorf("NewResource value mismatch: got %v, want %v", r.Value, BaseValue(ResourceWood))
	}
}

func TestValidate(t *testing.T) {
	// Valid resource
	valid := Resource{
		Type:     ResourceGrain,
		Quantity: 5.0,
	}
	if !valid.Validate() {
		t.Errorf("Valid resource should pass validation")
	}
	// Zero quantity is valid
	zeroQty := Resource{
		Type:     ResourceGrain,
		Quantity: 0.0,
	}
	if !zeroQty.Validate() {
		t.Errorf("Zero quantity should be valid")
	}
	// Invalid type
	invalidType := Resource{
		Type:     ResourceType("invalid"),
		Quantity: 5.0,
	}
	if invalidType.Validate() {
		t.Errorf("Invalid type should fail validation")
	}
	// Negative quantity
	negativeQuantity := Resource{
		Type:     ResourceGrain,
		Quantity: -1.0,
	}
	if negativeQuantity.Validate() {
		t.Errorf("Negative quantity should fail validation")
	}
}

func TestComputeValue(t *testing.T) {
	// Test all quality tiers
	type qualityMultiplier struct {
		quality QualityTier
		mult    float64
	}
	multipliers := []qualityMultiplier{
		{QualityPoor, 0.8},
		{QualityNormal, 1.0},
		{QualityGood, 1.2},
		{QualityExcellent, 1.5},
		{QualityMasterwork, 2.0},
	}
	base := BaseValue(ResourceIron)
	quantity := 10.0
	for _, qm := range multipliers {
		r := Resource{
			Type:     ResourceIron,
			Quantity: quantity,
			Quality:  qm.quality,
			Value:    base,
		}
		got := r.ComputeValue()
		expected := quantity * base * qm.mult
		if got != expected {
			t.Errorf("ComputeValue with quality %v = %f, expected %f", qm.quality, got, expected)
		}
	}
	// Edge case: zero quantity
	r := Resource{
		Type:     ResourceIron,
		Quantity: 0.0,
		Quality:  QualityNormal,
		Value:    base,
	}
	if r.ComputeValue() != 0 {
		t.Errorf("Zero quantity should produce zero value")
	}
}

func TestGameDateZero(t *testing.T) {
	var d GameDate
	if d.Year != 0 || d.Week != 0 {
		t.Errorf("Zero GameDate should have zero fields")
	}
}

func TestQualityToFloat(t *testing.T) {
	if QualityToFloat(QualityPoor) != 0.0 {
		t.Errorf("QualityPoor should map to 0.0")
	}
	if QualityToFloat(QualityNormal) != 0.25 {
		t.Errorf("QualityNormal should map to 0.25")
	}
	if QualityToFloat(QualityGood) != 0.5 {
		t.Errorf("QualityGood should map to 0.5")
	}
	if QualityToFloat(QualityExcellent) != 0.75 {
		t.Errorf("QualityExcellent should map to 0.75")
	}
	if QualityToFloat(QualityMasterwork) != 1.0 {
		t.Errorf("QualityMasterwork should map to 1.0")
	}
}

func TestFloatToQuality(t *testing.T) {
	if FloatToQuality(0.0) != QualityPoor {
		t.Errorf("0.0 should map to QualityPoor")
	}
	if FloatToQuality(0.25) != QualityNormal {
		t.Errorf("0.25 should map to QualityNormal")
	}
	if FloatToQuality(0.5) != QualityGood {
		t.Errorf("0.5 should map to QualityGood")
	}
	if FloatToQuality(0.75) != QualityExcellent {
		t.Errorf("0.75 should map to QualityExcellent")
	}
	if FloatToQuality(1.0) != QualityMasterwork {
		t.Errorf("1.0 should map to QualityMasterwork")
	}
	// Test clamping
	if FloatToQuality(-0.1) != QualityPoor {
		t.Errorf("Negative should clamp to QualityPoor")
	}
	if FloatToQuality(1.5) != QualityMasterwork {
		t.Errorf(">1.0 should clamp to QualityMasterwork")
	}
}
