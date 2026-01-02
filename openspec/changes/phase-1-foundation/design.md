# Phase 1 Foundation - Technical Design

## Context

Toutā Phase 1 establishes the foundational infrastructure for a Go-based web framework that emphasizes:
- **Interface-first design** - Everything abstracted behind interfaces
- **Message-passing architecture** - Smalltalk-inspired OOP
- **Pluggable components** - Swap implementations without code changes
- **Dependency injection** - No direct instantiation
- **Developer experience** - Simple CLI and hot reload

This design document outlines technical decisions and patterns for implementing these core components.

## Goals

### Primary Goals
1. Establish stable interface contracts for all core systems
2. Implement working DI container with auto-wiring
3. Create message bus with pub/sub and routing
4. Provide HTTP abstraction with Chi router default
5. Build CLI framework for project management
6. Enable hot-reload development workflow

### Non-Goals
- Advanced template dialects (Phase 2)
- Package system (Phase 3)
- Recipe system (Phase 4)
- Database adapters (Phase 5)
- WebSocket support (Phase 5)
- Production deployment tools (Phase 5)

## Architecture Decisions

### Decision 1: Interface-Only Dependencies

**Decision**: All component dependencies MUST be interfaces, never concrete types.

**Rationale**:
- Enables true pluggability - swap implementations via config
- Simplifies testing - mock any dependency
- Prevents coupling between components
- Supports runtime flexibility

**Implementation**:
```go
// ✅ Correct - interface dependency
type UserHandler struct {
    storage StorageAdapter  // Interface
    bus     MessageBus      // Interface
}

// ❌ Wrong - concrete dependency
type UserHandler struct {
    storage *FilesystemStorage  // Concrete type - NOT ALLOWED
}
```

**Enforcement**:
- Lint rule to detect concrete type fields
- Code review checklist
- Automated tests for interface compliance

---

### Decision 2: Dependency Injection Pattern

**Decision**: Use constructor injection with auto-wiring via struct tags.

**Rationale**:
- Explicit dependencies (constructor parameters)
- Auto-wiring reduces boilerplate (struct tags)
- No global state or singletons (except via DI)
- Testable - inject mocks easily

**Alternatives Considered**:
1. **Global variables** - Rejected: Creates hidden dependencies, hard to test
2. **Service locator** - Rejected: Hides dependencies, runtime errors
3. **Manual wiring only** - Rejected: Too verbose for large projects

**Implementation**:
```go
// Option 1: Constructor injection (explicit)
func NewUserHandler(storage StorageAdapter, bus MessageBus) *UserHandler {
    return &UserHandler{
        storage: storage,
        bus:     bus,
    }
}

// Option 2: Auto-wiring (convenience)
type UserHandler struct {
    storage StorageAdapter `inject:""`
    bus     MessageBus     `inject:""`
}

// Container resolves and injects
handler := &UserHandler{}
container.AutoWire(handler)
```

---

### Decision 3: Message Bus Architecture

**Decision**: Channel-based message queue with concurrent handler execution.

**Rationale**:
- Go channels are idiomatic for async communication
- Goroutines provide cheap concurrency
- Decouples message publishers from handlers
- Supports both sync and async dispatch

**Alternatives Considered**:
1. **Direct function calls** - Rejected: Couples components
2. **Event emitter pattern** - Considered: Channel-based is more idiomatic in Go
3. **External message queue (Redis/RabbitMQ)** - Rejected for Phase 1: Adds complexity, not needed initially

**Implementation**:
```go
type MessageBus interface {
    Publish(msg Message) error
    Subscribe(msgType string, handler Handler) error
}

type defaultMessageBus struct {
    subscribers map[string][]Handler
    messages    chan Message
    mu          sync.RWMutex
}

// Concurrent processing
func (mb *defaultMessageBus) process() {
    for msg := range mb.messages {
        handlers := mb.subscribers[msg.Type()]
        for _, h := range handlers {
            go h.Handle(msg)  // Each handler in own goroutine
        }
    }
}
```

**Trade-offs**:
- ✅ Pro: Simple, performant, idiomatic Go
- ✅ Pro: No external dependencies
- ⚠️ Con: Message ordering not guaranteed (acceptable for Phase 1)
- ⚠️ Con: No persistence (acceptable for Phase 1, Phase 5 adds this)

---

### Decision 4: Router Interface Abstraction

**Decision**: Define minimal router interface with Chi as default implementation.

**Rationale**:
- Most frameworks lock you into a router (Gin, Echo, Fiber)
- Interface abstraction allows swapping routers
- Chi is lightweight, compatible with stdlib
- Can add Fiber, Gin adapters later

**Interface Design**:
```go
type Router interface {
    GET(path string, handler HandlerFunc)
    POST(path string, handler HandlerFunc)
    PUT(path string, handler HandlerFunc)
    DELETE(path string, handler HandlerFunc)
    PATCH(path string, handler HandlerFunc)
    Group(prefix string) Router
    Use(middleware ...MiddlewareFunc)
    Listen(addr string) error
}
```

**Why Chi**:
- Compatible with stdlib `http.Handler`
- Lightweight (no heavy dependencies)
- Good middleware ecosystem
- Battle-tested in production

**Adapter Pattern**:
```go
type ChiRouter struct {
    mux *chi.Mux
}

func (r *ChiRouter) GET(path string, h HandlerFunc) {
    r.mux.Get(path, r.adapt(h))
}

func (r *ChiRouter) adapt(h HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, req *http.Request) {
        ctx := NewContext(w, req, r.container)
        h(ctx)
    }
}
```

---

### Decision 5: Configuration Format (YAML with Frontmatter)

**Decision**: Use YAML with frontmatter for all configuration files.

**Rationale**:
- YAML is human-readable and popular in Go ecosystem
- Frontmatter provides metadata separation
- Already using github.com/adrg/frontmatter library
- Familiar to developers (Jekyll, Hugo, etc.)

**Example**:
```yaml
---
name: my-project
version: 1.0.0
---
framework:
  mode: development
  port: 8080

server:
  websocket:
    enabled: true
```

**Alternatives Considered**:
1. **TOML** - Considered: Good option, may add in Phase 5
2. **JSON** - Rejected: Not human-friendly for config
3. **Custom DSL** - Rejected: Unnecessary complexity

**Pluggability**: Config loader is an interface, other formats can be added later.

---

### Decision 6: CLI Framework (Cobra)

**Decision**: Use spf13/cobra for CLI framework.

**Rationale**:
- Industry standard (kubectl, hugo, etc.)
- Excellent subcommand support
- Flag parsing built-in
- Auto-generated help
- Extensible command system

**Command Structure**:
```
touta                    # Root command
├── new <project>        # Create project
├── init                 # Initialize existing dir
├── serve                # Start dev server
└── version              # Show version
```

**Extensibility**:
Packages can add commands:
```
touta
├── auth:create-user     # From auth package
├── blog:seed            # From blog recipe
└── deploy               # From project config
```

---

### Decision 7: Template Renderer (html/template wrapper)

**Decision**: Wrap Go's stdlib `html/template` with an interface for Phase 1.

**Rationale**:
- Stdlib is battle-tested and secure (XSS protection)
- No external dependencies
- Good enough for Phase 1
- Interface allows custom dialect in Phase 2

**Interface**:
```go
type TemplateRenderer interface {
    Render(name string, data interface{}) ([]byte, error)
    RegisterFunction(name string, fn interface{})
    Parse(pattern string) error
}
```

**Phase 2** will introduce custom `<box:*>` dialect, but stdlib provides foundation.

---

### Decision 8: Hot Reload (cosmtrek/air)

**Decision**: Integrate cosmtrek/air for development hot reload.

**Rationale**:
- Most popular Go hot reload tool
- Watches file changes and auto-restarts
- Configurable via `.air.toml`
- Used by many Go projects (Fiber, etc.)

**Integration**:
```bash
# .air.toml
[build]
  cmd = "go build -o ./tmp/main cmd/myproject/main.go"
  bin = "tmp/main"
  
[watch]
  include_ext = ["go", "yaml", "html"]
```

**Developer Experience**:
```bash
$ touta serve
# Uses air under the hood
# Auto-reloads on file changes
```

---

## Data Models

### Message Structure

```go
// Base message interface
type Message interface {
    Slug() string        // Unique identifier (e.g., "user.registered")
    Type() string        // Category (e.g., "event", "command")
    Metadata() map[string]interface{}
}

// Base implementation
type BaseMessage struct {
    MessageSlug string                 `yaml:"slug" validate:"required"`
    MessageType string                 `yaml:"type" validate:"required"`
    Meta        map[string]interface{} `yaml:"metadata,omitempty"`
}

// Concrete message example
type UserRegistered struct {
    BaseMessage
    Email    string `validate:"required,email"`
    Username string `validate:"required,min=3"`
}
```

### Container Binding

```go
type binding struct {
    concrete interface{}                       // Concrete implementation
    factory  func(Container) (interface{}, error) // Factory function
    shared   bool                              // Is singleton?
}

type Container struct {
    bindings   map[string]binding        // Type -> binding
    singletons map[string]interface{}    // Type -> instance
    mu         sync.RWMutex
}
```

### Router Context

```go
type Context interface {
    Request() *http.Request
    Response() http.ResponseWriter
    Param(key string) string
    Get(key string) interface{}
    Set(key string, value interface{})
    Container() Container
}

type defaultContext struct {
    req       *http.Request
    res       http.ResponseWriter
    params    map[string]string
    container Container
    data      map[string]interface{}
}
```

---

## Risks and Mitigations

### Risk 1: Reflection Performance Overhead

**Risk**: Auto-wiring uses reflection, which is slower than direct assignment.

**Mitigation**:
- Auto-wiring happens at startup, not per-request
- Cache reflection results
- Provide manual wiring alternative
- Benchmark and optimize hot paths

**Acceptable**: Startup time increase is acceptable for better DX.

---

### Risk 2: Interface Complexity

**Risk**: Too many interfaces may overwhelm new contributors.

**Mitigation**:
- Clear documentation with examples
- Godoc comments on all interfaces
- Example implementations for reference
- Quick start guide with common patterns

---

### Risk 3: Message Bus Ordering

**Risk**: Concurrent handler execution means message ordering not guaranteed.

**Mitigation**:
- Document non-deterministic ordering
- Add optional ordering support in Phase 5 if needed
- For Phase 1, ordering not critical
- Use message metadata for sequencing if needed

**Acceptable**: Most use cases don't require strict ordering.

---

### Risk 4: Chi Router Lock-In

**Risk**: Despite interface, might accumulate Chi-specific patterns.

**Mitigation**:
- Strict adherence to interface abstraction
- No Chi types in public APIs
- Document adapter pattern clearly
- Create second adapter (Fiber or Gin) in Phase 2 to validate abstraction

---

## Migration Plan

**N/A** - This is the initial implementation, no migration needed.

---

## Testing Strategy

### Unit Tests
- Each component in isolation
- Mock all dependencies
- Table-driven tests for variations
- Target: >80% coverage

### Integration Tests
- Message flow end-to-end
- HTTP request → handler → response
- Config loading and hot reload
- DI resolution with real components

### Benchmarks
- DI container resolution speed
- Message routing throughput
- Template rendering performance
- HTTP request latency baseline

### Example Tests
```go
func TestContainer_AutoWire(t *testing.T) {
    container := di.NewContainer()
    container.Bind((*StorageAdapter)(nil), &MemoryStorage{})
    
    handler := &UserHandler{}
    err := container.AutoWire(handler)
    
    assert.NoError(t, err)
    assert.NotNil(t, handler.storage)
}

func TestMessageBus_PublishAndSubscribe(t *testing.T) {
    bus := message.NewBus()
    
    received := false
    handler := &mockHandler{
        handleFunc: func(msg Message) (Message, error) {
            received = true
            return nil, nil
        },
    }
    
    bus.Subscribe("user.registered", handler)
    bus.Publish(&UserRegistered{
        BaseMessage: BaseMessage{
            MessageSlug: "user.registered",
            MessageType: "event",
        },
        Email: "test@example.com",
    })
    
    time.Sleep(10 * time.Millisecond)  // Allow async processing
    assert.True(t, received)
}
```

---

## Open Questions

1. **Q**: Should we support synchronous message dispatch in Phase 1?
   **A**: Yes, add optional `PublishSync()` method for blocking dispatch.

2. **Q**: How to handle handler errors in message bus?
   **A**: Log errors, optionally publish error messages to error handlers.

3. **Q**: Should DI container support named bindings?
   **A**: Yes, allow tagging for multiple implementations of same interface.

4. **Q**: Config validation - strict or lenient?
   **A**: Lenient in Phase 1, add schema validation in Phase 2.

---

## Performance Targets

### Baseline Targets (Phase 1)
- **DI Resolution**: <1ms for typical dependency graph
- **Message Routing**: >10,000 msgs/sec on modest hardware
- **HTTP Requests**: <5ms p99 latency for simple handler
- **Template Rendering**: <10ms for typical page
- **Memory Usage**: <50MB for framework core

### Measurement
- Benchmark tests in CI
- Profiling with pprof
- Document baseline results
- Track regressions

---

## Dependencies

### Go Modules
```go
require (
    github.com/go-chi/chi/v5 v5.0.10
    github.com/spf13/cobra v1.8.0
    github.com/adrg/frontmatter v0.2.0
    github.com/go-playground/validator/v10 v10.16.0
    gopkg.in/yaml.v3 v3.0.1
)

require (
    github.com/cosmtrek/air v1.49.0 // dev
    github.com/stretchr/testify v1.8.4 // test
)
```

### External Tools
- **golangci-lint** - Linting
- **air** - Hot reload
- **GitHub Actions** - CI/CD

---

## Timeline and Milestones

### Week 1-2: Foundation
- Project structure
- Core interfaces defined
- DI container implementation

### Week 3-4: Message System
- Message bus implementation
- Routing (config + code)
- Handler pattern

### Week 5-6: HTTP and Config
- Router abstraction (Chi)
- Context implementation
- Config loader (YAML)

### Week 7-8: CLI and Polish
- CLI commands (new, init, serve)
- Template renderer
- Hot reload integration
- Documentation and examples

---

## Success Metrics

1. ✅ All core interfaces defined and documented
2. ✅ DI container resolves dependencies correctly
3. ✅ Message bus routes messages to handlers
4. ✅ HTTP server handles requests via Chi
5. ✅ Config loads from YAML files
6. ✅ CLI creates projects and starts dev server
7. ✅ Hot reload works in development
8. ✅ Test coverage >80%
9. ✅ Example "Hello World" project works
10. ✅ Documentation complete and clear

---

**Review and Approval**: This design should be reviewed before implementation begins. Any changes to core interfaces after Phase 1 ships will be breaking changes.
