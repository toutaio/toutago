# Phase 1 Implementation Tasks

## 1. Project Foundation

- [ ] 1.1 Initialize Go module structure
  - [ ] Create `go.mod` with Go 1.21+
  - [ ] Set up `internal/`, `pkg/`, `cmd/` directories
  - [ ] Configure linting (golangci-lint)
  - [ ] Set up CI/CD (GitHub Actions)
  
- [ ] 1.2 Define core interfaces
  - [ ] Create `pkg/touta/interfaces.go` with all core interfaces
  - [ ] Document interface contracts
  - [ ] Add interface validation tests

## 2. Dependency Injection Container

- [ ] 2.1 Define Container interface
  - [ ] `Bind()`, `Singleton()`, `Factory()` methods
  - [ ] `Make()`, `MakeWith()` for resolution
  - [ ] `Has()`, `Tagged()` for introspection
  - [ ] Interface documentation with examples

- [ ] 2.2 Implement default container
  - [ ] Create `internal/di/container.go`
  - [ ] Implement binding storage (map-based)
  - [ ] Implement singleton caching
  - [ ] Implement factory pattern support
  - [ ] Thread-safe operations (sync.RWMutex)

- [ ] 2.3 Auto-wiring with reflection
  - [ ] Parse struct tags (`inject:""`)
  - [ ] Resolve dependencies recursively
  - [ ] Handle circular dependency detection
  - [ ] Error reporting for missing dependencies

- [ ] 2.4 Service providers
  - [ ] Define ServiceProvider interface
  - [ ] Implement provider registration
  - [ ] Implement boot sequence
  - [ ] Add core service providers

- [ ] 2.5 Testing
  - [ ] Unit tests for binding/resolution
  - [ ] Unit tests for singletons
  - [ ] Unit tests for auto-wiring
  - [ ] Integration tests for service providers
  - [ ] Benchmark dependency resolution

## 3. Message Bus System

- [ ] 3.1 Define Message interfaces
  - [ ] Base `Message` interface (Slug, Type, Metadata)
  - [ ] `BaseMessage` struct implementation
  - [ ] Message validation with struct tags
  - [ ] Example message types

- [ ] 3.2 Define MessageBus interface
  - [ ] `Publish()`, `Subscribe()`, `Unsubscribe()` methods
  - [ ] Handler interface definition
  - [ ] Async/sync dispatch options
  - [ ] Error handling patterns

- [ ] 3.3 Implement default message bus
  - [ ] Create `internal/message/bus.go`
  - [ ] Channel-based message queue
  - [ ] Subscriber registry (type-based)
  - [ ] Concurrent message processing (goroutines)
  - [ ] Error handling and logging

- [ ] 3.4 Config-based routing
  - [ ] Parse `routes.yaml` format
  - [ ] Route messages to handlers by slug/type
  - [ ] Support conditional routing
  - [ ] Support async routing flags

- [ ] 3.5 Code-based routing
  - [ ] Programmatic subscription API
  - [ ] Handler registration helpers
  - [ ] Unsubscription support

- [ ] 3.6 Testing
  - [ ] Unit tests for message publishing
  - [ ] Unit tests for subscription
  - [ ] Integration tests for routing
  - [ ] Concurrent dispatch tests
  - [ ] Performance benchmarks

## 4. Router Abstraction

- [ ] 4.1 Define Router interface
  - [ ] HTTP methods (GET, POST, PUT, DELETE, PATCH)
  - [ ] `Group()` for route prefixes
  - [ ] `Use()` for middleware
  - [ ] `Listen()` to start server
  - [ ] `Native()` to access underlying implementation

- [ ] 4.2 Define Context interface
  - [ ] Request/Response access
  - [ ] Parameter extraction
  - [ ] Get/Set for request-scoped data
  - [ ] Container access
  - [ ] Message access

- [ ] 4.3 Implement Chi router adapter
  - [ ] Create `internal/router/chi_router.go`
  - [ ] Wrap Chi router with interface
  - [ ] Implement handler adaptation
  - [ ] Implement middleware adaptation
  - [ ] Implement context creation

- [ ] 4.4 HTTP server abstraction
  - [ ] Graceful shutdown support
  - [ ] TLS/SSL configuration
  - [ ] Timeout configuration
  - [ ] Request logging

- [ ] 4.5 Testing
  - [ ] Unit tests for route registration
  - [ ] Integration tests for HTTP requests
  - [ ] Middleware chain tests
  - [ ] Context tests
  - [ ] Performance benchmarks

## 5. Configuration System

- [ ] 5.1 Define ConfigLoader interface
  - [ ] `Load()` method for parsing
  - [ ] `Watch()` for hot reload
  - [ ] `Validate()` for schema validation
  - [ ] Support for nested config

- [ ] 5.2 Define Config structure
  - [ ] Framework settings
  - [ ] Router settings
  - [ ] Server settings
  - [ ] Package settings
  - [ ] Environment-specific overrides

- [ ] 5.3 Implement YAML frontmatter loader
  - [ ] Create `internal/config/yaml_loader.go`
  - [ ] Parse frontmatter metadata
  - [ ] Parse YAML body
  - [ ] Merge config sources
  - [ ] Environment variable substitution

- [ ] 5.4 Hot reload support
  - [ ] File watching (fsnotify)
  - [ ] Config reload triggers
  - [ ] Callback notification system
  - [ ] Graceful reconfiguration

- [ ] 5.5 Testing
  - [ ] Unit tests for parsing
  - [ ] Unit tests for file watching
  - [ ] Integration tests with real files
  - [ ] Invalid config handling tests

## 6. CLI Framework

- [ ] 6.1 Define Command interface
  - [ ] `Name()`, `Description()` methods
  - [ ] `Execute()` with CommandContext
  - [ ] `Flags()` for command flags
  - [ ] Subcommand support

- [ ] 6.2 Implement Cobra integration
  - [ ] Create `cmd/touta/main.go`
  - [ ] Root command setup
  - [ ] Command registration system
  - [ ] Flag parsing integration

- [ ] 6.3 Core commands
  - [ ] `touta new <name>` - Create new project
  - [ ] `touta init` - Initialize in existing directory
  - [ ] `touta serve` - Start development server
  - [ ] `touta version` - Show version info

- [ ] 6.4 Command discovery
  - [ ] Auto-discover commands from packages
  - [ ] Command registration API
  - [ ] Namespace support (package:command)
  - [ ] Help generation

- [ ] 6.5 Testing
  - [ ] Unit tests for each command
  - [ ] Integration tests for CLI flow
  - [ ] Flag parsing tests
  - [ ] Output capture tests

## 7. Component Registry

- [ ] 7.1 Define ComponentRegistry interface
  - [ ] `Register()` for components
  - [ ] `Get()` for retrieval
  - [ ] `List()` for enumeration
  - [ ] `LoadFromManifest()` for package.yaml

- [ ] 7.2 Package manifest parser
  - [ ] Parse `package.yaml` format
  - [ ] Extract component definitions
  - [ ] Extract dependencies
  - [ ] Validate manifest schema

- [ ] 7.3 Component registration
  - [ ] Store component metadata
  - [ ] Link to message handlers
  - [ ] Link to templates
  - [ ] Link to routes

- [ ] 7.4 Testing
  - [ ] Unit tests for registration
  - [ ] Unit tests for manifest parsing
  - [ ] Integration tests with real manifests

## 8. Template Renderer

- [ ] 8.1 Define TemplateRenderer interface
  - [ ] `Render()` method
  - [ ] `RegisterFunction()` for helpers
  - [ ] `Parse()` for template loading
  - [ ] `Execute()` for rendering

- [ ] 8.2 Implement html/template wrapper
  - [ ] Create `internal/template/html_renderer.go`
  - [ ] Wrap Go html/template
  - [ ] Template caching
  - [ ] Hot reload support

- [ ] 8.3 Template functions
  - [ ] Built-in helpers (url, asset, etc.)
  - [ ] Custom function registration
  - [ ] Context-aware functions

- [ ] 8.4 Testing
  - [ ] Unit tests for rendering
  - [ ] Unit tests for functions
  - [ ] Integration tests with real templates

## 9. Development Tools

- [ ] 9.1 Hot reload integration
  - [ ] Integrate cosmtrek/air
  - [ ] Configure file watching
  - [ ] Auto-restart on changes
  - [ ] Config file (`.air.toml`)

- [ ] 9.2 Development server
  - [ ] Implement serve command
  - [ ] Port configuration
  - [ ] Live reload support
  - [ ] Request logging

- [ ] 9.3 Project scaffolding
  - [ ] Template for new projects
  - [ ] Default directory structure
  - [ ] Example touta.yaml
  - [ ] Sample component

## 10. Testing Infrastructure

- [ ] 10.1 Test utilities
  - [ ] Mock implementations for all interfaces
  - [ ] Test helpers for DI
  - [ ] Test helpers for messages
  - [ ] Test server setup

- [ ] 10.2 Unit tests
  - [ ] All core components >80% coverage
  - [ ] Table-driven tests where applicable
  - [ ] Error case testing
  - [ ] Edge case testing

- [ ] 10.3 Integration tests
  - [ ] End-to-end message flow
  - [ ] HTTP request handling
  - [ ] Config loading and hot reload
  - [ ] Component registration

- [ ] 10.4 Benchmarks
  - [ ] DI resolution performance
  - [ ] Message routing performance
  - [ ] Template rendering performance
  - [ ] HTTP request performance

## 11. Documentation

- [ ] 11.1 API documentation
  - [ ] Godoc comments for all public APIs
  - [ ] Interface contract documentation
  - [ ] Usage examples in comments

- [ ] 11.2 Developer guides
  - [ ] Quick start guide
  - [ ] Creating your first handler
  - [ ] Understanding the message bus
  - [ ] Working with dependency injection

- [ ] 11.3 Example project
  - [ ] Hello World message handler
  - [ ] Form submission example
  - [ ] Multiple handlers example
  - [ ] Config-based routing example

## 12. Release Preparation

- [ ] 12.1 Version tagging
  - [ ] Semantic versioning setup
  - [ ] Git tags for v0.1.0
  - [ ] CHANGELOG.md

- [ ] 12.2 Distribution
  - [ ] Go module publication
  - [ ] Binary releases (GitHub Releases)
  - [ ] Installation instructions

- [ ] 12.3 Demo and validation
  - [ ] Create demo video
  - [ ] Blog post announcement
  - [ ] Validate all success criteria
  - [ ] Performance baseline documentation
