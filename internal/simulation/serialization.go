package simulation

import "encoding/json"

// SaveRNG encodes the RNG state as JSON bytes.
func SaveRNG(rng *RNG) ([]byte, error) {
	return json.Marshal(rng)
}

// LoadRNG decodes RNG state from JSON bytes.
func LoadRNG(data []byte) (*RNG, error) {
	var rng RNG
	if err := json.Unmarshal(data, &rng); err != nil {
		return nil, err
	}
	return &rng, nil
}

// SaveRNGBinary encodes the RNG state as binary bytes.
func SaveRNGBinary(rng *RNG) ([]byte, error) {
	return rng.MarshalBinary()
}

// LoadRNGBinary decodes RNG state from binary bytes.
func LoadRNGBinary(data []byte) (*RNG, error) {
	var rng RNG
	if err := rng.UnmarshalBinary(data); err != nil {
		return nil, err
	}
	return &rng, nil
}
