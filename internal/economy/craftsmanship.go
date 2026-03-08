package economy

import (
	"math/rand"
)

// QualityFromSkill determines the achievable quality tier based on worker skill level.
// Skill is a float64 in range [0,1] where 0 is novice, 1 is master.
// Randomness is used for masterwork chance; provide a deterministic *rand.Rand for reproducible results.
// Novice (skill ≤ 0.3) can produce Poor to Normal quality.
// Expert (0.3 < skill < 0.8) can produce Good to Excellent quality.
// Master (skill ≥ 0.8) always produces Masterwork.
func QualityFromSkill(skill float64, rng *rand.Rand) QualityTier {
	// Clamp skill to valid range
	if skill < 0 {
		skill = 0
	}
	if skill > 1 {
		skill = 1
	}

	if skill <= 0.3 {
		// Novice: Poor up to 0.15, Normal from 0.15 to 0.3
		if skill < 0.15 {
			return QualityPoor
		}
		return QualityNormal
	} else if skill < 0.8 {
		// Expert: Good up to 0.55, Excellent from 0.55 to 0.8
		if skill < 0.55 {
			return QualityGood
		}
		return QualityExcellent
	} else {
		// Master: always Masterwork
		return QualityMasterwork
	}
}
