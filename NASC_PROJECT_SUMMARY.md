# Nasc Dependency Injector - Complete Implementation Plan

## Overview

**Nasc** (Old Irish: "link" or "bond") is an independent, powerful dependency injection container for Go being developed as a separate project from Toutā.

## Project Details

- **Name:** Nasc Dependency Injector
- **Repository:** https://github.com/toutaio/toutago-nasc-dependency-injector
- **Local Path:** `~/Proyects/toutago-nasc-dependency-injector`
- **Language:** Go 1.21+
- **License:** MIT (recommended)
- **Status:** Planning Complete ✅

## Why Separate Project?

1. **Reusability** - Can be used in any Go project
2. **Independence** - Develops at its own pace
3. **Clean Architecture** - No coupling with Toutā internals
4. **Versioning** - Independent semantic versioning
5. **Community** - Own contributor base

## Key Features

### Core Capabilities
- ✅ **Interface-based binding** - Bind interfaces to implementations
- ✅ **Multiple lifetimes** - Transient, Singleton, Scoped, Factory
- ✅ **Auto-wiring** - Automatic field injection via struct tags
- ✅ **Constructor injection** - Resolve constructor parameters automatically
- ✅ **Service providers** - Organize related bindings (Register + Boot phases)
- ✅ **Tagged services** - Group and resolve services by tags
- ✅ **Named bindings** - Multiple implementations of same interface
- ✅ **Conditional resolution** - Context-based binding selection
- ✅ **Circular dependency detection** - Clear error messages with resolution chain
- ✅ **Scoping** - Per-request or custom scope lifetimes
- ✅ **Thread-safe** - Concurrent access with RWMutex
- ✅ **Performance optimized** - Reflection caching, minimal allocations

### Developer Experience
- Fluent binding API
- Clear error messages with resolution chains
- Debug mode with verbose logging
- Comprehensive documentation
- Rich examples

## API Examples

### Basic Usage
```go
import "github.com/toutaio/toutago-nasc-dependency-injector"

// Create container
nasc := nasc.New()

// Bind interface to implementation
nasc.Bind((*Logger)(nil), &ConsoleLogger{})

// Singleton
nasc.Singleton((*Database)(nil), &PostgresDB{})

// Resolve
logger := nasc.Make((*Logger)(nil)).(Logger)
```

### Auto-Wiring
```go
type UserService struct {
    Logger   Logger   `inject:""`
    Database Database `inject:""`
    Cache    Cache    `inject:"optional"`
    FileLog  Logger   `inject:"name=file"`
}

service := &UserService{}
nasc.AutoWire(service)
```

### Service Providers
```go
type DatabaseProvider struct{}

func (p *DatabaseProvider) Register(nasc *Nasc) error {
    nasc.Singleton((*Database)(nil), &PostgresDB{})
    return nil
}

func (p *DatabaseProvider) Boot(nasc *Nasc) error {
    db := nasc.Make((*Database)(nil)).(Database)
    return db.Connect()
}

nasc.RegisterProvider(&DatabaseProvider{})
nasc.Boot()
```

### Scoping
```go
// Create scope (e.g., per HTTP request)
scope := nasc.CreateScope()
defer scope.Dispose()

// Scoped instances reused within scope
uow := scope.Make((*UnitOfWork)(nil)).(UnitOfWork)
```

## Project Structure

```
~/Proyects/toutago-nasc-dependency-injector/
├── README.md
├── IMPLEMENTATION_PLAN.md
├── LICENSE
├── go.mod
│
├── nasc.go                  # Main container
├── binding.go               # Binding types
├── lifetime.go              # Lifetime management
├── errors.go                # Error types
│
├── resolver/                # Resolution engine
├── registry/                # Binding registry
├── scope/                   # Scoping support
├── provider/                # Service providers
├── cache/                   # Performance caching
│
├── examples/                # Usage examples
├── docs/                    # Documentation
└── benchmarks/              # Performance tests
```

## Implementation Timeline

### Phase 1: Core Container (Week 1)
- Basic container structure
- Binding registry
- Simple bind/resolve

### Phase 2: Lifetimes (Week 2)
- Transient, Singleton, Scoped, Factory
- Lifetime management

### Phase 3: Auto-Wiring (Week 3)
- Struct tag parsing
- Field injection
- Optional dependencies

### Phase 4: Constructor Injection (Week 4)
- Constructor support
- Parameter resolution
- Error handling

### Phase 5: Service Providers (Week 5)
- Register/Boot phases
- Provider ordering
- Dependency resolution

### Phase 6: Advanced Features (Week 6)
- Tagged services
- Named bindings
- Conditional resolution

### Phase 7: Error Handling (Week 7)
- Circular dependency detection
- Clear error messages
- Debug mode

### Phase 8: Scoping & Cleanup (Week 8)
- Scope management
- Resource disposal
- Cleanup callbacks

### Phase 9: Performance (Week 9)
- Resolution caching
- Reflection optimization
- Benchmarking

### Phase 10: Documentation & Integration (Week 10)
- Complete docs
- Toutā integration
- Examples

**Target:** v1.0.0 in ~10 weeks

## OpenSpec Documentation

All requirements, tasks, and design decisions documented in:

```
/home/nestor/Proyects/toutago/openspec/changes/create-nasc-dependency-injector/
├── proposal.md              # Why and what
├── tasks.md                 # 93 implementation tasks
├── design.md                # Technical architecture
└── specs/
    └── nasc-di/
        └── spec.md          # 26 detailed requirements
```

## Features Detail

### Binding Modes

**Transient** - New instance every time
```go
nasc.Bind((*Service)(nil), &MyService{})
```

**Singleton** - Same instance always
```go
nasc.Singleton((*DB)(nil), &PostgresDB{})
```

**Scoped** - Same instance per scope
```go
nasc.Scoped((*UnitOfWork)(nil), &DbUnitOfWork{})
```

**Factory** - Custom creation logic
```go
nasc.Factory((*Connection)(nil), func(n *Nasc) (interface{}, error) {
    return NewConnection(config)
})
```

### Advanced Features

**Named Bindings**
```go
nasc.Bind((*Logger)(nil), &FileLogger{}).Named("file")
logger := nasc.MakeNamed((*Logger)(nil), "file")
```

**Tagged Services**
```go
nasc.Bind((*Middleware)(nil), &AuthMiddleware{}).Tag("http")
middlewares := nasc.Tagged("http")
```

**Conditional Resolution**
```go
nasc.Bind((*Cache)(nil), &RedisCache{}).When(isProduction)
```

## Integration with Toutā

Once Nasc is ready, Toutā will:

1. **Import as module:**
   ```go
   import "github.com/toutaio/toutago-nasc-dependency-injector"
   ```

2. **Replace internal/di:**
   - Migrate existing bindings
   - Update framework code
   - Remove old DI implementation

3. **Enhance capabilities:**
   - Use service providers for framework modules
   - Leverage scoping for request handling
   - Implement cleanup for resources

## Performance Goals

- **Singleton resolution:** <100ns
- **Transient resolution:** <1μs (simple types)
- **Auto-wire:** <10μs (typical struct)
- **Memory overhead:** <100 bytes per binding
- **Thread-safe** with minimal lock contention

## Testing Strategy

- **Unit tests:** >90% coverage
- **Integration tests:** Real-world scenarios
- **Concurrent tests:** Race detector enabled
- **Benchmarks:** Performance tracking

## Next Steps

### Immediate (Today)
1. ✅ Initialize Git repository
2. ✅ Create project structure
3. ✅ Write implementation plan
4. ✅ Create OpenSpec documentation

### This Week
1. Set up Go module
2. Implement core container
3. Add basic binding
4. Write first tests

### Next Month
1. Complete Phase 1-3 (Core, Lifetimes, Auto-wiring)
2. Start Phase 4 (Constructor injection)

## Resources

- **OpenSpec:** `/home/nestor/Proyects/toutago/openspec/changes/create-nasc-dependency-injector/`
- **Implementation Plan:** `~/Proyects/toutago-nasc-dependency-injector/IMPLEMENTATION_PLAN.md`
- **Inspiration:** dig, wire, fx (Go DI libraries)
- **Patterns:** Service Locator, Constructor Injection, Property Injection

## How to Start

```bash
# Navigate to project
cd ~/Proyects/toutago-nasc-dependency-injector

# Read implementation plan
cat IMPLEMENTATION_PLAN.md

# Check OpenSpec docs
cd /home/nestor/Proyects/toutago
openspec show create-nasc-dependency-injector

# Start coding!
cd ~/Proyects/toutago-nasc-dependency-injector
go mod init github.com/toutaio/toutago-nasc-dependency-injector
```

## Summary

Nasc is a carefully planned, independent DI container that will:

- ✅ Provide complete dependency injection capabilities
- ✅ Maintain Celtic-themed identity
- ✅ Work as standalone Go library
- ✅ Integrate seamlessly with Toutā
- ✅ Support all common DI patterns
- ✅ Offer excellent performance

**Status:** Ready to begin implementation! 🔗

---

**Created:** 2025-12-28
**Planning Status:** Complete ✅  
**Implementation Status:** Not started  
**Target v1.0.0:** ~10 weeks
