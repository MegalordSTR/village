package economy

// Mine represents a mineral deposit that can be depleted.
type Mine struct {
	TotalOre      float64 // total ore originally present
	RemainingOre  float64 // current remaining ore
	DepletionRate float64 // fraction of extracted ore that is lost to depletion (0-1)
}

// Deplete removes the given given amount of ore from the mine.
// The amount removed reduces RemainingOre, but cannot go below zero.
func (m *Mine) Deplete(amount float64) {
	if amount > m.RemainingOre {
		amount = m.RemainingOre
	}
	m.RemainingOre -= amount
}

// YieldModifier returns a multiplier for production yield based on remaining ore.
// Rich mine (full) yields 1.5x, depleted mine yields 0.3x, scaling linearly.
func (m *Mine) YieldModifier() float64 {
	if m.TotalOre <= 0 {
		return 0.3
	}
	fraction := m.RemainingOre / m.TotalOre
	return 0.3 + 1.2*fraction
}

// Forest represents a forest area that can be logged and regrow.
type Forest struct {
	Health       float64 // current health (0 = dead, 1 = fully healthy)
	RegrowthRate float64 // fraction of health regained per year (0-1)
}

// Regrow increases forest health by the regrowth rate per year.
// Health cannot exceed 1.0.
func (f *Forest) Regrow(years float64) {
	f.Health += f.RegrowthRate * years
	if f.Health > 1.0 {
		f.Health = 1.0
	}
}

// Log reduces forest health by the given intensity (0-1).
// Health cannot go below 0.0.
func (f *Forest) Log(intensity float64) {
	f.Health -= intensity
	if f.Health < 0.0 {
		f.Health = 0.0
	}
}
