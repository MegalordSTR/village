package simulation

import (
	"github.com/vano44/village/internal/economy"
	"testing"
)

func TestEconomyImplementsSystem(t *testing.T) {
	var _ System = (*EconomicSystem)(nil)
}

func TestNewEconomy(t *testing.T) {
	econ := NewEconomicSystem()
	if econ == nil {
		t.Fatal("NewEconomicSystem returned nil")
	}
}

func TestEconomyUpdateDeterministic(t *testing.T) {
	econ := NewEconomicSystem()

	// Create two states with same seed
	state1 := NewGameState("test1", 123)
	state2 := NewGameState("test2", 123)

	// Add some residents and resources for testing
	state1.AddResident(Resident{ID: "r1", Name: "Alice", Age: 30})
	if err := state1.AddResource(Resource{Type: economy.ResourceGrain, Quantity: 100, Quality: 1.0}); err != nil {
		t.Fatalf("AddResource failed: %v", err)
	}
	state2.AddResident(Resident{ID: "r1", Name: "Alice", Age: 30})
	if err := state2.AddResource(Resource{Type: economy.ResourceGrain, Quantity: 100, Quality: 1.0}); err != nil {
		t.Fatalf("AddResource failed: %v", err)
	}

	// Process first week
	events1 := econ.Update(1, state1, state1.RNG.Rand())
	events2 := econ.Update(1, state2, state2.RNG.Rand())

	// Should generate same events
	if len(events1) != len(events2) {
		t.Errorf("event count mismatch: %d vs %d", len(events1), len(events2))
	}

	// Resource quantities should change deterministically
	// For now just ensure no panic
}

func TestEconomyUpdateConsumesFood(t *testing.T) {
	econ := NewEconomicSystem()
	state := NewGameState("test", 456)
	state.AddResident(Resident{ID: "r1", Name: "Bob", Age: 25})
	if err := state.AddResource(Resource{Type: economy.ResourceGrain, Quantity: 50, Quality: 1.0}); err != nil {
		t.Fatalf("AddResource failed: %v", err)
	}

	initialFood := 0
	for _, r := range state.Resources {
		if r.Type == economy.ResourceGrain {
			initialFood = r.Quantity
		}
	}

	events := econ.Update(state.Calendar.Week, state, state.RNG.Rand())

	// Should consume some food
	finalFood := 0
	for _, r := range state.Resources {
		if r.Type == economy.ResourceGrain {
			finalFood = r.Quantity
		}
	}
	if finalFood >= initialFood {
		t.Errorf("food consumption didn't happen: initial %d, final %d", initialFood, finalFood)
	}

	// Should generate consumption event
	foundConsumption := false
	for _, ev := range events {
		if ev.Type == "consumption" {
			foundConsumption = true
		}
	}
	if !foundConsumption {
		t.Error("no consumption event generated")
	}
}

func TestEconomyUpdateBuildingMaintenance(t *testing.T) {
	econ := NewEconomicSystem()
	state := NewGameState("test", 789)
	state.AddBuilding(Building{Type: "house", Location: "north", Level: 1})
	if err := state.AddResource(Resource{Type: "wood", Quantity: 20, Quality: 1.0}); err != nil {
		t.Fatalf("AddResource failed: %v", err)
	}
	if err := state.AddResource(Resource{Type: "stone", Quantity: 10, Quality: 1.0}); err != nil {
		t.Fatalf("AddResource failed: %v", err)
	}

	initialWood := 0
	initialStone := 0
	for _, r := range state.Resources {
		if r.Type == "wood" {
			initialWood = r.Quantity
		}
		if r.Type == "stone" {
			initialStone = r.Quantity
		}
	}

	events := econ.Update(state.Calendar.Week, state, state.RNG.Rand())

	finalWood := 0
	finalStone := 0
	for _, r := range state.Resources {
		if r.Type == "wood" {
			finalWood = r.Quantity
		}
		if r.Type == "stone" {
			finalStone = r.Quantity
		}
	}

	// Maintenance should consume some wood/stone
	if finalWood >= initialWood && finalStone >= initialStone {
		t.Errorf("maintenance didn't consume resources: wood %d->%d, stone %d->%d",
			initialWood, finalWood, initialStone, finalStone)
	}

	// Should generate maintenance event
	foundMaintenance := false
	for _, ev := range events {
		if ev.Type == "maintenance" {
			foundMaintenance = true
		}
	}
	if !foundMaintenance {
		t.Error("no maintenance event generated")
	}
}

func TestEconomyFoodSpoilage(t *testing.T) {
	econ := NewEconomicSystem()
	state := NewGameState("test", 999)
	if err := state.AddResource(Resource{Type: economy.ResourceGrain, Quantity: 100, Quality: 1.0}); err != nil {
		t.Fatalf("AddResource failed: %v", err)
	}

	initialFood := 0
	for _, r := range state.Resources {
		if r.Type == economy.ResourceGrain {
			initialFood = r.Quantity
		}
	}

	// Run multiple weeks to allow spoilage
	for week := 1; week <= 4; week++ {
		_ = econ.Update(week, state, state.RNG.Rand())
	}

	finalFood := 0
	for _, r := range state.Resources {
		if r.Type == economy.ResourceGrain {
			finalFood = r.Quantity
		}
	}

	// Food should spoil over time
	if finalFood >= initialFood {
		t.Errorf("food spoilage didn't happen: initial %d, final %d", initialFood, finalFood)
	}
}

func TestEconomyStorageLimits(t *testing.T) {
	econ := NewEconomicSystem()
	state := NewGameState("test", 111)
	state.AddBuilding(Building{Type: "warehouse", Location: "store", Level: 1})
	if err := state.AddResource(Resource{Type: economy.ResourceGrain, Quantity: 200, Quality: 1.0}); err != nil {
		t.Fatalf("AddResource failed: %v", err)
	}
	// capacity base 100 + warehouse 50 = 150
	events := econ.Update(state.Calendar.Week, state, state.RNG.Rand())
	// food should be reduced to 150
	finalFood := 0
	for _, r := range state.Resources {
		if r.Type == economy.ResourceGrain {
			finalFood = r.Quantity
		}
	}
	if finalFood > 150 {
		t.Errorf("food exceeds capacity: %d", finalFood)
	}
	// should generate storage event
	foundStorage := false
	for _, ev := range events {
		if ev.Type == "storage" {
			foundStorage = true
		}
	}
	if !foundStorage {
		t.Error("no storage event generated")
	}
}

func TestEconomyWealthAccumulation(t *testing.T) {
	econ := NewEconomicSystem()
	state := NewGameState("test", 222)
	// Add resident with skill and assign to building
	resident := Resident{
		ID:     "r1",
		Name:   "Worker",
		Age:    30,
		Skills: []Skill{{ID: "farming", Name: "Farming", Level: 3}},
		Needs:  []Need{{ID: "happiness", Name: "Happiness", Level: 0.5}},
	}
	state.AddResident(resident)
	building := Building{
		Type:     "farm",
		Location: "south",
		Level:    1,
		Workers:  []string{"r1"},
	}
	state.AddBuilding(building)
	// Run a few weeks
	for week := 1; week <= 4; week++ {
		_ = econ.Update(week, state, state.RNG.Rand())
	}
	// Check happiness need decreased (wealth makes happier)
	found := false
	for _, r := range state.Residents {
		if r.ID == "r1" {
			for _, n := range r.Needs {
				if n.ID == "happiness" && n.Level < 0.5 {
					found = true
				}
			}
		}
	}
	if !found {
		t.Error("happiness need did not decrease despite wealth accumulation")
	}
}

func TestEconomyDeterministicWealth(t *testing.T) {
	// Same seed should produce same wealth effect on happiness
	econ1 := NewEconomicSystem()
	state1 := NewGameState("test1", 333)
	resident := Resident{
		ID:     "r1",
		Name:   "Worker",
		Age:    30,
		Skills: []Skill{{ID: "mining", Name: "Mining", Level: 2}},
		Needs:  []Need{{ID: "happiness", Name: "Happiness", Level: 0.6}},
	}
	state1.AddResident(resident)
	state1.AddBuilding(Building{Type: "mine", Location: "north", Level: 1, Workers: []string{"r1"}})
	// Run 2 weeks
	_ = econ1.Update(1, state1, state1.RNG.Rand())
	_ = econ1.Update(2, state1, state1.RNG.Rand())
	happiness1 := 0.0
	for _, r := range state1.Residents {
		for _, n := range r.Needs {
			if n.ID == "happiness" {
				happiness1 = n.Level
			}
		}
	}
	// Second run with same seed
	econ2 := NewEconomicSystem()
	state2 := NewGameState("test2", 333)
	resident2 := Resident{
		ID:     "r1",
		Name:   "Worker",
		Age:    30,
		Skills: []Skill{{ID: "mining", Name: "Mining", Level: 2}},
		Needs:  []Need{{ID: "happiness", Name: "Happiness", Level: 0.6}},
	}
	state2.AddResident(resident2)
	state2.AddBuilding(Building{Type: "mine", Location: "north", Level: 1, Workers: []string{"r1"}})
	_ = econ2.Update(1, state2, state2.RNG.Rand())
	_ = econ2.Update(2, state2, state2.RNG.Rand())
	happiness2 := 0.0
	for _, r := range state2.Residents {
		for _, n := range r.Needs {
			if n.ID == "happiness" {
				happiness2 = n.Level
			}
		}
	}
	if happiness1 != happiness2 {
		t.Errorf("wealth effect not deterministic: happiness %f vs %f", happiness1, happiness2)
	}
}

func TestEconomyWealthEdgeCases(t *testing.T) {
	// Test 1: resident without happiness need
	econ := NewEconomicSystem()
	state := NewGameState("test-no-happiness", 444)
	resident := Resident{
		ID:    "r1",
		Name:  "NoHappiness",
		Age:   30,
		Needs: []Need{{ID: "hunger", Name: "Hunger", Level: 0.5}}, // no happiness
	}
	state.AddResident(resident)
	// Should not panic
	_ = econ.Update(1, state, state.RNG.Rand())
	// Test 2: wealth effect clamping
	econ2 := NewEconomicSystem()
	state2 := NewGameState("test-clamp", 555)
	resident2 := Resident{
		ID:     "r2",
		Name:   "Worker",
		Age:    30,
		Skills: []Skill{{ID: "crafting", Name: "Crafting", Level: 10}},
		Needs:  []Need{{ID: "happiness", Name: "Happiness", Level: 0.5}},
	}
	state2.AddResident(resident2)
	state2.AddBuilding(Building{Type: "workshop", Location: "west", Level: 1, Workers: []string{"r2"}})
	// Run many weeks to accumulate wealth >20
	for week := 1; week <= 30; week++ {
		_ = econ2.Update(week, state2, state2.RNG.Rand())
	}
	// Check happiness need decreased but not below 0
	for _, r := range state2.Residents {
		for _, n := range r.Needs {
			if n.ID == "happiness" {
				if n.Level < 0 {
					t.Errorf("happiness need below 0: %f", n.Level)
				}
				if n.Level > 0.5 {
					t.Errorf("happiness need increased (should decrease): %f", n.Level)
				}

			}
		}
	}
}
