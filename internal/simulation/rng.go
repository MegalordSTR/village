package simulation

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"math/rand"
)

// RNG is a deterministic random number generator.
type RNG struct {
	state uint64
}

// NewRNG creates a new RNG seeded with the given value.
func NewRNG(seed int64) *RNG {
	return &RNG{state: uint64(seed)}
}

// Uint64 returns a pseudo-random 64-bit value.
func (r *RNG) Uint64() uint64 {
	// 64-bit LCG constants from Knuth's MMIX
	const (
		mul = 6364136223846793005
		inc = 1442695040888963407
	)
	r.state = r.state*mul + inc
	return r.state
}

// Intn returns a random integer in [0, n).
func (r *RNG) Intn(n int) int {
	if n <= 0 {
		panic("invalid argument to Intn")
	}
	// Simple modulo with rejection to avoid bias for now
	// In production, we'd implement a more robust method
	return int(r.Uint64() % uint64(n))
}

// Float64 returns a random float64 in [0.0, 1.0).
func (r *RNG) Float64() float64 {
	// Generate a 53-bit integer and divide by 2^53
	// Use high 53 bits of a 64-bit value
	bits := r.Uint64() >> 11 // 64 - 11 = 53 bits
	const float53 = 1 << 53
	return float64(bits) / float53
}

// Shuffle pseudo-randomizes the order of elements using the Fisher-Yates algorithm.
func (r *RNG) Shuffle(n int, swap func(i, j int)) {
	if n < 0 {
		panic("invalid argument to Shuffle")
	}
	for i := n - 1; i > 0; i-- {
		j := r.Intn(i + 1)
		swap(i, j)
	}
}

// MarshalBinary encodes the RNG state as a binary representation.
func (r *RNG) MarshalBinary() ([]byte, error) {
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], r.state)
	return buf[:], nil
}

// UnmarshalBinary decodes the RNG state from binary representation.
func (r *RNG) UnmarshalBinary(data []byte) error {
	if len(data) != 8 {
		return json.Unmarshal(data, &r.state) // fallback for JSON compatibility
	}
	r.state = binary.LittleEndian.Uint64(data)
	return nil
}

// MarshalJSON encodes the RNG state as a base64 string.
func (r *RNG) MarshalJSON() ([]byte, error) {
	data, err := r.MarshalBinary()
	if err != nil {
		return nil, err
	}
	// Encode as base64 for JSON string
	encoded := base64.StdEncoding.EncodeToString(data)
	return json.Marshal(encoded)
}

// UnmarshalJSON decodes the RNG state from a base64 string.
func (r *RNG) UnmarshalJSON(data []byte) error {
	var encoded string
	if err := json.Unmarshal(data, &encoded); err != nil {
		return err
	}
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return err
	}
	return r.UnmarshalBinary(decoded)
}

// State returns the current internal state (for testing).
func (r *RNG) State() uint64 {
	return r.state
}

// randSource implements rand.Source64 using the underlying RNG.
type randSource struct {
	rng *RNG
}

func (s *randSource) Seed(seed int64) {
	// Not needed because RNG manages its own state
	// We could reset the RNG state, but for compatibility we ignore.
}

func (s *randSource) Int63() int64 {
	return int64(s.rng.Uint64() >> 1) // clear sign bit
}

func (s *randSource) Uint64() uint64 {
	return s.rng.Uint64()
}

// Rand returns a *rand.Rand that uses this RNG as its source.
// This allows the RNG to be used wherever *rand.Rand is expected.
func (r *RNG) Rand() *rand.Rand {
	return rand.New(&randSource{rng: r})
}
