# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.2.0] - 2026-01-02

### BREAKING CHANGES

**The main framework now uses production-ready component libraries instead of internal implementations.**

- **Removed** `internal/di/container.go` - replaced by `toutago-nasc-dependency-injector`
- **Removed** `internal/router/chi_router.go` - replaced by `toutago-cosan-router`
- **Removed** `internal/template/html_renderer.go` - replaced by `toutago-fith-renderer`

### Added

- **Integration Layer** - New `pkg/touta/integration` package with adapters for component libraries
  - `NascContainerAdapter` - Adapts nasc.Nasc to touta.Container interface
  - `CosanRouterAdapter` - Adapts cosan.Router to touta.Router interface
  - `FithRendererAdapter` - Adapts fith.Engine to touta.TemplateRenderer interface
  - `NewDataMapper()` - Factory for creating datamapper instances
  - `NewMigrator()` - Factory for creating database migrators

- **Component Dependencies**:
  - `github.com/toutaio/toutago-nasc-dependency-injector` v1.0.9 - DI container
  - `github.com/toutaio/toutago-cosan-router` v1.0.5 - HTTP router
  - `github.com/toutaio/toutago-fith-renderer` v1.0.6 - Template engine
  - `github.com/toutaio/toutago-datamapper` v1.0.8 - Database abstraction
  - `github.com/toutaio/toutago-sil-migrator` v1.0.5 - Migration tool

- **Examples**:
  - `examples/basic/` - Basic server with DI and routing
  - `examples/with-templates/` - Template rendering example

- **Comprehensive Integration Tests** - All adapters have test coverage

### Changed

- Framework architecture now demonstrates pluggable component design
- Reduced main framework LOC by ~500+ lines
- Improved separation of concerns via adapter pattern
- Component libraries can now be used standalone in any Go project

### Migration Guide

**Before (v0.1.x):**
```go
import "github.com/toutaio/toutago/internal/di"
import "github.com/toutaio/toutago/internal/router"

container := di.NewContainer()
router := router.NewChiRouter(container)
```

**After (v0.2.0):**
```go
import "github.com/toutaio/toutago/pkg/touta/integration"

container := integration.NewContainer()
router := integration.NewRouter(container)
```

### Benefits

- **Best-in-class components**: Each library is mature, tested (80%+ coverage), production-ready
- **Flexibility**: Easy to swap implementations via interfaces
- **Ecosystem**: Components work standalone in any Go project
- **SOLID principles**: Clear separation of concerns
- **Reduced complexity**: Less code to maintain in the main framework

## [0.1.0] - 2024-12-21

### Added
- Core dependency injection container with auto-wiring support
- Message bus system with async/sync publish and subscription
- HTTP router abstraction with Chi as default implementation
- YAML configuration loader with frontmatter support
- CLI framework with commands: new, init, serve, version
- Template renderer using html/template
- Component registry for package management
- Project scaffolding functionality
- **Complete test suite with 85.9% coverage** ✅
- **Hot reload system for development** ✅
- Example project generation

### Testing
- 78 comprehensive tests across all components
- Config: 90.8% coverage (16 tests)
- DI Container: 80.6% coverage (24 tests)
- Message Bus: 80.2% coverage (4 tests)
- Registry: 95.1% coverage (11 tests)
- Router: 92.7% coverage (17 tests)
- Template: 89.3% coverage (6 tests)

### Architecture
- Interface-first design pattern established
- Message-passing architecture implemented
- No global state - pure dependency injection
- Pluggable component system

### Developer Experience
- `touta new` command to create new projects
- `touta init` command to initialize existing directories
- `touta serve` command with **integrated hot reload** ✅
- Auto-generated project structure with examples
- File watching for *.go, *.yaml, *.html files
- Automatic process restart on changes

### Documentation
- Complete README with examples
- QUICKSTART.md reference guide
- IMPLEMENTATION_SUMMARY.md technical details
- PHASE1_FINAL.md completion report
- Godoc comments on all public APIs

## [Unreleased]

### Planned for Phase 1 Completion
- Complete test coverage (>80%)
- Hot reload integration with cosmtrek/air
- Performance benchmarks
- More comprehensive examples
- Integration tests
- GitHub Actions CI/CD pipeline

### Planned for Phase 2
- Custom template dialect with `<box:*>` tags
- Enhanced template features
- Template hot-reloading
- Template caching optimizations

### Planned for Phase 3
- Package system implementation
- Package discovery and loading
- Package dependencies
- Package CLI commands

### Planned for Phase 4
- Recipe system
- Pre-built component recipes
- Recipe marketplace

### Planned for Phase 5
- Production deployment tools
- Database adapters
- WebSocket support
- Advanced middleware
- Monitoring and observability

[0.1.0]: https://github.com/toutaio/toutago/releases/tag/v0.1.0
