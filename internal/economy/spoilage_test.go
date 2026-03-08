package economy

import (
	"testing"
)

func TestWeeksSince(t *testing.T) {
	tests := []struct {
		current, produced GameDate
		want              int
	}{
		{GameDate{Year: 1, Week: 1}, GameDate{Year: 1, Week: 1}, 0},
		{GameDate{Year: 1, Week: 5}, GameDate{Year: 1, Week: 1}, 4},
		{GameDate{Year: 2, Week: 1}, GameDate{Year: 1, Week: 1}, 52},
		{GameDate{Year: 2, Week: 10}, GameDate{Year: 1, Week: 50}, 12}, // cross year
		{GameDate{Year: 1, Week: 1}, GameDate{Year: 1, Week: 5}, 0},    // produced after current
	}
	for _, tt := range tests {
		got := WeeksSince(tt.current, tt.produced)
		if got != tt.want {
			t.Errorf("WeeksSince(%v, %v) = %d, want %d", tt.current, tt.produced, got, tt.want)
		}
	}
}

func TestApplySpoilage_Grain_Granary_Normal(t *testing.T) {
	storage := &StorageBuilding{
		Type:    StorageGranary,
		Quality: StorageQualityNormal,
		Level:   1,
	}
	resource := NewResource(ResourceGrain, 100.0)
	resource.Produced = GameDate{Year: 1, Week: 1}
	current := GameDate{Year: 1, Week: 10} // 9 weeks later
	spoiled, degraded := ApplySpoilage(&resource, storage, current)
	// base rate 0.01 per week, storage multiplier 1.0, age 9 weeks => spoiled fraction = 0.09
	// original quantity 100, spoiled = 9
	expectedSpoiled := 9.0
	if spoiled < expectedSpoiled-0.001 || spoiled > expectedSpoiled+0.001 {
		t.Errorf("spoiled = %f, want %f", spoiled, expectedSpoiled)
	}
	if degraded {
		t.Error("quality should not degrade at 9% spoilage")
	}
	if resource.Quantity != 100.0-spoiled {
		t.Errorf("quantity after spoilage = %f, want %f", resource.Quantity, 100.0-spoiled)
	}
	if resource.Spoiled != spoiled {
		t.Errorf("resource.Spoiled = %f, want %f", resource.Spoiled, spoiled)
	}
}

func TestApplySpoilage_Vegetables_OutdoorPile_Poor(t *testing.T) {
	storage := &StorageBuilding{
		Type:    StorageOutdoorPile,
		Quality: StorageQualityPoor,
		Level:   1,
	}
	resource := NewResource(ResourceVegetables, 200.0)
	resource.Produced = GameDate{Year: 1, Week: 1}
	current := GameDate{Year: 1, Week: 5} // 4 weeks later
	spoiled, degraded := ApplySpoilage(&resource, storage, current)
	// base rate 0.02 per week, storage multiplier 2.0 (poor), age 4 weeks => totalRate = 0.04, fraction = 0.16
	// original quantity 200, spoiled = 32
	expectedSpoiled := 32.0
	if spoiled < expectedSpoiled-0.001 || spoiled > expectedSpoiled+0.001 {
		t.Errorf("spoiled = %f, want %f", spoiled, expectedSpoiled)
	}
	if degraded {
		t.Error("quality should not degrade at 16% spoilage")
	}
}

func TestApplySpoilage_QualityDegradation(t *testing.T) {
	storage := &StorageBuilding{
		Type:    StorageGranary,
		Quality: StorageQualityNormal,
		Level:   1,
	}
	resource := NewResource(ResourceGrain, 100.0)
	resource.Quality = QualityGood
	resource.Produced = GameDate{Year: 1, Week: 1}
	current := GameDate{Year: 1, Week: 30} // 29 weeks later
	spoiled, degraded := ApplySpoilage(&resource, storage, current)
	// base rate 0.01, multiplier 1, age 29 => fraction 0.29 > 0.2 threshold
	if !degraded {
		t.Error("quality should degrade when spoilage exceeds 20%")
	}
	if resource.Quality != QualityNormal {
		t.Errorf("quality after degradation = %v, want %v", resource.Quality, QualityNormal)
	}
	if spoiled <= 0 {
		t.Error("spoiled should be positive")
	}
}

func TestApplySpoilageToInventory(t *testing.T) {
	storageReg := NewStorageRegistry()
	granary := &StorageBuilding{
		ID:      "granary1",
		Type:    StorageGranary,
		Quality: StorageQualityNormal,
		Level:   1,
	}
	storageReg.AddBuilding(granary)
	inv := NewInventoryWithStorage(storageReg)
	// Add some grain
	grain := NewResource(ResourceGrain, 500.0)
	grain.Produced = GameDate{Year: 1, Week: 1}
	inv.AddResource("granary1", grain)
	// Add vegetables elsewhere with no storage building (should not spoil)
	veg := NewResource(ResourceVegetables, 300.0)
	veg.Produced = GameDate{Year: 1, Week: 1}
	inv.AddResource("outdoor1", veg)
	current := GameDate{Year: 1, Week: 10}
	totals := ApplySpoilageToInventory(inv, storageReg, current)
	// Only grain should spoil
	grainSpoiled, ok := totals["grain"]
	if !ok {
		t.Error("expected spoilage for grain")
	}
	if grainSpoiled <= 0 {
		t.Error("grain spoilage should be positive")
	}
	// Vegetables should have no entry
	if _, ok := totals["vegetables"]; ok {
		t.Error("vegetables should not spoil (no storage building)")
	}
	// Check inventory quantities updated
	avail := inv.GetAvailable("granary1", ResourceGrain)
	if avail >= 500.0 {
		t.Errorf("grain quantity should be reduced, got %f", avail)
	}
}
