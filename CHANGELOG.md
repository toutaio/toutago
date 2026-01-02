# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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
