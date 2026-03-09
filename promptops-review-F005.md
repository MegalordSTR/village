# PromptOps Review for Feature F005

## SCOPE:
- sdp/prompts/skills/*/SKILL.md (24 skills)
- sdp/prompts/agents/*.md (10 agents)
- sdp/prompts/commands/*.md (2 commands)
- sdp/AGENTS.md (global agent instructions)

## RISK MAP:
1. **Language-specific assumptions** - Quality gates and skill notes assume Go language, reducing reusability across other language projects.
2. **CLI command inconsistencies** - Mixed references to `sdp-orchestrate` vs `sdp orchestrate` could cause confusion.
3. **Handoff list anti-pattern** - Potential for vague handoff lists; review skill includes a handoff block but it's a specific command.

## EVIDENCE:
1. Language-specific:
   - `sdp/AGENTS.md:20-22` - Quality gates: `go build ./...`, `go test ./...`, `go vet ./...`
   - `sdp/prompts/skills/feature/SKILL.md:16` - "Phase 0: This skill targets Go projects"
   - `sdp/prompts/skills/oneshot/SKILL.md:16` - "Run orchestrate: Either `sdp-orchestrate` on PATH, or from project root: `go run ./cmd/sdp-orchestrate`" (assumes Go installed)

2. CLI inconsistencies:
   - `sdp/prompts/skills/oneshot/SKILL.md:14` - references `sdp-orchestrate` (hyphenated)
   - `sdp/prompts/skills/build/SKILL.md:62` - references `sdp guard activate` (valid)
   - Actual command is `sdp orchestrate` (no hyphen) as verified.

3. Handoff lists:
   - `sdp/prompts/skills/review/SKILL.md:70-78` - Handoff block when verdict=CHANGES_REQUESTED with specific command `@design phase4-remediation`. This is a clear command, not a vague list.
   - `sdp/AGENTS.md:49` - "Hand off - Provide context for next session" generic step.

## SEVERITY:
1. Language-specific assumptions: **P2** (maintenance debt) - reduces reusability but not blocking for current Go project.
2. CLI inconsistencies: **P3** (style) - minor inconsistency, unlikely to cause failures.
3. Handoff lists: **PASS** - no problematic handoff lists found.

Max severity: P2.

## VERDICT:
PASS (max severity ≤ P2)

## FINDINGS_CREATED:
village-rvy village-oyx village-f8n

## PromptOps Checks:
```json
[
  {
    "check_id": "language-agnostic",
    "status": "fail",
    "note": "Quality gates in AGENTS.md are Go-specific (go build/go test/go vet). Feature skill mentions targeting Go projects. Reduces reusability across languages."
  },
  {
    "check_id": "no-phantom-cli",
    "status": "pass",
    "note": "All referenced CLI commands exist (sdp, bd, sdp guard, sdp quality, sdp orchestrate). sdp-orchestrate binary not in PATH but source exists and can be run via go run."
  },
  {
    "check_id": "no-handoff-lists",
    "status": "pass",
    "note": "No vague handoff lists found. The handoff block in review skill is a specific command."
  },
  {
    "check_id": "skill-size",
    "status": "pass",
    "note": "All 24 skill files are under 200 lines (max 191 lines)."
  }
]
```