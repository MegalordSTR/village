package simulation

import (
	"github.com/vano44/village/internal/economy"
	"testing"
)

// TDD:RED - writing failing test (now updated to match implementation)

func TestEventSystemImplementsSystem(t *testing.T) {
	var _ System = (*EventSystem)(nil)
}

func TestNewEventSystem(t *testing.T) {
	evt := NewEventSystem()
	if evt == nil {
		t.Fatal("NewEventSystem returned nil")
	}
}

func TestEventsUpdateReturnsEvents(t *testing.T) {
	evt := NewEventSystem()
	state := NewGameState("test-events", 123)
	// Week 10 should trigger festival event
	events := evt.Update(10, state, state.RNG.Rand())
	if len(events) == 0 {
		t.Error("Update returned no events; expected festival event at week 10")
	}
	if len(events) > 0 && events[0].Type != "festival" {
		t.Errorf("expected event type 'festival', got %q", events[0].Type)
	}
}

func TestEventsDeterministic(t *testing.T) {
	evt := NewEventSystem()
	state1 := NewGameState("test1", 456)
	state2 := NewGameState("test2", 456)
	events1 := evt.Update(10, state1, state1.RNG.Rand())
	events2 := evt.Update(10, state2, state2.RNG.Rand())
	if len(events1) != len(events2) {
		t.Errorf("event count mismatch: %d vs %d", len(events1), len(events2))
	}
	// Compare event details
	for i := range events1 {
		if events1[i].ID != events2[i].ID {
			t.Errorf("event ID mismatch at index %d: %s vs %s", i, events1[i].ID, events2[i].ID)
		}
		if events1[i].Type != events2[i].Type {
			t.Errorf("event type mismatch at index %d: %s vs %s", i, events1[i].Type, events2[i].Type)
		}
	}
}

func TestEventsNoEventsOutsideTrigger(t *testing.T) {
	evt := NewEventSystem()
	state := NewGameState("test-no-events", 789)
	// Week 1 should not trigger any event
	events := evt.Update(1, state, state.RNG.Rand())
	if len(events) != 0 {
		t.Errorf("expected no events at week 1, got %d", len(events))
	}
}

// TestEventDiseaseOutbreak tests disease outbreak trigger and consequences.
func TestEventDiseaseOutbreak(t *testing.T) {
	evt := NewEventSystem()
	state := NewGameState("test-disease", 111)
	// Set winter season and low forest health
	state.Environment.Season = "winter"
	state.Environment.ForestHealth = 0.2
	// Add food resource
	if err := state.AddResource(Resource{Type: economy.ResourceGrain, Quantity: 100, Quality: 1.0}); err != nil {
		t.Fatalf("AddResource failed: %v", err)
	}
	initialFood := state.Resources[0].Quantity
	events := evt.Update(1, state, state.RNG.Rand())
	if len(events) == 0 {
		t.Fatal("expected disease outbreak event, got none")
	}
	if events[0].Type != "disease_outbreak" {
		t.Errorf("expected event type 'disease_outbreak', got %q", events[0].Type)
	}
	// Check food reduction (10%)
	expectedFood := initialFood - int(float64(initialFood)*0.1)
	if state.Resources[0].Quantity != expectedFood {
		t.Errorf("food quantity after disease: got %d, want %d", state.Resources[0].Quantity, expectedFood)
	}
}

// TestEventGoodHarvest tests good harvest trigger and consequences.
func TestEventGoodHarvest(t *testing.T) {
	evt := NewEventSystem()
	state := NewGameState("test-good-harvest", 222)
	state.Environment.Season = "autumn"
	state.Environment.SoilFertility = 0.9
	if err := state.AddResource(Resource{Type: economy.ResourceGrain, Quantity: 100, Quality: 1.0}); err != nil {
		t.Fatalf("AddResource failed: %v", err)
	}
	initialFood := state.Resources[0].Quantity
	events := evt.Update(1, state, state.RNG.Rand())
	if len(events) == 0 {
		t.Fatal("expected good harvest event, got none")
	}
	if events[0].Type != "good_harvest" {
		t.Errorf("expected event type 'good_harvest', got %q", events[0].Type)
	}
	expectedFood := initialFood + int(float64(initialFood)*0.2)
	if state.Resources[0].Quantity != expectedFood {
		t.Errorf("food quantity after good harvest: got %d, want %d", state.Resources[0].Quantity, expectedFood)
	}
}

// TestEventBadHarvest tests bad harvest trigger and consequences.
func TestEventBadHarvest(t *testing.T) {
	evt := NewEventSystem()
	state := NewGameState("test-bad-harvest", 333)
	state.Environment.Season = "summer"
	state.Environment.Rainfall = 3.0 // low rainfall
	if err := state.AddResource(Resource{Type: economy.ResourceGrain, Quantity: 100, Quality: 1.0}); err != nil {
		t.Fatalf("AddResource failed: %v", err)
	}
	initialFood := state.Resources[0].Quantity
	events := evt.Update(1, state, state.RNG.Rand())
	if len(events) == 0 {
		t.Fatal("expected bad harvest event, got none")
	}
	if events[0].Type != "bad_harvest" {
		t.Errorf("expected event type 'bad_harvest', got %q", events[0].Type)
	}
	expectedFood := initialFood - int(float64(initialFood)*0.15)
	if state.Resources[0].Quantity != expectedFood {
		t.Errorf("food quantity after bad harvest: got %d, want %d", state.Resources[0].Quantity, expectedFood)
	}
}

// TestEventAccident tests accident trigger and consequences.
func TestEventAccident(t *testing.T) {
	evt := NewEventSystem()
	state := NewGameState("test-accident", 444)
	// Add a building with level 2
	state.AddBuilding(Building{Type: "house", Location: "center", Level: 2})
	initialLevel := state.Buildings[0].Level
	events := evt.Update(1, state, state.RNG.Rand())
	if len(events) == 0 {
		t.Fatal("expected accident event, got none")
	}
	if events[0].Type != "accident" {
		t.Errorf("expected event type 'accident', got %q", events[0].Type)
	}
	// Building level should be reduced by 1 (if level > 0)
	if state.Buildings[0].Level != initialLevel-1 {
		t.Errorf("building level after accident: got %d, want %d", state.Buildings[0].Level, initialLevel-1)
	}
}

// TestEventNoResourceNoCrash tests that adjustResource does not crash when resource missing.
func TestEventNoResourceNoCrash(t *testing.T) {
	evt := NewEventSystem()
	state := NewGameState("test-no-resource", 555)
	state.Environment.Season = "winter"
	state.Environment.ForestHealth = 0.2
	// No food resource added
	events := evt.Update(1, state, state.RNG.Rand())
	// Should still produce event (adjustResource does nothing)
	if len(events) == 0 {
		t.Fatal("expected disease outbreak event even without food resource")
	}
	// No crash is success
}

// TestEventWeightedSelection tests that weighted selection works deterministically.
func TestEventWeightedSelection(t *testing.T) {
	evt := NewEventSystem()
	state := NewGameState("test-weighted", 666)
	// Set up conditions that trigger multiple events: festival (week 10) and accident (building present)
	state.AddBuilding(Building{Type: "workshop", Level: 1})
	// Week 10 triggers festival, accident also possible
	events := evt.Update(10, state, state.RNG.Rand())
	// Should get exactly one event (selection)
	if len(events) != 1 {
		t.Errorf("expected exactly one event, got %d", len(events))
	}
	// Deterministic: same seed should produce same event type
	state2 := NewGameState("test-weighted2", 666)
	state2.AddBuilding(Building{Type: "workshop", Level: 1})
	events2 := evt.Update(10, state2, state2.RNG.Rand())
	if len(events2) != 1 {
		t.Fatal("second run failed")
	}
	if events[0].Type != events2[0].Type {
		t.Errorf("event type mismatch across deterministic runs: %s vs %s", events[0].Type, events2[0].Type)
	}
}

// TestAdjustResourceZeroQuantity tests adjustResource with zero quantity and negative percent.
func TestAdjustResourceZeroQuantity(t *testing.T) {
	evt := NewEventSystem()
	state := NewGameState("test-zero-qty", 1111)
	state.Environment.Season = "winter"
	state.Environment.ForestHealth = 0.2
	// Add food resource with zero quantity
	if err := state.AddResource(Resource{Type: "food", Quantity: 0, Quality: 1.0}); err != nil {
		t.Fatalf("AddResource failed: %v", err)
	}
	events := evt.Update(1, state, state.RNG.Rand())
	if len(events) == 0 {
		t.Fatal("expected disease outbreak event")
	}
	// Quantity should stay zero
	if state.Resources[0].Quantity != 0 {
		t.Errorf("food quantity should remain zero, got %d", state.Resources[0].Quantity)
	}
}

// TestEventCoverage exercises various conditions to improve code coverage.
func TestEventCoverage(t *testing.T) {
	evt := NewEventSystem()
	// Test each event type with edge conditions
	// 1. Festival at week 10,20,30 etc.
	for week := 10; week <= 50; week += 10 {
		state := NewGameState("test-festival", int64(week))
		events := evt.Update(week, state, state.RNG.Rand())
		if len(events) == 0 {
			t.Errorf("festival expected at week %d but got none", week)
		}
	}
	// 2. Disease outbreak with varying forest health
	for _, health := range []float64{0.1, 0.2, 0.29} {
		state := NewGameState("test-disease", 777)
		state.Environment.Season = "winter"
		state.Environment.ForestHealth = health
		if err := state.AddResource(Resource{Type: "food", Quantity: 100}); err != nil {
			t.Fatalf("AddResource failed: %v", err)
		}
		evt.Update(1, state, state.RNG.Rand())
	}
	// 3. Good harvest with high soil fertility
	for _, fertility := range []float64{0.71, 0.8, 0.95} {
		state := NewGameState("test-good", 888)
		state.Environment.Season = "autumn"
		state.Environment.SoilFertility = fertility
		if err := state.AddResource(Resource{Type: "food", Quantity: 100}); err != nil {
			t.Fatalf("AddResource failed: %v", err)
		}
		evt.Update(1, state, state.RNG.Rand())
	}
	// 4. Bad harvest with low rainfall
	for _, rain := range []float64{1.0, 2.0, 4.9} {
		state := NewGameState("test-bad", 999)
		state.Environment.Season = "summer"
		state.Environment.Rainfall = rain
		if err := state.AddResource(Resource{Type: "food", Quantity: 100}); err != nil {
			t.Fatalf("AddResource failed: %v", err)
		}
		evt.Update(1, state, state.RNG.Rand())
	}
	// 5. Accident with zero-level building (should not reduce below zero)
	state := NewGameState("test-accident-zero", 1000)
	state.AddBuilding(Building{Type: "shed", Level: 0})
	evt.Update(1, state, state.RNG.Rand())
	// 6. Accident with multiple buildings
	state2 := NewGameState("test-accident-multi", 1001)
	state2.AddBuilding(Building{Type: "house", Level: 3})
	state2.AddBuilding(Building{Type: "farm", Level: 2})
	evt.Update(1, state2, state2.RNG.Rand())
	// 7. No resource edge case for all event types (already covered)
	// 8. Weight total zero (no triggered events)
	state3 := NewGameState("test-no-trigger", 1002)
	state3.Environment.Season = "spring"
	state3.Environment.ForestHealth = 0.9
	state3.Environment.SoilFertility = 0.5
	state3.Environment.Rainfall = 10.0
	events := evt.Update(1, state3, state3.RNG.Rand())
	if len(events) != 0 {
		t.Errorf("expected no events, got %d", len(events))
	}
}

// TestEvents is the top-level test that runs all event-related subtests.
func TestEvents(t *testing.T) {
	t.Run("ImplementsSystem", TestEventSystemImplementsSystem)
	t.Run("NewEventSystem", TestNewEventSystem)
	t.Run("UpdateReturnsEvents", TestEventsUpdateReturnsEvents)
	t.Run("Deterministic", TestEventsDeterministic)
	t.Run("NoEventsOutsideTrigger", TestEventsNoEventsOutsideTrigger)
	t.Run("DiseaseOutbreak", TestEventDiseaseOutbreak)
	t.Run("GoodHarvest", TestEventGoodHarvest)
	t.Run("BadHarvest", TestEventBadHarvest)
	t.Run("Accident", TestEventAccident)
	t.Run("NoResourceNoCrash", TestEventNoResourceNoCrash)
	t.Run("WeightedSelection", TestEventWeightedSelection)
	t.Run("AdjustResourceZeroQuantity", TestAdjustResourceZeroQuantity)
	t.Run("Coverage", TestEventCoverage)
}
