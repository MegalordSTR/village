package simulation

import (
	"github.com/vano44/village/internal/economy"
)

// GameState represents the entire simulation state.
type GameState struct {
	ID          string             `json:"id"`
	Version     int                `json:"version"`
	Seed        int64              `json:"seed"`
	RNG         *RNG               `json:"rng"`
	Calendar    Calendar           `json:"calendar"`
	Village     Village            `json:"village"`
	Residents   []Resident         `json:"residents"`
	Resources   []Resource         `json:"resources"`
	Buildings   []Building         `json:"buildings"`
	History     []Event            `json:"history"`
	Policies    []Policy           `json:"policies"`
	Environment Environment        `json:"environment"`
	Inventory   *economy.Inventory `json:"-"`
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
		Resources: []Resource{},
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
	}
}

// AddResident adds a resident to the state.
func (gs *GameState) AddResident(r Resident) {
	gs.Residents = append(gs.Residents, r)
}

// AddResource adds a resource to the state.
func (gs *GameState) AddResource(r Resource) {
	gs.Resources = append(gs.Resources, r)
	// If inventory exists, add the resource there as well (at default "global" location)
	if gs.Inventory != nil {
		er := ToEconomyResource(r)
		er.Location = "global"
		// Ignore error for now; in production we might want to handle it
		_ = gs.Inventory.AddResource("global", er)
	}
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

// SyncInventory ensures Inventory is populated from Resources.
// If Inventory is nil, a new Inventory is created and resources are loaded.
// If StorageRegistry is needed, it must be attached separately.
func (gs *GameState) SyncInventory() {
	if gs.Inventory == nil {
		gs.Inventory = economy.NewInventory()
		// Load existing resources into inventory at default location "global"
		LoadInventoryFromGameState(gs.Inventory, gs.Resources, "global")
	}
}

// SyncResources updates the Resources slice from Inventory.
// This should be called after inventory modifications to keep Resources in sync.
func (gs *GameState) SyncResources() {
	if gs.Inventory == nil {
		return
	}
	gs.Resources = ExportInventoryToGameState(gs.Inventory)
}
