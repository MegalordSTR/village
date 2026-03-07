# Social Systems

**Vision:** Transform villagers from economic units into believable medieval people with complex social lives—hierarchies, beliefs, conflicts, and culture that shape village dynamics as much as resource production.

**Status:** Draft  
**Created:** March 7, 2026  
**Priority:** P1 (Post-MVP Depth)

---

## Executive Summary

Social systems add the human layer to economic simulation: residents have social status, religious beliefs, political opinions, and cultural traditions that influence their behavior, productivity, and interactions. A blacksmith's daughter might marry a farmer's son despite class differences, a charismatic resident could inspire others to work harder or protest unfair conditions, and religious festivals become community touchstones. These systems create emergent stories of love, ambition, faith, and conflict.

**Success Definition:** Players discuss villagers' social lives ("The miller's feud with the blacksmith is affecting tool production!") and feel they're managing a society, not just an economy.

---

## User Stories

### US-001: Social Classes & Status
**As:** A resident in a medieval village  
**I want:** My social status to affect my opportunities and treatment  
**So that:** The village feels historically authentic

**Acceptance Criteria:**
- [ ] Four social classes: peasants, craftsmen, clergy, minor nobility
- [ ] Status affects: marriage options, work assignments, resource access
- [ ] Social mobility through exceptional skill, wealth, or marriage
- [ ] Class tensions affecting productivity and happiness
- [ ] Visible status indicators in UI (clothing, housing quality)

### US-002: Religion & Beliefs
**As:** A medieval villager  
**I want:** Religious beliefs that shape my daily life and community  
**So that:** The village reflects medieval spirituality

**Acceptance Criteria:**
- [ ] Multiple belief systems: Christianity (Catholic/heresy), paganism, skepticism
- [ ] Religious practices: weekly services, holidays, rites of passage
- [ ] Faith level affects: happiness, event interpretation, health during plague
- [ ] Religious conflicts between believers of different faiths
- [ ] Clergy as distinct social role with political influence

### US-003: Crime & Justice
**As:** A village leader  
**I want:** To handle theft, violence, and other crimes appropriately  
**So that:** Law and order affect village stability

**Acceptance Criteria:**
- [ ] Crime system: theft, assault, murder, heresy accusations
- [ ] Justice options: fines, imprisonment, exile, execution
- [ ] Crime causes: desperation (starvation), personality (aggressive traits)
- [ ] Justice consequences: deterrence vs. resentment, vigilante justice
- [ ] Reputation system affecting trade and immigration

### US-004: Education & Knowledge
**As:** A player investing in village future  
**I want:** To educate residents and preserve knowledge  
**So that:** Skills and traditions survive across generations

**Acceptance Criteria:**
- [ ] Literacy as rare skill (clergy, nobility, exceptional craftsmen)
- [ ] Apprenticeship system: masters train apprentices over years
- [ ] Knowledge loss when skilled residents die without passing on skills
- [ ] Library/scriptorium building for knowledge preservation
- [ ] Oral tradition for illiterate majority (stories, songs, techniques)

### US-005: Culture & Traditions
**As:** A villager participating in community life  
**I want:** Cultural traditions that define our village identity  
**So that:** Each village develops unique character

**Acceptance Criteria:**
- [ ] Village-specific traditions: founding story, local saints, unique festivals
- [ ] Cultural drift over generations (traditions change, new ones emerge)
- [ ] Cultural exchange with visitors/traders
- [ ] Material culture: distinct architectural styles, clothing, crafts
- [ ] Oral history tracking major village events across generations

### US-006: Politics & Leadership
**As:** A player managing village governance  
**I want:** Political dynamics that challenge my authority  
**So that:** Leadership feels earned, not automatic

**Acceptance Criteria:**
- [ ] Factions within village: traditionalists vs. innovators, young vs. old
- [ ] Leadership challenges: popular residents might oppose player decisions
- [ ] Decision legitimacy: unpopular decisions may be ignored or protested
- [ ] Succession planning: player's chosen successor may not be accepted
- [ ] External political forces: lord's demands, church authority, king's laws

---

## Technical Specifications

### Social Class System
```go
type SocialClass int

const (
    ClassPeasant SocialClass = iota    // 80% of population
    ClassCraftsman                    // 15% of population
    ClassClergy                       // 3% of population  
    ClassNobility                     // 2% of population (minor gentry)
)

type SocialStatus struct {
    Class          SocialClass   `json:"class"`
    Prestige       float64       `json:"prestige"`       // 0-100 within class
    Wealth         float64       `json:"wealth"`         // Relative to class norms
    Reputation     float64       `json:"reputation"`     // -100 to 100
    Lineage        FamilyLine    `json:"lineage"`        // Family history
    Titles         []Title       `json:"titles"`         // Earned or inherited
}

// Class effects
var ClassEffects = map[SocialClass]ClassModifiers{
    ClassPeasant: {
        WorkEfficiency:    1.0,   // Base
        PoliticalPower:    0.1,   // Little influence
        ResourceAccess:    0.5,   // Limited access to luxury goods
        EducationChance:   0.05,  // 5% chance to become literate
    },
    ClassCraftsman: {
        WorkEfficiency:    1.2,   // Skilled work
        PoliticalPower:    0.3,   // Some influence through guilds
        ResourceAccess:    0.8,   // Access to tools/materials
        EducationChance:   0.3,   // 30% literacy
    },
    // ... similar for clergy and nobility
}
```

### Religion System
```go
type Religion struct {
    ID           string        `json:"id"`           // "catholic", "lollard", "pagan"
    Name         string        `json:"name"`         // "Catholic Christianity"
    Description  string        `json:"description"`
    
    // Belief tenets
    Tenets       []Tenet       `json:"tenets"`       // Core beliefs
    Practices    []Practice    `json:"practices"`    // Rituals, requirements
    Taboos       []Taboo       `json:"taboos"`       // Forbidden actions
    
    // Organizational structure
    Hierarchy    []Role        `json:"hierarchy"`    // Priest, bishop, etc.
    Buildings    []BuildingType `json:"buildings"`   // Church, shrine, etc.
}

type Belief struct {
    ReligionID   string        `json:"religionId"`
    Faith        float64       `json:"faith"`        // 0-100 belief strength
    Doubt        float64       `json:"doubt"`        // 0-100 skepticism
    Practices    []string      `json:"practices"`    // Which rituals followed
    Conversion   float64       `json:"conversion"`   // Likelihood to convert others
}

// Weekly religious update
func UpdateReligion(resident *Resident, village *Village) {
    belief := resident.Belief
    
    // Faith naturally decays without reinforcement
    belief.Faith *= 0.995  // 0.5% decay per week
    
    // Participate in religious practices
    if village.HasChurch && resident.Class >= ClassPeasant {
        // Weekly service increases faith
        belief.Faith = min(belief.Faith + 2, 100)
        
        // Social pressure in highly religious villages
        villageAverageFaith := village.CalculateAverageFaith()
        if villageAverageFaith > 70 {
            belief.Faith = min(belief.Faith + 5, 100)  // Conformity boost
        }
    }
    
    // Critical events affect faith
    for _, event := range resident.RecentEvents {
        switch event.Type {
        case EventTypeDeath:
            // Death often strengthens or weakens faith
            if rng.Float64() < 0.3 {
                belief.Faith = min(belief.Faith + 10, 100)  // Seek comfort
            } else if rng.Float64() < 0.2 {
                belief.Faith = max(belief.Faith - 15, 0)    // Question faith
            }
        case EventTypeMiracle:
            belief.Faith = min(belief.Faith + 25, 100)
        case EventTypeTragedy:
            belief.Faith *= 0.8  // 20% faith loss
        }
    }
    
    resident.Belief = belief
}
```

### Crime & Justice System
```go
type CrimeType int

const (
    CrimeTheft CrimeType = iota
    CrimeAssault
    CrimeMurder
    CrimeHeresy
    CrimeTreason
)

type Crime struct {
    ID           string        `json:"id"`
    Type         CrimeType     `json:"type"`
    Perpetrator  string        `json:"perpetrator"`  // Resident ID
    Victim       string        `json:"victim"`       // Resident ID or "village"
    Severity     float64       `json:"severity"`     // 0-100
    Evidence     float64       `json:"evidence"`     // 0-100 certainty
    Motive       string        `json:"motive"`       // "starvation", "revenge", "greed"
    Reported     bool          `json:"reported"`     // Was crime reported?
}

type JusticeOption struct {
    Name         string        `json:"name"`         // "Fine", "Imprison", "Execute"
    Description  string        `json:"description"`
    Severity     float64       `json:"severity"`     // 0-100 harshness
    Effectiveness float64      `json:"effectiveness"` // Deterrence value
    Requirements []Requirement `json:"requirements"` // Resources, buildings needed
    Consequences []Consequence `json:"consequences"` // Social effects
}

// Crime generation
func GenerateCrimes(village *Village, week int) []Crime {
    var crimes []Crime
    
    // Check each resident for crime likelihood
    for _, resident := range village.Residents {
        if !resident.Alive {
            continue
        }
        
        crimeChance := CalculateCrimeChance(resident, village)
        
        if rng.Float64() < crimeChance {
            crime := CreateCrime(resident, village)
            crimes = append(crimes, crime)
        }
    }
    
    return crimes
}

func CalculateCrimeChance(resident *Resident, village *Village) float64 {
    baseChance := 0.01  // 1% per week baseline
    
    // Economic pressure
    if resident.Hunger > 80 {
        baseChance *= 3.0  // Starvation triples crime chance
    }
    
    // Personality traits
    if resident.HasTrait("Aggressive") {
        baseChance *= 2.0
    }
    if resident.HasTrait("Honest") {
        baseChance *= 0.5
    }
    
    // Social bonds
    if resident.HasFamilyInVillage(village) {
        baseChance *= 0.7  // Family responsibility reduces crime
    }
    
    // Deterrence from previous justice
    deterrence := village.CalculateDeterrence()
    baseChance *= (1.0 - deterrence)
    
    return min(baseChance, 0.5)  // Cap at 50%
}
```

### Education & Knowledge System
```go
type KnowledgeType int

const (
    KnowledgeLiteracy KnowledgeType = iota
    KnowledgeMathematics
    KnowledgeMedicine
    KnowledgeAgriculture
    KnowledgeCraft
    KnowledgeHistory
)

type Knowledge struct {
    Type         KnowledgeType `json:"type"`
    Level        float64       `json:"level"`      // 0-100 mastery
    Source       string        `json:"source"`     // "apprenticeship", "self-taught", "inherited"
    Teachers     []string      `json:"teachers"`   // Resident IDs who taught
    Students     []string      `json:"students"`   // Resident IDs taught
    Written      bool          `json:"written"`    // Is knowledge documented?
}

type Apprenticeship struct {
    MasterID     string        `json:"masterId"`
    ApprenticeID string        `json:"apprenticeId"`
    Skill        SkillType     `json:"skill"`
    StartWeek    int           `json:"startWeek"`
    Duration     int           `json:"duration"`   // Weeks required
    Progress     float64       `json:"progress"`   // 0-100 completion
    Quality      float64       `json:"quality"`    // 0-100 teaching quality
}

// Knowledge transmission
func TransmitKnowledge(village *Village) {
    // Apprenticeships
    for i := range village.Apprenticeships {
        app := &village.Apprenticeships[i]
        
        // Both must be alive and in village
        master := village.GetResident(app.MasterID)
        apprentice := village.GetResident(app.ApprenticeID)
        
        if master == nil || apprentice == nil || 
           !master.Alive || !apprentice.Alive {
            continue
        }
        
        // Weekly progress
        weeklyProgress := 1.0
        weeklyProgress *= app.Quality / 100.0
        weeklyProgress *= master.Skills[app.Skill].Level / 100.0
        weeklyProgress *= apprentice.Skills[app.Skill].Aptitude
        
        app.Progress = min(app.Progress + weeklyProgress, 100)
        
        // Skill transfer when master and apprentice work together
        if apprentice.Assignment != nil && 
           apprentice.Assignment.BuildingID == master.Assignment.BuildingID {
           
           transferAmount := 0.1 * weeklyProgress
           apprentice.Skills[app.Skill].Level = min(
               apprentice.Skills[app.Skill].Level + transferAmount,
               master.Skills[app.Skill].Level,  // Can't exceed master
           )
        }
        
        // Completion
        if app.Progress >= 100 {
            // Apprentice becomes journeyman
            apprentice.Skills[app.Skill].Level = min(
                master.Skills[app.Skill].Level * 0.8,  // 80% of master's skill
                100,
            )
            // Remove apprenticeship
            village.Apprenticeships = append(
                village.Apprenticeships[:i],
                village.Apprenticeships[i+1:]...,
            )
        }
    }
    
    // Knowledge loss when experts die
    for _, resident := range village.Residents {
        if !resident.Alive && resident.WasExpert() {
            // Check if knowledge was passed on
            knowledgePreserved := resident.KnowledgePreservationRate(village)
            
            // Some knowledge is lost forever
            for knowledgeType, level := range resident.Knowledge {
                preservedLevel := level * knowledgePreserved
                village.DistributeKnowledge(knowledgeType, preservedLevel)
            }
        }
    }
}
```

### UI Interface Requirements

#### Social Class Dashboard
- **Class distribution pie chart:** Peasants, craftsmen, clergy, nobility
- **Social mobility tracker:** Residents who changed class, by generation
- **Class tensions heatmap:** Areas of conflict between classes
- **Wealth inequality graph:** Gini coefficient over time

#### Religion Interface
- **Faith map:** Geographic distribution of beliefs within village
- **Religious calendar:** Upcoming holidays and observances
- **Conversion tracking:** Who converted whom, when, and why
- **Religious satisfaction:** Faith levels by demographic (age, class, gender)

#### Crime & Justice Panel
- **Crime reports:** List of recent crimes with details
- **Justice history:** Past sentences and their outcomes
- **Deterrence effectiveness:** Crime rate vs. punishment severity
- **Prison management:** Current inmates, sentences, rehabilitation

#### Education & Knowledge Views
- **Literacy map:** Who can read/write (highlighted on village map)
- **Skill transmission graph:** Master-apprentice relationships
- **Knowledge preservation:** Documented vs. oral tradition tracking
- **Library contents:** Index of written knowledge in village

#### Cultural Traditions Tracker
- **Village identity panel:** Unique traditions, founding story, symbols
- **Tradition evolution timeline:** How customs have changed over generations
- **Cultural exchange log:** Ideas adopted from visitors/traders
- **Material culture showcase:** Distinctive architecture, clothing, crafts

#### Political Power Interface
- **Faction alignment chart:** Where residents stand on key issues
- **Leadership legitimacy meter:** Player's authority level
- **Succession planning:** Potential heirs with support levels
- **External pressures:** Demands from lord, church, king

---

## Integration Points

### With Resident Management System
- Social class affects skill development opportunities
- Religious beliefs modify need calculations (faith as comfort)
- Education level affects learning speed and teaching ability
- Political views influence work motivation and compliance

### With Resource Economy System
- Social class determines resource access (luxury goods for nobility)
- Religious requirements consume resources (church maintenance, festival costs)
- Crime affects resource security (theft reduces stores)
- Education requires resource investment (books, teacher time)

### With Seasonal & Event System
- Religious holidays integrated into seasonal calendar
- Crime waves triggered by economic stress events
- Cultural traditions evolve through seasonal repetition
- Political crises triggered by external events (war, famine)

### With Building System
- Class-specific housing (peasant huts vs. manor houses)
- Religious buildings (church, shrine, monastery)
- Justice buildings (gaol, stocks, gallows)
- Education buildings (school, library, guild hall)

---

## Balancing & Tuning

### Class Distribution (Realistic Medieval)
- **Peasants:** 75-85% of population
- **Craftsmen:** 10-15% of population  
- **Clergy:** 2-5% of population
- **Nobility:** 1-3% of population

### Religious Belief Distribution (England c. 1350)
- **Catholic:** 85-90% (nominal)
- **Devout Catholic:** 30-40% (regular practice)
- **Heretical groups:** 5-10% (Lollards, etc.)
- **Pagan survivals:** 5-10% (folk practices)
- **Atheist/skeptical:** 1-2% (rare)

### Crime Rates (Estimated Medieval)
- **Minor theft:** 2-5% chance per poor resident per year
- **Violent crime:** 0.5-1% chance per aggressive resident per year
- **Murder:** 0.1-0.2% chance per resident per year
- **Heresy accusations:** 0.5% chance during religious tension periods

### Literacy Rates (Medieval Europe)
- **Nobility:** 30-50%
- **Clergy:** 70-90%
- **Craftsmen:** 10-20%
- **Peasants:** 1-5%
- **Women:** Half of male rate for same class

---

## Performance Considerations

### Optimization Strategies
- **Class calculations:** Batch update weekly, not daily
- **Religion updates:** Only for residents with changing faith
- **Crime detection:** Probabilistic sampling for large populations
- **Knowledge tracking:** Aggregate similar knowledge, prune old data

### Scaling Targets
| Population | Social Update Time | Memory per Resident |
|------------|-------------------|---------------------|
| 50 residents | < 50ms | ~1 KB social data |
| 100 residents | < 100ms | ~1 KB social data |
| 200 residents | < 200ms | ~1 KB social data |

### Memory Management
- Old crime records archived after 10 years
- Detailed relationship history aggregated annually
- Cultural tradition details loaded on demand
- Political faction data compressed when inactive

---

## Testing Strategy

### Unit Tests
- Class mobility calculations under different conditions
- Religious faith changes from events and practices
- Crime generation probabilities based on resident state
- Knowledge transmission efficiency

### Integration Tests
- Complete social lifecycle (birth → education → work → death)
- Religious conflict scenarios (heresy accusations, conversions)
- Crime waves and justice system responses
- Cultural evolution over multiple generations

### Playtesting Focus
- Social dynamics feeling organic, not mechanical
- Religious beliefs affecting gameplay meaningfully
- Crime and justice creating moral dilemmas
- Education system providing long-term benefits
- Political challenges to player authority feeling fair

---

## Risks & Mitigations

### Technical Risks
1. **Social network calculations becoming exponentially complex**  
   **Mitigation:** Limit active relationships per resident, use efficient graph algorithms

2. **Save file bloat from detailed social histories**  
   **Mitigation:** Aggregate social data, optional detailed history, compression

3. **AI behavior complexity causing performance issues**  
   **Mitigation:** Behavior trees with caching, simplified decision making for minor NPCs

### Design Risks
1. **Social systems feeling like spreadsheet management**  
   **Mitigation:** Personal stories, unique resident personalities, emotional connections

2. **Historical accuracy making game depressing (rigid class system)**  
   **Mitigation:** Social mobility opportunities, player ability to challenge norms

3. **Religious content offending players**  
   **Mitigation:** Historical context explanations, optional content filters, sensitivity review

4. **Crime system encouraging punitive playstyles**  
   **Mitigation:** Rehabilitation options, systemic solutions to crime causes, moral complexity

### Mitigation Strategies
- **Progressive implementation:** Basic class system first, complex politics later
- **Player choice emphasis:** Multiple approaches to social management
- **Historical disclaimer:** Clear separation of game mechanics from endorsement
- **Community feedback:** Early access for historical accuracy review

---

## Success Metrics

### Technical Metrics
- [ ] Social system updates <100ms for 100 residents
- [ ] Crime detection accuracy >90% (correct perpetrator identification)
- [ ] Knowledge preservation >70% when experts die with apprentices
- [ ] Memory usage <50MB for complete social data

### Player Experience Metrics
- [ ] 70% of players can name village's dominant religion and social tensions
- [ ] Players report emotional investment in resident social lives
- [ ] Social management feels like integral part of gameplay, not tacked on
- [ ] Historical authenticity rated >4/5 by history enthusiast players

### Gameplay Metrics
- [ ] Multiple viable social strategies (authoritarian, egalitarian, theocratic)
- [ ] Social dynamics create emergent stories players want to share
- [ ] Education system visibly improves village capabilities over generations
- [ ] Crime rates respond to player policies (deterrence vs. prevention)

---

## Dependencies

### Required First
- Resident relationship system (basic)
- Event system with social hooks
- UI framework for social interfaces

### Dependent Features
- Advanced politics (requires class system)
- Religious reforms (requires religion system)
- Legal code development (requires crime system)
- University/scholarship (requires education system)

---

## Open Questions

1. **Historical sensitivity:** How to handle religion without offending modern believers?
2. **Social mobility:** How much upward mobility should be possible in medieval setting?
3. **Crime punishment:** Include historically accurate but harsh punishments?
4. **Education access:** Balance historical accuracy (limited education) with player agency?

---

*Feature Version 1.0 · Owner: Social Systems Team · Estimated Effort: 5-7 sprints*