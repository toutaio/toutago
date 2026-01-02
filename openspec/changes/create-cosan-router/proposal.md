# Change: Create Cosan - Independent HTTP Router Component

## Why

The toutā framework currently relies on pluggable router interfaces, but lacks a canonical router implementation that embodies the project's principles of extreme decoupling, SOLID design, and message-centric architecture. Creating an independent, standalone router component (Cosan - Irish for "pathway") will:

1. Provide a reference implementation demonstrating toutā's architectural principles
2. Offer a production-ready, framework-agnostic router usable in any Go project
3. Establish patterns for creating independent, reusable toutā ecosystem components
4. Enable testing and validation of router interface design before broader adoption

Cosan will be developed in its own repository (`https://github.com/toutaio/toutago-cosan-router`) at `/home/nestor/Proyects/toutago-cosan-router`, making it truly independent and reusable across projects beyond toutā.

## What Changes

This creates a new independent component with the following characteristics:

**Core Features:**
- Interface-driven HTTP router implementation in Go
- SOLID principles compliance (Single Responsibility, Open/Closed, Liskov Substitution, Interface Segregation, Dependency Inversion)
- Middleware chain support with clean composition patterns
- Route parameter extraction and validation
- Method-based routing (GET, POST, PUT, DELETE, PATCH, OPTIONS, HEAD)
- Route grouping and nested routing
- Context-aware request handling
- Performance optimizations inspired by leading routers

**Architectural Principles:**
- Zero framework dependencies (standard library + minimal dependencies)
- Complete interface abstraction for all major components
- Pluggable route matching strategies (trie, radix tree, etc.)
- Message-passing friendly design
- Middleware as composable, testable units
- Dependency injection support

**Quality Attributes:**
- Comprehensive test coverage (>90%)
- Performance benchmarks against Chi, Gin, Echo, Fiber
- Documentation with examples and best practices
- Clean API design following Go idioms
- Minimal memory allocations in hot paths

**Repository Structure:**
- Independent GitHub repository under toutaio organization
- Standard Go module layout
- CI/CD pipeline (GitHub Actions)
- Semantic versioning
- MIT or Apache 2.0 license

**Development Phases:**
- **Phase 1:** Core routing engine and basic middleware
- **Phase 2:** Advanced features (route groups, parameters, validation)
- **Phase 3:** Performance optimization and benchmarking
- **Phase 4:** Documentation, examples, and ecosystem integration

## Impact

**New Capabilities:**
- `cosan-router` - New independent router component specification

**Affected Components:**
- None in toutā core (this is a new independent component)
- Future: toutā router interface will reference Cosan as canonical implementation

**Repository Structure:**
- Creates new repository at `/home/nestor/Proyects/toutago-cosan-router`
- GitHub repository at `https://github.com/toutaio/toutago-cosan-router`

**Ecosystem Impact:**
- Establishes pattern for creating independent toutā ecosystem components
- Demonstrates interface-first design principles
- Provides benchmark for future component development
- Enables community contributions to router ecosystem

**Migration Path:**
- No migration needed (new component)
- toutā framework can optionally adopt Cosan as default router
- Existing Chi/Gin/Echo integrations remain unchanged

**Breaking Changes:**
- None (new independent component)
