# Fíth Renderer Project - Complete Implementation Plan

## Overview

**Fíth** (Old Irish: "weaving patterns") is an independent, powerful template engine for Go being developed as a separate project from Toutā.

## Project Details

- **Name:** Fíth Renderer
- **Repository:** https://github.com/toutaio/toutago-fith-renderer
- **Local Path:** `~/Proyects/toutago-fith-renderer`
- **Language:** Go 1.21+
- **License:** MIT (recommended)
- **Status:** Planning Complete ✅

## Why Separate Project?

1. **Reusability** - Can be used in any Go project, not just Toutā
2. **Independence** - Develops at its own pace
3. **Clean Architecture** - No coupling with Toutā internals
4. **Versioning** - Independent semantic versioning
5. **Community** - Can have its own contributor base

## Key Features

### Core Capabilities
- ✅ **Slug-based template resolution** with path separators
- ✅ **Multiple data sources** (maps, structs)
- ✅ **Control flow** (if/else, loops with special variables)
- ✅ **Template composition** (includes with parameters)
- ✅ **Layout inheritance** (parent/child with blocks)
- ✅ **Function system** (built-in + custom)
- ✅ **Filter pipeline** (chainable data transformations)
- ✅ **Multi-format** output (HTML, CSV, markdown, text, etc.)
- ✅ **Auto-escaping** (configurable, context-aware)
- ✅ **Template caching** with hot-reload in dev mode
- ✅ **Embed.FS support** for embedded templates

### Developer Experience
- Clear error messages with line numbers
- Thread-safe concurrent rendering
- Extensive documentation
- Rich examples
- Migration guide from html/template

## Template Syntax Examples

### Basic Usage
```go
fith := fith.New(fith.Config{
    TemplateDir: "templates",
})

data := map[string]interface{}{
    "Title": "Welcome",
    "User":  user,
    "Items": items,
}

output, err := fith.Render("home", data)
```

### Template Examples

**Variables:**
```html
<h1>{{.Title}}</h1>
<p>Welcome, {{.User.Name}}!</p>
```

**Conditionals:**
```html
{{if .User.IsAdmin}}
  <a href="/admin">Admin Panel</a>
{{else}}
  <p>Regular User</p>
{{end}}
```

**Loops:**
```html
{{range .Items}}
  <li>{{@index}}: {{.Name}}</li>
  {{if @first}}<strong>First!</strong>{{end}}
{{end}}
```

**Filters:**
```html
<h1>{{.Title | upper}}</h1>
<p>{{.Bio | truncate 100}}</p>
<time>{{.Created | date "Jan 2, 2006"}}</time>
```

**Includes:**
```html
{{include "header"}}
{{include "card" title="Hello" content=.Message}}
{{include "user-profile" .User}}
```

**Layouts:**
```html
<!-- child.fith -->
{{extends "layouts/main"}}

{{block "title"}}My Page{{end}}

{{block "content"}}
  <p>Page content</p>
{{end}}

<!-- layouts/main.fith -->
<html>
  <head><title>{{block "title"}}Default{{end}}</title></head>
  <body>{{block "content"}}{{end}}</body>
</html>
```

## Project Structure

```
~/Proyects/toutago-fith-renderer/
├── README.md                    # Project overview
├── IMPLEMENTATION_PLAN.md       # Detailed implementation guide
├── LICENSE                      # MIT license
├── go.mod                       # Go module
│
├── fith.go                      # Main API
├── config.go                    # Configuration types
├── errors.go                    # Error types
│
├── lexer/                       # Tokenization
├── parser/                      # AST generation
├── compiler/                    # Template compilation
├── runtime/                     # Execution engine
├── loader/                      # Template loading
├── builtins/                    # Built-in functions
│
├── examples/                    # Usage examples
├── docs/                        # Documentation
└── benchmarks/                  # Performance tests
```

## Implementation Timeline

### Phase 1: Foundation (Weeks 1-2)
- Lexer and tokenization
- Parser and AST generation
- Basic variable interpolation

### Phase 2: Control Flow (Week 3)
- If/else conditionals
- Range loops
- Loop variables

### Phase 3: Rendering (Week 4)
- Runtime execution
- Expression evaluation
- Error reporting

### Phase 4: Functions & Filters (Week 5)
- Built-in functions
- Filter pipeline
- Custom function registration

### Phase 5: Composition (Week 6)
- Include directive
- Layout inheritance
- Block system

### Phase 6: Loading (Week 7)
- Template loader
- Caching system
- Hot-reload

### Phase 7: Optimization (Week 8)
- Performance tuning
- Benchmarking
- Memory optimization

### Phase 8-10: Documentation & Polish
- Complete documentation
- Examples
- Integration guides

**Target:** v1.0.0 in ~12 weeks

## OpenSpec Documentation

All requirements, tasks, and design decisions are documented in:

```
/home/nestor/Proyects/toutago/openspec/changes/create-fith-renderer/
├── proposal.md              # Why and what
├── tasks.md                 # 85 implementation tasks
├── design.md                # Technical architecture
└── specs/
    └── fith-renderer/
        └── spec.md          # 21 detailed requirements
```

## Built-in Functions (Planned)

### String Functions
- `upper`, `lower`, `title`, `trim`
- `replace`, `split`, `substring`
- `truncate`, `wordwrap`

### Array/Slice Functions
- `join`, `first`, `last`, `slice`
- `reverse`, `sort`, `unique`
- `length`, `contains`

### Formatting
- `date` - Date formatting
- `number` - Number formatting
- `currency` - Currency formatting
- `bytes` - Byte size formatting

### Logic
- `default` - Default value if nil
- `coalesce` - First non-nil value
- `ternary` - Inline if/else

## Integration with Toutā

Once Fíth is ready, Toutā will:

1. Import as a Go module:
   ```go
   import "github.com/toutaio/toutago-fith-renderer"
   ```

2. Create Toutā-specific adapter:
   ```go
   type ToutaRenderer struct {
       fith *fith.Fíth
   }
   ```

3. Provide configuration helpers:
   ```go
   renderer := touta.NewFithRenderer(config)
   ```

4. Document migration path from current template system

## Performance Goals

- **Parsing:** <1ms per template
- **Rendering:** <100μs for simple templates
- **Memory:** <10KB per cached template
- **Comparison:** Within 2x of Go's html/template

## Testing Strategy

- **Unit tests:** >90% coverage
- **Integration tests:** Real-world scenarios
- **Benchmarks:** Performance tracking
- **Examples:** Serve as tests

## Next Steps

### Immediate (Today)
1. ✅ Initialize Git repository
2. ✅ Create project structure
3. ✅ Write implementation plan
4. ✅ Create OpenSpec documentation

### This Week
1. Set up Go module
2. Design token types
3. Start lexer implementation
4. Write first tests

### Next Month
1. Complete Phase 1 (Foundation)
2. Complete Phase 2 (Control Flow)
3. Start Phase 3 (Rendering)

## Resources

- **OpenSpec Documentation:** `/home/nestor/Proyects/toutago/openspec/changes/create-fith-renderer/`
- **Implementation Plan:** `~/Proyects/toutago-fith-renderer/IMPLEMENTATION_PLAN.md`
- **Inspiration:** Jinja2, Twig, Liquid
- **Go Templates:** stdlib text/template, html/template

## How to Start

```bash
# Navigate to project
cd ~/Proyects/toutago-fith-renderer

# Read the implementation plan
cat IMPLEMENTATION_PLAN.md

# Check OpenSpec docs
cd /home/nestor/Proyects/toutago
openspec show create-fith-renderer

# Start coding!
cd ~/Proyects/toutago-fith-renderer
# Create go.mod, start with lexer/
```

## Summary

Fíth is a carefully planned, independent template engine that will:

- ✅ Provide modern templating features
- ✅ Maintain Celtic-themed identity
- ✅ Work as standalone Go library
- ✅ Integrate seamlessly with Toutā
- ✅ Support multiple output formats
- ✅ Offer excellent developer experience

**Status:** Ready to begin implementation! 🚀

---

**Created:** 2025-12-27  
**Planning Status:** Complete ✅  
**Implementation Status:** Not started  
**Target v1.0.0:** ~12 weeks
