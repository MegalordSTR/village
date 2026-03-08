package economy

// Skill IDs for economic activities.
const (
	SkillFarming   = "farming"
	SkillMining    = "mining"
	SkillSmithing  = "smithing"
	SkillWeaving   = "weaving"
	SkillBaking    = "baking"
	SkillMilling   = "milling"
	SkillCarpentry = "carpentry"
	SkillTailoring = "tailoring"
)

// Skill represents a worker's proficiency in a particular economic activity.
type Skill struct {
	ID    string  // skill identifier
	Level float64 // 0..1, where 0 is novice, 1 is master
	XP    float64 // experience points accumulated
}

// IsValidSkill returns true if the given skill ID is one of the defined constants.
func IsValidSkill(id string) bool {
	switch id {
	case SkillFarming, SkillMining, SkillSmithing, SkillWeaving,
		SkillBaking, SkillMilling, SkillCarpentry, SkillTailoring:
		return true
	default:
		return false
	}
}

// Practice adds experience points to the skill.
// The amount is typically proportional to time spent or tasks completed.
func (s *Skill) Practice(amount float64) {
	if amount < 0 {
		amount = 0
	}
	s.XP += amount
}

// UpdateLevel recomputes Level based on XP.
// The formula: Level = XP / 100, capped at 1.0.
// Teachers can provide bonus XP (not implemented yet).
func (s *Skill) UpdateLevel() {
	level := s.XP / 100.0
	if level > 1.0 {
		level = 1.0
	}
	s.Level = level
}

// NewSkill creates a new skill with zero XP and Level 0.
func NewSkill(id string) *Skill {
	return &Skill{
		ID:    id,
		Level: 0.0,
		XP:    0.0,
	}
}
