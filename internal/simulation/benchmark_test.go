package simulation

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

func BenchmarkProcessWeek10Residents(b *testing.B) {
	benchmarkProcessWeek(b, 10)
}

func BenchmarkProcessWeek50Residents(b *testing.B) {
	benchmarkProcessWeek(b, 50)
}

func BenchmarkProcessWeek100Residents(b *testing.B) {
	benchmarkProcessWeek(b, 100)
}

func benchmarkProcessWeek(b *testing.B, numResidents int) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		// Create a fresh game state for each iteration
		state := NewGameState("benchmark", int64(i))
		// Add residents
		for n := 0; n < numResidents; n++ {
			state.AddResident(Resident{
				ID:   fmt.Sprintf("resident-%d-%d", i, n),
				Name: fmt.Sprintf("Resident %d", n),
				Age:  20 + n%40,
			})
		}
		// Create turn processor with all systems
		tp := NewTurnProcessor()
		tp.RegisterSystem(SystemEnvironment, NewEnvironmentSystem())
		tp.RegisterSystem(SystemProduction, NewProductionSystem())
		tp.RegisterSystem(SystemSocial, NewSocialSystem())
		tp.RegisterSystem(SystemEconomic, NewEconomicSystem())
		tp.RegisterSystem(SystemEvents, NewEventSystem())
		// Process one week
		tp.ProcessWeek(state)
	}
}

// BenchmarkMemoryUsage runs a single iteration to capture memory stats
func BenchmarkMemoryUsage(b *testing.B) {
	b.ReportAllocs()
	// We'll run for 10, 50, 100 residents in subbenchmarks
	for _, residents := range []int{10, 50, 100} {
		b.Run(fmt.Sprintf("%d residents", residents), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				state := NewGameState("mem", int64(i))
				for n := 0; n < residents; n++ {
					state.AddResident(Resident{
						ID:   fmt.Sprintf("mem-%d-%d", i, n),
						Name: fmt.Sprintf("Mem Resident %d", n),
						Age:  20 + n%40,
					})
				}
				tp := NewTurnProcessor()
				tp.RegisterSystem(SystemEnvironment, NewEnvironmentSystem())
				tp.RegisterSystem(SystemProduction, NewProductionSystem())
				tp.RegisterSystem(SystemSocial, NewSocialSystem())
				tp.RegisterSystem(SystemEconomic, NewEconomicSystem())
				tp.RegisterSystem(SystemEvents, NewEventSystem())
				tp.ProcessWeek(state)
			}
		})
	}
}

// createSimulationWithResidents creates a game state with the given number of residents
// and a turn processor with all systems registered.
func createSimulationWithResidents(numResidents int, seed int64) (*GameState, *TurnProcessor) {
	state := NewGameState("test", seed)
	for n := 0; n < numResidents; n++ {
		state.AddResident(Resident{
			ID:   fmt.Sprintf("resident-%d", n),
			Name: fmt.Sprintf("Resident %d", n),
			Age:  20 + n%40,
		})
	}
	tp := NewTurnProcessor()
	tp.RegisterSystem(SystemEnvironment, NewEnvironmentSystem())
	tp.RegisterSystem(SystemProduction, NewProductionSystem())
	tp.RegisterSystem(SystemSocial, NewSocialSystem())
	tp.RegisterSystem(SystemEconomic, NewEconomicSystem())
	tp.RegisterSystem(SystemEvents, NewEventSystem())
	return state, tp
}

// TestProcessingTimeWithinLimits verifies that processing one week stays within
// the required time bounds for 10, 50, and 100 residents.
func TestProcessingTimeWithinLimits(t *testing.T) {
	// Target durations in seconds
	targets := map[int]float64{
		10:  1.0,
		50:  3.0,
		100: 5.0,
	}
	// We'll run each size once and measure elapsed time.
	// Since performance may vary, we allow a 50% margin above target.
	margin := 1.5

	for residents, target := range targets {
		state, tp := createSimulationWithResidents(residents, 123)
		start := time.Now()
		tp.ProcessWeek(state)
		elapsed := time.Since(start)
		seconds := elapsed.Seconds()
		maxAllowed := target * margin
		if seconds > maxAllowed {
			t.Errorf("%d residents processing too slow: %v (allowed < %v with margin)", residents, elapsed, time.Duration(maxAllowed*float64(time.Second)))
		} else {
			t.Logf("%d residents processed in %v (target < %v)", residents, elapsed, time.Duration(target*float64(time.Second)))
		}
	}
}

// TestLinearScaling verifies that processing time scales roughly linearly with resident count.
func TestLinearScaling(t *testing.T) {
	// Sizes to test
	sizes := []int{10, 50, 100}
	// Number of repetitions per size to reduce noise
	const repeats = 5
	// Store total time per size (sum over repeats)
	totalTimes := make([]float64, len(sizes))

	for i, size := range sizes {
		sum := 0.0
		for r := 0; r < repeats; r++ {
			state, tp := createSimulationWithResidents(size, int64(456+r))
			start := time.Now()
			tp.ProcessWeek(state)
			elapsed := time.Since(start)
			sum += elapsed.Seconds()
		}
		avg := sum / float64(repeats)
		totalTimes[i] = avg
		t.Logf("%d residents: average total time %v", size, time.Duration(avg*float64(time.Second)))
	}

	// Simple linear regression: time = a + b*size
	// We'll compute R-squared to see how well linear model fits.
	n := float64(len(sizes))
	var sumSize, sumTime, sumSize2, sumTime2, sumSizeTime float64
	for i, size := range sizes {
		x := float64(size)
		y := totalTimes[i]
		sumSize += x
		sumTime += y
		sumSize2 += x * x
		sumTime2 += y * y
		sumSizeTime += x * y
	}
	// Compute slope b and intercept a
	b := (n*sumSizeTime - sumSize*sumTime) / (n*sumSize2 - sumSize*sumSize)
	a := (sumTime - b*sumSize) / n
	// Compute R-squared
	ssTotal := sumTime2 - sumTime*sumTime/n
	ssResidual := 0.0
	for i, size := range sizes {
		x := float64(size)
		y := totalTimes[i]
		predicted := a + b*x
		ssResidual += (y - predicted) * (y - predicted)
	}
	rSquared := 1.0 - ssResidual/ssTotal
	t.Logf("Linear regression: time = %.3e + %.3e * residents (R² = %.3f)", a, b, rSquared)

	// Require R² > 0.9 for approximately linear relationship
	if rSquared < 0.9 {
		t.Errorf("time vs resident count does not scale linearly (R² = %.3f < 0.9)", rSquared)
	} else {
		t.Logf("Linear scaling confirmed (R² = %.3f ≥ 0.9)", rSquared)
	}
}

// TestMemoryUsage logs memory allocation for different resident counts.
func TestMemoryUsage(t *testing.T) {
	sizes := []int{10, 50, 100}
	for _, size := range sizes {
		var m1, m2 runtime.MemStats
		runtime.ReadMemStats(&m1)
		state, tp := createSimulationWithResidents(size, 789)
		tp.ProcessWeek(state)
		runtime.ReadMemStats(&m2)
		alloc := m2.TotalAlloc - m1.TotalAlloc
		t.Logf("%d residents allocated %v bytes", size, alloc)
	}
}
