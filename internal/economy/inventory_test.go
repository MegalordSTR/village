package economy

import (
	"github.com/vano44/village/internal/simulation"
	"testing"
)

func TestInventory_AddResource(t *testing.T) {
	inv := NewInventory()
	location := "granary1"
	resource := NewResource(ResourceGrain, 100.0)

	err := inv.AddResource(location, resource)
	if err != nil {
		t.Fatalf("AddResource failed: %v", err)
	}

	available := inv.GetAvailable(location, ResourceGrain)
	if available != 100.0 {
		t.Errorf("Expected 100.0 grain available, got %f", available)
	}
}

func TestInventory_RemoveResource(t *testing.T) {
	inv := NewInventory()
	location := "warehouse1"
	resource := NewResource(ResourceWood, 50.0)

	// Add first
	err := inv.AddResource(location, resource)
	if err != nil {
		t.Fatalf("AddResource failed: %v", err)
	}

	// Remove some
	removed, err := inv.RemoveResource(location, ResourceWood, 30.0)
	if err != nil {
		t.Fatalf("RemoveResource failed: %v", err)
	}
	if removed != 30.0 {
		t.Errorf("Expected removed 30.0, got %f", removed)
	}

	available := inv.GetAvailable(location, ResourceWood)
	if available != 20.0 {
		t.Errorf("Expected 20.0 wood remaining, got %f", available)
	}

	// Remove more than available should fail
	_, err = inv.RemoveResource(location, ResourceWood, 30.0)
	if err == nil {
		t.Errorf("Expected error when removing more than available")
	}
}

func TestInventory_TransferResource(t *testing.T) {
	inv := NewInventory()
	src := "granary1"
	dst := "granary2"
	resource := NewResource(ResourceVegetables, 80.0)

	err := inv.AddResource(src, resource)
	if err != nil {
		t.Fatalf("AddResource failed: %v", err)
	}

	transferred, err := inv.TransferResource(src, dst, ResourceVegetables, 50.0)
	if err != nil {
		t.Fatalf("TransferResource failed: %v", err)
	}
	if transferred != 50.0 {
		t.Errorf("Expected transferred 50.0, got %f", transferred)
	}

	srcAvail := inv.GetAvailable(src, ResourceVegetables)
	if srcAvail != 30.0 {
		t.Errorf("Expected src remaining 30.0, got %f", srcAvail)
	}
	dstAvail := inv.GetAvailable(dst, ResourceVegetables)
	if dstAvail != 50.0 {
		t.Errorf("Expected dst received 50.0, got %f", dstAvail)
	}
}

func TestInventory_GetAvailableEmpty(t *testing.T) {
	inv := NewInventory()
	available := inv.GetAvailable("nowhere", ResourceStone)
	if available != 0.0 {
		t.Errorf("Expected 0.0 for unknown location/resource, got %f", available)
	}
}

func TestInventory_CapacityLimit(t *testing.T) {
	storage := NewStorageRegistry()
	granary := &StorageBuilding{
		ID:    "granary1",
		Type:  StorageGranary,
		Level: 1,
	}
	err := storage.AddBuilding(granary)
	if err != nil {
		t.Fatalf("Failed to add building: %v", err)
	}
	inv := NewInventoryWithStorage(storage)

	// Add grain up to capacity (1000 kg)
	resource := NewResource(ResourceGrain, 1000.0)
	err = inv.AddResource("granary1", resource)
	if err != nil {
		t.Fatalf("AddResource failed: %v", err)
	}
	// Try to add one more kg, should fail
	extra := NewResource(ResourceGrain, 1.0)
	err = inv.AddResource("granary1", extra)
	if err == nil {
		t.Errorf("Expected error when exceeding capacity")
	}
	// Verify only 1000 kg stored
	available := inv.GetAvailable("granary1", ResourceGrain)
	if available != 1000.0 {
		t.Errorf("Expected 1000.0 grain stored, got %f", available)
	}
}

func TestInventory_StorageTypeRestriction(t *testing.T) {
	storage := NewStorageRegistry()
	granary := &StorageBuilding{
		ID:    "granary1",
		Type:  StorageGranary,
		Level: 1,
	}
	err := storage.AddBuilding(granary)
	if err != nil {
		t.Fatalf("Failed to add building: %v", err)
	}
	inv := NewInventoryWithStorage(storage)

	// Grain should be allowed
	grain := NewResource(ResourceGrain, 100.0)
	err = inv.AddResource("granary1", grain)
	if err != nil {
		t.Errorf("Grain should be allowed in granary: %v", err)
	}
	// Wood should be rejected
	wood := NewResource(ResourceWood, 100.0)
	err = inv.AddResource("granary1", wood)
	if err == nil {
		t.Errorf("Wood should not be allowed in granary")
	}
}

func TestInventory_CheckAlerts(t *testing.T) {
	inv := NewInventory()
	// Add some resources across locations
	err := inv.AddResource("loc1", NewResource(ResourceGrain, 50.0))
	if err != nil {
		t.Fatal(err)
	}
	err = inv.AddResource("loc2", NewResource(ResourceGrain, 30.0))
	if err != nil {
		t.Fatal(err)
	}
	err = inv.AddResource("loc1", NewResource(ResourceWood, 10.0))
	if err != nil {
		t.Fatal(err)
	}
	thresholds := AlertThreshold{
		ResourceGrain: 100.0, // total grain across locations = 80 < 100 -> alert
		ResourceWood:  5.0,   // wood total 10 > 5, no alert
	}
	alerts := inv.CheckAlerts(thresholds)
	// Expect single alert for grain (total 80 < 100)
	if len(alerts) != 1 {
		t.Errorf("Expected 1 alert, got %d", len(alerts))
	}
	for _, alert := range alerts {
		if alert.Resource != ResourceGrain {
			t.Errorf("Alert resource should be grain, got %v", alert.Resource)
		}
		if alert.Current >= alert.Threshold {
			t.Errorf("Alert current %f should be below threshold %f", alert.Current, alert.Threshold)
		}
		if alert.Location != "" {
			t.Errorf("Alert location should be empty for global total, got %s", alert.Location)
		}
	}
}

func TestInventory_GameStateIntegration(t *testing.T) {
	inv := NewInventory()
	// Create simulation resources
	simResources := []simulation.Resource{
		{Type: "grain", Quantity: 20, Quality: 0.25},
		{Type: "wood", Quantity: 15, Quality: 0.5},
	}
	err := LoadInventoryFromGameState(inv, simResources, "default")
	if err != nil {
		t.Fatalf("LoadInventoryFromGameState failed: %v", err)
	}
	// Check inventory has those resources
	grainAvail := inv.GetAvailable("default", ResourceGrain)
	if grainAvail != 20.0 {
		t.Errorf("Expected 20.0 grain, got %f", grainAvail)
	}
	woodAvail := inv.GetAvailable("default", ResourceWood)
	if woodAvail != 15.0 {
		t.Errorf("Expected 15.0 wood, got %f", woodAvail)
	}
	// Export back
	exported := ExportInventoryToGameState(inv)
	if len(exported) != 2 {
		t.Errorf("Expected 2 exported resources, got %d", len(exported))
	}
	// Check types and quantities
	// Order not guaranteed
	grainFound, woodFound := false, false
	for _, sr := range exported {
		if sr.Type == "grain" && sr.Quantity == 20 {
			grainFound = true
		}
		if sr.Type == "wood" && sr.Quantity == 15 {
			woodFound = true
		}
	}
	if !grainFound || !woodFound {
		t.Error("Exported resources missing expected items")
	}
}
