# Fíth Renderer - Technical Design

## Context

Toutā needs a powerful template engine that aligns with its philosophy of clean architecture, Celtic-inspired identity, and developer experience. While Go's standard `html/template` is solid, it lacks features like template inheritance, advanced composition, and extensible filters. Creating Fíth as an independent library provides maximum reusability and maintains Toutā's unique identity.

## Goals

- Create a standalone, reusable Go template engine
- Support modern template features (inheritance, includes, filters)
- Handle any text format (HTML, CSV, markdown, etc.)
- Provide excellent developer experience
- Maintain high performance
- Enable easy integration with Toutā and other Go projects

## Non-Goals

- Not a web framework (just templates)
- Not a replacement for all Go template use cases
- Not focused on backward compatibility with html/template syntax
- Not a JavaScript/CSS preprocessor

## Architecture

### High-Level Components

```
┌─────────────────────────────────────────────────────────┐
│                    Fíth Renderer                        │
├─────────────────────────────────────────────────────────┤
│                                                          │
│  ┌────────────┐  ┌──────────┐  ┌──────────────┐       │
│  │   Lexer    │→ │  Parser  │→ │  Compiler    │       │
│  └────────────┘  └──────────┘  └──────────────┘       │
│                                                          │
│  ┌────────────────────────────────────────────────┐    │
│  │              Template Loader                    │    │
│  │  • Directory Loader                             │    │
│  │  • Embed.FS Loader                              │    │
│  │  • Slug Resolution                              │    │
│  │  • Caching                                      │    │
│  └────────────────────────────────────────────────┘    │
│                                                          │
│  ┌────────────────────────────────────────────────┐    │
│  │               Renderer Engine                   │    │
│  │  • Context Management                           │    │
│  │  • Function Registry                            │    │
│  │  • Filter Pipeline                              │    │
│  │  • Layout Resolution                            │    │
│  └────────────────────────────────────────────────┘    │
│                                                          │
│  ┌────────────────────────────────────────────────┐    │
│  │           Built-in Functions                    │    │
│  │  • String: upper, lower, title, trim...        │    │
│  │  • Array: join, slice, first, last...          │    │
│  │  • Format: date, number, currency...           │    │
│  │  • Logic: default, coalesce, ternary...        │    │
│  └────────────────────────────────────────────────┘    │
│                                                          │
└─────────────────────────────────────────────────────────┘
```

### Template Syntax Design

#### Variables
```
{{.Variable}}              - Simple variable
{{.User.Name}}             - Nested access
{{.Items[0]}}              - Array/slice access
{{.Map.key}}               - Map access
```

#### Control Flow
```
{{if .Condition}}...{{end}}
{{if .X}}...{{else}}...{{end}}
{{if .X}}...{{else if .Y}}...{{else}}...{{end}}

{{range .Items}}
  {{.}} or {{@item}}
  {{@index}} {{@first}} {{@last}}
{{end}}

{{range $index, $item := .Items}}
  {{$index}}: {{$item}}
{{end}}
```

#### Functions and Filters
```
{{upper .Name}}                    - Function call
{{.Name | upper}}                  - Filter
{{.Name | upper | trim}}           - Chained filters
{{truncate .Text 100}}             - Function with args
{{.Text | truncate 100}}           - Filter with args
{{date .Created "2006-01-02"}}    - Formatting
```

#### Template Composition
```
{{include "header"}}                          - Simple include
{{include "card" title="Hello" .User}}       - Include with params
{{include "item" .}}                         - Include with context

{{extends "layouts/main"}}                   - Extend layout
{{block "content"}}...{{end}}               - Define block
{{parent}}                                  - Call parent block
```

#### Comments
```
{{# This is a comment }}
{{#
  Multi-line comment
  Second line
#}}
```

### Core Data Structures

```go
// AST Node types
type Node interface {
    Type() NodeType
    Position() Position
}

type TextNode struct {
    Content string
}

type VariableNode struct {
    Path []string  // ["User", "Name"]
}

type IfNode struct {
    Condition Expression
    Then      []Node
    Else      []Node
}

type RangeNode struct {
    Variable  string
    Index     string
    Source    Expression
    Body      []Node
}

type IncludeNode struct {
    Template string
    Params   map[string]Expression
    Context  Expression
}

type BlockNode struct {
    Name    string
    Content []Node
}

type FunctionNode struct {
    Name string
    Args []Expression
}

// Template structure
type Template struct {
    Name    string
    Extends string
    Blocks  map[string]*BlockNode
    Nodes   []Node
    Funcs   FuncMap
}

// Renderer configuration
type Config struct {
    TemplateDir   string
    LeftDelim     string
    RightDelim    string
    AutoEscape    bool
    CacheEnabled  bool
    DevMode       bool
    FuncMap       FuncMap
}

// Main renderer
type Fíth struct {
    config    *Config
    loader    TemplateLoader
    cache     TemplateCache
    funcMap   FuncMap
}
```

### Template Loading Strategy

```go
// Loader interface
type TemplateLoader interface {
    Load(slug string) (*Template, error)
    Exists(slug string) bool
}

// Directory loader
type DirectoryLoader struct {
    baseDir   string
    extension string
    cache     map[string]*Template
}

// Slug resolution:
// "home" → templates/home.fith
// "admin/users/list" → templates/admin/users/list.fith
// "layouts/main" → templates/layouts/main.fith
```

### Function System

```go
// Function signature types
type Func0 func() interface{}
type Func1 func(interface{}) interface{}
type Func2 func(interface{}, interface{}) interface{}
type FuncVariadic func(...interface{}) interface{}

// Function registry
type FuncMap map[string]interface{}

// Built-in functions
var builtins = FuncMap{
    // String
    "upper":  strings.ToUpper,
    "lower":  strings.ToLower,
    "title":  strings.Title,
    "trim":   strings.TrimSpace,
    
    // Array
    "join":   strings.Join,
    "first":  arrayFirst,
    "last":   arrayLast,
    "slice":  arraySlice,
    
    // Format
    "date":   formatDate,
    "number": formatNumber,
    
    // Logic
    "default": defaultValue,
    "ternary": ternary,
}

// Custom function registration
func (f *Fíth) RegisterFunction(name string, fn interface{}) error {
    // Validate function signature
    // Add to funcMap
}
```

### Execution Model

```
1. Parse Phase (once per template):
   Template Text → Lexer → Tokens → Parser → AST

2. Compile Phase (once per template):
   AST → Optimizer → Compiled Template

3. Execute Phase (per render):
   Compiled Template + Data → Executor → Output
```

### Layout Resolution

```
Child template:
  {{extends "layouts/main"}}
  {{block "title"}}My Page{{end}}
  {{block "content"}}...{{end}}

Parent template (layouts/main.fith):
  <html>
    <head><title>{{block "title"}}Default{{end}}</title></head>
    <body>{{block "content"}}{{end}}</body>
  </html>

Resolution:
  1. Parse child template
  2. Identify extends directive
  3. Load parent template
  4. Merge blocks (child overrides parent)
  5. Render merged template
```

## Implementation Phases

### Phase 1: Core Parser (Weeks 1-2)
- Lexer implementation
- Parser for basic syntax
- AST generation
- Variable interpolation

### Phase 2: Control Flow (Week 3)
- If/else conditionals
- Range loops
- Loop variables

### Phase 3: Functions & Filters (Week 4)
- Function system
- Built-in functions
- Filter pipeline
- Custom function registration

### Phase 4: Template Composition (Week 5)
- Include directive
- Parameter passing
- Layout inheritance
- Block system

### Phase 5: Performance & Polish (Week 6)
- Template caching
- Optimization
- Benchmarking
- Documentation

## Decisions

### Decision: Custom Syntax vs html/template Compatibility
**Choice:** Custom syntax inspired by Jinja2/Twig
**Rationale:**
- More intuitive for developers
- Better support for modern features
- Cleaner syntax for filters and inheritance
- No backward compatibility constraints

**Alternatives Considered:**
- Extend html/template: Too limiting, syntax awkward
- Use existing library: Doesn't fit Toutā philosophy

### Decision: Compilation Strategy
**Choice:** AST-based with optional compilation to bytecode
**Rationale:**
- Fast parsing with caching
- Room for future optimization
- Clear separation of concerns

**Alternatives Considered:**
- Direct interpretation: Too slow
- Generate Go code: Complex, slow compile times

### Decision: Auto-Escaping Default
**Choice:** Auto-escape HTML by default, configurable per template type
**Rationale:**
- Security by default
- Explicit opt-out for raw content
- Context-aware escaping for HTML

### Decision: Template Loading
**Choice:** Slug-based with multiple loaders
**Rationale:**
- Clean API
- Flexible for different use cases
- Supports embedded templates

## Risks and Mitigations

### Risk: Performance vs stdlib html/template
**Mitigation:**
- Benchmark early and often
- Implement caching aggressively
- Consider compilation optimizations
- Accept slight performance trade-off for features

### Risk: Breaking changes during development
**Mitigation:**
- Use semantic versioning strictly
- Comprehensive test suite
- Keep v0.x for experimentation
- Document breaking changes clearly

### Risk: Complexity in layout resolution
**Mitigation:**
- Start with simple parent/child model
- Add features incrementally
- Test edge cases thoroughly
- Provide clear error messages

## Migration Plan

### From html/template
1. Identify template patterns
2. Map to Fíth equivalent syntax
3. Test rendering output matches
4. Gradually migrate templates

### Integration with Toutā
1. Create adapter package
2. Update template renderer interface
3. Provide migration examples
4. Update documentation

## Success Metrics

- ✅ 90%+ test coverage
- ✅ Performance within 2x of html/template
- ✅ Comprehensive documentation
- ✅ Working examples for all features
- ✅ Zero critical bugs in 1.0 release
- ✅ Positive developer feedback

## Open Questions

1. Should we support async/parallel rendering?
2. Should we add template precompilation CLI tool?
3. How to handle very large templates (streaming)?
4. Should we support template macros/mixins?
5. Integration with hot-reload during development?

## Timeline

- Week 1-2: Core parser and variable interpolation
- Week 3: Control flow structures
- Week 4: Functions and filters
- Week 5: Template composition
- Week 6: Performance, testing, documentation
- Week 7-8: Integration with Toutā, examples, polish

Target: v0.1.0 release in ~8 weeks
