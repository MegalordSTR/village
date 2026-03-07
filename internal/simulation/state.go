package simulation

// GameState represents the entire simulation state.
type GameState struct {
	ID        string     `json:"id"`
	Version   int        `json:"version"`
	Seed      int64      `json:"seed"`
	Calendar  Calendar   `json:"calendar"`
	Village   Village    `json:"village"`
	Residents []Resident `json:"residents"`
	Resources []Resource `json:"resources"`
	Buildings []Building `json:"buildings"`
	History   []Event    `json:"history"`
	Policies  []Policy   `json:"policies"`
}

// NewGameState creates a new GameState with default values.
func NewGameState(id string, seed int64) *GameState {
	return &GameState{
		ID:      id,
		Version: 1,
		Seed:    seed,
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
	}
}

// AddResident adds a resident to the state.
func (gs *GameState) AddResident(r Resident) {
	gs.Residents = append(gs.Residents, r)
}

// AddResource adds a resource to the state.
func (gs *GameState) AddResource(r Resource) {
	gs.Resources = append(gs.Resources, r)
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
