# Project Context

## Purpose

Toutā is a Go-based web framework designed for extreme decoupling, extensibility, and maintainability. The project aims to create a system that:

- Encourages separation of concerns and single responsibility through independent nemetons
- Provides both configuration-based and code-based flow options
- Makes deployment, customization, and extension as simple as possible
- Supports complete solution delivery through a "ritual" system (e.g., blog, wiki, eCommerce)
- Implements message-passing OOP philosophy inspired by Smalltalk

## Tech Stack

### Core Technologies
- **Language:** Go 1.21+
- **Router:** Interface-based (default: Chi implementation)
- **WebSocket:** Interface-based (default: gorilla/websocket)
- **Template Engine:** Interface-based (default: custom wrapping html/template)
- **Config Parser:** Interface-based (default: github.com/adrg/frontmatter for YAML)
- **Validation:** github.com/go-playground/validator (Pydantic-like validation)
- **DI Container:** Interface-based dependency injection
- **CLI Tool:** Cobra-based extensible ogam system
- **Hot Reload (dev):** air or cosmtrek/air

### Supporting Libraries
- Chi router (pluggable)
- Gorilla WebSocket (pluggable)
- Go Playground Validator
- YAML frontmatter parser

## Project Conventions

### Code Style

#### Interface-First Design (CRITICAL)
- **Everything is an interface unless there's a compelling reason otherwise**
- Never code against concrete implementations—always interfaces
- All dependencies must be injected (via constructor or DI container)
- Never use `new()` or direct instantiation for dependencies

#### Naming Conventions
- Interfaces: Noun describing capability (e.g., `StorageAdapter`, `MessageBus`, `Router`)
- Implementations: Implementation detail + Interface name (e.g., `ChiRouter`, `FilesystemAdapter`)
- Messages: Past tense + Message (e.g., `UserRegistered`, `FormSubmitted`)
- Handlers: Domain + Handler (e.g., `UserHandler`, `AuthHandler`)

#### File Organization
```
nemeton-name/
├── nemeton.yaml        # Nemeton metadata
├── components/         # Go handlers/components
├── messages/          # Message type definitions
├── templates/         # Frontend templates
├── routes.yaml        # Message routing config
├── ogam/          # CLI ogam
└── assets/            # Static resources
```

### Architecture Patterns

#### 1. Message-Centric Architecture
All inter-component communication happens through typed, validated messages:
- Messages implement the `Message` interface with `Slug()`, `Type()`, `Metadata()`
- Messages extend `BaseMessage` for core metadata
- Message Bus coordinates all message routing (config-based or code-based)
- Handlers listen for specific message types/slugs

#### 2. Dependency Injection Pattern
All dependencies resolved through DI container:
- Bind interfaces to implementations
- Support singleton and factory patterns
- Auto-wiring with struct tags (`inject:""`)
- Service providers organize related bindings

#### 3. Interface-Driven Design
Every major component is abstracted behind interfaces:
- `Router` - HTTP routing (swappable: Chi, Fiber, Gin, etc.)
- `Container` - Dependency injection
- `ConfigLoader` - Configuration parsing
- `TemplateRenderer` - Template rendering
- `StorageAdapter` - Data persistence
- `MessageBus` - Message routing
- `Ogam` - CLI ogam

#### 4. Nemeton System
Components organized into independent, reusable nemetons:
- Local nemetons (in-repo development)
- External nemetons (imported via nemeton manager)
- Nemetons can be moved between local and external
- Centralized dependency management (Composer-style for Go modules)

#### 5. Ritual System
Complete solutions packaged as rituals:
- Rituals compose multiple nemetons
- Rituals can extend other rituals
- Examples: blog, wiki, eCommerce, custom sites

### Testing Strategy

#### Framework Testing
- Unit tests for all core components
- Integration tests for message flow
- E2E tests for template rendering
- Benchmark tests for performance baselines

#### Nemeton Testing
Each nemeton must include:
- Component unit tests
- Message validation tests
- Template rendering tests
- Handler integration tests

#### Ritual Testing
- Deployment smoke tests
- Critical path testing
- Performance baseline validation

#### Test Conventions
- Use table-driven tests where applicable
- Mock all dependencies using interfaces
- Test message validation separately
- Benchmark critical paths

### Git Workflow

**Status:** Early development (no commits yet)

Planned workflow:
- **Main branch:** Stable releases only
- **Development branch:** Integration branch for features
- **Feature branches:** Individual features/nemetons
- **Commit messages:** Conventional commits format
  - `feat:` New features
  - `fix:` Bug fixes
  - `docs:` Documentation changes
  - `refactor:` Code refactoring
  - `test:` Test additions/updates
  - `chore:` Build/tooling changes

## Domain Context

### Message-Passing OOP Philosophy
Toutā follows Smalltalk's message-passing approach:
- Objects (components) don't call methods on each other directly
- All communication happens through message objects
- Message Bus routes messages to appropriate handlers
- Handlers respond with new messages

### Nemeton Independence
Nemetons are designed to be truly independent:
- No direct dependencies between nemetons (except via messages)
- Each nemeton declares its own dependencies in `nemeton.yaml`
- Project-level dependency resolution ensures compatibility
- Nemetons can be developed, tested, and distributed separately

### Hybrid Rendering Strategy
Inspired by Next.js:
- **SSR (Server-Side Rendering):** Default for initial page loads
- **CSR (Client-Side Rendering):** For dynamic updates
- **Postback Mechanism:** Components can trigger backend code via WebSocket/HTTP
- Templates specify rendering strategy per component

### Configuration Options
Two approaches for component relationships:
1. **Configuration-based:** YAML files define message routing
2. **Code-based:** Direct handler subscription in Go code
Developers choose based on decoupling needs

## Important Constraints

### Technical Constraints
- **Database Agnostic:** No database in core; storage adapters for persistence
- **Interface-Only Dependencies:** No concrete type coupling allowed
- **Message-Based Communication:** Required for inter-component interaction
- **Go Version:** Minimum Go 1.21+
- **Stateless First:** Architecture leans stateless with optional stateful support

### Design Constraints
- **Extreme Decoupling:** Nemetons must be independently distributable
- **Pluggable Everything:** All major components must be swappable via interfaces
- **Zero Lock-In:** No vendor lock-in for routers, databases, config formats, etc.
- **Convention Over Configuration:** Sensible defaults with full customization

### Development Constraints
- **Interface-First:** All code must use interfaces for dependencies
- **DI Required:** No direct instantiation of dependencies
- **Message Validation:** All messages must use struct validation tags
- **Nemeton Isolation:** Nemetons cannot directly call other nemetons' code

## External Dependencies

### Required Dependencies
- **github.com/go-chi/chi:** Default HTTP router (pluggable)
- **github.com/gorilla/websocket:** WebSocket support (pluggable)
- **github.com/go-playground/validator:** Message validation
- **github.com/adrg/frontmatter:** YAML frontmatter parsing
- **github.com/spf13/cobra:** CLI framework

### Optional Dependencies
- **cosmtrek/air:** Hot reload for development
- Alternative routers (Fiber, Gin, Echo) - user choice
- Alternative DI containers - user choice
- Database drivers - via storage adapter nemetons

### Future Nemeton Ecosystem
Planned optional nemetons:
- `touta-postgres` - PostgreSQL storage adapter
- `touta-mysql` - MySQL storage adapter
- `touta-mongodb` - MongoDB storage adapter
- `touta-sqlite` - SQLite storage adapter
- Ritual nemetons (blog, wiki, eCommerce, etc.)
- Community-contributed nemetons

---

**Version:** 0.1.0 (Planning Phase)  
**Last Updated:** 2025-12-20  
**Project Status:** Architecture design and planning stage
