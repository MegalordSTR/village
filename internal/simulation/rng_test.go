package simulation

import (
	"encoding/json"
	"testing"
)

var _ = json.Marshal

func TestRNGStruct(t *testing.T) {
	// This test will fail until RNG is defined
	rng := NewRNG(42)
	if rng == nil {
		t.Fatal("NewRNG returned nil")
	}
	_ = rng.State() // ensure State method is covered
}

func TestRNGIntnDeterministic(t *testing.T) {
	rng1 := NewRNG(123)
	rng2 := NewRNG(123)

	// Generate sequence of 10 numbers
	seq1 := make([]int, 10)
	for i := range seq1 {
		seq1[i] = rng1.Intn(100)
	}

	seq2 := make([]int, 10)
	for i := range seq2 {
		seq2[i] = rng2.Intn(100)
	}

	for i := range seq1 {
		if seq1[i] != seq2[i] {
			t.Errorf("position %d: seq1=%d, seq2=%d", i, seq1[i], seq2[i])
		}
	}
}

func TestRNGFloat64Deterministic(t *testing.T) {
	rng1 := NewRNG(456)
	rng2 := NewRNG(456)

	seq1 := make([]float64, 10)
	for i := range seq1 {
		seq1[i] = rng1.Float64()
	}

	seq2 := make([]float64, 10)
	for i := range seq2 {
		seq2[i] = rng2.Float64()
	}

	for i := range seq1 {
		if seq1[i] != seq2[i] {
			t.Errorf("position %d: seq1=%f, seq2=%f", i, seq1[i], seq2[i])
		}
	}
}

func TestRNGShuffleDeterministic(t *testing.T) {
	rng1 := NewRNG(789)
	rng2 := NewRNG(789)

	slice1 := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	slice2 := make([]int, len(slice1))
	copy(slice2, slice1)

	rng1.Shuffle(len(slice1), func(i, j int) {
		slice1[i], slice1[j] = slice1[j], slice1[i]
	})

	rng2.Shuffle(len(slice2), func(i, j int) {
		slice2[i], slice2[j] = slice2[j], slice2[i]
	})

	for i := range slice1 {
		if slice1[i] != slice2[i] {
			t.Errorf("position %d: slice1=%d, slice2=%d", i, slice1[i], slice2[i])
		}
	}
}

func TestRNGMarshalUnmarshal(t *testing.T) {
	rng := NewRNG(999)
	// Generate some numbers to advance state
	for i := 0; i < 5; i++ {
		rng.Intn(100)
	}

	data, err := rng.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON failed: %v", err)
	}

	rng2 := NewRNG(0) // different seed
	if err := rng2.UnmarshalJSON(data); err != nil {
		t.Fatalf("UnmarshalJSON failed: %v", err)
	}

	// Both should produce same sequence now
	val1 := rng.Intn(1000)
	val2 := rng2.Intn(1000)
	if val1 != val2 {
		t.Errorf("after unmarshal: val1=%d, val2=%d", val1, val2)
	}
}

func TestRNGDeterminismSequence1000(t *testing.T) {
	rng1 := NewRNG(42)
	rng2 := NewRNG(42)

	for i := 0; i < 1000; i++ {
		val1 := rng1.Intn(1000000)
		val2 := rng2.Intn(1000000)
		if val1 != val2 {
			t.Fatalf("position %d: val1=%d, val2=%d", i, val1, val2)
		}
	}
}

func TestRNGGameStateIntegration(t *testing.T) {
	gs := NewGameState("test", 12345)
	// Check that RNG field is initialized
	if gs.RNG == nil {
		t.Fatal("GameState.RNG field is nil")
	}
	rng := gs.RNG

	// Should be deterministic based on seed
	rng2 := NewRNG(12345)
	val1 := rng.Intn(100)
	val2 := rng2.Intn(100)
	if val1 != val2 {
		t.Errorf("RNG from GameState mismatch: %d vs %d", val1, val2)
	}
}

func TestRNGSerialization(t *testing.T) {
	rng := NewRNG(777)
	// advance state
	rng.Intn(100)
	rng.Float64()

	// Test JSON serialization
	data, err := SaveRNG(rng)
	if err != nil {
		t.Fatalf("SaveRNG failed: %v", err)
	}
	rng2, err := LoadRNG(data)
	if err != nil {
		t.Fatalf("LoadRNG failed: %v", err)
	}
	if rng.State() != rng2.State() {
		t.Errorf("state mismatch after JSON serialization")
	}

	// Test binary serialization
	bin, err := SaveRNGBinary(rng)
	if err != nil {
		t.Fatalf("SaveRNGBinary failed: %v", err)
	}
	rng3, err := LoadRNGBinary(bin)
	if err != nil {
		t.Fatalf("LoadRNGBinary failed: %v", err)
	}
	if rng.State() != rng3.State() {
		t.Errorf("state mismatch after binary serialization")
	}
}

func TestRNGEdgeCases(t *testing.T) {
	t.Run("Intn panics on n <= 0", func(t *testing.T) {
		rng := NewRNG(42)
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for n <= 0")
			}
		}()
		rng.Intn(0)
	})

	t.Run("Intn panics on negative n", func(t *testing.T) {
		rng := NewRNG(42)
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for negative n")
			}
		}()
		rng.Intn(-5)
	})

	t.Run("Shuffle panics on negative n", func(t *testing.T) {
		rng := NewRNG(42)
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for negative n")
			}
		}()
		rng.Shuffle(-1, func(i, j int) {})
	})

	t.Run("Shuffle n=0 does nothing", func(t *testing.T) {
		rng := NewRNG(42)
		// Should not panic
		rng.Shuffle(0, func(i, j int) {
			t.Error("swap should not be called for n=0")
		})
	})

	t.Run("Shuffle n=1 does nothing", func(t *testing.T) {
		rng := NewRNG(42)
		called := false
		rng.Shuffle(1, func(i, j int) {
			called = true
		})
		if called {
			t.Error("swap should not be called for n=1")
		}
	})

	t.Run("UnmarshalBinary with invalid length", func(t *testing.T) {
		var rng RNG
		err := rng.UnmarshalBinary([]byte{1, 2, 3}) // length 3
		if err == nil {
			t.Error("expected error for invalid length")
		}
	})

	t.Run("UnmarshalJSON invalid base64", func(t *testing.T) {
		var rng RNG
		invalidJSON := []byte(`"not-base64!"`)
		err := rng.UnmarshalJSON(invalidJSON)
		if err == nil {
			t.Error("expected error for invalid base64")
		}
	})

	t.Run("UnmarshalJSON invalid JSON", func(t *testing.T) {
		var rng RNG
		err := rng.UnmarshalJSON([]byte(`{invalid}`))
		if err == nil {
			t.Error("expected error for invalid JSON")
		}
	})

	t.Run("LoadRNG invalid data", func(t *testing.T) {
		_, err := LoadRNG([]byte(`{invalid}`))
		if err == nil {
			t.Error("expected error for invalid JSON")
		}
	})

	t.Run("LoadRNGBinary invalid data", func(t *testing.T) {
		_, err := LoadRNGBinary([]byte{1, 2, 3})
		if err == nil {
			t.Error("expected error for invalid binary length")
		}
	})
}
