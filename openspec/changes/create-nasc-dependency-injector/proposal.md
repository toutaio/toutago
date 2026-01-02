# Change: Create Nasc Dependency Injector - Independent Project

## Why
Dependency injection is fundamental to Toutā's architecture, enabling interface-first design, testability, and loose coupling. Rather than embedding DI functionality in the core framework or relying on external libraries, creating Nasc as an independent, reusable Go library provides maximum flexibility and follows Toutā's philosophy. This allows the DI container to be used in any Go project while maintaining the framework's unique Celtic-inspired identity.

## What Changes
- Create a new independent Go project: `toutago-nasc-dependency-injector`
- Repository: https://github.com/toutaio/toutago-nasc-dependency-injector
- Location: `~/Proyects/toutago-nasc-dependency-injector`
- Implement a full-featured dependency injection container with:
  - Interface-based binding and resolution
  - Multiple binding modes (transient, singleton, scoped, factory)
  - Auto-wiring with struct tags
  - Constructor injection support
  - Lazy resolution and circular dependency detection
  - Service providers for organizing bindings
  - Tagged services and conditional resolution
  - Thread-safe concurrent access
  - Clear error messages and debugging support
- Design as a standalone library that can be imported into Toutā or any Go project

## Impact
- Affected specs: nasc-di (new capability - separate project)
- Affected code:
  - New repository: toutago-nasc-dependency-injector
  - Toutā can optionally integrate via Go modules
  - Replace current internal/di with Nasc
- Benefits:
  - Independent, reusable DI container
  - Clean separation from core framework
  - Can be used in any Go project
  - Follows Toutā's nemeton philosophy
  - Provides powerful, flexible dependency injection
  - Supports advanced DI patterns
  - Extensible with service providers
  - Celtic-themed identity consistent with Toutā
  - Better testing and maintainability
