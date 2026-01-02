## ADDED Requirements

### Requirement: Independent Go Library
The Fíth renderer SHALL be implemented as a standalone Go library that can be used independently of the Toutā framework.

#### Scenario: Import in any Go project
- **WHEN** a developer imports `github.com/toutaio/toutago-fith-renderer`
- **THEN** they can use Fíth for templating without any Toutā dependencies

#### Scenario: Semantic versioning
- **WHEN** the library is released
- **THEN** it follows semantic versioning for stability guarantees

### Requirement: Slug-Based Template Resolution
The renderer SHALL resolve templates using slug-based naming with support for path separators.

#### Scenario: Simple slug resolution
- **WHEN** rendering with slug "home"
- **THEN** the renderer loads `templates/home.fith` or `templates/home.html`

#### Scenario: Nested slug with path separators
- **WHEN** rendering with slug "admin/users/list"
- **THEN** the renderer loads `templates/admin/users/list.fith`

#### Scenario: Configurable template directory
- **WHEN** the template directory is set to "views"
- **THEN** the renderer looks for templates in the "views" directory

### Requirement: Data Source Support
The renderer SHALL accept both maps and structs as data sources for template rendering.

#### Scenario: Render with map data
- **WHEN** template receives `map[string]interface{}{"name": "John", "age": 30}`
- **THEN** `{{.name}}` renders as "John" and `{{.age}}` renders as "30"

#### Scenario: Render with struct data
- **WHEN** template receives a struct with fields Name and Age
- **THEN** `{{.Name}}` and `{{.Age}}` access the struct fields

#### Scenario: Nested data access
- **WHEN** data contains nested maps or structs
- **THEN** dot notation like `{{.User.Profile.Name}}` accesses nested values

### Requirement: Control Flow - Loops
The renderer SHALL support iteration over collections using loop constructs.

#### Scenario: Loop over array
- **WHEN** template contains `{{range .Items}}{{.}}{{end}}`
- **THEN** each item in the array is rendered

#### Scenario: Loop with index
- **WHEN** template uses `{{range $index, $item := .Items}}{{$index}}: {{$item}}{{end}}`
- **THEN** both index and item are available in the loop

#### Scenario: Loop variables
- **WHEN** inside a range loop
- **THEN** special variables like `@first`, `@last`, `@index` are available

### Requirement: Control Flow - Conditionals
The renderer SHALL support conditional rendering with if/else statements.

#### Scenario: Simple if condition
- **WHEN** template contains `{{if .IsLoggedIn}}Welcome{{end}}`
- **THEN** content renders only when condition is true

#### Scenario: If-else condition
- **WHEN** template uses `{{if .IsAdmin}}Admin{{else}}User{{end}}`
- **THEN** appropriate branch renders based on condition

#### Scenario: Else-if chains
- **WHEN** template has `{{if .Role == "admin"}}Admin{{else if .Role == "moderator"}}Mod{{else}}User{{end}}`
- **THEN** the first matching condition renders

### Requirement: Template Includes
The renderer SHALL support including sub-templates with parameter passing.

#### Scenario: Include another template
- **WHEN** template contains `{{include "header"}}`
- **THEN** the header template is rendered inline

#### Scenario: Include with parameters
- **WHEN** template uses `{{include "card" title="Welcome" content=.Message}}`
- **THEN** the card template receives title and content parameters

#### Scenario: Include with data context
- **WHEN** template uses `{{include "user-card" .User}}`
- **THEN** the included template receives .User as its root context

### Requirement: Layout System
The renderer SHALL support parent/child layout inheritance for consistent page structure.

#### Scenario: Extend parent layout
- **WHEN** template declares `{{extends "layout"}}`
- **THEN** it inherits the parent layout structure

#### Scenario: Define content blocks
- **WHEN** child template defines `{{block "content"}}...{{end}}`
- **THEN** the content replaces the corresponding block in parent

#### Scenario: Multiple blocks
- **WHEN** layout has multiple blocks (header, content, footer)
- **THEN** child can override any or all blocks

#### Scenario: Default block content
- **WHEN** child doesn't override a block
- **THEN** the parent's default block content is used

### Requirement: Custom Functions
The renderer SHALL allow registration of custom functions for use in templates.

#### Scenario: Register custom function
- **WHEN** developer registers a function `fith.RegisterFunction("greet", greetFunc)`
- **THEN** templates can use `{{greet .Name}}`

#### Scenario: Function with multiple arguments
- **WHEN** function is registered with signature `func(string, string) string`
- **THEN** templates can call `{{myFunc "arg1" "arg2"}}`

### Requirement: Built-in Functions
The renderer SHALL provide commonly used built-in functions for string manipulation, formatting, and data operations.

#### Scenario: String functions
- **WHEN** template uses `{{upper .Name}}`
- **THEN** the name is rendered in uppercase

#### Scenario: Array/slice functions
- **WHEN** template uses `{{join .Tags ", "}}`
- **THEN** array elements are joined with separator

#### Scenario: Date formatting
- **WHEN** template uses `{{date .CreatedAt "2006-01-02"}}`
- **THEN** date is formatted according to Go time layout

### Requirement: Filter Pipeline
The renderer SHALL support chaining filters for data transformation.

#### Scenario: Single filter
- **WHEN** template uses `{{.Name | upper}}`
- **THEN** the filter is applied to the value

#### Scenario: Chained filters
- **WHEN** template uses `{{.Name | lower | title}}`
- **THEN** filters are applied left-to-right

#### Scenario: Filter with arguments
- **WHEN** template uses `{{.Text | truncate 100}}`
- **THEN** the filter receives the argument

### Requirement: Multi-Format Output
The renderer SHALL support rendering any text-based format, not just HTML.

#### Scenario: Render HTML
- **WHEN** rendering HTML templates
- **THEN** output is valid HTML with optional auto-escaping

#### Scenario: Render CSV
- **WHEN** rendering CSV templates
- **THEN** output is properly formatted CSV

#### Scenario: Render markdown
- **WHEN** rendering markdown templates
- **THEN** output is valid markdown syntax

#### Scenario: Render plain text
- **WHEN** rendering text templates
- **THEN** output is plain text without escaping

### Requirement: Auto-Escaping
The renderer SHALL provide configurable auto-escaping for safe HTML rendering.

#### Scenario: HTML auto-escape enabled
- **WHEN** auto-escaping is enabled and data contains `<script>`
- **THEN** output is escaped as `&lt;script&gt;`

#### Scenario: Raw output
- **WHEN** template uses `{{.HTML | raw}}`
- **THEN** content is rendered without escaping

#### Scenario: Context-aware escaping
- **WHEN** escaping is context-aware
- **THEN** JavaScript strings, HTML attributes, and HTML content use appropriate escaping

### Requirement: Template Caching
The renderer SHALL cache parsed templates for performance optimization.

#### Scenario: Cache enabled in production
- **WHEN** cache is enabled
- **THEN** templates are parsed once and reused

#### Scenario: Hot-reload in development
- **WHEN** development mode is enabled
- **THEN** templates are reloaded when files change

#### Scenario: Cache invalidation
- **WHEN** a template file is modified
- **THEN** the cache entry is invalidated in development mode

### Requirement: Error Reporting
The renderer SHALL provide clear, actionable error messages with line numbers and context.

#### Scenario: Syntax error
- **WHEN** template has syntax error on line 15
- **THEN** error message includes line number and excerpt

#### Scenario: Undefined variable
- **WHEN** template references undefined variable
- **THEN** error identifies the variable and location

#### Scenario: Template not found
- **WHEN** included template doesn't exist
- **THEN** error clearly states which template is missing

### Requirement: Embedded Filesystem Support
The renderer SHALL support loading templates from Go's embed.FS for embedded templates.

#### Scenario: Load from embed.FS
- **WHEN** templates are embedded using `//go:embed`
- **THEN** renderer loads templates from the embedded filesystem

#### Scenario: Fallback to disk
- **WHEN** template not found in embed.FS
- **THEN** renderer can optionally fall back to disk

### Requirement: Thread Safety
The renderer SHALL be thread-safe for concurrent rendering operations.

#### Scenario: Concurrent renders
- **WHEN** multiple goroutines render simultaneously
- **THEN** all renders complete successfully without data races

#### Scenario: Shared cache
- **WHEN** multiple renders access cached templates
- **THEN** cache operations are synchronized

### Requirement: Performance
The renderer SHALL be optimized for production use with minimal overhead.

#### Scenario: Compilation to bytecode
- **WHEN** template is parsed
- **THEN** it's compiled to an efficient intermediate representation

#### Scenario: Streaming output
- **WHEN** rendering large templates
- **THEN** output can be streamed to io.Writer without buffering everything

#### Scenario: Memory efficiency
- **WHEN** rendering many templates
- **THEN** memory allocations are minimized through pooling

### Requirement: Extensibility
The renderer SHALL provide extension points for custom behavior.

#### Scenario: Custom delimiter syntax
- **WHEN** developer configures delimiters as `<% %>` instead of `{{ }}`
- **THEN** templates use the custom delimiters

#### Scenario: Custom template loader
- **WHEN** developer implements custom TemplateLoader interface
- **THEN** renderer uses the custom loading logic

#### Scenario: Hooks and middleware
- **WHEN** developer registers render hooks
- **THEN** hooks are called before/after rendering

### Requirement: Documentation
The renderer SHALL have comprehensive documentation for all features and APIs.

#### Scenario: GoDoc coverage
- **WHEN** viewing package documentation
- **THEN** all exported types and functions have clear documentation

#### Scenario: Syntax reference
- **WHEN** developer reads the syntax guide
- **THEN** all template syntax features are documented with examples

#### Scenario: Migration guide
- **WHEN** migrating from html/template
- **THEN** migration guide shows equivalent syntax and patterns
