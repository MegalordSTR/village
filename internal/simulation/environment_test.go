package simulation

import (
	"math/rand"
	"testing"
)

func TestEnvironmentSystemImplementsSystem(t *testing.T) {
	var _ System = (*EnvironmentSystem)(nil)
}

func TestNewEnvironmentSystem(t *testing.T) {
	env := NewEnvironmentSystem()
	if env == nil {
		t.Fatal("NewEnvironmentSystem returned nil")
	}
}

func TestEnvironmentUpdateDeterministic(t *testing.T) {
	env := NewEnvironmentSystem()

	// Create two states with same seed
	state1 := NewGameState("test1", 123)
	state2 := NewGameState("test2", 123)

	// Process first week
	events1 := env.Update(1, state1, state1.RNG.Rand())
	events2 := env.Update(1, state2, state2.RNG.Rand())

	// Should generate same events
	if len(events1) != len(events2) {
		t.Errorf("event count mismatch: %d vs %d", len(events1), len(events2))
	}

	// Environment should change deterministically
	if state1.Environment.Temperature != state2.Environment.Temperature {
		t.Errorf("temperature mismatch: %f vs %f",
			state1.Environment.Temperature, state2.Environment.Temperature)
	}
	if state1.Environment.Rainfall != state2.Environment.Rainfall {
		t.Errorf("rainfall mismatch: %f vs %f",
			state1.Environment.Rainfall, state2.Environment.Rainfall)
	}
	if state1.Environment.Season != state2.Environment.Season {
		t.Errorf("season mismatch: %s vs %s",
			state1.Environment.Season, state2.Environment.Season)
	}

	// Season should match calculated season
	expectedSeason := calculateSeason(1)
	if state1.Environment.Season != expectedSeason {
		t.Errorf("season %s doesn't match calculated season %s",
			state1.Environment.Season, expectedSeason)
	}
}

func TestEnvironmentUpdateChangesAllFields(t *testing.T) {
	env := NewEnvironmentSystem()
	state := NewGameState("test", 456)
	initial := state.Environment

	events := env.Update(state.Calendar.Week, state, state.RNG.Rand())

	// Check all fields were updated (may stay same due to random seed, but likely change)
	// We'll check at least some fields changed
	changed := false
	if state.Environment.Temperature != initial.Temperature {
		changed = true
	}
	if state.Environment.Rainfall != initial.Rainfall {
		changed = true
	}
	if state.Environment.SoilFertility != initial.SoilFertility {
		changed = true
	}
	if state.Environment.ForestHealth != initial.ForestHealth {
		changed = true
	}
	if state.Environment.MineQuality != initial.MineQuality {
		changed = true
	}
	if state.Environment.WildlifePopulation != initial.WildlifePopulation {
		changed = true
	}
	if !changed {
		t.Error("no environment fields changed (unlikely with random seed)")
	}

	// Season should match calculated season
	expectedSeason := calculateSeason(state.Calendar.Week)
	if state.Environment.Season != expectedSeason {
		t.Errorf("season %s doesn't match calculated season %s",
			state.Environment.Season, expectedSeason)
	}

	// Should generate some events (maybe)
	_ = events
}

func TestCalculateSeason(t *testing.T) {
	tests := []struct {
		week   int
		season string
	}{
		{1, "spring"},
		{13, "spring"},
		{14, "summer"},
		{26, "summer"},
		{27, "autumn"},
		{39, "autumn"},
		{40, "winter"},
		{52, "winter"},
		{53, "winter"}, // week 53 should wrap to winter (though shouldn't happen)
		{0, "winter"},
		{-1, "winter"},
	}

	for _, tt := range tests {
		got := calculateSeason(tt.week)
		if got != tt.season {
			t.Errorf("calculateSeason(%d) = %s, want %s", tt.week, got, tt.season)
		}
	}
}

func TestCalculateWeatherDeterministic(t *testing.T) {
	// Test each season
	seasons := []string{"spring", "summer", "autumn", "winter"}
	for _, season := range seasons {
		t.Run(season, func(t *testing.T) {
			// Use same seed for both calls
			rng1 := rand.New(rand.NewSource(42))
			temp1, rain1 := calculateWeather(season, rng1)

			rng2 := rand.New(rand.NewSource(42))
			temp2, rain2 := calculateWeather(season, rng2)

			if temp1 != temp2 {
				t.Errorf("temperature not deterministic: %f vs %f", temp1, temp2)
			}
			if rain1 != rain2 {
				t.Errorf("rainfall not deterministic: %f vs %f", rain1, rain2)
			}

			// Basic sanity checks
			switch season {
			case "spring":
				if temp1 < -5 || temp1 > 25 {
					t.Errorf("spring temperature %f out of reasonable range", temp1)
				}
				if rain1 < 0 || rain1 > 30 {
					t.Errorf("spring rainfall %f out of reasonable range", rain1)
				}
			case "summer":
				if temp1 < 10 || temp1 > 35 {
					t.Errorf("summer temperature %f out of reasonable range", temp1)
				}
			case "autumn":
				if temp1 < 5 || temp1 > 25 {
					t.Errorf("autumn temperature %f out of reasonable range", temp1)
				}
			case "winter":
				if temp1 < -10 || temp1 > 15 {
					t.Errorf("winter temperature %f out of reasonable range", temp1)
				}
			}
		})
	}
}

func TestUpdateSoilFertility(t *testing.T) {
	// Test with different rainfall levels
	tests := []struct {
		name     string
		current  float64
		rainfall float64
	}{
		{"very dry", 0.5, 2.0},
		{"dry", 0.5, 7.0},
		{"ideal", 0.5, 15.0},
		{"heavy", 0.5, 25.0},
		{"very heavy", 0.5, 40.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rng := rand.New(rand.NewSource(123))
			got := updateSoilFertility(tt.current, tt.rainfall, rng)

			// Should stay in [0,1] range
			if got < 0 || got > 1.0 {
				t.Errorf("updateSoilFertility(%f, %f) = %f out of [0,1] range",
					tt.current, tt.rainfall, got)
			}

			// Different rainfall should produce different results
			// (not testing exact values due to randomness)
		})
	}

	// Test clamping at boundaries
	t.Run("clamp at 0", func(t *testing.T) {
		rng := rand.New(rand.NewSource(456))
		got := updateSoilFertility(0.0, 2.0, rng) // very dry reduces fertility
		if got < 0 {
			t.Errorf("updateSoilFertility should clamp at 0, got %f", got)
		}
	})

	t.Run("clamp at 1", func(t *testing.T) {
		rng := rand.New(rand.NewSource(789))
		got := updateSoilFertility(1.0, 15.0, rng) // ideal improves fertility
		if got > 1.0 {
			t.Errorf("updateSoilFertility should clamp at 1.0, got %f", got)
		}
	})
}

func TestRegenerateResource(t *testing.T) {
	rng := rand.New(rand.NewSource(456))

	// Test regeneration
	initial := 0.5
	got := regenerateResource(initial, rng, 0.01)

	// Should change slightly
	if got == initial {
		t.Error("regenerateResource should change value")
	}

	// Should stay in [0,1] range
	if got < 0 || got > 1.0 {
		t.Errorf("regenerateResource result %f out of [0,1] range", got)
	}

	// Test clamping at 0
	rng = rand.New(rand.NewSource(456))
	low := regenerateResource(0.0, rng, -0.1) // negative rate
	if low < 0 {
		t.Errorf("regenerateResource should clamp at 0, got %f", low)
	}

	// Test clamping at 1
	rng = rand.New(rand.NewSource(456))
	high := regenerateResource(1.0, rng, 0.1)
	if high > 1.0 {
		t.Errorf("regenerateResource should clamp at 1.0, got %f", high)
	}
}

func TestUpdateWildlife(t *testing.T) {
	rng := rand.New(rand.NewSource(789))

	// Test with different forest health values
	tests := []struct {
		name         string
		current      float64
		forestHealth float64
	}{
		{"healthy forest", 0.5, 0.8},
		{"poor forest", 0.5, 0.2},
		{"average forest", 0.5, 0.5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := updateWildlife(tt.current, tt.forestHealth, rng)
			if got == tt.current {
				t.Error("updateWildlife should change value")
			}
			if got < 0 || got > 1.0 {
				t.Errorf("updateWildlife result %f out of [0,1] range", got)
			}
		})
	}
}

func TestGenerateEventID(t *testing.T) {
	// Test deterministic IDs
	id1 := generateEventID("rain", 1)
	id2 := generateEventID("rain", 1)
	if id1 != id2 {
		t.Errorf("generateEventID not deterministic: %s vs %s", id1, id2)
	}

	// Different base or week should produce different IDs
	id3 := generateEventID("snow", 1)
	if id3 == id1 {
		t.Error("different base should produce different ID")
	}

	id4 := generateEventID("rain", 2)
	if id4 == id1 {
		t.Error("different week should produce different ID")
	}

	// Should contain base and some suffix
	if len(id1) <= len("rain") {
		t.Errorf("ID too short: %s", id1)
	}
}

func TestFormatWeekTimestamp(t *testing.T) {
	tests := []struct {
		year int
		week int
		want string
	}{
		{1, 1, "1-W1"},
		{2025, 52, "2025-W52"},
		{0, 0, "0-W0"},
	}

	for _, tt := range tests {
		got := formatWeekTimestamp(tt.year, tt.week)
		if got != tt.want {
			t.Errorf("formatWeekTimestamp(%d, %d) = %s, want %s", tt.year, tt.week, got, tt.want)
		}
	}
}

func TestGenerateWeatherEvents(t *testing.T) {
	tests := []struct {
		name       string
		temp       float64
		rain       float64
		week       int
		year       int
		wantEvents int
	}{
		{"normal", 15.0, 10.0, 1, 1, 0},
		{"heavy rain", 15.0, 25.0, 1, 1, 1},
		{"freezing", -1.0, 5.0, 1, 1, 1},
		{"both extreme", -5.0, 30.0, 1, 1, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			events := GenerateWeatherEvents(tt.temp, tt.rain, tt.week, tt.year)
			if len(events) != tt.wantEvents {
				t.Errorf("GenerateWeatherEvents(%f, %f, %d, %d) = %d events, want %d",
					tt.temp, tt.rain, tt.week, tt.year, len(events), tt.wantEvents)
			}
			// Check event types
			for _, ev := range events {
				if ev.Type != "weather" {
					t.Errorf("event type %s, want weather", ev.Type)
				}
				if tt.rain > 20.0 && ev.Data["severity"] == "heavy" {
					// good
				} else if tt.temp < 0.0 && ev.Data["severity"] == "freezing" {
					// good
				} else {
					t.Errorf("unexpected event data: %v", ev.Data)
				}
			}
		})
	}
}

func TestEnvironmentUpdateManyWeeks(t *testing.T) {
	env := NewEnvironmentSystem()
	state := NewGameState("test-many-weeks", 12345)

	// Run through a full year (52 weeks)
	for week := 1; week <= 52; week++ {
		state.Calendar.Week = week
		events := env.Update(week, state, state.RNG.Rand())
		_ = events // ensure no panic
		// After each update, environment should reflect the week
		if state.Environment.Season != calculateSeason(week) {
			t.Errorf("week %d: season %s, expected %s", week, state.Environment.Season, calculateSeason(week))
		}
	}

	// Ensure all environmental values remain in valid ranges
	if state.Environment.Temperature < -50 || state.Environment.Temperature > 50 {
		t.Errorf("temperature out of plausible range: %f", state.Environment.Temperature)
	}
	if state.Environment.Rainfall < 0 || state.Environment.Rainfall > 100 {
		t.Errorf("rainfall out of plausible range: %f", state.Environment.Rainfall)
	}
	if state.Environment.SoilFertility < 0 || state.Environment.SoilFertility > 1.0 {
		t.Errorf("soil fertility out of range: %f", state.Environment.SoilFertility)
	}
	if state.Environment.ForestHealth < 0 || state.Environment.ForestHealth > 1.0 {
		t.Errorf("forest health out of range: %f", state.Environment.ForestHealth)
	}
	if state.Environment.MineQuality < 0 || state.Environment.MineQuality > 1.0 {
		t.Errorf("mine quality out of range: %f", state.Environment.MineQuality)
	}
	if state.Environment.WildlifePopulation < 0 || state.Environment.WildlifePopulation > 1.0 {
		t.Errorf("wildlife population out of range: %f", state.Environment.WildlifePopulation)
	}
}
