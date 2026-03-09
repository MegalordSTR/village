package simulation

import (
	"encoding/json"
	"testing"
)

func TestGameStateStruct(t *testing.T) {
	// This test will fail until GameState is defined
	gs := GameState{
		ID:        "test-id",
		Version:   1,
		Seed:      12345,
		Calendar:  Calendar{},
		Village:   Village{},
		Residents: []Resident{},
		Resources: []Resource{},
		Buildings: []Building{},
		History:   []Event{},
		Policies:  []Policy{},
	}

	// Check field values
	if gs.ID != "test-id" {
		t.Errorf("expected ID 'test-id', got %q", gs.ID)
	}
	if gs.Version != 1 {
		t.Errorf("expected Version 1, got %d", gs.Version)
	}
	if gs.Seed != 12345 {
		t.Errorf("expected Seed 12345, got %d", gs.Seed)
	}

	// Marshal to JSON to verify tags
	data, err := json.Marshal(gs)
	if err != nil {
		t.Fatalf("failed to marshal GameState: %v", err)
	}

	// Parse back to map to verify JSON field names
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}

	expectedFields := []string{"id", "version", "seed", "calendar", "village", "residents", "resources", "buildings", "history", "policies"}
	for _, field := range expectedFields {
		if _, ok := m[field]; !ok {
			t.Errorf("missing JSON field %q", field)
		}
	}
}

func TestResidentStruct(t *testing.T) {
	r := Resident{
		ID:            "res-1",
		Name:          "Alice",
		Age:           30,
		Skills:        []Skill{},
		Needs:         []Need{},
		Relationships: []Relationship{},
	}

	if r.ID != "res-1" {
		t.Errorf("expected ID 'res-1', got %q", r.ID)
	}
	if r.Name != "Alice" {
		t.Errorf("expected Name 'Alice', got %q", r.Name)
	}
	if r.Age != 30 {
		t.Errorf("expected Age 30, got %d", r.Age)
	}

	data, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("failed to marshal Resident: %v", err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}
	expectedFields := []string{"id", "name", "age", "skills", "needs", "relationships"}
	for _, field := range expectedFields {
		if _, ok := m[field]; !ok {
			t.Errorf("missing JSON field %q", field)
		}
	}
}

func TestBuildingStruct(t *testing.T) {
	b := Building{
		Type:       "house",
		Location:   "A1",
		Level:      1,
		Workers:    []string{},
		Production: []Production{},
	}

	if b.Type != "house" {
		t.Errorf("expected Type 'house', got %q", b.Type)
	}
	if b.Location != "A1" {
		t.Errorf("expected Location 'A1', got %q", b.Location)
	}
	if b.Level != 1 {
		t.Errorf("expected Level 1, got %d", b.Level)
	}

	data, err := json.Marshal(b)
	if err != nil {
		t.Fatalf("failed to marshal Building: %v", err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}
	expectedFields := []string{"type", "location", "level", "workers", "production"}
	for _, field := range expectedFields {
		if _, ok := m[field]; !ok {
			t.Errorf("missing JSON field %q", field)
		}
	}
}

func TestResourceStruct(t *testing.T) {
	r := Resource{
		Type:     "wood",
		Quantity: 100,
		Quality:  1.0,
	}

	if r.Type != "wood" {
		t.Errorf("expected Type 'wood', got %q", r.Type)
	}
	if r.Quantity != 100 {
		t.Errorf("expected Quantity 100, got %d", r.Quantity)
	}
	if r.Quality != 1.0 {
		t.Errorf("expected Quality 1.0, got %f", r.Quality)
	}

	data, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("failed to marshal Resource: %v", err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}
	expectedFields := []string{"type", "quantity", "quality"}
	for _, field := range expectedFields {
		if _, ok := m[field]; !ok {
			t.Errorf("missing JSON field %q", field)
		}
	}
}

func TestNewGameState(t *testing.T) {
	gs := NewGameState("test", 42)
	if gs.ID != "test" {
		t.Errorf("expected ID 'test', got %q", gs.ID)
	}
	if gs.Version != 1 {
		t.Errorf("expected Version 1, got %d", gs.Version)
	}
	if gs.Seed != 42 {
		t.Errorf("expected Seed 42, got %d", gs.Seed)
	}
	if gs.Calendar.Year != 1 {
		t.Errorf("expected Calendar.Year 1, got %d", gs.Calendar.Year)
	}
	if gs.Village.Name != "New Village" {
		t.Errorf("expected Village.Name 'New Village', got %q", gs.Village.Name)
	}
	if len(gs.Residents) != 0 {
		t.Errorf("expected empty Residents, got %d", len(gs.Residents))
	}
	if len(gs.Resources) != 0 {
		t.Errorf("expected empty Resources, got %d", len(gs.Resources))
	}
	if len(gs.Buildings) != 0 {
		t.Errorf("expected empty Buildings, got %d", len(gs.Buildings))
	}
	if len(gs.History) != 0 {
		t.Errorf("expected empty History, got %d", len(gs.History))
	}
	if len(gs.Policies) != 0 {
		t.Errorf("expected empty Policies, got %d", len(gs.Policies))
	}
}

func TestAddResident(t *testing.T) {
	gs := NewGameState("test", 42)
	r := Resident{ID: "r1", Name: "Bob", Age: 25}
	gs.AddResident(r)
	if len(gs.Residents) != 1 {
		t.Fatalf("expected 1 resident, got %d", len(gs.Residents))
	}
	if gs.Residents[0].ID != "r1" {
		t.Errorf("expected resident ID 'r1', got %q", gs.Residents[0].ID)
	}
}

func TestAddResource(t *testing.T) {
	gs := NewGameState("test", 42)
	res := Resource{Type: "wood", Quantity: 10, Quality: 0.8}
	if err := gs.AddResource(res); err != nil {
		t.Fatalf("AddResource failed: %v", err)
	}
	if len(gs.Resources) != 1 {
		t.Fatalf("expected 1 resource, got %d", len(gs.Resources))
	}
	if gs.Resources[0].Type != "wood" {
		t.Errorf("expected resource type 'wood', got %q", gs.Resources[0].Type)
	}
}

func TestAddBuilding(t *testing.T) {
	gs := NewGameState("test", 42)
	b := Building{Type: "house", Location: "A1", Level: 1}
	gs.AddBuilding(b)
	if len(gs.Buildings) != 1 {
		t.Fatalf("expected 1 building, got %d", len(gs.Buildings))
	}
	if gs.Buildings[0].Type != "house" {
		t.Errorf("expected building type 'house', got %q", gs.Buildings[0].Type)
	}
}

func TestAddEvent(t *testing.T) {
	gs := NewGameState("test", 42)
	e := Event{ID: "e1", Type: "festival"}
	gs.AddEvent(e)
	if len(gs.History) != 1 {
		t.Fatalf("expected 1 event, got %d", len(gs.History))
	}
	if gs.History[0].ID != "e1" {
		t.Errorf("expected event ID 'e1', got %q", gs.History[0].ID)
	}
}

func TestAddPolicy(t *testing.T) {
	gs := NewGameState("test", 42)
	p := Policy{ID: "p1", Name: "Tax", Active: true}
	gs.AddPolicy(p)
	if len(gs.Policies) != 1 {
		t.Fatalf("expected 1 policy, got %d", len(gs.Policies))
	}
	if gs.Policies[0].ID != "p1" {
		t.Errorf("expected policy ID 'p1', got %q", gs.Policies[0].ID)
	}
}

func TestGameStateJSON(t *testing.T) {
	gs := NewGameState("json-test", 99)
	gs.AddResident(Resident{ID: "r1", Name: "Charlie"})
	if err := gs.AddResource(Resource{Type: "stone", Quantity: 5}); err != nil {
		t.Fatalf("AddResource failed: %v", err)
	}
	gs.AddBuilding(Building{Type: "market", Location: "B2"})

	data, err := json.Marshal(gs)
	if err != nil {
		t.Fatalf("failed to marshal GameState: %v", err)
	}

	var decoded GameState
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal GameState: %v", err)
	}
	if decoded.ID != gs.ID {
		t.Errorf("decoded ID mismatch: got %q, want %q", decoded.ID, gs.ID)
	}
	if len(decoded.Residents) != 1 {
		t.Errorf("decoded Residents length mismatch: got %d, want 1", len(decoded.Residents))
	}
}

func TestCore(t *testing.T) {
	c := NewCore()
	if c == nil {
		t.Fatal("NewCore returned nil")
	}
	err := c.ProcessWeek()
	if err != nil {
		t.Errorf("ProcessWeek returned error: %v", err)
	}
}

func TestAddResource_Invalid(t *testing.T) {
	gs := NewGameState("test", 42)
	// Invalid resource type
	err := gs.AddResource(Resource{Type: "unknown", Quantity: 10})
	if err == nil {
		t.Error("expected error for unknown resource type")
	}
	// Negative quantity
	err = gs.AddResource(Resource{Type: "wood", Quantity: -5})
	if err == nil {
		t.Error("expected error for negative quantity")
	}
	// With inventory present (should also propagate error)
	gs.SyncInventory()
	err = gs.AddResource(Resource{Type: "unknown", Quantity: 10})
	if err == nil {
		t.Error("expected error for unknown resource type with inventory")
	}
}
