# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
### Changed
### Fixed

## [0.1.2] - 2026-03-09

### Added
- **Go Modernization (F006) - Final Workstream**: Enforce mandatory linters and fix violations across codebase
  - Integrated `govet`, `staticcheck`, `gosimple`, `unused`, `stylecheck` linters
  - Fixed all linter violations across production code
  - Added mandatory linter enforcement in CI pipeline

### Changed
- Updated `.golangci.yml` configuration with stricter linter settings
- Improved code quality with TODO comments and documentation

### Fixed
- None in this release

## [0.1.1] - 2026-03-09

### Added
- **Go Modernization (F006)**: Comprehensive Go tooling and infrastructure updates with 5 workstreams:
  - Update Go Version and CI Configuration: Consistent Go 1.25.7 across CI/CD, Docker, and go.mod
  - Add Modern Go Linting Tools: Integrated golangci-lint with automated CI workflow
  - Code Audit for Modern Go Best Practices: Analysis of error handling, generics, concurrency, API design
  - Update Development Tools and Security Audit: Node.js LTS updates, npm dependency security fixes
  - Integration Tests and Build Verification: Go version matrix testing, race detection, benchmarks

### Changed
- Updated GitHub Actions workflows for Go 1.25 and improved linting
- Enhanced security scanning and dependency management
- Improved test infrastructure with integration tests and benchmarks

### Fixed
- None in this release

## [0.1.0] - 2026-03-08

### Added
- **Resource Economy System (F002)**: Complete resource economy simulation with 8 workstreams:
  - Resource Definition System: 15 core resource types with quality tiers, spoilage rates, and base values
  - Production Chain System: Conversion of raw materials to processed goods via production buildings
  - Inventory & Storage System: Storage buildings with capacity limits and spatial organization
  - Seasonal Economic Cycle: Seasonal price fluctuations and weather impact on production
  - Craftsmanship & Quality System: Quality tiers affecting value and production success
  - Resource Scarcity & Substitution: Scarcity-based price adjustments and substitution logic
  - UI Integration (Angular): Frontend components for resource management
  - Deployment Setup: Docker Compose configuration and CI/CD pipeline

### Changed
- Project structure with new internal/economy package
- Updated Docker Compose configuration for production deployment
- Enhanced CI/CD pipeline with automated testing and deployment

### Fixed
- None in this release

[Unreleased]: https://github.com/MegalordSTR/village/compare/v0.1.2...HEAD
[0.1.2]: https://github.com/MegalordSTR/village/releases/tag/v0.1.2
[0.1.1]: https://github.com/MegalordSTR/village/releases/tag/v0.1.1
[0.1.0]: https://github.com/MegalordSTR/village/releases/tag/v0.1.0