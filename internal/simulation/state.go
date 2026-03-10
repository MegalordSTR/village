package simulation

import (
	"encoding/json"
	"errors"
	"github.com/vano44/village/internal/economy"
	"log"
)

// GameState represents the entire simulation state.
type GameState struct {
	ID        string     `json:"id"`
	Version   int        `json:"version"`
	Seed      int64      `json:"seed"`
	RNG       *RNG       `json:"rng"`
	Calendar  Calendar   `json:"calendar"`
	Village   Village    `json:"village"`
	Residents []Resident `json:"residents"`

	Buildings   []Building         `json:"buildings"`
	History     []Event            `json:"history"`
	Policies    []Policy           `json:"policies"`
	Environment Environment        `json:"environment"`
	Inventory   *economy.Inventory `json:"inventory"`
}

// NewGameState creates a new GameState with default values.
func NewGameState(id string, seed int64) *GameState {
	return &GameState{
		ID:      id,
		Version: 1,
		Seed:    seed,
		RNG:     NewRNG(seed),
		Calendar: Calendar{
			Year:  1,
			Month: 1,
			Week:  1,
			Day:   1,
		},
		Village: Village{
			ID:      id + "-village",
			Name:    "New Village",
			Region:  "temperate",
			Climate: "mild",
		},
		Residents: []Resident{},
		Buildings: []Building{},
		History:   []Event{},
		Policies:  []Policy{},
		Environment: Environment{
			Temperature:        15.0, // mild temperature
			Rainfall:           10.0, // mm per week
			Season:             "spring",
			SoilFertility:      0.7,
			ForestHealth:       0.8,
			MineQuality:        0.6,
			WildlifePopulation: 0.5,
		},
		Inventory: economy.NewInventory(),
	}
}

// AddResident adds a resident to the state.
func (gs *GameState) AddResident(r Resident) {
	gs.Residents = append(gs.Residents, r)
}

// AddResource adds a resource to the state.
func (gs *GameState) AddResource(r Resource) error {
	// Reject unknown resource types (including unmapped legacy strings)
	if !IsKnownType(string(r.Type)) {
		log.Printf("WARN: operation=GameState.AddResource resource=%s quantity=%d error=\"invalid resource type\"", r.Type, r.Quantity)
		return errors.New("invalid resource type")
	}
	// Convert to economy resource and validate quantity
	er := ToEconomyResource(r)
	if !er.Validate() {
		log.Printf("WARN: operation=GameState.AddResource resource=%s quantity=%d error=\"invalid resource\"", r.Type, r.Quantity)
		return errors.New("invalid resource")
	}
	// Add to inventory (Inventory is always present after NewGameState)
	er.Location = "global"
	err := gs.Inventory.AddResource("global", er)
	if err != nil {
		// Inventory.AddResource already logs the error
		return err
	}
	log.Printf("INFO: operation=GameState.AddResource resource=%s quantity=%d", r.Type, r.Quantity)
	return nil
}

// AddBuilding adds a building to the state.
func (gs *GameState) AddBuilding(b Building) {
	gs.Buildings = append(gs.Buildings, b)
}

// AddEvent records an event in history.
func (gs *GameState) AddEvent(e Event) {
	gs.History = append(gs.History, e)
}

// AddPolicy adds a policy.
func (gs *GameState) AddPolicy(p Policy) {
	gs.Policies = append(gs.Policies, p)
}

// MarshalJSON customizes JSON serialization for GameState.
func (gs *GameState) MarshalJSON() ([]byte, error) {
	resources := ExportInventoryToGameState(gs.Inventory)
	return json.Marshal(struct {
		ID          string      `json:"id"`
		Version     int         `json:"version"`
		Seed        int64       `json:"seed"`
		RNG         *RNG        `json:"rng"`
		Calendar    Calendar    `json:"calendar"`
		Village     Village     `json:"village"`
		Residents   []Resident  `json:"residents"`
		Buildings   []Building  `json:"buildings"`
		History     []Event     `json:"history"`
		Policies    []Policy    `json:"policies"`
		Environment Environment `json:"environment"`
		Resources   []Resource  `json:"resources,omitempty"`
	}{
		ID:          gs.ID,
		Version:     gs.Version,
		Seed:        gs.Seed,
		RNG:         gs.RNG,
		Calendar:    gs.Calendar,
		Village:     gs.Village,
		Residents:   gs.Residents,
		Buildings:   gs.Buildings,
		History:     gs.History,
		Policies:    gs.Policies,
		Environment: gs.Environment,
		Resources:   resources,
	})
}

// UnmarshalJSON customizes JSON deserialization for GameState.
func (gs *GameState) UnmarshalJSON(data []byte) error {
	var aux struct {
		ID          string      `json:"id"`
		Version     int         `json:"version"`
		Seed        int64       `json:"seed"`
		RNG         *RNG        `json:"rng"`
		Calendar    Calendar    `json:"calendar"`
		Village     Village     `json:"village"`
		Residents   []Resident  `json:"residents"`
		Buildings   []Building  `json:"buildings"`
		History     []Event     `json:"history"`
		Policies    []Policy    `json:"policies"`
		Environment Environment `json:"environment"`
		Resources   []Resource  `json:"resources"`
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	gs.ID = aux.ID
	gs.Version = aux.Version
	gs.Seed = aux.Seed
	gs.RNG = aux.RNG
	gs.Calendar = aux.Calendar
	gs.Village = aux.Village
	gs.Residents = aux.Residents
	gs.Buildings = aux.Buildings
	gs.History = aux.History
	gs.Policies = aux.Policies
	gs.Environment = aux.Environment
	gs.Inventory = economy.NewInventory()
	if aux.Resources != nil {
		produced := economy.GameDate{Year: gs.Calendar.Year, Week: gs.Calendar.Week}
		if err := LoadInventoryFromGameState(gs.Inventory, aux.Resources, "global", produced); err != nil {
			return err
		}
	}
	return nil
}
