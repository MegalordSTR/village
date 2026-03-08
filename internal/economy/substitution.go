package economy

// Alternative represents a possible substitute for a required resource.
type Alternative struct {
	Type              ResourceType // substitute resource type
	TimeMultiplier    float64      // multiplier for production time (e.g., 2.0 for 2x time)
	QualityMultiplier float64      // multiplier for output quality (e.g., 0.8 for 80%)
}

// substitutionRules maps a required resource type to a list of possible alternatives.
var substitutionRules = make(map[ResourceType][]Alternative)

// RegisterSubstitution adds an alternative for a required resource type.
func RegisterSubstitution(required ResourceType, alt Alternative) {
	substitutionRules[required] = append(substitutionRules[required], alt)
}

// GetAlternatives returns the list of alternatives for a required resource type.
// Returns nil if no alternatives defined.
func GetAlternatives(required ResourceType) []Alternative {
	return substitutionRules[required]
}

// ResetSubstitutionRules clears all registered substitution rules.
// Useful for testing.
func ResetSubstitutionRules() {
	substitutionRules = make(map[ResourceType][]Alternative)
}
