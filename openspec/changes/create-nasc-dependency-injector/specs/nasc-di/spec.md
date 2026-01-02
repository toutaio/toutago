## ADDED Requirements

### Requirement: Independent Go Library
The Nasc DI container SHALL be implemented as a standalone Go library that can be used independently of the Toutā framework.

#### Scenario: Import in any Go project
- **WHEN** a developer imports `github.com/toutaio/toutago-nasc-dependency-injector`
- **THEN** they can use Nasc for DI without any Toutā dependencies

#### Scenario: Semantic versioning
- **WHEN** the library is released
- **THEN** it follows semantic versioning for stability guarantees

### Requirement: Interface-Based Binding
The container SHALL support binding interfaces to concrete implementations.

#### Scenario: Bind interface to implementation
- **WHEN** developer calls `nasc.Bind((*Logger)(nil), &ConsoleLogger{})`
- **THEN** the interface is mapped to the concrete type

#### Scenario: Resolve interface to instance
- **WHEN** developer calls `nasc.Make((*Logger)(nil))`
- **THEN** an instance of ConsoleLogger is returned

#### Scenario: Multiple implementations
- **WHEN** multiple implementations exist for an interface
- **THEN** the most recently bound or tagged implementation is resolved

### Requirement: Binding Modes - Transient
The container SHALL support transient binding where a new instance is created on each resolution.

#### Scenario: Transient returns new instance
- **WHEN** interface is bound as transient
- **THEN** each call to Make returns a different instance

#### Scenario: No instance reuse
- **WHEN** same type is resolved twice
- **THEN** the instances are independent

### Requirement: Binding Modes - Singleton
The container SHALL support singleton binding where the same instance is returned for all resolutions.

#### Scenario: Singleton returns same instance
- **WHEN** interface is bound as singleton
- **THEN** all calls to Make return the same instance

#### Scenario: Lazy initialization
- **WHEN** singleton is bound but not yet resolved
- **THEN** instance is created on first resolution

#### Scenario: Thread-safe singleton creation
- **WHEN** multiple goroutines resolve singleton simultaneously
- **THEN** only one instance is created

### Requirement: Binding Modes - Scoped
The container SHALL support scoped binding where instance lifetime is tied to a scope.

#### Scenario: Scoped instance per scope
- **WHEN** a scope is created
- **THEN** scoped bindings return the same instance within that scope

#### Scenario: Different scopes have different instances
- **WHEN** two separate scopes exist
- **THEN** each scope has its own instance of scoped bindings

#### Scenario: Scope disposal
- **WHEN** a scope is disposed
- **THEN** all scoped instances are cleaned up

### Requirement: Binding Modes - Factory
The container SHALL support factory binding where a custom function creates instances.

#### Scenario: Factory function called
- **WHEN** interface is bound with factory function
- **THEN** the factory is invoked to create instances

#### Scenario: Factory receives container
- **WHEN** factory function is called
- **THEN** it receives the container for resolving dependencies

#### Scenario: Factory error handling
- **WHEN** factory returns an error
- **THEN** the error is propagated to the caller

### Requirement: Auto-Wiring
The container SHALL support automatic dependency injection via struct tags.

#### Scenario: Auto-wire tagged fields
- **WHEN** struct has `inject:""` tags on fields
- **THEN** fields are automatically populated with resolved dependencies

#### Scenario: Optional dependencies
- **WHEN** field has `inject:"optional"` tag
- **THEN** nil is allowed if dependency not found

#### Scenario: Named dependencies
- **WHEN** field has `inject:"name=logger"` tag
- **THEN** the named binding is resolved

#### Scenario: Auto-wire nested structs
- **WHEN** injected dependency has its own inject tags
- **THEN** nested dependencies are also resolved

### Requirement: Constructor Injection
The container SHALL support constructor-based dependency injection.

#### Scenario: Constructor with dependencies
- **WHEN** type is registered with constructor `func(Logger, Database) *Service`
- **THEN** dependencies are resolved and passed to constructor

#### Scenario: Constructor errors
- **WHEN** constructor returns error
- **THEN** error is propagated with context

#### Scenario: Variadic constructors
- **WHEN** constructor accepts variadic parameters
- **THEN** all matching dependencies are resolved

### Requirement: Service Providers
The container SHALL support service providers for organizing related bindings.

#### Scenario: Register phase
- **WHEN** provider's Register method is called
- **THEN** bindings are registered but not resolved

#### Scenario: Boot phase
- **WHEN** provider's Boot method is called
- **THEN** provider can resolve dependencies and perform initialization

#### Scenario: Provider dependencies
- **WHEN** one provider depends on another
- **THEN** providers are booted in dependency order

### Requirement: Tagged Services
The container SHALL support tagging services for selective resolution.

#### Scenario: Tag binding
- **WHEN** binding is registered with tag `nasc.Bind(...).Tag("logging")`
- **THEN** service is associated with that tag

#### Scenario: Resolve by tag
- **WHEN** calling `nasc.Tagged("logging")`
- **THEN** all services with that tag are returned

#### Scenario: Multiple tags per binding
- **WHEN** binding has multiple tags
- **THEN** it appears in results for all tags

### Requirement: Circular Dependency Detection
The container SHALL detect and report circular dependencies.

#### Scenario: Direct circular dependency
- **WHEN** type A depends on B and B depends on A
- **THEN** clear error is returned with dependency chain

#### Scenario: Indirect circular dependency
- **WHEN** A → B → C → A dependency chain exists
- **THEN** error shows full circular path

#### Scenario: Lazy resolution breaks cycles
- **WHEN** circular dependency uses factory or lazy injection
- **THEN** resolution succeeds

### Requirement: Error Handling
The container SHALL provide clear, actionable error messages.

#### Scenario: Missing dependency error
- **WHEN** dependency cannot be resolved
- **THEN** error states which type is missing and where it's needed

#### Scenario: Resolution chain in errors
- **WHEN** nested dependency resolution fails
- **THEN** error shows full resolution chain (A → B → C)

#### Scenario: Type mismatch error
- **WHEN** bound type doesn't implement interface
- **THEN** clear error explains the mismatch

### Requirement: Thread Safety
The container SHALL be thread-safe for concurrent access.

#### Scenario: Concurrent binding
- **WHEN** multiple goroutines register bindings
- **THEN** all bindings are stored correctly without data races

#### Scenario: Concurrent resolution
- **WHEN** multiple goroutines resolve dependencies
- **THEN** all resolutions succeed without races

#### Scenario: Singleton thread safety
- **WHEN** multiple goroutines resolve singleton simultaneously
- **THEN** only one instance is created safely

### Requirement: Scope Management
The container SHALL support creating and managing scopes for scoped lifetimes.

#### Scenario: Create scope
- **WHEN** `scope := nasc.CreateScope()` is called
- **THEN** a new isolated scope is created

#### Scenario: Resolve in scope
- **WHEN** resolving scoped binding within a scope
- **THEN** same instance is returned within that scope

#### Scenario: Dispose scope
- **WHEN** `scope.Dispose()` is called
- **THEN** all scoped instances are cleaned up

#### Scenario: Nested scopes
- **WHEN** scopes are nested
- **THEN** each level has its own scoped instances

### Requirement: Named Bindings
The container SHALL support named bindings for multiple implementations of same interface.

#### Scenario: Bind with name
- **WHEN** binding is registered as `nasc.Bind(...).Named("console")`
- **THEN** binding is associated with that name

#### Scenario: Resolve by name
- **WHEN** calling `nasc.MakeNamed((*Logger)(nil), "console")`
- **THEN** the named implementation is returned

#### Scenario: Default binding
- **WHEN** no name is specified in resolution
- **THEN** unnamed or default binding is used

### Requirement: Conditional Resolution
The container SHALL support conditional binding based on runtime context.

#### Scenario: Conditional binding
- **WHEN** binding has condition `nasc.Bind(...).When(ctx.IsProduction)`
- **THEN** binding is only used when condition is true

#### Scenario: Fallback binding
- **WHEN** conditional binding doesn't match
- **THEN** next matching binding is used

### Requirement: Binding Decorators
The container SHALL support decorating resolved instances.

#### Scenario: Decorator wraps instance
- **WHEN** decorator is registered for interface
- **THEN** resolved instance is wrapped before returning

#### Scenario: Multiple decorators
- **WHEN** multiple decorators exist
- **THEN** they are applied in order

### Requirement: Instance Binding
The container SHALL support binding pre-created instances.

#### Scenario: Bind existing instance
- **WHEN** `nasc.Instance((*Logger)(nil), existingLogger)` is called
- **THEN** that instance is always returned

#### Scenario: Instance as singleton
- **WHEN** instance is bound
- **THEN** it behaves as a singleton

### Requirement: Parameterized Resolution
The container SHALL support passing parameters during resolution.

#### Scenario: Resolve with parameters
- **WHEN** `nasc.MakeWith((*Service)(nil), param1, param2)` is called
- **THEN** parameters are passed to constructor or factory

#### Scenario: Override dependencies
- **WHEN** parameters match dependency types
- **THEN** provided values override container bindings

### Requirement: Disposal and Cleanup
The container SHALL support proper cleanup of resources.

#### Scenario: Disposable interface
- **WHEN** instance implements Dispose() method
- **THEN** Dispose is called when container/scope is disposed

#### Scenario: Cleanup callbacks
- **WHEN** binding registers cleanup callback
- **THEN** callback is invoked on disposal

#### Scenario: Container disposal
- **WHEN** container.Dispose() is called
- **THEN** all singletons and resources are cleaned up

### Requirement: Debug Mode
The container SHALL provide debugging capabilities for troubleshooting.

#### Scenario: Enable debug mode
- **WHEN** debug mode is enabled
- **THEN** verbose resolution logging is available

#### Scenario: Resolution tracing
- **WHEN** resolving dependencies in debug mode
- **THEN** full resolution tree is logged

#### Scenario: Binding inspection
- **WHEN** querying container in debug mode
- **THEN** all bindings can be listed and inspected

### Requirement: Performance
The container SHALL be optimized for production use.

#### Scenario: Resolution caching
- **WHEN** type is resolved multiple times
- **THEN** reflection metadata is cached

#### Scenario: Minimal allocations
- **WHEN** resolving dependencies
- **THEN** memory allocations are minimized

#### Scenario: Fast singleton access
- **WHEN** accessing singleton
- **THEN** resolution is O(1) after first access

### Requirement: Extensibility
The container SHALL provide extension points for custom behavior.

#### Scenario: Custom binding strategies
- **WHEN** custom binding strategy is registered
- **THEN** it's used for resolution

#### Scenario: Resolution middleware
- **WHEN** middleware is registered
- **THEN** it intercepts resolution pipeline

### Requirement: Documentation
The container SHALL have comprehensive documentation.

#### Scenario: GoDoc coverage
- **WHEN** viewing package documentation
- **THEN** all exported types and functions are documented

#### Scenario: Usage examples
- **WHEN** reading documentation
- **THEN** examples cover all major features

#### Scenario: Best practices guide
- **WHEN** developer needs guidance
- **THEN** documentation includes patterns and anti-patterns
