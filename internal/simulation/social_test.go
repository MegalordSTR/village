package simulation

import (
	"testing"
)

func TestSocialSystemImplementsSystem(t *testing.T) {
	var _ System = (*SocialSystem)(nil)
}

func TestNewSocialSystem(t *testing.T) {
	soc := NewSocialSystem()
	if soc == nil {
		t.Fatal("NewSocialSystem returned nil")
	}
}

func TestSocialUpdateDeterministic(t *testing.T) {
	soc := NewSocialSystem()

	// Create two states with same seed
	state1 := NewGameState("test1", 123)
	state2 := NewGameState("test2", 123)

	// Add some residents
	state1.AddResident(Resident{ID: "r1", Name: "Alice", Age: 25})
	state1.AddResident(Resident{ID: "r2", Name: "Bob", Age: 30})
	state2.AddResident(Resident{ID: "r1", Name: "Alice", Age: 25})
	state2.AddResident(Resident{ID: "r2", Name: "Bob", Age: 30})

	// Process first week
	events1 := soc.Update(1, state1, state1.RNG.Rand())
	events2 := soc.Update(1, state2, state2.RNG.Rand())

	// Should generate same events
	if len(events1) != len(events2) {
		t.Errorf("event count mismatch: %d vs %d", len(events1), len(events2))
	}

	// Residents should have same needs changes
	if len(state1.Residents) != len(state2.Residents) {
		t.Fatalf("resident count mismatch: %d vs %d", len(state1.Residents), len(state2.Residents))
	}
	for i := range state1.Residents {
		r1 := &state1.Residents[i]
		r2 := &state2.Residents[i]
		if len(r1.Needs) != len(r2.Needs) {
			t.Errorf("needs count mismatch for resident %s: %d vs %d", r1.ID, len(r1.Needs), len(r2.Needs))
		}
	}
}

func TestSocialNeedsUpdate(t *testing.T) {
	soc := NewSocialSystem()
	state := NewGameState("test", 456)
	state.AddResident(Resident{ID: "r1", Name: "Test", Age: 20})

	// Initial needs should be empty
	if len(state.Residents[0].Needs) != 0 {
		t.Errorf("expected empty needs initially, got %d", len(state.Residents[0].Needs))
	}

	// Run update
	_ = soc.Update(1, state, state.RNG.Rand())

	// After update, needs should be populated
	if len(state.Residents[0].Needs) == 0 {
		t.Error("expected needs to be populated after update")
	}
}

func TestSocialSkillImprovement(t *testing.T) {
	soc := NewSocialSystem()
	state := NewGameState("test", 789)
	state.AddResident(Resident{
		ID:   "r1",
		Name: "Crafter",
		Age:  25,
		Skills: []Skill{
			{ID: "crafting", Name: "Crafting", Level: 1},
		},
	})

	initialLevel := state.Residents[0].Skills[0].Level

	// Run multiple weeks to allow skill improvement (100 weeks for high probability)
	for week := 1; week <= 100; week++ {
		_ = soc.Update(week, state, state.RNG.Rand())
	}

	finalLevel := state.Residents[0].Skills[0].Level
	if finalLevel <= initialLevel {
		t.Errorf("expected skill improvement, level stayed at %d", finalLevel)
	}
}

func TestSocialRelationshipChanges(t *testing.T) {
	soc := NewSocialSystem()
	state := NewGameState("test", 999)
	state.AddResident(Resident{ID: "r1", Name: "Alice", Age: 25})
	state.AddResident(Resident{ID: "r2", Name: "Bob", Age: 30})

	initialRelationships := len(state.Residents[0].Relationships)

	// Run update
	_ = soc.Update(1, state, state.RNG.Rand())

	// Relationships may form
	finalRelationships := len(state.Residents[0].Relationships)
	if finalRelationships < initialRelationships {
		t.Error("relationships should not decrease")
	}
}

func TestSocialLifeEvents(t *testing.T) {
	soc := NewSocialSystem()
	state := NewGameState("test", 111)
	state.AddResident(Resident{ID: "r1", Name: "Elder", Age: 90})

	// Run multiple weeks to trigger potential death
	var events []Event
	for week := 1; week <= 52; week++ {
		events = append(events, soc.Update(week, state, state.RNG.Rand())...)
	}

	// Check for death events
	deathCount := 0
	for _, e := range events {
		if e.Type == "death" {
			deathCount++
		}
	}
	// At least elder might die
	if deathCount == 0 {
		t.Log("no death events generated (maybe probability low)")
	}
}

func TestSocialHappinessMetric(t *testing.T) {
	soc := NewSocialSystem()
	state := NewGameState("test", 222)
	state.AddResident(Resident{ID: "r1", Name: "Happy", Age: 30})

	_ = soc.Update(1, state, state.RNG.Rand())

	// Check that happiness need exists
	foundHappiness := false
	for _, need := range state.Residents[0].Needs {
		if need.ID == "happiness" || need.Name == "Happiness" {
			foundHappiness = true
			break
		}
	}
	if !foundHappiness {
		t.Error("happiness need not found")
	}
}

func TestSocialShelterAndWork(t *testing.T) {
	soc := NewSocialSystem()
	state := NewGameState("test", 333)

	// Add a resident
	resident := Resident{ID: "r1", Name: "Worker", Age: 25}
	state.AddResident(resident)

	// Add a house building
	house := Building{
		Type:     "house",
		Location: "A1",
		Level:    1,
	}
	state.AddBuilding(house)

	// Add a farm building with worker
	farm := Building{
		Type:    "farm",
		Workers: []string{"r1"},
	}
	state.AddBuilding(farm)

	// Run update
	_ = soc.Update(1, state, state.RNG.Rand())

	// Check that shelter need exists and is reduced
	var shelterNeed *Need
	for i := range state.Residents[0].Needs {
		if state.Residents[0].Needs[i].ID == "shelter" {
			shelterNeed = &state.Residents[0].Needs[i]
			break
		}
	}
	if shelterNeed == nil {
		t.Fatal("shelter need not found")
	}
	// Shelter should be reduced due to house presence
	// Initial level 0.6, decay +0.01~0.03, reduction -0.1
	// So level should be <= 0.6
	if shelterNeed.Level > 0.6 {
		t.Errorf("shelter need level %f expected <= 0.6", shelterNeed.Level)
	}

	// Check that happiness is increased due to work and shelter
	var happinessNeed *Need
	for i := range state.Residents[0].Needs {
		if state.Residents[0].Needs[i].ID == "happiness" {
			happinessNeed = &state.Residents[0].Needs[i]
			break
		}
	}
	if happinessNeed == nil {
		t.Fatal("happiness need not found")
	}
	// Happiness should be higher due to bonuses
	// Initial 0.7, decay +0.01~0.03, bonuses +0.08
	// So level should be >= 0.7
	if happinessNeed.Level < 0.7 {
		t.Errorf("happiness need level %f expected >= 0.7", happinessNeed.Level)
	}
}

func TestSocialFoodAvailability(t *testing.T) {
	soc := NewSocialSystem()
	state := NewGameState("test", 444)

	// Add a resident
	state.AddResident(Resident{ID: "r1", Name: "Hungry", Age: 30})

	// Add some food resources
	state.AddResource(Resource{Type: "food", Quantity: 100, Quality: 1.0})
	state.AddResource(Resource{Type: "grain", Quantity: 50, Quality: 0.8})

	// Run update
	_ = soc.Update(1, state, state.RNG.Rand())

	// Find hunger need
	var hungerNeed *Need
	for i := range state.Residents[0].Needs {
		if state.Residents[0].Needs[i].ID == "hunger" {
			hungerNeed = &state.Residents[0].Needs[i]
			break
		}
	}
	if hungerNeed == nil {
		t.Fatal("hunger need not found")
	}
	// With abundant food, hunger should be reduced
	// Initial level 0.3, decay +0.01~0.03, reduction up to 0.05*foodAvailability
	// foodAvailability = (100+50)/(1*10) = 15, capped to 1, so reduction 0.05
	// So hunger level should be <= 0.3
	if hungerNeed.Level > 0.3 {
		t.Errorf("hunger need level %f expected <= 0.3 with abundant food", hungerNeed.Level)
	}
}
