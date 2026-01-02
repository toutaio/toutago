# Design: Cosan Router Architecture

## Context

Cosan (Irish for "pathway") is an independent HTTP router component designed to embody SOLID principles and demonstrate toutā's architectural philosophy while being usable in any Go project. The router must balance performance, simplicity, extensibility, and maintainability.

### Background
- Most Go routers couple routing logic with specific patterns (Chi's middleware, Gin's context)
- Existing routers either prioritize performance (Fiber) or developer experience (Echo)
- Few routers strictly follow SOLID principles or interface-driven design
- The Go ecosystem lacks a canonical example of message-passing compatible routing

### Constraints
- Must work with standard `net/http` package
- Zero framework dependencies beyond standard library
- Must be performant (within 10% of fastest routers)
- API must be simple and Go-idiomatic
- Must support both imperative and declarative configuration

### Stakeholders
- toutā framework developers and users
- Go community seeking clean router implementations
- Developers building message-centric architectures
- Teams requiring fully testable, mockable routing

## Goals / Non-Goals

### Goals
1. **SOLID Compliance**: Demonstrate all five SOLID principles in practice
2. **Interface-First**: Every component behind well-defined interfaces
3. **Performance**: Competitive with top routers (Chi, Gin, Echo)
4. **Simplicity**: Clean API with minimal learning curve
5. **Extensibility**: Plugin architecture for matchers, middleware, context
6. **Testability**: 100% mockable, no global state
7. **Message-Friendly**: Easy integration with message-passing systems
8. **Production-Ready**: Comprehensive error handling, logging hooks, observability

### Non-Goals
1. **HTTP/3 or QUIC support** - Focus on HTTP/1.1 and HTTP/2
2. **Built-in validation** - Validation is a separate concern
3. **Built-in serialization** - Keep encoding pluggable
4. **WebSocket handling** - Delegate to specialized libraries
5. **Rate limiting** - Implement as middleware, not core feature
6. **Authentication** - Implement as middleware, not core feature

## Decisions

### Decision 1: Radix Tree for Route Matching
**What:** Use a radix tree (compressed trie) for route matching algorithm

**Why:**
- O(k) lookup time where k is path length (vs O(n) linear search)
- Memory efficient with common prefixes
- Industry proven (used by Chi, Echo, Gin)
- Supports parameter extraction efficiently

**Alternatives Considered:**
- **Hash map:** Fast but doesn't support path parameters well
- **Linear search:** Simple but O(n) performance
- **Full trie:** More memory usage than radix tree

**Trade-offs:**
- More complex implementation than linear search
- Requires careful concurrency management
- Worth it for performance at scale

### Decision 2: Interface Segregation for Context
**What:** Split context into multiple focused interfaces rather than one large interface

**Interfaces:**
```go
type ParamReader interface {
    Param(key string) string
    Params() map[string]string
}

type QueryReader interface {
    Query(key string) string
    QueryAll(key string) []string
}

type BodyReader interface {
    Bind(v interface{}) error
    BodyBytes() ([]byte, error)
}

type ResponseWriter interface {
    JSON(code int, v interface{}) error
    String(code int, format string, args ...interface{})
    Status(code int)
}

// Context composes all interfaces
type Context interface {
    ParamReader
    QueryReader
    BodyReader
    ResponseWriter
    Request() *http.Request
    ResponseWriter() http.ResponseWriter
}
```

**Why:**
- Interface Segregation Principle: Clients depend only on methods they use
- Easier mocking in tests
- Clear separation of concerns
- Allows partial implementations

**Alternatives Considered:**
- **Single interface:** Simpler but violates ISP
- **No context abstraction:** More flexible but less convenient

### Decision 3: Functional Options Pattern for Router Configuration
**What:** Use functional options for router configuration

```go
type Option func(*Router)

func WithMatcher(m Matcher) Option {
    return func(r *Router) { r.matcher = m }
}

func New(opts ...Option) *Router {
    r := &Router{
        matcher: NewRadixMatcher(),
        // defaults...
    }
    for _, opt := range opts {
        opt(r)
    }
    return r
}
```

**Why:**
- Open/Closed Principle: Extensible without modifying constructor
- Optional parameters without multiple constructors
- Self-documenting API
- Go community standard pattern

**Alternatives Considered:**
- **Config struct:** Less discoverable, requires documentation
- **Builder pattern:** More verbose in Go

### Decision 4: Middleware as Interface + Adapter
**What:** Define middleware as interface with adapter for functions

```go
type Middleware interface {
    Process(next HandlerFunc) HandlerFunc
}

type MiddlewareFunc func(HandlerFunc) HandlerFunc

func (mw MiddlewareFunc) Process(next HandlerFunc) HandlerFunc {
    return mw(next)
}
```

**Why:**
- Supports both functional and object-oriented middleware
- Easy to test (mock interface)
- Composable and reusable
- Follows Liskov Substitution Principle

**Alternatives Considered:**
- **Function only:** Less testable, can't hold state elegantly
- **Interface only:** Less convenient for simple middleware

### Decision 5: Dependency Injection via Constructor
**What:** All dependencies injected via constructor, never via globals

```go
type Router struct {
    matcher    Matcher
    logger     Logger
    errorHandler ErrorHandler
}

func New(matcher Matcher, logger Logger, handler ErrorHandler) *Router {
    return &Router{
        matcher: matcher,
        logger: logger,
        errorHandler: handler,
    }
}
```

**Why:**
- Dependency Inversion Principle: Depend on abstractions
- Fully testable with mocks
- No global state
- Clear dependencies

**Alternatives Considered:**
- **Service locator:** Hidden dependencies, testing harder
- **Global instances:** Not testable, concurrency issues

### Decision 6: Immutable Route Registration
**What:** Routes immutable after registration; router compilation step freezes routes

```go
func (r *Router) Compile() error {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    if r.compiled {
        return ErrAlreadyCompiled
    }
    
    // Build optimized route tree
    r.tree = buildTree(r.routes)
    r.compiled = true
    return nil
}
```

**Why:**
- Thread-safe serving without locks in hot path
- Explicit compilation allows optimizations
- Prevents runtime route modifications
- Clear separation: configuration vs execution

**Alternatives Considered:**
- **Mutable routes:** Requires locking on every request
- **Automatic compilation:** Less explicit, harder to test

## Architecture Diagrams

### Component Structure
```
┌─────────────────────────────────────────────┐
│              Router (main)                   │
│  - Route registration                        │
│  - Middleware management                     │
│  - ServeHTTP implementation                  │
└─────────────────┬───────────────────────────┘
                  │
        ┌─────────┼─────────┐
        │         │         │
        ▼         ▼         ▼
   ┌────────┐ ┌──────┐ ┌─────────┐
   │Matcher │ │Logger│ │ErrorHdlr│
   │(Radix) │ │      │ │         │
   └────────┘ └──────┘ └─────────┘
        │
        ▼
   ┌──────────────┐
   │  RouteTree   │
   │ (compiled)   │
   └──────────────┘
```

### Request Flow
```
HTTP Request
     │
     ▼
┌─────────────┐
│ ServeHTTP   │
└─────┬───────┘
      │
      ▼
┌─────────────┐
│  Matcher    │ ──── Match route & extract params
└─────┬───────┘
      │
      ▼
┌─────────────┐
│ Middleware  │ ──── Chain execution (outer to inner)
│   Chain     │
└─────┬───────┘
      │
      ▼
┌─────────────┐
│   Handler   │ ──── Business logic
└─────┬───────┘
      │
      ▼
┌─────────────┐
│  Response   │
└─────────────┘
```

### SOLID Principles Applied

**Single Responsibility:**
- `Router`: Route management and HTTP serving
- `Matcher`: Route matching only
- `Context`: Request/response abstraction only
- `Middleware`: Request transformation only

**Open/Closed:**
- Functional options for router configuration
- Middleware interface for extensibility
- Custom matcher implementations

**Liskov Substitution:**
- All matcher implementations interchangeable
- All middleware implementations interchangeable
- Context implementations fully substitutable

**Interface Segregation:**
- Context split into ParamReader, QueryReader, BodyReader, ResponseWriter
- Clients depend only on needed interfaces

**Dependency Inversion:**
- Router depends on Matcher interface, not concrete implementation
- Logger interface instead of concrete logger
- All dependencies injected, never instantiated internally

## Risks / Trade-offs

### Risk 1: Performance vs Abstraction
**Risk:** Interface overhead might impact performance

**Mitigation:**
- Benchmark early and often
- Use interface inlining where possible
- Keep hot paths (matching, serving) minimal
- Consider sync.Pool for high-allocation objects

**Trade-off:** Accept small performance cost for testability and flexibility

### Risk 2: API Complexity
**Risk:** Too many interfaces might confuse users

**Mitigation:**
- Provide sensible defaults for everything
- Comprehensive examples and docs
- Simple use cases work out-of-box
- Advanced features clearly documented

**Trade-off:** Power users get full control, beginners use defaults

### Risk 3: Compatibility with Standard Library
**Risk:** Custom context might not integrate well with existing middleware

**Mitigation:**
- Provide adapter for standard http.Handler
- Keep context.Request() and context.ResponseWriter() accessible
- Document integration patterns

**Trade-off:** Some verbosity in adapter code

## Migration Plan

### Phase 1: Standalone Development
1. Develop Cosan in isolation at `/home/nestor/Proyects/toutago-cosan-router`
2. Publish as independent Go module
3. Build community adoption outside toutā

### Phase 2: Toutā Integration
1. Create adapter implementing toutā's Router interface
2. Update toutā examples to use Cosan
3. Document migration from Chi/other routers
4. Maintain backward compatibility

### Phase 3: Ecosystem Growth
1. Build middleware collection
2. Create integration examples
3. Support community contributions
4. Version according to semantic versioning

### Rollback Strategy
- Cosan is independent; no impact on toutā if issues arise
- toutā maintains support for Chi, Gin, Echo
- Users can switch routers via configuration

## Performance Targets

### Benchmarks
- **Route matching:** < 500ns per request (simple routes)
- **Parameter extraction:** < 100ns additional overhead
- **Memory allocation:** < 2 allocations per request in hot path
- **Concurrent requests:** Linear scaling to 100k requests/sec

### Comparison Targets
- Within 10% of Chi (current toutā default)
- Within 15% of Gin (fastest popular router)
- Better than stdlib ServeMux on complex routing

## Open Questions

1. **Should we support HTTP method wildcards (ANY)?**
   - Leaning: Yes, for flexibility, but document anti-patterns

2. **How to handle trailing slashes?**
   - Options: Strict matching, auto-redirect, configuration option
   - Leaning: Configuration option with sensible default (auto-redirect)

3. **Support for route versioning (v1, v2 prefixes)?**
   - Leaning: Handle via route groups, not special syntax

4. **Custom error responses format?**
   - Leaning: User-provided ErrorHandler interface

5. **Observability hooks (metrics, tracing)?**
   - Leaning: Middleware-based, provide examples

---

**Version:** 1.0 (Initial Design)  
**Status:** Proposed  
**Last Updated:** 2025-12-29
