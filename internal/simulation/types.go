package simulation

// Resident represents a village inhabitant.
type Resident struct {
	ID            string         `json:"id"`
	Name          string         `json:"name"`
	Age           int            `json:"age"`
	Skills        []Skill        `json:"skills"`
	Needs         []Need         `json:"needs"`
	Relationships []Relationship `json:"relationships"`
}

// Building represents a constructed building in the village.
type Building struct {
	Type       string                 `json:"type"`
	Location   string                 `json:"location"`
	Level      int                    `json:"level"`
	Workers    []string               `json:"workers"` // resident IDs
	Production []Production           `json:"production"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}

// Resource represents a material resource.
type Resource struct {
	Type     string  `json:"type"`
	Quantity int     `json:"quantity"`
	Quality  float64 `json:"quality"`
}

// Skill represents a resident's ability.
type Skill struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Level int    `json:"level"`
}

// Need represents a resident's need.
type Need struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Level float64 `json:"level"` // 0.0 to 1.0
}

// Relationship represents a social connection between residents.
type Relationship struct {
	From     string  `json:"from"`
	To       string  `json:"to"`
	Type     string  `json:"type"`
	Strength float64 `json:"strength"`
}

// Production represents a building's output.
type Production struct {
	ResourceType string  `json:"resource_type"`
	Amount       int     `json:"amount"`
	Quality      float64 `json:"quality"`
}

// Calendar tracks game time.
type Calendar struct {
	Year  int `json:"year"`
	Month int `json:"month"`
	Week  int `json:"week"`
	Day   int `json:"day"`
}

// Village represents the village location and map.
type Village struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Region  string `json:"region"`
	Climate string `json:"climate"`
}

// Event represents a historical event.
type Event struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Timestamp string                 `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}

// Policy represents a village policy.
type Policy struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Active      bool   `json:"active"`
}

// Environment represents environmental state.
type Environment struct {
	Temperature float64 `json:"temperature"` // in Celsius
	Rainfall    float64 `json:"rainfall"`    // mm per week
	Season      string  `json:"season"`      // "spring", "summer", "autumn", "winter"
	// Soil fertility per region (simplified: single value for now)
	SoilFertility float64 `json:"soil_fertility"` // 0.0 to 1.0
	// Natural resources
	ForestHealth       float64 `json:"forest_health"`       // 0.0 to 1.0
	MineQuality        float64 `json:"mine_quality"`        // 0.0 to 1.0
	WildlifePopulation float64 `json:"wildlife_population"` // 0.0 to 1.0
}
