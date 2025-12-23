# SOLID Principles Specification

> **Note**: *These principles are fundamental to creating maintainable, flexible, and robust object-oriented software. They are not rules, but guiding lights that help you make better design decisions.* — Robert C. Martin

## Purpose

This document establishes the specifications and guidelines for ensuring any software project adheres to the SOLID principles of object-oriented design. Use this as a reference when designing, reviewing, or refactoring code.

---

## The Five SOLID Principles

### 1. Single Responsibility Principle (SRP)

**Definition**: *A class should have one, and only one, reason to change.*

#### Specifications

- **Each class must have a single, well-defined responsibility**
  - The responsibility should be entirely encapsulated by the class
  - If you need to use "and" to describe what a class does, it likely violates SRP
  
- **Change triggers must be isolated**
  - Identify the actors or stakeholders who would request changes
  - Each class should serve only one actor/stakeholder
  - Changes from different sources should not affect the same class

- **Cohesion indicators**
  - All methods should utilize most or all of the class's instance variables
  - Methods should work together toward a unified purpose
  - Low cohesion indicates SRP violation

#### Anti-patterns to Avoid

- **God Classes**: Classes that do everything
- **Anemic Domain Models**: Classes with no behavior, only data
- **Mixed concerns**: Business logic mixed with presentation, persistence, or infrastructure

#### Implementation Checklist

- [ ] Each class name clearly reflects its single responsibility
- [ ] Class has a focused, coherent set of methods
- [ ] No mixing of business logic with infrastructure concerns
- [ ] Changes to different features require modifying different classes
- [ ] Class documentation can describe purpose in one sentence

---

### 2. Open/Closed Principle (OCP)

**Definition**: *Software entities (classes, modules, functions) should be open for extension, but closed for modification.*

#### Specifications

- **Extension without modification**
  - New functionality should be added by writing new code, not changing existing code
  - Existing, working code should remain untouched when requirements change
  - Use abstraction to allow multiple implementations

- **Design for variability**
  - Identify points of variation in your system
  - Abstract these variations behind interfaces or abstract classes
  - Concrete implementations should be pluggable

- **Protection mechanisms**
  - Use polymorphism, not conditional logic (if/switch statements)
  - Employ design patterns: Strategy, Template Method, Decorator
  - Dependency injection for swappable components

#### Anti-patterns to Avoid

- **Shotgun surgery**: One change requires modifications across multiple classes
- **Conditional bloat**: Excessive if/else or switch statements for type checking
- **Direct instantiation**: Creating concrete dependencies with `new` keyword

#### Implementation Checklist

- [ ] Core business logic depends on abstractions, not concretions
- [ ] New features added through new classes, not modifying existing ones
- [ ] Plugin architecture allows easy extension
- [ ] Minimal use of type-checking conditionals
- [ ] Changes are localized to specific modules

---

### 3. Liskov Substitution Principle (LSP)

**Definition**: *Derived classes must be substitutable for their base classes without altering the correctness of the program.*

#### Specifications

- **Behavioral compatibility**
  - Subtypes must honor the contract of their parent types
  - Preconditions cannot be strengthened in derived classes
  - Postconditions cannot be weakened in derived classes
  - Invariants must be preserved

- **Reasonable expectations**
  - Clients using a base class should work correctly with any derived class
  - No surprising behavior changes in subclasses
  - Derived classes should not throw unexpected exceptions

- **Design by contract**
  - Base class establishes contracts (preconditions, postconditions, invariants)
  - Derived classes must respect and honor these contracts
  - Override methods must accept what base accepts and return what base returns

#### Anti-patterns to Avoid

- **Refused bequest**: Subclass doesn't support parent's interface fully
- **Conditional type checking**: Using `instanceof` or type checks
- **Empty overrides**: Throwing exceptions or doing nothing in overridden methods
- **Classic example**: Square extending Rectangle (mathematical "is-a" doesn't mean OO "is-a")

#### Implementation Checklist

- [ ] Derived classes don't remove or disable base class functionality
- [ ] No type-checking or casting when using polymorphic references
- [ ] Substituting derived class doesn't break client code
- [ ] Preconditions and postconditions are honored
- [ ] Unit tests for base class pass with all derived classes

---

### 4. Interface Segregation Principle (ISP)

**Definition**: *No client should be forced to depend on methods it does not use.*

#### Specifications

- **Client-specific interfaces**
  - Design interfaces based on how clients use them
  - Keep interfaces small, focused, and cohesive
  - Multiple specific interfaces are better than one general-purpose interface

- **Role-based design**
  - Create interfaces that represent roles or capabilities
  - Clients depend only on the roles they need
  - Classes can implement multiple focused interfaces

- **Decoupling through segregation**
  - Changes to one interface don't affect clients of other interfaces
  - Reduces coupling between components
  - Prevents ripple effects of changes

#### Anti-patterns to Avoid

- **Fat interfaces**: Interfaces with too many methods
- **Interface pollution**: Adding methods to interface that only some implementers need
- **Header file hell**: Changes to unused methods forcing recompilation
- **Dummy implementations**: Implementing interface methods with empty bodies or exceptions

#### Implementation Checklist

- [ ] Interfaces are small and focused (typically 1-5 methods)
- [ ] No client implements interface methods it doesn't use
- [ ] Interfaces named after client capabilities (IReadable, IComparable)
- [ ] Adapter pattern used to wrap fat interfaces when unavoidable
- [ ] Changes to one interface don't cascade to unrelated clients

---

### 5. Dependency Inversion Principle (DIP)

**Definition**: *High-level modules should not depend on low-level modules. Both should depend on abstractions. Abstractions should not depend on details. Details should depend on abstractions.*

#### Specifications

- **Depend on abstractions**
  - Code to interfaces or abstract classes, not concrete implementations
  - High-level policy should not be contaminated by low-level details
  - Source code dependencies point toward abstractions

- **Inversion of control**
  - Dependencies are injected, not created internally
  - Use dependency injection (constructor, setter, or interface injection)
  - Frameworks control flow, not your code

- **Plugin architecture**
  - Low-level details plug into high-level abstractions
  - Business rules isolated from framework and database details
  - Main component composes the system by injecting dependencies

#### Anti-patterns to Avoid

- **New keyword abuse**: Directly instantiating dependencies
- **Concrete coupling**: Classes depending on concrete implementations
- **Control freak**: Classes creating and managing their own dependencies
- **Framework marriage**: Business logic tightly coupled to frameworks

#### Implementation Checklist

- [ ] High-level modules define abstractions for their needs
- [ ] Low-level modules implement those abstractions
- [ ] Dependencies injected through constructors or setters
- [ ] No direct instantiation of volatile dependencies
- [ ] Core business logic has no framework dependencies
- [ ] Dependency injection container configured in composition root

---

## Project Implementation Guidelines

### Architecture Patterns That Support SOLID

1. **Clean Architecture / Hexagonal Architecture**
   - Business logic at the center
   - Dependencies point inward
   - Infrastructure at the edges

2. **Dependency Injection**
   - Use DI containers (Spring, .NET Core DI, etc.)
   - Configure dependencies at composition root
   - Constructor injection as default

3. **Design Patterns**
   - **Strategy Pattern**: OCP, DIP
   - **Factory Pattern**: OCP, DIP
   - **Adapter Pattern**: ISP, LSP
   - **Template Method**: OCP, LSP
   - **Decorator Pattern**: OCP, SRP

### Code Review Checklist

Use this checklist during code reviews to ensure SOLID compliance:

#### Single Responsibility
- [ ] Can the class's purpose be described in one sentence without "and"?
- [ ] Would changes from different business requirements affect this class?
- [ ] Are there methods that don't use most instance variables?

#### Open/Closed
- [ ] Can new features be added without modifying existing code?
- [ ] Are there type checks or switch statements that might grow?
- [ ] Are extension points identified and abstracted?

#### Liskov Substitution
- [ ] Can derived classes be used interchangeably with base classes?
- [ ] Are there methods that throw NotImplementedException?
- [ ] Do tests for the base class pass with all derived classes?

#### Interface Segregation
- [ ] Are any interface implementations empty or throw exceptions?
- [ ] Does the interface have more than 5-7 methods?
- [ ] Do clients depend on interface methods they don't use?

#### Dependency Inversion
- [ ] Are dependencies injected rather than instantiated?
- [ ] Does high-level code depend on low-level details?
- [ ] Are volatile dependencies hidden behind interfaces?

### Testing Strategy

SOLID principles make code more testable:

- **SRP**: Each class has focused, easy-to-test responsibilities
- **OCP**: Mock implementations can be injected for testing
- **LSP**: Tests for base classes verify all implementations
- **ISP**: Test doubles implement only needed methods
- **DIP**: Dependencies are mockable through injection

#### Testing Requirements

- [ ] Unit tests exist for each class with business logic
- [ ] Dependencies are mocked/stubbed through interfaces
- [ ] Tests verify behavior, not implementation
- [ ] No tests require database or external services
- [ ] Test coverage > 80% for business logic

---

## Metrics and Indicators

### Code Quality Metrics

Track these to measure SOLID compliance:

1. **Cyclomatic Complexity**
   - Target: < 10 per method
   - High complexity indicates SRP, OCP violations

2. **Coupling Metrics**
   - Afferent Coupling (Ca): Number of classes depending on this class
   - Efferent Coupling (Ce): Number of classes this class depends on
   - Target: Low Ce, appropriate Ca

3. **Abstractness**
   - Ratio of abstract types to total types
   - Target: 20-30% for business logic layers

4. **Instability (I = Ce / (Ce + Ca))**
   - Target: High for volatile components, low for stable components

5. **Distance from Main Sequence (D = |A + I - 1|)**
   - Target: < 0.25

### Refactoring Triggers

Refactor when you observe:

- Methods longer than 20 lines (SRP violation)
- Classes with more than 7 methods (SRP violation)
- Switch statements on type (OCP violation)
- Instance checks or casts (LSP violation)
- Interfaces with unused methods (ISP violation)
- Direct instantiation of dependencies (DIP violation)

---

## Migration Strategy

### For Existing Projects

1. **Assess Current State**
   - Run static analysis tools
   - Identify most violated principles
   - Map dependencies between components

2. **Establish Testing**
   - Add characterization tests for existing code
   - Create test harness before refactoring
   - Ensure tests pass before and after changes

3. **Incremental Refactoring**
   - Start with new features (apply SOLID from the start)
   - Refactor when you touch existing code
   - Don't rewrite everything at once

4. **Create Abstractions**
   - Identify seams in existing code
   - Extract interfaces for dependencies
   - Introduce dependency injection gradually

5. **Measure Progress**
   - Track metrics over time
   - Review code regularly
   - Celebrate improvements

### For New Projects

1. **Architecture First**
   - Choose architecture pattern (Clean, Hexagonal, etc.)
   - Define layer boundaries
   - Establish dependency rules

2. **Start With Abstractions**
   - Define interfaces for external dependencies
   - Create domain models independent of frameworks
   - Plan extension points

3. **Configure DI Early**
   - Set up dependency injection container
   - Define composition root
   - Document dependency graph

4. **Continuous Enforcement**
   - Set up linting rules
   - Configure static analysis
   - Make SOLID part of code review process

---

## Common Pitfalls and Solutions

### Over-engineering

**Problem**: Applying SOLID too rigidly, creating unnecessary abstractions

**Solution**: 
- Apply principles when complexity warrants it
- Start simple, refactor when you see duplication or rigidity
- Three strikes rule: Wait for three similar cases before abstracting

### Premature Abstraction

**Problem**: Creating interfaces "just in case" without clear need

**Solution**:
- YAGNI (You Aren't Gonna Need It)
- Create abstractions when you have actual variation
- Refactor to abstractions when second implementation appears

### Analysis Paralysis

**Problem**: Spending too much time designing, not enough coding

**Solution**:
- Make it work, make it right, make it fast (in that order)
- Refactor continuously as you learn more
- Embrace iterative design

---

## Resources and References

### Essential Reading

- **"Clean Code"** by Robert C. Martin
- **"Agile Software Development, Principles, Patterns, and Practices"** by Robert C. Martin
- **"Design Patterns"** by Gang of Four
- **"Refactoring"** by Martin Fowler

### Tools

- **Static Analysis**: SonarQube, NDepend, ReSharper
- **Dependency Analysis**: Structure101, Lattix
- **Metrics**: CodeClimate, CodeCov
- **DI Containers**: Spring (Java), Autofac (.NET), Inversify (TypeScript)

---

## Conclusion

The SOLID principles are not laws to be followed dogmatically, but guidelines that help you reason about good design. They work together to create systems that are:

- **Maintainable**: Easy to change and extend
- **Testable**: Dependencies can be mocked and isolated
- **Flexible**: Adapt to changing requirements
- **Understandable**: Clear responsibilities and boundaries

Remember: *The goal is not perfect compliance with principles, but code that is easy to understand, easy to change, and hard to break.*

---

**Document Version**: 1.0  
**Last Updated**: December 2025  
**Maintained By**: Software Architecture Team  
**Review Cycle**: Quarterly

*"The only way to go fast is to go well."* — Robert C. Martin
