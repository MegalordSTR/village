# Workstreams Index

## Feature F006: Go Modernization

**Goal**: Update Go version, tooling, and ensure code follows modern best practices.

### Workstreams

| WS ID | Title | Status | Dependencies | Scope | Estimated LOC |
|-------|-------|--------|--------------|-------|---------------|
| [WS-006-01](backlog/00-006-01.md) | Update Go Version and CI Configuration | backlog | Independent | SMALL | ~50 |
| [WS-006-02](backlog/00-006-02.md) | Add Modern Go Linting Tools | backlog | WS-006-01 | MEDIUM | ~100 |
| [WS-006-03](backlog/00-006-03.md) | Code Audit for Modern Go Best Practices | backlog | WS-006-02 | MEDIUM | Analysis only |
| [WS-006-04](backlog/00-006-04.md) | Update Development Tools and Security Audit | backlog | Independent | MEDIUM | ~100 |
| [WS-006-05](backlog/00-006-05.md) | Integration Tests and Build Verification | backlog | WS-006-01, WS-006-02, WS-006-04 | MEDIUM | ~200 |
| [WS-006-06](backlog/00-006-06.md) | Enforce Mandatory Linters and Fix Violations | backlog | WS-006-02 | MEDIUM | ~150 |

### Dependency Graph

```mermaid
graph TD
    A[WS-006-01: Update Go Version] --> B[WS-006-02: Linting Tools]
    B --> C[WS-006-03: Code Audit]
    B --> F[WS-006-06: Linter Enforcement]
    A --> E[WS-006-05: Integration Tests]
    B --> E
    D[WS-006-04: Dev Tools] --> E
```

### Execution Order

1. **WS-006-01**: Update Go version in CI and configurations (independent)
2. **WS-006-04**: Update development tools and security audit (parallel with #1)
3. **WS-006-02**: Add linting tools after Go version is updated
4. **WS-006-03**: Perform code audit using linting tools
5. **WS-006-06**: Enforce mandatory linters and fix violations (depends on #3)
6. **WS-006-05**: Final integration tests and build verification

### Notes (F006)

- WS-006-03 produces an audit report, not code changes. Subsequent workstreams may be created based on its recommendations.
- All workstreams are scoped SMALL or MEDIUM (<500 LOC changes).
- Feature ID F006 corresponds to "Go Modernization" initiative.

---

## Feature F004: Resource Economy System Integration

**Goal**: Integrate the advanced Resource Economy System (F002) with the simulation economic system, replacing duplicate resource handling with proper economy package types and functions.

### Workstreams

| WS ID | Title | Status | Dependencies | Scope | Estimated LOC |
|-------|-------|--------|--------------|-------|---------------|
| [WS-004-01](backlog/00-004-01.md) | Eliminate Circular Dependencies | backlog | Independent | SMALL | ~30 |
| [WS-004-02](backlog/00-004-02.md) | Update Simulation Types to Use Economy Types | backlog | WS-004-01 | MEDIUM | ~100 |
| [WS-004-03](backlog/00-004-03.md) | Integrate Inventory into GameState | backlog | WS-004-02 | MEDIUM | ~80 |
| [WS-004-04](backlog/00-004-04.md) | Update Economic System to Use Inventory | backlog | WS-004-03 | MEDIUM | ~150 |
| [WS-004-05](backlog/00-004-05.md) | Update Production System to Use Economy Types | backlog | WS-004-02 | MEDIUM | ~120 |
| [WS-004-06](backlog/00-004-06.md) | Update Tests and Ensure Backward Compatibility | backlog | WS-004-04, WS-004-05 | MEDIUM | ~100 |

### Dependency Graph

```mermaid
graph TD
    A[WS-004-01: Eliminate Circular Dependencies] --> B[WS-004-02: Update Simulation Types]
    B --> C[WS-004-03: Integrate Inventory]
    B --> E[WS-004-05: Update Production System]
    C --> D[WS-004-04: Update Economic System]
    D --> F[WS-004-06: Tests & Compatibility]
    E --> F
```

### Execution Order

1. **WS-004-01**: Remove circular imports between economy and simulation packages
2. **WS-004-02**: Update simulation types (Resource, Production) to use economy types
3. **WS-004-03**: Add Inventory field to GameState and implement sync with Resources
4. **WS-004-05**: Update production system to use economy resource types (parallel with #3)
5. **WS-004-04**: Rewrite economic system to use Inventory and StorageRegistry
6. **WS-004-06**: Update all tests, ensure backward compatibility and deterministic behavior

### Notes

- All workstreams are scoped SMALL or MEDIUM (<500 LOC changes).
- Feature ID F004 corresponds to "Resource Economy System Integration" initiative.
- Depends on completed F002 (Resource Economy System).

---

