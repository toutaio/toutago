# Nasc Dependency Injector - Technical Design

## Context

Dependency injection is fundamental to Toutā's philosophy of interface-first design, loose coupling, and testability. While the current `internal/di` package works, creating Nasc as an independent, feature-rich library provides better reusability, maintainability, and aligns with Toutā's vision of composable, independent components.

## Goals

- Create a standalone, production-ready DI container for Go
- Support all common DI patterns (transient, singleton, scoped, factory)
- Provide excellent developer experience with auto-wiring and service providers
- Maintain high performance with minimal reflection overhead
- Enable easy integration with Toutā and other Go projects
- Follow Go idioms and best practices

## Non-Goals

- Not a full IoC framework (just dependency injection)
- Not focused on XML/YAML configuration
- Not attempting backward compatibility with other DI libraries
- Not supporting code generation (pure runtime)

## Architecture

### High-Level Components

```
┌─────────────────────────────────────────────────────────┐
│                    Nasc Container                       │
├─────────────────────────────────────────────────────────┤
│                                                          │
│  ┌────────────────────────────────────────────────┐    │
│  │          Binding Registry                       │    │
│  │  • Interface → Implementation mappings          │    │
│  │  • Lifetime management (transient/singleton)    │    │
│  │  • Named bindings                               │    │
│  │  • Tagged bindings                              │    │
│  └────────────────────────────────────────────────┘    │
│                                                          │
│  ┌────────────────────────────────────────────────┐    │
│  │           Resolution Engine                     │    │
│  │  • Type reflection                              │    │
│  │  • Constructor injection                        │    │
│  │  • Auto-wiring                                  │    │
│  │  • Circular dependency detection                │    │
│  └────────────────────────────────────────────────┘    │
│                                                          │
│  ┌────────────────────────────────────────────────┐    │
│  │          Lifecycle Management                   │    │
│  │  • Singleton cache                              │    │
│  │  • Scope management                             │    │
│  │  • Disposal handling                            │    │
│  └────────────────────────────────────────────────┘    │
│                                                          │
│  ┌────────────────────────────────────────────────┐    │
│  │         Service Providers                       │    │
│  │  • Register phase                               │    │
│  │  • Boot phase                                   │    │
│  │  • Dependency ordering                          │    │
│  └────────────────────────────────────────────────┘    │
│                                                          │
└─────────────────────────────────────────────────────────┘
```

### Core Data Structures

```go
// Main container
type Nasc struct {
    mu        sync.RWMutex
    bindings  map[reflect.Type]*bindingRegistry
    singletons map[reflect.Type]interface{}
    providers  []ServiceProvider
    parent    *Nasc  // For scoped containers
}

// Binding represents a registration
type Binding struct {
    interfaceType reflect.Type
    concreteType  reflect.Type
    lifetime      Lifetime
    factory       FactoryFunc
    instance      interface{}
    tags          []string
    name          string
    condition     Condition
}

// Lifetime types
type Lifetime int
const (
    Transient Lifetime = iota
    Singleton
    Scoped
)

// Service provider interface
type ServiceProvider interface {
    Register(nasc *Nasc) error
    Boot(nasc *Nasc) error
}

// Factory function signature
type FactoryFunc func(*Nasc) (interface{}, error)

// Constructor function signature  
type ConstructorFunc interface{}
```

### Binding Fluent API

```go
// Fluent binding API
nasc := NewNasc()

// Basic binding
nasc.Bind((*Logger)(nil), &ConsoleLogger{})

// Singleton
nasc.Singleton((*Database)(nil), &PostgresDB{})

// Scoped
nasc.Scoped((*UnitOfWork)(nil), &DbUnitOfWork{})

// Factory
nasc.Factory((*Connection)(nil), func(n *Nasc) (interface{}, error) {
    config := n.Make((*Config)(nil)).(*Config)
    return NewConnection(config.DSN)
})

// Named binding
nasc.Bind((*Logger)(nil), &FileLogger{}).Named("file")

// Tagged binding
nasc.Bind((*Middleware)(nil), &AuthMiddleware{}).Tag("http")

// Conditional binding
nasc.Bind((*Cache)(nil), &RedisCache{}).When(isProduction)

// Instance binding
nasc.Instance((*Config)(nil), existingConfig)
```

### Resolution API

```go
// Basic resolution
logger := nasc.Make((*Logger)(nil)).(Logger)

// Type-safe helper
logger, err := MakeTyped[Logger](nasc)

// Named resolution
fileLogger := nasc.MakeNamed((*Logger)(nil), "file").(Logger)

// Tagged resolution
middlewares := nasc.Tagged("http")

// With parameters
service := nasc.MakeWith((*Service)(nil), param1, param2)
```

### Auto-Wiring

```go
type UserService struct {
    Logger   Logger   `inject:""`
    Database Database `inject:""`
    Cache    Cache    `inject:"optional"`
    FileLog  Logger   `inject:"name=file"`
}

// Auto-wire instance
service := &UserService{}
err := nasc.AutoWire(service)

// Or use constructor with auto-resolution
nasc.Bind((*UserService)(nil), &UserService{}).AutoWire()
```

### Service Providers

```go
type DatabaseProvider struct{}

func (p *DatabaseProvider) Register(nasc *Nasc) error {
    // Register bindings
    nasc.Singleton((*Database)(nil), &PostgresDB{})
    nasc.Scoped((*Transaction)(nil), &DbTransaction{})
    return nil
}

func (p *DatabaseProvider) Boot(nasc *Nasc) error {
    // Resolve and initialize
    db := nasc.Make((*Database)(nil)).(Database)
    return db.Connect()
}

// Register providers
nasc.RegisterProvider(&DatabaseProvider{})
nasc.RegisterProvider(&CacheProvider{})
nasc.Boot() // Calls all providers in order
```

### Scoping

```go
// Create scope (e.g., per HTTP request)
scope := nasc.CreateScope()
defer scope.Dispose()

// Resolve in scope
unitOfWork := scope.Make((*UnitOfWork)(nil)).(UnitOfWork)

// Scoped instances are reused within this scope
same := scope.Make((*UnitOfWork)(nil)).(UnitOfWork)
// unitOfWork == same

// But different in another scope
scope2 := nasc.CreateScope()
different := scope2.Make((*UnitOfWork)(nil)).(UnitOfWork)
// unitOfWork != different
```

## Implementation Strategy

### Phase 1: Core Container (Week 1)
- Basic container structure
- Binding registry
- Simple transient resolution
- Interface to concrete mapping

### Phase 2: Lifetimes (Week 2)
- Singleton implementation
- Scoped lifetime
- Factory functions
- Instance binding

### Phase 3: Auto-Wiring (Week 3)
- Struct tag parsing
- Field injection
- Optional dependencies
- Named injection

### Phase 4: Constructor Injection (Week 4)
- Constructor function support
- Parameter resolution
- Error handling
- Variadic constructors

### Phase 5: Advanced Features (Week 5)
- Service providers
- Tagged services
- Named bindings
- Conditional resolution

### Phase 6: Circular Dependencies (Week 6)
- Detection algorithm
- Clear error messages
- Lazy resolution support
- Resolution chain tracking

### Phase 7: Scoping & Cleanup (Week 7)
- Scope creation and management
- Disposal pattern
- Cleanup callbacks
- Resource management

### Phase 8: Performance (Week 8)
- Resolution caching
- Reflection optimization
- Benchmark suite
- Memory profiling

## Key Design Decisions

### Decision: Reflection vs Code Generation
**Choice:** Runtime reflection
**Rationale:**
- Simpler to use (no build step)
- More flexible
- Acceptable performance with caching
- Go idioms favor reflection for DI

**Alternatives Considered:**
- Code generation: Too complex, poor DX
- Build tags: Limited flexibility

### Decision: Interface-First Binding
**Choice:** Bind interfaces to implementations
**Rationale:**
- Enforces interface-first design
- Aligns with Toutā philosophy
- Better testability
- Clear contracts

### Decision: Fluent API
**Choice:** Chainable method calls
**Rationale:**
- Better discoverability
- Cleaner syntax
- Self-documenting
- Common pattern in DI libraries

### Decision: Service Providers
**Choice:** Two-phase (Register + Boot)
**Rationale:**
- Clear separation of registration and initialization
- Enables dependency ordering
- Better error handling
- Familiar pattern from Laravel/ASP.NET

### Decision: Thread Safety
**Choice:** Full thread-safety with RWMutex
**Rationale:**
- Essential for production use
- Concurrent web servers need it
- Minimal performance impact with RW locks
- Safe singleton creation

## Error Handling

### Clear Error Messages

```go
// Missing dependency
"cannot resolve type *UserService: dependency Logger not registered"

// Circular dependency
"circular dependency detected: *UserService → *Logger → *UserService"

// Type mismatch
"type *FileLogger does not implement interface Logger"

// Resolution chain
"failed to resolve *Controller → *Service → *Repository: missing Database binding"
```

### Error Types

```go
type BindingError struct {
    Type    reflect.Type
    Message string
}

type ResolutionError struct {
    Type  reflect.Type
    Chain []reflect.Type
    Cause error
}

type CircularDependencyError struct {
    Chain []reflect.Type
}
```

## Performance Considerations

### Optimization Strategies

1. **Reflection Caching**
   - Cache type info after first use
   - Reuse reflection results
   - Minimize Type() calls

2. **Singleton Fast Path**
   - Direct map lookup for singletons
   - No reflection after first creation
   - Lock-free reads with sync.Map

3. **Auto-Wire Metadata Cache**
   - Parse struct tags once
   - Store field injection metadata
   - Reuse across instances

4. **Memory Pooling**
   - Pool transient instances (optional)
   - Reduce GC pressure
   - Configurable pool sizes

### Benchmark Targets

- **Singleton resolution:** <100ns
- **Transient resolution:** <1μs (simple types)
- **Auto-wire:** <10μs (typical struct)
- **Memory:** <100 bytes overhead per binding

## Testing Strategy

### Unit Tests
- Test each binding mode in isolation
- Test auto-wiring edge cases
- Test circular dependency detection
- Test thread-safety with race detector

### Integration Tests
- Test real-world scenarios
- Test with actual Toutā code
- Test complex dependency graphs
- Test service provider ordering

### Benchmarks
- Resolution performance
- Memory allocations
- Concurrent access
- Compare with wire, dig, fx

## Migration from internal/di

### Compatibility Layer
```go
// Adapter for current Toutā code
type LegacyContainer struct {
    nasc *Nasc
}

func (c *LegacyContainer) Bind(iface, impl interface{}) error {
    return c.nasc.Bind(iface, impl).Error()
}

// Gradual migration path
```

### Migration Steps
1. Create Nasc package
2. Implement core features
3. Build compatibility layer
4. Migrate Toutā code incrementally
5. Remove old internal/di
6. Update documentation

## Integration with Toutā

```go
// In Toutā framework
import "github.com/toutaio/toutago-nasc-dependency-injector"

type Application struct {
    container *nasc.Nasc
}

func NewApplication() *Application {
    n := nasc.New()
    
    // Register core providers
    n.RegisterProvider(&RoutingProvider{})
    n.RegisterProvider(&ConfigProvider{})
    n.RegisterProvider(&MessageBusProvider{})
    
    // Boot all providers
    n.Boot()
    
    return &Application{container: n}
}
```

## Success Metrics

- ✅ >90% test coverage
- ✅ <100ns singleton resolution
- ✅ Zero data races
- ✅ Comprehensive documentation
- ✅ Toutā integration complete
- ✅ Community feedback positive

## Timeline

- Week 1: Core container
- Week 2: Lifetimes
- Week 3: Auto-wiring
- Week 4: Constructor injection
- Week 5: Advanced features
- Week 6: Circular dependency handling
- Week 7: Scoping and cleanup
- Week 8: Performance optimization
- Week 9: Documentation
- Week 10: Toutā integration

Target: v1.0.0 in ~10 weeks
