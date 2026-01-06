# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed
- Project template now uses PORT environment variable instead of hardcoded port 8080

### Added
- Comprehensive test coverage for integration package
  - Coverage improved from 52.6% to 58.3%
  - Added tests for PUT, DELETE, PATCH HTTP methods
  - Added tests for context parameter and query extraction
  - Added tests for redirect functionality
  - Added tests for request/response access

### Fixed
- Fixed missing go.sum entries for spf13/pflag dependency

## [0.3.0] - 2026-01-04

### BREAKING CHANGES

**Complete integration of Scéla message bus with new interface design.**

- **Removed** `internal/message/` - Old internal message bus implementation completely removed
- **Changed** Message interfaces - Updated to align with Scéla's API design
  - `Message.Topic()` replaces `Message.Slug()` and `Message.Type()`
  - `Handler.Handle(ctx, msg)` returns `error` only (no response message)
  - `Bus.Subscribe()` returns `Subscription` instead of `error` only
  - `Bus.Close()` replaces `Bus.Start()` and `Bus.Stop()`

### Added

- **Scéla Message Bus Integration**:
  - `github.com/toutaio/toutago-scela-bus` v1.4.0 - Production-ready message bus
  - `integration.NewScelaBus()` - Factory for creating Scéla bus instances
  - `integration.NewScelaBusWithMiddleware()` - Factory with global middleware
  - `ScelaBus.PublishSync()` - Synchronous message publishing
  - `ScelaBus.PublishWithPriority()` - Priority message publishing
  - `ScelaBus.Use()` - Add middleware to the bus
  - Pattern matching support (`user.*`, `app.**`)
  - Priority message processing (Low, Normal, High, Urgent)
  - Optional persistence (Redis, Filesystem)
  - Dead letter queue support
  - Retry logic with exponential backoff

- **Documentation**:
  - `docs/message-bus.md` - Comprehensive message bus guide
  - `examples/with-scela/README.md` - Detailed Scéla example documentation
  - Enhanced README with Scéla usage examples
  - Updated QUICKSTART with modern message bus patterns

- **Examples**:
  - Enhanced `examples/with-scela/` with middleware, patterns, and priorities
  - Added `examples/basic/` message bus integration
  - Comprehensive integration test suite

### Changed

- **Simplified Scéla Adapter** - More direct exposure of Scéla features
- **Updated all examples** - Now use Scéla bus instead of old internal bus
- **Improved integration layer** - Better separation of concerns

### Removed

- `internal/message/bus.go` - Replaced by Scéla integration
- `internal/message/` directory - No longer needed
- Old `BaseMessage` type - Use Scéla's Message interface
- Legacy message bus interfaces - Replaced with Scéla-aligned design

### Migration Guide

**Message Publishing - Before (v0.2.x):**
```go
import "github.com/toutaio/toutago/internal/message"

bus := message.NewBus()
bus.Start(context.Background())

bus.Publish(ctx, &UserCreated{
    BaseMessage: message.BaseMessage{
        MessageSlug: "user.created",
        MessageType: "event",
    },
    UserID: "123",
})
```

**After (v0.3.0):**
```go
import "github.com/toutaio/toutago/pkg/touta/integration"

bus := integration.NewScelaBus()
defer bus.Close()

bus.Publish(ctx, "user.created", map[string]interface{}{
    "id": "123",
})
```

**Message Subscription - Before:**
```go
type MyHandler struct{}

func (h *MyHandler) Handle(ctx context.Context, msg touta.Message) (touta.Message, error) {
    // Handle message
    return nil, nil
}

bus.Subscribe("user.created", &MyHandler{})
```

**After:**
```go
bus.Subscribe("user.created", touta.HandlerFunc(
    func(ctx context.Context, msg touta.Message) error {
        // Handle message
        return nil
    },
))
```

### Benefits

- **Production-ready** - Built on battle-tested Scéla bus (80%+ test coverage)
- **Feature-rich** - Pattern matching, priorities, persistence, middleware
- **Simpler API** - More intuitive interface design
- **Better performance** - Optimized worker pool and message queue
- **Flexible** - Easy to configure for different use cases

### See Also

- [Scéla Bus Repository](https://github.com/toutaio/toutago-scela-bus)
- [Message Bus Guide](docs/message-bus.md)
- [Examples](examples/with-scela/)

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
