# Change: Phase 1 Foundation - Core Infrastructure

## Why

Toutā is currently in the planning phase with no implementation. Phase 1 establishes the foundational infrastructure required for all future development. This includes the core abstractions (interfaces), dependency injection container, message bus, router, configuration system, and CLI framework.

Without these foundational components, we cannot:
- Implement packages or recipes
- Create reusable components
- Support message-based communication
- Provide pluggable implementations
- Build developer tools

Phase 1 delivers the minimum viable framework core that enables all subsequent phases.

## What Changes

This change introduces the entire core framework infrastructure from scratch:

### Core Components
- **Dependency Injection Container**: Interface-based DI with auto-wiring, singleton support, and service providers
- **Message Bus**: Pub/sub system for message-based communication with config and code-based routing
- **Router Abstraction**: Interface-based HTTP router with Chi as default implementation
- **Configuration System**: YAML frontmatter loader with hot-reload support
- **CLI Framework**: Cobra-based extensible command system with discovery mechanism
- **Component Registry**: Package manifest parsing and component registration
- **Template Renderer**: Basic html/template wrapper with interface abstraction
- **Development Tools**: Hot reload integration and development server

### Developer Experience
- CLI commands: `touta new`, `touta init`, `touta serve`
- Project scaffolding with sensible defaults
- Development server with hot reload
- Clear interface-driven API

### Testing & Documentation
- Unit tests for all core components
- Integration tests for message flow
- Example implementations
- Developer documentation

## Impact

### Affected Specs
- **Creates** `core/dependency-injection` - DI container specification
- **Creates** `core/message-bus` - Message system specification
- **Creates** `core/router` - HTTP router specification
- **Creates** `core/config` - Configuration system specification
- **Creates** `core/cli` - CLI framework specification
- **Creates** `core/component-registry` - Component registration specification
- **Creates** `core/template-renderer` - Template rendering specification

### Affected Code
- **Creates** entire `internal/` directory structure
- **Creates** `cmd/touta/` CLI entry point
- **Creates** `pkg/touta/` public framework APIs
- **Creates** example implementations and tests

### Timeline
- **Estimated Duration**: 6-8 weeks
- **Dependencies**: None (first implementation)
- **Blockers**: None

### Success Criteria
1. ✅ DI container resolves dependencies via interfaces
2. ✅ Message bus routes messages to handlers
3. ✅ HTTP server responds to requests using Chi router
4. ✅ Config loads from YAML with frontmatter
5. ✅ CLI creates new projects and starts dev server
6. ✅ All core interfaces defined and documented
7. ✅ >80% test coverage for core components
8. ✅ Developer can create "Hello World" message handler
9. ✅ Hot reload works in development mode
10. ✅ Example project demonstrates core features

### Breaking Changes
None (initial implementation)

### Migration Guide
N/A (no existing code to migrate)
