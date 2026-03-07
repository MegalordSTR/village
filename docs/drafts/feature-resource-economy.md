# Resource Economy System

**Vision:** Create a historically-plausible medieval economy where 15 interconnected resources flow through production chains, with scarcity, spoilage, and craftsmanship affecting village survival and prosperity.

**Status:** Draft  
**Created:** March 7, 2026  
**Priority:** P0 (MVP Core)

---

## Executive Summary

The resource economy system simulates the production, consumption, and exchange of goods in a medieval village. Starting with 15 core resources arranged in multi-step production chains (raw → processed → advanced), the system creates meaningful trade-offs where players must balance immediate needs against long-term investment. Resources spoil, quality matters, and stockpiling strategies become essential for surviving winter and unexpected events.

**Success Definition:** Players feel the weight of economic decisions—choosing between planting more wheat vs. building a mill, balancing tool production against iron mining, and preparing for winter scarcity.

---

## User Stories

### US-001: Resource Production Chains
**As:** A player managing village industry  
**I want:** To see how raw materials become finished goods through multiple steps  
**So that:** I can plan efficient production pipelines

**Acceptance Criteria:**
- [ ] Clear visualization of production chains (wheat → flour → bread)
- [ ] Each production step shows: input resources, worker requirements, time, output
- [ ] Unlockable production buildings (mill, smithy, bakery) with prerequisites
- [ ] Production efficiency affected by worker skill, tool quality, building level
- [ ] Resource quality system (poor/normal/good/excellent) affecting output

### US-002: Inventory Management
**As:** A player stocking resources for winter  
**I want:** To manage storage capacity and track resource quantities  
**So that:** I can prevent spoilage and ensure adequate supplies

**Acceptance Criteria:**
- [ ] Central inventory view showing all resources with quantities
- [ ] Storage buildings with type-specific capacities (granary for food, warehouse for goods)
- [ ] Spoilage system: food decays over time, faster in poor conditions
- [ ] Stock alerts when resources fall below threshold
- [ ] Transfer system to move resources between storage locations

### US-003: Seasonal Economic Cycle
**As:** A player experiencing medieval agriculture  
**I want:** Production to follow realistic seasonal patterns  
**So that:** I must plan ahead for harvests and winter scarcity

**Acceptance Criteria:**
- [ ] Planting season (spring), growth (summer), harvest (autumn)
- [ ] Weather affects crop yields (drought reduces, rain increases)
- [ ] Winter: no agriculture, increased fuel consumption for heating
- [ ] Seasonal resource price fluctuations (food cheap after harvest, expensive in spring)
- [ ] Storage management critical for surviving winter

### US-004: Craftsmanship & Quality
**As:** A player investing in skilled craftsmen  
**I want:** Higher quality tools and goods to provide tangible benefits  
**So that:** I'm motivated to develop specialist workers

**Acceptance Criteria:**
- [ ] Resource quality tiers: poor, normal, good, excellent, masterwork
- [ ] Quality affects: tool durability, production speed, product value
- [ ] Craftsman skill determines achievable quality level
- [ ] Special materials (fine iron, seasoned wood) enable higher quality
- [ ] Quality visible in UI with icons and tooltips

### US-005: Resource Scarcity & Substitution
**As:** A player facing material shortages  
**I want:** To find alternative solutions when preferred resources are unavailable  
**So that:** I can adapt to changing conditions

**Acceptance Criteria:**
- [ ] Multiple paths to similar outcomes (stone vs. wood buildings)
- [ ] Resource substitution with efficiency penalties
- [ ] Scavenging system to find limited resources in environment
- [ ] Trade-off between using scarce resources now vs. saving for future
- [ ] Resource exhaustion (mines deplete, forests regrow slowly)

---

## Technical Specifications

### Resource Definition
```go
// Core resource types for MVP
type ResourceType string

const (
    // Raw materials (6)
    ResourceGrain      ResourceType = "grain"
    ResourceVegetables ResourceType = "vegetables"
    ResourceWood       ResourceType = "wood"
    ResourceStone      ResourceType = "stone"
    ResourceIronOre    ResourceType = "iron_ore"
    ResourceWool       ResourceType = "wool"
    
    // Processed goods (5)
    ResourceFlour      ResourceType = "flour"
    ResourceBread      ResourceType = "bread"
    ResourcePlanks     ResourceType = "planks"
    ResourceIron       ResourceType = "iron"
    ResourceCloth      ResourceType = "cloth"
    
    // Advanced goods (4)
    ResourceTools      ResourceType = "tools"
    ResourceFurniture  ResourceType = "furniture"
    ResourceWeapons    ResourceType = "weapons"
    ResourceClothing   ResourceType = "clothing"
)

type Resource struct {
    Type        ResourceType `json:"type"`
    Quantity    float64      `json:"quantity"`    // Units (kg, pieces, etc.)
    Quality     QualityTier  `json:"quality"`     // 0-5 (poor to masterwork)
    Location    string       `json:"location"`    // Building ID where stored
    Produced    GameDate     `json:"produced"`    // Production date (for spoilage)
    Value       float64      `json:"value"`       // Base value in abstract units
}
```

### Production Chain Definitions
```go
// Production recipe for transforming inputs to outputs
type Recipe struct {
    ID          string                     `json:"id"`
    Name        string                     `json:"name"`        // "Bake Bread"
    Building    BuildingType               `json:"building"`    // Required workplace
    Skill       SkillType                  `json:"skill"`       // Required skill
    
    Inputs      []ResourceRequirement     `json:"inputs"`      // Required resources
    Outputs     []ResourceOutput          `json:"outputs"`     // Produced resources
    
    Time        int                        `json:"time"`        // Days to complete
    BaseYield   float64                    `json:"baseYield"`   // Base units per week
    MaxWorkers  int                        `json:"maxWorkers"`  // Workers per building
}

// Example recipe: Flour to Bread
var BakeBreadRecipe = Recipe{
    ID: "bake_bread",
    Name: "Bake Bread",
    Building: BuildingBakery,
    Skill: SkillBaking,
    Inputs: []ResourceRequirement{
        {Type: ResourceFlour, Quantity: 5, Quality: QualityNormal},
        {Type: ResourceWood, Quantity: 1, Quality: AnyQuality}, // Fuel
    },
    Outputs: []ResourceOutput{
        {Type: ResourceBread, Quantity: 10, QualityInheritance: 0.8}, // 80% of input quality
    },
    Time: 2,        // 2 days per batch
    BaseYield: 35,  // 35 bread per week with 1 skilled worker
    MaxWorkers: 3,  // Up to 3 bakers per bakery
}
```

### Production Chains (MVP)

#### Food Chain
```
Field Work → Grain → Mill → Flour → Bakery → Bread
                    ↓
              Vegetables (direct consumption)
```

#### Wood Chain
```
Forestry → Wood → Sawmill → Planks → Carpenter → Furniture
                                    ↓
                              Tools (with iron)
```

#### Metal Chain  
```
Mining → Iron Ore → Smelter → Iron → Smithy → Tools/Weapons
```

#### Textile Chain
```
Sheep → Wool → Spinner → Yarn → Weaver → Cloth → Tailor → Clothing
```

### Production Calculation Algorithm
```go
func CalculateProduction(recipe Recipe, workers []*Resident, building *Building) ProductionResult {
    result := ProductionResult{}
    
    // Check input availability
    for _, input := range recipe.Inputs {
        available := inventory.GetAvailable(input.Type, input.Quality)
        if available < input.Quantity {
            // Partial production possible
            productionRatio := available / input.Quantity
            result.Success = false
            result.Reason = fmt.Sprintf("Insufficient %s", input.Type)
            result.PartialRatio = productionRatio
            return result
        }
    }
    
    // Calculate total yield based on workers
    totalYield := 0.0
    qualityMultiplier := 1.0
    
    for _, worker := range workers {
        skill := worker.Skills[recipe.Skill]
        efficiency := skill.Level / 100.0
        
        // Quality affects output
        if skill.Level > 75 {
            qualityMultiplier = math.Max(qualityMultiplier, 1.2)
        }
        
        // Needs affect productivity
        needsPenalty := worker.CalculateNeedsPenalty()
        efficiency *= needsPenalty
        
        // Building condition affects work
        buildingBonus := building.Condition / 100.0
        
        workerYield := recipe.BaseYield * efficiency * buildingBonus
        totalYield += workerYield
    }
    
    // Apply weekly time scaling (recipe.Time is days per batch)
    batchesPerWeek := 7.0 / float64(recipe.Time)
    totalYield *= batchesPerWeek
    
    // Consume inputs
    for _, input := range recipe.Inputs {
        inventory.Consume(input.Type, input.Quantity, input.Quality)
    }
    
    // Produce outputs
    for _, output := range recipe.Outputs {
        finalQuantity := totalYield * output.Quantity
        finalQuality := CalculateOutputQuality(output, workers)
        
        inventory.Add(output.Type, finalQuantity, finalQuality)
        result.Produced = append(result.Produced, Resource{
            Type: output.Type,
            Quantity: finalQuantity,
            Quality: finalQuality,
        })
    }
    
    result.Success = true
    result.TotalYield = totalYield
    return result
}
```

### Spoilage System
```go
func UpdateSpoilage(resources []Resource, week int, storageConditions map[string]float64) {
    for i := range resources {
        r := &resources[i]
        
        // Only certain resources spoil
        if !r.Type.Spoils() {
            continue
        }
        
        // Calculate spoilage rate
        baseRate := r.Type.SpoilageRate()  // e.g., 0.01 per week for grain
        storageModifier := storageConditions[r.Location]  // 0.5 for good storage, 2.0 for poor
        
        // Quality affects spoilage
        qualityModifier := 1.0
        if r.Quality >= QualityGood {
            qualityModifier = 0.7  // Better preservation
        }
        
        // Age of resource
        weeksOld := (week - r.Produced.Week)
        if weeksOld > 52 {
            qualityModifier *= 1.5  // Accelerated spoilage for very old items
        }
        
        // Calculate loss
        lossPercent := baseRate * storageModifier * qualityModifier
        lossAmount := r.Quantity * lossPercent
        
        r.Quantity -= lossAmount
        
        // Update quality if significant loss
        if lossPercent > 0.1 {
            r.Quality = max(r.Quality-1, QualityPoor)
        }
        
        // Remove if quantity negligible
        if r.Quantity < 0.01 {
            r.Quantity = 0
        }
    }
}
```

### UI Interface Requirements

#### Production Chain Visualization
- Flowchart showing resource transformations
- Current bottlenecks highlighted (insufficient inputs)
- Production rates displayed at each stage
- Unlock requirements for advanced production

#### Inventory Dashboard
- Grid view of all resources with icons
- Sort by: quantity, value, spoilage rate
- Filter by: category (food, material, goods)
- Stock level indicators (low/medium/high)
- Spoilage warnings for perishables

#### Production Management
- Building view showing assigned workers
- Input/output requirements for selected recipe
- Production queue for multi-step items
- Efficiency metrics (output per worker hour)
- Quality distribution of produced goods

#### Economic Overview
- Resource flow diagram (inputs → outputs)
- Consumption rates vs. production rates
- Stockpile trends over time (graphs)
- Seasonal preparation indicators

---

## Integration Points

### With Resident Management System
- Worker skills affect production quality and quantity
- Resident needs affect work efficiency
- Assignment system places workers in production buildings

### With Building System
- Buildings enable specific production recipes
- Building condition affects production efficiency
- Storage buildings determine capacity and spoilage rates

### With Seasonal System
- Agriculture only possible in growing season
- Weather affects crop yields
- Winter increases fuel consumption for heating
- Seasonal events affect resource availability

### With Event System
- Resource discovery events (new mine, abundant harvest)
- Loss events (fire, spoilage, theft)
- Trade events (caravan offers/exchange)
- Crafting breakthroughs (quality improvement)

---

## Balancing & Tuning

### Production Rates (per skilled worker per week)
| Resource | Base Rate | Notes |
|----------|-----------|-------|
| Grain | 50 kg | Varies by field size |
| Wood | 20 units | Depletes forest over time |
| Stone | 15 units | Mine depletion |
| Iron Ore | 10 units | Rare, slow to mine |
| Flour | 40 kg | From 50 kg grain |
| Bread | 35 loaves | From 5 kg flour |
| Planks | 25 units | From 30 wood |
| Iron | 8 units | From 15 iron ore + fuel |
| Tools | 5 units | From 2 iron + 3 planks |

### Consumption Rates (per resident per week)
- **Food:** 5 kg grain equivalent (2.5 kg bread)
- **Fuel:** 3 units wood (winter: 8 units)
- **Tool wear:** 0.1 tools (craftsmen: 0.2 tools)

### Storage & Spoilage
- **Granary:** 1000 kg food, spoilage rate 0.5% weekly
- **Warehouse:** 500 units goods, spoilage rate 0.1% weekly
- **Outdoor pile:** 200 units, spoilage rate 2.0% weekly
- **Food spoilage multipliers:** Grain 1x, Vegetables 2x, Meat 3x, Bread 0.5x

---

## Performance Considerations

### Optimization Strategies
- **Batch processing:** Update all resources in single pass
- **Lazy evaluation:** Only calculate spoilage for perishables
- **Delta compression:** Send only changed quantities to frontend
- **Caching:** Production calculations cached until inputs change

### Scaling Targets
| Resource Types | Update Time | Memory Usage |
|----------------|-------------|--------------|
| 15 resources   | < 10ms      | < 1 MB       |
| 50 resources   | < 30ms      | < 3 MB       |
| 200 resources  | < 100ms     | < 10 MB      |

### Data Structures
- **Resource map:** O(1) lookup by type and quality
- **Production graph:** Adjacency list for chain traversal
- **Inventory index:** Spatial hash for location-based queries
- **Spoilage queue:** Priority queue by expiration date

---

## Testing Strategy

### Unit Tests
- Production calculations with varying inputs
- Spoilage rates under different conditions
- Quality inheritance from inputs to outputs
- Resource substitution logic

### Integration Tests
- Complete production chains from raw to finished
- Seasonal agriculture cycle
- Storage capacity and overflow handling
- Economic collapse scenarios (resource exhaustion)

### Playtesting Focus
- Intuitive understanding of production chains
- Satisfaction with production planning
- Tension around resource scarcity
- Reward from efficient economic management

---

## Risks & Mitigations

### Technical Risks
1. **Complex production graphs causing performance issues**  
   **Mitigation:** Limit chain depth, optimize graph traversal, use caching

2. **Floating-point errors in resource quantities**  
   **Mitigation:** Fixed-point arithmetic, tolerance thresholds, validation

3. **Save file size from detailed resource tracking**  
   **Mitigation:** Aggregate identical resources, compress history, optional detail

### Design Risks
1. **Overwhelming complexity for new players**  
   **Mitigation:** Progressive disclosure, tutorial chains, simplified view

2. **Optimal strategies becoming obvious/repetitive**  
   **Mitigation:** Random events, multiple viable strategies, evolving challenges

3. **Historical accuracy conflicting with fun gameplay**  
   **Mitigation:** Configurable realism settings, optional hardship modes

### Mitigation Strategies
- **Early economic testing:** Playtest economy before adding combat/social
- **Modular difficulty:** Separate tuning for production rates, spoilage, scarcity
- **Player feedback loops:** Clear cause-effect for economic decisions
- **Alternative playstyles:** Support different economic strategies (trade, self-sufficiency)

---

## Success Metrics

### Technical Metrics
- [ ] Production calculations <5ms for 50 active recipes
- [ ] Spoilage updates <2ms for 1000 resource items
- [ ] Memory usage <5MB for complete economy state
- [ ] Save/load economy state in <100ms

### Player Experience Metrics
- [ ] 80% of players understand basic production chains within 1 hour
- [ ] Players report meaningful trade-offs in resource allocation
- [ ] Economic preparation for winter feels necessary and rewarding
- [ ] Quality system provides visible benefits worth pursuing

### Gameplay Metrics
- [ ] Multiple viable economic strategies emerge in playtesting
- [ ] Resource scarcity creates interesting decisions, not frustration
- [ ] Production optimization provides satisfying progression
- [ ] Economic collapse possible but avoidable with good management

---

## Dependencies

### Required First
- Basic resource type definitions
- Inventory storage system
- Building placement system

### Dependent Features
- Crafting quality system (requires resource quality)
- Trade system (requires resource valuation)
- Event system (requires resource availability checks)
- Technology tree (requires production building unlocks)

---

## Open Questions

1. **Abstraction level:** Should resources be measured in abstract units or realistic weights/volumes?
2. **Transportation:** Should moving resources between buildings require time/labor?
3. **Market dynamics:** Should prices fluctuate based on village supply/demand?
4. **Specialization:** Should villages develop economic specialties (mining vs. farming)?

---

*Feature Version 1.0 · Owner: Game Economy Team · Estimated Effort: 5-7 sprints*