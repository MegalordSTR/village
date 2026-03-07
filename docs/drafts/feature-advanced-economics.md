# Advanced Economics System

**Vision:** Transform the village from a self-sufficient subsistence economy into a complex medieval market with trade networks, currency, price fluctuations, banking, and economic crises that challenge even experienced players.

**Status:** Draft  
**Created:** March 7, 2026  
**Priority:** P1 (Post-MVP Depth)

---

## Executive Summary

The advanced economics system adds layers of complexity to the MVP's basic resource management: villages can specialize and trade, currency replaces barter, prices fluctuate based on supply and demand, banking enables investment and debt, and economic bubbles and crashes create new challenges. Players must now understand not just production chains but market dynamics, comparative advantage, inflation, and financial risk.

**Success Definition:** Players discuss economic strategies in market terms ("I'm short on grain because prices are low this season, but I've invested in iron futures") and feel they're managing a dynamic economy, not just a production spreadsheet.

---

## User Stories

### US-001: Trade Routes & Caravans
**As:** A player in a village with surplus goods  
**I want:** To trade with other villages via traveling merchants  
**So that:** I can acquire resources we cannot produce locally

**Acceptance Criteria:**
- [ ] Trade routes to 3-5 neighboring villages (simulated, not directly controlled)
- [ ] Caravan system: merchants arrive seasonally with goods for trade
- [ ] Trade goods availability varies by origin village specialization
- [ ] Transportation costs: distance affects trade profitability
- [ ] Trade route security: bandits can raid caravans, requiring guards

### US-002: Currency & Monetary System
**As:** A player tired of complex barter calculations  
**I want:** A currency system to simplify trade and valuation  
**So that:** I can think in prices rather than resource-for-resource exchanges

**Acceptance Criteria:**
- [ ] Coinage system: pennies, shillings, pounds (£sd system)
- [ ] Minting: convert precious metals (silver, gold) into coins
- [ ] Prices in currency for all goods and services
- [ ] Wage system: pay residents in coins rather than just food/shelter
- [ ] Currency debasement: lowering silver content causes inflation

### US-003: Market Price Dynamics
**As:** A player participating in medieval markets  
**I want:** Prices to fluctuate based on supply, demand, and events  
**So that:** I can profit from good timing and suffer from bad timing

**Acceptance Criteria:**
- [ ] Local price system: village supply/demand affects prices
- [ ] Regional price differences: goods cheaper where produced
- [ ] Seasonal price cycles: grain cheap after harvest, expensive in spring
- [ ] Event-driven price spikes: plague increases medicine prices
- [ ] Price memory: past prices affect current price expectations

### US-004: Banking & Credit System
**As:** A player needing capital for large projects  
**I want:** To borrow money and pay interest  
**So that:** I can invest in infrastructure without immediate resources

**Acceptance Criteria:**
- [ ] Moneylenders: residents with excess wealth lend at interest
- [ ] Loans: borrow coins for investments, repay with interest
- [ ] Collateral: assets pledged against loans (buildings, future harvests)
- [ ] Default risk: failed repayments lead to asset seizure
- [ ] Compound interest: debts grow if not serviced

### US-005: Economic Specialization
**As:** A player optimizing village economy  
**I want:** To specialize in producing goods where we have advantage  
**So that:** We can trade for greater overall wealth

**Acceptance Criteria:**
- [ ] Comparative advantage: villages naturally better at certain productions
- [ ] Specialization bonuses: focused production increases efficiency
- [ ] Trade dependency: specialized villages vulnerable to supply disruptions
- [ ] Economic diversity slider: choose between self-sufficiency vs. specialization
- [ ] Knowledge specialization: unique recipes/techniques develop in specialized villages

### US-006: Economic Crises & Recovery
**As:** A player managing long-term village prosperity  
**I want:** Economic challenges that test my management skills  
**So that:** Economic success feels earned, not automatic

**Acceptance Criteria:**
- [ ] Inflation: too much currency chasing too few goods
- [ ] Deflation: currency shortage causing economic paralysis
- [ ] Bubbles: speculative frenzies in certain goods (tulip mania equivalent)
- [ ] Crashes: sudden price collapses with cascading failures
- [ ] Depression/recovery cycles: multi-year economic patterns

---

## Technical Specifications

### Trade System
```go
type TradeRoute struct {
    ID            string        `json:"id"`
    FromVillage   string        `json:"fromVillage"`   // Always "our" village
    ToVillage     string        `json:"toVillage"`     // Simulated external village
    Distance      float64       `json:"distance"`      // Days travel
    Security      float64       `json:"security"`      // 0-100, affects bandit risk
    Established   int           `json:"established"`   // Week route was established
    Volume        float64       `json:"volume"`        // Trade volume modifier
}

type Caravan struct {
    ID            string        `json:"id"`
    Route         string        `json:"route"`         // TradeRoute ID
    ArrivalWeek   int           `json:"arrivalWeek"`   // Expected arrival
    Status        CaravanStatus `json:"status"`        // Traveling, Arrived, Raided
    Goods         []TradeGood   `json:"goods"`         // What they're carrying
    Merchant      Merchant      `json:"merchant"`      // NPC merchant data
    GuardStrength float64       `json:"guardStrength"` // 0-100
}

type TradeGood struct {
    ResourceType  ResourceType  `json:"resourceType"`
    Quantity      float64       `json:"quantity"`
    AskingPrice   Price         `json:"askingPrice"`   // In coins or barter
    WillBuy       bool          `json:"willBuy"`       // Merchant wants to buy
    WillSell      bool          `json:"willSell"`      // Merchant wants to sell
}

// Trade negotiation
func NegotiateTrade(playerOffer []TradeGood, merchant *Merchant, village *Village) TradeResult {
    result := TradeResult{}
    
    // Calculate value of player's offer
    playerValue := 0.0
    for _, good := range playerOffer {
        price := village.Market.GetPrice(good.ResourceType)
        playerValue += good.Quantity * price
    }
    
    // Calculate value of merchant's offer  
    merchantValue := 0.0
    for _, good := range merchant.Goods {
        if good.WillSell {
            // Merchant's asking price
            merchantValue += good.Quantity * good.AskingPrice.Amount
        }
    }
    
    // Negotiation logic
    exchangeRatio := playerValue / merchantValue
    result.FairTrade = exchangeRatio >= 0.8 && exchangeRatio <= 1.2
    
    // Merchant personality affects negotiation
    switch merchant.Personality {
    case "Generous":
        result.FairTrade = exchangeRatio >= 0.7  // More lenient
    case "Greedy":
        result.FairTrade = exchangeRatio >= 0.9  // Stricter
    }
    
    // Reputation affects success
    reputationBonus := village.Reputation[merchant.ID] / 100.0
    result.Success = result.FairTrade || rng.Float64() < reputationBonus
    
    if result.Success {
        // Execute trade
        village.Resources.Consume(playerOffer)
        village.Resources.Add(merchant.Goods)
        
        // Update reputation
        if exchangeRatio >= 1.0 {
            village.Reputation[merchant.ID] = min(
                village.Reputation[merchant.ID] + 5, 
                100,
            )
        }
    }
    
    return result
}
```

### Currency & Price System
```go
type Currency struct {
    Name          string        `json:"name"`          // "Penny", "Shilling", "Pound"
    Metal         MetalType     `json:"metal"`         // Silver, Gold, Copper
    Purity        float64       `json:"purity"`        // 0-1, e.g., 0.925 for sterling
    Weight        float64       `json:"weight"`        // Grams
    ExchangeRate  float64       `json:"exchangeRate"`  // To base currency (pennies)
}

// Price in currency
type Price struct {
    Amount        float64       `json:"amount"`        // In pennies
    Currency      string        `json:"currency"`      // "pennies"
    ValidFrom     int           `json:"validFrom"`     // Week price set
    ValidTo       int           `json:"validTo"`       // Week price expires (0 = until changed)
}

type Market struct {
    VillageID     string                 `json:"villageId"`
    Prices        map[ResourceType]Price `json:"prices"`
    PriceHistory  map[ResourceType][]PricePoint `json:"priceHistory"`
    Supply        map[ResourceType]float64 `json:"supply"`      // Local supply
    Demand        map[ResourceType]float64 `json:"demand"`      // Local demand
    BasePrices    map[ResourceType]float64 `json:"basePrices"`  // Base value in pennies
}

// Update market prices weekly
func (m *Market) UpdatePrices(week int, village *Village) {
    for resourceType := range m.Prices {
        // Calculate new price
        newPrice := m.CalculatePrice(resourceType, village)
        
        // Record in history
        m.PriceHistory[resourceType] = append(
            m.PriceHistory[resourceType],
            PricePoint{
                Week: week,
                Price: newPrice,
                Supply: m.Supply[resourceType],
                Demand: m.Demand[resourceType],
            },
        )
        
        // Keep only recent history (52 weeks)
        if len(m.PriceHistory[resourceType]) > 52 {
            m.PriceHistory[resourceType] = m.PriceHistory[resourceType][-52:]
        }
        
        // Update current price
        m.Prices[resourceType] = Price{
            Amount: newPrice,
            Currency: "pennies",
            ValidFrom: week,
        }
    }
}

func (m *Market) CalculatePrice(resourceType ResourceType, village *Village) float64 {
    basePrice := m.BasePrices[resourceType]
    
    // Supply/demand effect
    supply := m.Supply[resourceType]
    demand := m.Demand[resourceType]
    
    if supply <= 0 {
        supply = 0.001 // Avoid division by zero
    }
    
    supplyDemandRatio := demand / supply
    priceMultiplier := math.Pow(supplyDemandRatio, 0.5) // Square root for diminishing effect
    
    // Seasonal effects
    seasonMultiplier := 1.0
    if resourceType.IsFood() {
        seasonMultiplier = village.Calendar.GetFoodPriceMultiplier()
    }
    
    // Event effects
    eventMultiplier := 1.0
    for _, event := range village.RecentEvents {
        if event.AffectsResource(resourceType) {
            eventMultiplier *= event.PriceEffect
        }
    }
    
    // Calculate final price
    price := basePrice * priceMultiplier * seasonMultiplier * eventMultiplier
    
    // Price stickiness: prices don't change too quickly
    oldPrice := m.Prices[resourceType].Amount
    maxChange := oldPrice * 0.1 // Max 10% change per week
    price = clamp(price, oldPrice-maxChange, oldPrice+maxChange)
    
    // Minimum price (production cost)
    minPrice := m.CalculateProductionCost(resourceType)
    price = math.Max(price, minPrice * 0.8) // Can fall to 80% of production cost
    
    return price
}
```

### Banking & Credit System
```go
type Loan struct {
    ID            string        `json:"id"`
    Borrower      string        `json:"borrower"`      // Village ID or resident ID
    Lender        string        `json:"lender"`        // Resident ID (moneylender)
    Principal     float64       `json:"principal"`     // Amount borrowed (pennies)
    InterestRate  float64       `json:"interestRate"`  // Annual percentage (e.g., 0.1 for 10%)
    Term          int           `json:"term"`          // Weeks until due
    StartWeek     int           `json:"startWeek"`
    Collateral    []Collateral  `json:"collateral"`    // Assets pledged
    Payments      []Payment     `json:"payments"`      // Payment history
    Status        LoanStatus    `json:"status"`        // Active, Paid, Defaulted
}

type Bank struct {
    VillageID     string        `json:"villageId"`
    Deposits      float64       `json:"deposits"`      // Coins deposited by residents
    Reserves      float64       `json:"reserves"`      // Coins kept in reserve
    Loans         []Loan        `json:"loans"`         // Active loans
    InterestRate  float64       `json:"interestRate"`  // Deposit interest rate
    ReserveRatio  float64       `json:"reserveRatio"`  // Required reserves (0-1)
}

// Weekly banking update
func (b *Bank) Update(week int, village *Village) {
    // Calculate interest on deposits
    weeklyInterestRate := b.InterestRate / 52.0
    interestEarned := b.Deposits * weeklyInterestRate
    b.Deposits += interestEarned
    
    // Update loans
    for i := range b.Loans {
        loan := &b.Loans[i]
        
        if loan.Status != Active {
            continue
        }
        
        // Calculate interest due
        weeksActive := week - loan.StartWeek
        annualInterest := loan.Principal * loan.InterestRate
        weeklyInterest := annualInterest / 52.0
        totalInterestDue := weeklyInterest * float64(weeksActive)
        
        // Check for default
        if weeksActive > loan.Term {
            // Loan is overdue
            b.HandleDefault(loan, village)
            continue
        }
        
        // Automatic payment if borrower has funds
        if loan.Borrower == village.ID {
            // Village loan
            weeklyPayment := loan.Principal / float64(loan.Term) + weeklyInterest
            
            if village.Coins >= weeklyPayment {
                village.Coins -= weeklyPayment
                loan.Payments = append(loan.Payments, Payment{
                    Week: week,
                    Amount: weeklyPayment,
                    Type: PaymentRegular,
                })
                
                loan.Principal -= weeklyPayment - weeklyInterest
                
                // Loan paid off
                if loan.Principal <= 0 {
                    loan.Status = Paid
                    village.Reputation[loan.Lender] = min(
                        village.Reputation[loan.Lender] + 10,
                        100,
                    )
                }
            }
        }
    }
    
    // Check reserve requirements
    requiredReserves := b.Deposits * b.ReserveRatio
    if b.Reserves < requiredReserves {
        // Reserve deficiency - may trigger crisis
        b.HandleReserveDeficiency(village)
    }
}
```

### Economic Crisis System
```go
type EconomicCrisis struct {
    Type          CrisisType    `json:"type"`
    Severity      float64       `json:"severity"`      // 0-100
    StartWeek     int           `json:"startWeek"`
    Duration      int           `json:"duration"`      // Weeks crisis lasts
    Causes        []string      `json:"causes"`        // What triggered it
    Effects       []CrisisEffect `json:"effects"`      // Active effects
    Resolution    CrisisResolution `json:"resolution"` // How it ended
}

type CrisisType int

const (
    CrisisInflation CrisisType = iota
    CrisisDeflation
    CrisisBankRun
    CrisisTradeCollapse
    CrisisBubbleBurst
    CrisisDepression
)

// Crisis detection and management
func CheckForCrises(village *Village, week int) []EconomicCrisis {
    var newCrises []EconomicCrisis
    
    // Inflation check
    inflationRate := village.CalculateInflationRate()
    if inflationRate > 0.1 { // 10% weekly inflation
        crisis := EconomicCrisis{
            Type: CrisisInflation,
            Severity: (inflationRate - 0.1) * 1000, // Scale to 0-100
            StartWeek: week,
            Duration: rng.Intn(26) + 26, // 26-52 weeks
            Causes: []string{
                "Excessive coin minting",
                "Poor harvest driving up food prices",
                "Trade disruption reducing supply",
            },
        }
        newCrises = append(newCrises, crisis)
    }
    
    // Bank run check
    if village.Bank.Reserves / village.Bank.Deposits < 0.1 {
        // Less than 10% reserves
        crisis := EconomicCrisis{
            Type: CrisisBankRun,
            Severity: (0.1 - village.Bank.Reserves/village.Bank.Deposits) * 1000,
            StartWeek: week,
            Duration: rng.Intn(13) + 13, // 13-26 weeks
            Causes: []string{
                "Loss of confidence in bank",
                "Rumors of insolvency",
                "External economic shock",
            },
        }
        newCrises = append(newCrises, crisis)
    }
    
    // Bubble detection (rapid price increase followed by correction)
    for resourceType, history := range village.Market.PriceHistory {
        if len(history) < 10 {
            continue
        }
        
        recentPrices := history[len(history)-10:]
        priceIncrease := (recentPrices[9].Price - recentPrices[0].Price) / recentPrices[0].Price
        
        if priceIncrease > 2.0 { // 200% increase in 10 weeks
            crisis := EconomicCrisis{
                Type: CrisisBubbleBurst,
                Severity: priceIncrease * 25, // Scale to 0-100
                StartWeek: week,
                Duration: rng.Intn(52) + 52, // 52-104 weeks
                Causes: []string{
                    fmt.Sprintf("Speculative frenzy in %s", resourceType),
                    "Irrational exuberance",
                    "Herding behavior among traders",
                },
            }
            newCrises = append(newCrises, crisis)
            break // Only one bubble at a time
        }
    }
    
    return newCrises
}

// Apply crisis effects
func ApplyCrisisEffects(village *Village, crisis EconomicCrisis) {
    switch crisis.Type {
    case CrisisInflation:
        // Prices rise faster
        for resourceType := range village.Market.Prices {
            currentPrice := village.Market.Prices[resourceType].Amount
            inflationEffect := crisis.Severity / 100.0 // 0-1 multiplier
            village.Market.Prices[resourceType].Amount = currentPrice * (1 + inflationEffect*0.1)
        }
        
        // Residents hoard goods, distrust currency
        for i := range village.Residents {
            resident := &village.Residents[i]
            if resident.Wealth > 1000 { // Wealthy residents
                // Convert coins to goods
                convertAmount := resident.Coins * 0.1
                resident.Coins -= convertAmount
                // Buy goods (simplified)
                village.Resources.ConsumeForResident(resident, convertAmount)
            }
        }
        
    case CrisisBankRun:
        // Depositors withdraw funds
        withdrawalRate := crisis.Severity / 100.0 * 0.05 // 0-5% per week
        withdrawals := village.Bank.Deposits * withdrawalRate
        
        // Bank may not have enough reserves
        if withdrawals > village.Bank.Reserves {
            // Bank failure
            village.Bank.Reserves = 0
            village.Bank.Deposits -= withdrawals
            
            // Only some depositors get their money
            recoveryRate := village.Bank.Reserves / withdrawals
            for i := range village.Residents {
                resident := &village.Residents[i]
                if resident.BankDeposit > 0 {
                    recovered := resident.BankDeposit * recoveryRate
                    resident.Coins += recovered
                    resident.BankDeposit -= recovered
                }
            }
        } else {
            // Bank survives run
            village.Bank.Reserves -= withdrawals
            village.Bank.Deposits -= withdrawals
        }
        
    case CrisisBubbleBurst:
        // Find the bubbled resource
        var bubbledResource ResourceType
        maxIncrease := 0.0
        
        for resourceType, history := range village.Market.PriceHistory {
            if len(history) < 10 {
                continue
            }
            increase := (history[len(history)-1].Price - history[len(history)-10].Price) / history[len(history)-10].Price
            if increase > maxIncrease {
                maxIncrease = increase
                bubbledResource = resourceType
            }
        }
        
        // Crash the price
        currentPrice := village.Market.Prices[bubbledResource].Amount
        crashSeverity := crisis.Severity / 100.0 // 0-1
        newPrice := currentPrice * (1 - crashSeverity*0.7) // Up to 70% crash
        
        village.Market.Prices[bubbledResource].Amount = newPrice
        
        // Investors lose wealth
        for i := range village.Residents {
            resident := &village.Residents[i]
            if resident.Investments[bubbledResource] > 0 {
                loss := resident.Investments[bubbledResource] * crashSeverity * 0.7
                resident.Wealth -= loss
                resident.Happiness = max(resident.Happiness - loss/100, 0)
            }
        }
    }
}
```

### UI Interface Requirements

#### Trade Interface
- **Caravan arrivals calendar:** Upcoming merchant visits
- **Trade route map:** Visual representation of routes with distances/risks
- **Negotiation screen:** Side-by-side comparison of offers with value calculations
- **Trade history:** Past trades with profit/loss calculations
- **Specialization advisor:** Suggests optimal trade goods based on comparative advantage

#### Market Dashboard
- **Price ticker:** Current prices for major goods with change indicators
- **Supply/demand graphs:** Historical data for key resources
- **Seasonal price calendar:** Forecast of expected price movements
- **Arbitrage opportunities:** Price differences between villages
- **Market depth:** Buy/sell orders at different price levels

#### Banking Interface
- **Account summary:** Deposits, loans, interest rates
- **Loan application:** Apply for loans with different terms/collateral
- **Investment portfolio:** Resident investments with current values
- **Risk assessment:** Village economic health indicators
- **Crisis management:** Tools for responding to economic crises

#### Economic Planning Tools
- **Specialization planner:** Compare self-sufficiency vs. trade strategies
- **Investment calculator:** ROI projections for different projects
- **Price forecasting:** Predict future prices based on trends
- **Economic model:** Simulate effects of policy changes
- **Wealth distribution:** Gini coefficient and poverty metrics

#### Crisis Management Interface
- **Early warning indicators:** Metrics approaching dangerous levels
- **Crisis response options:** Policy tools for different crisis types
- **Impact assessment:** Predicted effects of different responses
- **Recovery tracking:** Progress toward economic recovery
- **Lesson learned:** Analysis of what caused crisis and how to prevent recurrence

---

## Integration Points

### With Resource Economy System
- Prices affect production decisions (produce high-price goods)
- Currency enables wage labor instead of subsistence distribution
- Banking provides capital for production expansion
- Trade allows import of missing resources

### With Resident Management System
- Wages in coins affect resident wealth and class mobility
- Economic crises affect resident happiness and needs
- Investment opportunities for wealthy residents
- Debt can drive residents to crime or desperation

### With Social Systems
- Wealth inequality affects class tensions
- Economic specialization creates distinct social roles (merchants, bankers)
- Crises can trigger social unrest or innovation
- Education affects economic productivity and innovation

### With Seasonal & Event System
- Seasonal price cycles integrated with agricultural calendar
- Economic events (harvest failures, trade disruptions)
- Crisis events triggering or resulting from economic conditions
- Recovery events after crises

---

## Balancing & Tuning

### Price Ranges (Pennies per Unit)
| Resource | Base Price | Min Price | Max Price | Volatility |
|----------|------------|-----------|-----------|------------|
| Grain (kg) | 2 | 1 | 10 | High (seasonal) |
| Bread (loaf) | 1 | 0.5 | 4 | Medium |
| Iron (kg) | 20 | 15 | 60 | Low |
| Tools (each) | 50 | 30 | 150 | Medium |
| Land (acre) | 1000 | 500 | 5000 | Very Low |

### Trade Parameters
- **Caravan frequency:** 1-2 per season per trade route
- **Travel time:** 1-4 weeks depending on distance
- **Bandit risk:** 5-20% per trip, reducible with guards
- **Price margins:** Merchants buy at 70-90% of sell price
- **Trade volume limits:** Based on village size and route establishment

### Banking Parameters
- **Interest rates:** 5-20% annual for loans, 2-5% for deposits
- **Reserve requirement:** 10-20% of deposits
- **Loan terms:** 13-104 weeks (3 months to 2 years)
- **Default rates:** 5-15% depending on economic conditions
- **Wealth distribution:** Top 10% hold 40-60% of wealth (medieval realistic)

### Crisis Thresholds
- **Inflation crisis:** >10% weekly price increase
- **Deflation crisis:** >5% weekly price decrease for 4+ weeks
- **Bank run:** Reserve ratio <10%
- **Bubble:** >200% price increase in 10 weeks
- **Depression:** >50% price decrease across all goods

---

## Performance Considerations

### Optimization Strategies
- **Price updates:** Batch process all prices weekly
- **Trade calculations:** Cache comparative advantage calculations
- **Banking updates:** Only process active loans weekly
- **Crisis detection:** Sample-based rather than exhaustive checking

### Scaling Targets
| Economic Complexity | Update Time | Memory Usage |
|---------------------|-------------|--------------|
| Basic (MVP) | < 10ms | < 1 MB |
| Advanced (Trade) | < 50ms | < 5 MB |
| Full (Banking+Crises) | < 100ms | < 10 MB |

### Data Structures
- **Market prices:** Hash map for O(1) lookups
- **Price history:** Circular buffer for recent prices
- **Trade routes:** Graph structure for route optimization
- **Loan portfolio:** Priority queue by due date

---

## Testing Strategy

### Unit Tests
- Price calculation algorithms under various conditions
- Trade negotiation logic with different merchant personalities
- Interest calculation and compounding
- Crisis detection thresholds and triggers

### Integration Tests
- Complete trade cycle: production → pricing → trade → profit
- Banking lifecycle: deposit → loan → interest → repayment → withdrawal
- Economic crisis: trigger → effects → player response → recovery
- Long-term simulation: multiple economic cycles over decades

### Playtesting Focus
- Economic complexity feeling engaging, not overwhelming
- Trade decisions providing meaningful strategic choices
- Banking system usefulness for village development
- Crisis management creating dramatic tension without frustration
- Economic systems integration feeling cohesive

---

## Risks & Mitigations

### Technical Risks
1. **Floating-point economics causing rounding errors**  
   **Mitigation:** Fixed-point arithmetic for currency, tolerance thresholds

2. **Economic simulation becoming computationally expensive**  
   **Mitigation:** Simplified models for AI villages, caching, incremental updates

3. **Save file size explosion from economic history**  
   **Mitigation:** Aggregate economic data, optional detailed history, compression

4. **Economic systems creating unstable feedback loops**  
   **Mitigation:** Damping factors in calculations, bounds on extreme values

### Design Risks
1. **Economic complexity alienating casual players**  
   **Mitigation:** Progressive disclosure, simplified economic modes, automation options

2. **Historical accuracy conflicting with balanced gameplay**  
   **Mitigation:** Configurable realism settings, gameplay-balanced defaults

3. **Economic systems feeling disconnected from core gameplay**  
   **Mitigation:** Tight integration with resource production and resident needs

4. **Player frustration from economic crises outside their control**  
   **Mitigation:** Early warning systems, crisis preparation time, recovery tools

### Mitigation Strategies
- **Gradual implementation:** Basic trade first, then currency, then banking, then crises
- **Player choice emphasis:** Multiple economic strategies viable
- **Difficulty settings:** Adjustable economic complexity and crisis frequency
- **Community feedback:** Economic balance testing with strategy game enthusiasts

---

## Success Metrics

### Technical Metrics
- [ ] Economic updates <100ms per week for full system
- [ ] Price calculation accuracy >95% (matches expected formulas)
- [ ] Trade simulation stability (no infinite money glitches)
- [ ] Save/load economic state in <200ms

### Player Experience Metrics
- [ ] 70% of players engage with trade system within first 5 hours
- [ ] Economic complexity rated as "deep but understandable" by strategy gamers
- [ ] Crisis management creates memorable stories players want to share
- [ ] Banking system used by >50% of players for major projects

### Gameplay Metrics
- [ ] Multiple viable economic strategies (self-sufficient, trade, financial)
- [ ] Economic specialization occurs naturally in playtesting
- [ ] Crises create challenge without causing excessive player frustration
- [ ] Economic systems provide long-term engagement beyond basic production

---

## Dependencies

### Required First
- Basic resource production system
- Resident needs and wealth tracking
- UI framework for complex economic interfaces

### Dependent Features
- International trade (requires basic trade system)
- Stock market/exchanges (requires banking system)
- Economic policies/taxation (requires currency system)
- Economic history/analysis (requires price tracking)

---

## Open Questions

1. **Economic abstraction:** How abstract should currency be? Actual coin counts or simplified "wealth"?
2. **AI economics:** How sophisticated should simulated village economies be?
3. **Multiplayer economics:** If multiplayer added later, how handle inter-player trade?
4. **Historical monetary systems:** Use actual medieval currencies or simplified system?

---

*Feature Version 1.0 · Owner: Economic Systems Team · Estimated Effort: 6-8 sprints*