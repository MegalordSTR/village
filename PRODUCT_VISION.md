# Product Vision: Medieval Village Economic Simulator

## Core Concept
A historically-accurate, deep economic simulation of a medieval village with weekly turn-based gameplay, individual resident management, and complex resource production chains. Players make strategic decisions before each turn, then observe emergent outcomes influenced by seasonal cycles, social dynamics, and historical events.

## Target Audience
- Strategy enthusiasts who enjoy systemic complexity (Aurora 4x, Dwarf Fortress, Rimworld players)
- History buffs interested in medieval economics and social structures
- Players seeking long-term engagement with meaningful progression

## Unique Value Proposition
**Historical depth meets strategic complexity:** Unlike simplified city builders, this simulator models authentic medieval economics where every decision creates cascading consequences through production chains, social hierarchies, and environmental factors.

## Success Metrics
- Player retention: Average session length > 2 hours, return rate > 70% weekly
- Depth perception: Players discover new system interactions after 10+ hours
- Historical accuracy: Positive feedback from history-focused communities
- Performance: Complex village (50+ residents) processes in < 5 seconds per turn

## Expert Analysis

### Product Expert
**Summary:** This is a niche but passionate market segment that values depth over accessibility. The weekly turn-based system creates a "slow burn" engagement model perfect for thoughtful decision-making. Individual resident management adds emotional investment missing from abstract city builders.

**Key finding:** The MVP should focus on making the complex simulation feel tangible through clear cause-effect visualization—when a player assigns a resident to blacksmithing, they should immediately see how that affects tool production, then farming efficiency, then food stores.

**Risk/opportunity:** Risk of overwhelming new players with complexity; opportunity to create a "Eureka!" moment when systems click together.

### Market Expert
**Summary:** The medieval/strategy simulation market is underserved between casual mobile games and ultra-complex PC titles. Steam analytics show consistent 50K+ monthly active users for games like Banished, with strong modding communities extending lifespan 5-10x.

**Key finding:** Modding capability isn't a nice-to-have—it's essential for longevity. The most successful simulators in this space (Rimworld, Cities: Skylines) derive 80%+ of content from community mods.

**Risk/opportunity:** Risk of competing against established titles; opportunity to capture disenfranchised players seeking more historical accuracy than fantasy elements.

### Technical Expert
**Summary:** Go backend excels at concurrent simulation calculations but requires careful architecture for deterministic turn processing. PostgreSQL 18 with JSONB fields can handle complex resident state efficiently. Angular 21 frontend needs optimization for frequent UI updates during turn resolution.

**Key finding:** Simulation determinism is critical—the same inputs must produce identical results for save/load functionality. Implement a seed-based random system with serializable game state.

**Risk/opportunity:** Risk of frontend slowdown with 50+ animated residents; opportunity to use Web Workers for background simulation calculations during player planning phase.

### UX Expert
**Summary:** The tension between data richness and interface clarity defines this genre. Players need to understand complex relationships (resident A's farming skill affects yield, which affects food stores, which affects health, which affects productivity) without drowning in spreadsheets.

**Key finding:** Implement a "causal chain" visualization that shows how player decisions propagate through systems. When adjusting tax rates, show predicted effects on resident morale, productivity, and rebellion risk.

**Risk/opportunity:** Risk of analysis paralysis; opportunity to create satisfaction through understanding complex systems.

### Business Expert
**Summary:** Premium pricing ($19.99-29.99) with eventual Steam release aligns with market expectations. Development cost concentrated in simulation engine (backend) with frontend as visualization layer. Long-term revenue through DLC expanding historical periods (Roman, Renaissance) rather than microtransactions.

**Key finding:** Build community tools early—modding documentation, save game sharing, scenario creator. These create network effects that reduce customer acquisition costs.

**Risk/opportunity:** Risk of niche audience size limiting revenue; opportunity to become definitive historical simulator with expansion into educational markets.

### Growth Expert
**Summary:** Content creation drives discovery in simulation genres. Twitch/YouTube creators need interesting failure states and emergent stories to share. The weekly turn structure creates natural cliffhangers perfect for episodic content.

**Key finding:** Design for "shareable moments"—unexpected events (plague, rebellion, miraculous harvest) that players will want to screenshot or record. Include built-in screenshot tools with data overlays.

**Risk/opportunity:** Risk of slow initial growth; opportunity to build dedicated community that creates its own marketing through stories.

### Risk Expert
**Summary:** Technical risk: simulation complexity causing performance issues or bugs in cascading systems. Market risk: too niche to sustain development. Execution risk: scope creep from "one more system" mentality common in simulation games.

**Key finding:** Implement aggressive scoping for MVP—one village, 15 resources, basic seasons. Use modular architecture so additional systems (trade, warfare, religion) can be added without rewriting core.

**Risk/opportunity:** Risk of never feeling "complete" due to simulation depth; opportunity to create decade-long engagement through gradual expansion.

## Core Gameplay Loop
1. **Plan Phase:** Review village status, assign residents to tasks, set policies, make strategic decisions
2. **Execute Turn:** Simulation processes week of game time with deterministic outcomes
3. **Review Results:** Observe consequences, unexpected events, system interactions
4. **Adapt Strategy:** Adjust plans based on new information, manage crises, optimize systems
5. **Progress:** Unlock new buildings, technologies, or social structures over long-term play

## Technical Architecture
```
Frontend (Angular 21) ↔ REST API ↔ Backend (Go)
       ↑                           ↑
    Browser                    PostgreSQL 18
   (SPA/PWA)                   (Game State)
                                ↓
                           Simulation Engine
                           (Deterministic)
```

## Development Principles
1. **Determinism First:** Same inputs → same outputs for reliable save/load
2. **Modular Systems:** Each subsystem (economics, social, environmental) independently testable
3. **Data-Driven Design:** Balance values in configuration files, not code
4. **Progressive Disclosure:** Complexity revealed as player learns, not dumped upfront
5. **Historical Plausibility:** Systems based on medieval realities, not game balance alone

## Long-Term Vision
Become the definitive historical village simulator, eventually expanding to:
- Multiple historical periods (Ancient, Medieval, Renaissance)
- Different geographic regions (European, Middle Eastern, Asian villages)
- Multi-village interactions (trade networks, conflicts)
- Educational modes with historical commentary

---

*Version 1.0 · Created: March 7, 2026 · Next Review: MVP Launch*