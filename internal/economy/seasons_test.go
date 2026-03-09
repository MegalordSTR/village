package economy

import (
	"github.com/vano44/village/internal/simulation"
	"testing"
)

func TestSeasonConstants(t *testing.T) {
	// This test will fail until Season type and constants are defined
	var _ Season = Spring //nolint:staticcheck
	var _ Season = Summer //nolint:staticcheck
	var _ Season = Autumn //nolint:staticcheck
	var _ Season = Winter //nolint:staticcheck
}

func TestSeasonFromWeek(t *testing.T) {
	tests := []struct {
		week int
		want Season
	}{
		{1, Spring},
		{13, Spring},
		{14, Summer},
		{26, Summer},
		{27, Autumn},
		{39, Autumn},
		{40, Winter},
		{52, Winter},
	}
	for _, tt := range tests {
		got := SeasonFromWeek(tt.week)
		if got != tt.want {
			t.Errorf("SeasonFromWeek(%d) = %v, want %v", tt.week, got, tt.want)
		}
	}
}

func TestAgricultureYieldMultiplier(t *testing.T) {
	tests := []struct {
		season Season
		want   float64
	}{
		{Spring, 1.0},
		{Summer, 1.5},
		{Autumn, 1.0},
		{Winter, 0.0},
	}
	for _, tt := range tests {
		got := AgricultureYieldMultiplier(tt.season)
		if got != tt.want {
			t.Errorf("AgricultureYieldMultiplier(%v) = %v, want %v", tt.season, got, tt.want)
		}
	}
}

func TestFuelConsumptionMultiplier(t *testing.T) {
	tests := []struct {
		season Season
		want   float64
	}{
		{Spring, 1.0},
		{Summer, 1.0},
		{Autumn, 1.0},
		{Winter, 2.5},
	}
	for _, tt := range tests {
		got := FuelConsumptionMultiplier(tt.season)
		if got != tt.want {
			t.Errorf("FuelConsumptionMultiplier(%v) = %v, want %v", tt.season, got, tt.want)
		}
	}
}

func TestCropGrowthStages(t *testing.T) {
	// Test CropStage constants
	var _ CropStage = CropPlanted         //nolint:staticcheck
	var _ CropStage = CropGrowing         //nolint:staticcheck
	var _ CropStage = CropReadyForHarvest //nolint:staticcheck
	var _ CropStage = CropHarvested       //nolint:staticcheck

	// Test NewCrop
	crop := NewCrop(ResourceGrain, 1)
	if crop.Type != ResourceGrain {
		t.Errorf("NewCrop type = %v, want %v", crop.Type, ResourceGrain)
	}
	if crop.Stage != CropPlanted {
		t.Errorf("NewCrop stage = %v, want %v", crop.Stage, CropPlanted)
	}
	if crop.PlantedWeek != 1 {
		t.Errorf("NewCrop planted week = %v, want 1", crop.PlantedWeek)
	}
	if crop.GrowthProgress != 0.0 {
		t.Errorf("NewCrop progress = %v, want 0.0", crop.GrowthProgress)
	}
}

func TestWeatherPatterns(t *testing.T) {
	// Test WeatherPattern constants
	var _ WeatherPattern = WeatherDrought //nolint:staticcheck
	var _ WeatherPattern = WeatherNormal  //nolint:staticcheck
	var _ WeatherPattern = WeatherRain    //nolint:staticcheck

	// Test YieldMultiplier
	if got := YieldMultiplier(WeatherDrought); got != 0.5 {
		t.Errorf("YieldMultiplier(drought) = %v, want 0.5", got)
	}
	if got := YieldMultiplier(WeatherNormal); got != 1.0 {
		t.Errorf("YieldMultiplier(normal) = %v, want 1.0", got)
	}
	if got := YieldMultiplier(WeatherRain); got != 1.3 {
		t.Errorf("YieldMultiplier(rain) = %v, want 1.3", got)
	}
	// Test FloodingProbability
	if got := FloodingProbability(WeatherRain); got != 0.1 {
		t.Errorf("FloodingProbability(rain) = %v, want 0.1", got)
	}
	if got := FloodingProbability(WeatherDrought); got != 0.0 {
		t.Errorf("FloodingProbability(drought) = %v, want 0.0", got)
	}
}

func TestAgriculturalCalendar(t *testing.T) {
	// Planting window: spring weeks 1-4
	for week := 1; week <= 4; week++ {
		if !CanPlant(week) {
			t.Errorf("CanPlant(%d) = false, want true", week)
		}
	}
	// Outside planting window
	if CanPlant(5) {
		t.Errorf("CanPlant(5) = true, want false")
	}
	if CanPlant(14) {
		t.Errorf("CanPlant(14) = true, want false")
	}
	// Harvest window: autumn weeks 8-12 (weeks 35-39?)
	// autumn starts week 27, weeks 8-12 => weeks 34-38? Wait: week 27 + 7 = week 34? Let's compute.
	// autumn weeks: 27-39. weekInSeason = week - 26. So week 35 => weekInSeason 9. Should be within 8-12.
	for week := 34; week <= 38; week++ {
		if !CanHarvest(week) {
			t.Errorf("CanHarvest(%d) = false, want true", week)
		}
	}
	// Outside harvest window
	if CanHarvest(27) {
		t.Errorf("CanHarvest(27) = true, want false")
	}
	if CanHarvest(40) {
		t.Errorf("CanHarvest(40) = true, want false")
	}
}

func TestWinterChallenges(t *testing.T) {
	// No agriculture possible in winter
	crop := NewCrop(ResourceGrain, 1)
	changed := crop.UpdateGrowth(40, Winter, 1.0) // week 40 is winter
	if changed {
		t.Errorf("UpdateGrowth in winter changed stage, should not")
	}
	if crop.GrowthProgress != 0.0 {
		t.Errorf("GrowthProgress in winter should remain 0, got %v", crop.GrowthProgress)
	}
	// Fuel consumption multiplier already tested
}

func TestSeasonalPriceFluctuations(t *testing.T) {
	// Food resources
	foods := []ResourceType{ResourceGrain, ResourceVegetables, ResourceFlour, ResourceBread}
	// Non-food resource
	nonFood := ResourceWood
	// Test spring expensive
	for _, food := range foods {
		got := PriceMultiplier(food, Spring, 1)
		if got != 1.2 {
			t.Errorf("PriceMultiplier(%v, Spring, 1) = %v, want 1.2", food, got)
		}
	}
	// Test autumn weeks 8-12 cheap
	for week := 34; week <= 38; week++ {
		for _, food := range foods {
			got := PriceMultiplier(food, Autumn, week)
			if got != 0.8 {
				t.Errorf("PriceMultiplier(%v, Autumn, %d) = %v, want 0.8", food, week, got)
			}
		}
	}
	// Test non-food unaffected
	if got := PriceMultiplier(nonFood, Spring, 1); got != 1.0 {
		t.Errorf("PriceMultiplier(Wood, Spring, 1) = %v, want 1.0", got)
	}
}

func TestStoragePreparationIndicators(t *testing.T) {
	// Winter recommendations
	rec := RecommendedStockpiles(Winter)
	if rec == nil {
		t.Fatal("RecommendedStockpiles(Winter) returned nil")
	}
	if rec[ResourceGrain] != 100.0 {
		t.Errorf("Recommended grain = %v, want 100.0", rec[ResourceGrain])
	}
	if rec[ResourceVegetables] != 50.0 {
		t.Errorf("Recommended vegetables = %v, want 50.0", rec[ResourceVegetables])
	}
	if rec[ResourceWood] != 200.0 {
		t.Errorf("Recommended wood = %v, want 200.0", rec[ResourceWood])
	}
	// Other seasons return nil
	if RecommendedStockpiles(Spring) != nil {
		t.Error("RecommendedStockpiles(Spring) should return nil")
	}
}

func TestIntegrationWithEnvironment(t *testing.T) {
	env := &simulation.Environment{
		Season: "spring",
	}
	got := SeasonFromEnvironment(env)
	if got != Spring {
		t.Errorf("SeasonFromEnvironment(spring) = %v, want Spring", got)
	}
	env.Season = "summer"
	got = SeasonFromEnvironment(env)
	if got != Summer {
		t.Errorf("SeasonFromEnvironment(summer) = %v, want Summer", got)
	}
	env.Season = "autumn"
	got = SeasonFromEnvironment(env)
	if got != Autumn {
		t.Errorf("SeasonFromEnvironment(autumn) = %v, want Autumn", got)
	}
	env.Season = "winter"
	got = SeasonFromEnvironment(env)
	if got != Winter {
		t.Errorf("SeasonFromEnvironment(winter) = %v, want Winter", got)
	}
}
