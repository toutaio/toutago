# Toutā Framework Architecture

## Technology Stack

### Core Technologies
- **Language:** Go 1.21+
- **Router:** Interface-based (default: Chi implementation)
- **WebSocket:** Interface-based (default: gorilla/websocket)
- **Template Engine:** Interface-based (default: custom wrapping html/template)
- **Config Parser:** Interface-based (default: github.com/adrg/frontmatter for YAML)
- **Validation:** github.com/go-playground/validator (Pydantic-like validation)
- **DI Container:** Interface-based dependency injection (pluggable, default implementation TBD)
- **Dependency Management:** Centralized at project level (Composer-like for Go modules)
- **CLI Tool:** Extensible command system (cobra-based)
- **Hot Reload (dev):** air or cosmtrek/air

### Design Philosophy
- **Interface-Driven Design:** Everything uses interfaces, not concrete types (except rare exceptions)
- **Message Objects Over Classes:** All interactions through message interfaces, avoiding concrete class coupling
- **Database Independent:** No database dependency in core; storage adapters for future modules
- **Message-Centric:** Following Smalltalk's OOP philosophy—objects communicate via messages
- **Stateless First:** Lean towards stateless architecture with optional stateful support
- **Pluggable Everything:** Routers, DI containers, config loaders, template dialects, storage adapters all replaceable
- **Centralized Dependencies:** All nemeton dependencies resolved and installed at project level (Composer-style)
- **Extensible CLI:** Ogam can be added by framework, recipes, nemetons, and projects
- **Rendering Strategy:** Next.js-inspired SSR/CSR hybrid approach

## Architecture Overview

### Core Layers

```
┌─────────────────────────────────────────────────────────┐
│                    Recipe Layer                         │
│  (Complete applications: blog, wiki, ecommerce)         │
└─────────────────────────────────────────────────────────┘
                         ▼
┌─────────────────────────────────────────────────────────┐
│                   Nemeton Layer                         │
│  (Reusable components: auth, comments, media)           │
└─────────────────────────────────────────────────────────┘
                         ▼
┌─────────────────────────────────────────────────────────┐
│                    Core Framework                       │
│  ┌─────────────┬──────────────┬──────────────────────┐ │
│  │ DI          │ Message Bus  │ Component            │ │
│  │ Container   │ (Interface)  │ Registry             │ │
│  │ (Interface) │              │ (Interface)          │ │
│  └─────────────┴──────────────┴──────────────────────┘ │
│  ┌─────────────┬──────────────┬──────────────────────┐ │
│  │ HTTP/WS     │ Frontend     │ Template             │ │
│  │ Server      │ Engine       │ Renderer             │ │
│  │ (Interface) │ (Interface)  │ (Interface)          │ │
│  └─────────────┴──────────────┴──────────────────────┘ │
│  ┌─────────────┬──────────────┬──────────────────────┐ │
│  │ Router      │ Config       │ CLI Tool             │ │
│  │ (Interface) │ Loader       │ (Extensible)         │ │
│  │             │ (Interface)  │                      │ │
│  └─────────────┴──────────────┴──────────────────────┘ │
└─────────────────────────────────────────────────────────┘
                         ▼
┌─────────────────────────────────────────────────────────┐
│                  Storage Adapters                       │
│  (Interface-based: memory, file, future DB modules)     │
└─────────────────────────────────────────────────────────┘
```

## Message-Centric Architecture

### Core Concept
The framework is built on the principle that real OOP is about message passing, not objects. Every interaction in the system flows through typed, validated messages.

### Message Flow

```
User Request (HTTP/WS)
       ▼
┌─────────────────┐
│  HTTP Handler   │
└─────────────────┘
       ▼
┌─────────────────┐      ┌──────────────────────┐
│  Message Bus    │◄────►│  Routing Rules       │
│                 │      │  (YAML/Config based) │
└─────────────────┘      └──────────────────────┘
       ▼
┌─────────────────┐
│  Components     │ ◄─── Registered handlers listen
│  (Handlers)     │      for message types/slugs
└─────────────────┘
       ▼
┌─────────────────┐
│ Response Msg    │
└─────────────────┘
       ▼
┌─────────────────┐
│ Frontend Engine │ ◄─── SSR/CSR decision
│ (Template +     │      (like Next.js)
│  Postback JS)   │
└─────────────────┘
       ▼
   HTML + JS sent to client
```

### Base Message Structure

All messages extend from a base structure that provides core metadata:

```go
// Core message interface
type Message interface {
    Slug() string
    Type() string
    Metadata() map[string]interface{}
}

// Base message struct
type BaseMessage struct {
    MessageSlug string                 `yaml:"slug" validate:"required"`
    MessageType string                 `yaml:"type" validate:"required"`
    Meta        map[string]interface{} `yaml:"metadata,omitempty"`
}

func (m BaseMessage) Slug() string { return m.MessageSlug }
func (m BaseMessage) Type() string { return m.MessageType }
func (m BaseMessage) Metadata() map[string]interface{} { return m.Meta }
```

### Example Concrete Message

```go
type UserFormSubmitted struct {
    BaseMessage
    FormData    map[string]string `yaml:"data" validate:"required"`
    ComponentID string            `yaml:"component_id" validate:"required"`
    PageID      string            `yaml:"page_id"`
}
```

### Message Routing

**Configuration-based routing** (YAML):
```yaml
---
slug: user-form-submitted
type: form.submission
---
routes:
  - handler: auth.validateUser
    condition: data.email != ""
  - handler: user.saveProfile
  - handler: notification.sendWelcome
    async: true
```

**Code-based routing** (alternative):
```go
bus.Subscribe("user-form-submitted", userHandler.HandleFormSubmit)
```

## Nemeton System

### Nemeton Structure

```
nemetons/
└── auth/
    ├── nemeton.yaml        # Nemeton metadata and dependencies
    ├── components/         # Go code for handlers
    │   ├── login.go
    │   └── register.go
    ├── templates/          # Frontend templates
    │   ├── login.html
    │   └── register.html
    ├── messages/           # Message definitions
    │   └── auth_messages.go
    ├── routes.yaml         # Message routing config
    ├── migrations/         # Future: DB migrations
    └── assets/             # CSS, JS, images
        ├── styles/
        └── scripts/
```

### Nemeton Metadata (nemeton.yaml)

```yaml
---
name: auth
version: 1.0.0
description: Authentication and authorization nemeton
---
# Nemeton dependencies (Go modules)
# These will be resolved and installed at the project level
dependencies:
  go_modules:
    - github.com/golang-jwt/jwt/v5: "^5.0.0"
    - golang.org/x/crypto/bcrypt: "latest"
  
  # Other Toutā nemetons this nemeton depends on
  nemetons:
    - name: core/validation
      version: ">=1.0.0"

components:
  - name: login
    type: handler
    messages:
      - user.login.attempt
      - user.login.success
  
  - name: register
    type: handler
    messages:
      - user.register.attempt

templates:
  - login.html
  - register.html

routes:
  - path: /login
    handler: login
  - path: /register
    handler: register

# Optional: CLI ogam provided by this nemeton
cli_commands:
  - name: auth:create-admin
    description: Create an admin user
    handler: ogam.CreateAdminCommand
```

### Project Structure

```
project/
├── touta.yaml              # Main config with frontmatter
├── go.mod                  # Centralized Go module dependencies
├── go.sum                  # Go module checksums
├── touta.lock              # Resolved nemeton dependencies (like composer.lock)
├── recipes/                # Recipe definitions
│   └── blog/
│       ├── recipe.yaml
│       └── nemetons/       # Recipe-specific nemetons
├── nemetons/               # Local development nemetons
│   ├── auth/
│   │   ├── nemeton.yaml
│   │   ├── components/
│   │   ├── templates/
│   │   ├── messages/
│   │   ├── ogam/       # CLI ogam
│   │   └── routes.yaml
│   └── comments/
├── vendor/                 # Go vendor directory (standard Go)
├── storage/                # File-based storage (initial)
├── assets/                 # Global assets
├── cmd/
│   └── touta/             # CLI entry point
└── .touta/                # Framework internals
    ├── cache/
    └── registry/          # Nemeton registry cache
```

### Dependency Management (Composer-Style)

Toutā uses a centralized dependency management system inspired by PHP's Composer:

1. **Nemeton Dependencies**: Each nemeton declares its Go module dependencies in `nemeton.yaml`
2. **Dependency Resolution**: The `touta install` command collects all dependencies from:
   - Project's `touta.yaml`
   - All local nemetons
   - All imported nemetons
   - Active recipe dependencies
3. **Centralized Installation**: All Go modules are resolved and added to the project's `go.mod`
4. **Lock File**: `touta.lock` stores exact resolved versions (reproducible builds)

**Example workflow:**
```bash
# Install all dependencies
touta install

# Add a nemeton (automatically installs its dependencies)
touta nemeton add github.com/user/touta-auth

# Update all dependencies
touta update
```

## Recipe System

### Recipe Concept
A recipe is a complete, deployable application composed of nemetons. Recipes can extend other recipes.

### Recipe Composition

```yaml
---
name: blog
version: 1.0.0
base: core/web
---
description: A simple blog application

extends:
  - recipe: core/web
    version: "^1.0.0"

nemetons:
  - auth
  - posts
  - comments
  - media

configuration:
  site:
    name: My Blog
    theme: default
  
  features:
    comments: true
    social_sharing: true

customization:
  templates:
    - override: posts/index.html
      with: custom/blog-index.html
```

### Recipe Repository
- Recipes stored in Git repositories
- Version controlled with semantic versioning
- Can be forked and customized
- Pull-based distribution (no marketplace initially)

## Frontend Engine

### Rendering Strategy (Next.js-inspired)

1. **SSR (Server-Side Rendering)**: Default for initial page loads
2. **CSR (Client-Side Rendering)**: For dynamic updates and interactions
3. **Hybrid**: Templates can specify rendering strategy per component

### Template Dialect (Custom)

```html
<!-- Example template syntax -->
<box:page title="User Profile">
  <box:component name="user-form" id="profile-form">
    <box:input 
      name="username" 
      value="{{ .User.Name }}"
      @change="validateUsername"
      @blur="saveField"
    />
    
    <box:button @click="submitForm">
      Save Profile
    </box:button>
  </box:component>
</box:page>
```

### Postback Mechanism

**Primary:** WebSocket connection for real-time updates
**Fallback:** HTTP POST requests

```javascript
// Auto-generated postback handler
boxIt.postback('profile-form', 'validateUsername', {
  field: 'username',
  value: inputValue
}).then(response => {
  // Handle validation response
  updateDOM(response.updates);
});
```

### Component Lifecycle

```
Client Event (click, change, etc)
       ▼
Postback to Server (WS or HTTP)
       ▼
Message Bus receives event message
       ▼
Component Handler processes
       ▼
Response with DOM updates/validation
       ▼
Client applies changes
```

## Configuration System

### Pluggable Config Loaders

```go
type ConfigLoader interface {
    Load(path string) (Config, error)
    Watch(path string, callback func(Config)) error
}

// Default YAML loader
type YAMLFrontmatterLoader struct {}

// Future: TOML, JSON, custom DSL loaders
```

### Main Configuration (touta.yaml)

```yaml
---
name: my-project
version: 1.0.0
---
framework:
  mode: development  # development, production
  port: 8080
  
config_loader:
  type: yaml-frontmatter  # Pluggable
  
template_engine:
  type: touta-default     # Pluggable
  hot_reload: true
  
message_bus:
  router: config-based    # config-based or code-based
  
server:
  websocket:
    enabled: true
    fallback: http
  
storage:
  adapter: filesystem     # Pluggable: memory, filesystem, future DB
  path: ./storage

recipes:
  active: blog
  
nemetons:
  local:
    - ./nemetons/auth
    - ./nemetons/posts
  external:
    - github.com/user/touta-comments@v1.0.0
```

## Storage Adapters (Database Agnostic)

### Adapter Interface

```go
type StorageAdapter interface {
    Get(key string) (interface{}, error)
    Set(key string, value interface{}) error
    Delete(key string) error
    Query(criteria Criteria) ([]interface{}, error)
}

// Initial implementations
type MemoryAdapter struct {}      // In-memory (development)
type FilesystemAdapter struct {}  // File-based persistence
```

### Future Database Modules

Database support will be added via optional nemetons:
- `touta-postgres`
- `touta-mysql`
- `touta-mongodb`
- `touta-sqlite`

Each implementing the `StorageAdapter` interface.

## Router Interface (Pluggable)

### Core Router Interface

The router is completely abstracted behind an interface, allowing the underlying implementation to be swapped without changing application code.

```go
// Router interface - all routing implementations must satisfy this
type Router interface {
    // HTTP methods
    GET(path string, handler HandlerFunc)
    POST(path string, handler HandlerFunc)
    PUT(path string, handler HandlerFunc)
    DELETE(path string, handler HandlerFunc)
    PATCH(path string, handler HandlerFunc)
    
    // Group routes with prefix
    Group(prefix string) Router
    
    // Middleware
    Use(middleware ...MiddlewareFunc)
    
    // Start server
    Listen(addr string) error
    
    // Get underlying implementation (for advanced use)
    Native() interface{}
}

// Handler function signature
type HandlerFunc func(Context) error

// Middleware function signature
type MiddlewareFunc func(HandlerFunc) HandlerFunc
```

### Default Implementation (Chi)

```go
// ChiRouter wraps Chi router
type ChiRouter struct {
    mux *chi.Mux
}

func NewChiRouter() Router {
    return &ChiRouter{
        mux: chi.NewRouter(),
    }
}

func (r *ChiRouter) GET(path string, handler HandlerFunc) {
    r.mux.Get(path, r.adaptHandler(handler))
}

// Convert Toutā handler to Chi handler
func (r *ChiRouter) adaptHandler(h HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, req *http.Request) {
        ctx := NewContext(w, req, r.di)
        if err := h(ctx); err != nil {
            // Error handling
        }
    }
}
```

### Alternative Implementations

Users can create implementations for other routers:
- Fiber
- Gin
- Echo
- Standard library ServeMux

```go
// In touta.yaml
framework:
  router:
    type: chi  # or: fiber, gin, echo, stdlib
```

## Dependency Injection Container

### DI Container Interface

```go
// Container interface - completely pluggable
type Container interface {
    // Bind a concrete implementation to an interface
    Bind(abstract interface{}, concrete interface{}) error
    
    // Bind a singleton
    Singleton(abstract interface{}, concrete interface{}) error
    
    // Bind with a factory function
    Factory(abstract interface{}, factory func(Container) (interface{}, error)) error
    
    // Resolve a dependency
    Make(abstract interface{}) (interface{}, error)
    
    // Resolve with parameters
    MakeWith(abstract interface{}, params map[string]interface{}) (interface{}, error)
    
    // Check if binding exists
    Has(abstract interface{}) bool
    
    // Get all bindings of a type
    Tagged(tag string) ([]interface{}, error)
}
```

### Default Implementation

```go
// SimpleContainer - default DI container
type SimpleContainer struct {
    bindings   map[string]binding
    singletons map[string]interface{}
    mu         sync.RWMutex
}

type binding struct {
    concrete interface{}
    factory  func(Container) (interface{}, error)
    shared   bool
}

func NewContainer() Container {
    return &SimpleContainer{
        bindings:   make(map[string]binding),
        singletons: make(map[string]interface{}),
    }
}

// Usage example
func (app *Application) Bootstrap() {
    // Bind interfaces to implementations
    app.Container.Bind((*Router)(nil), NewChiRouter)
    app.Container.Singleton((*MessageBus)(nil), NewMessageBus)
    app.Container.Bind((*StorageAdapter)(nil), NewFilesystemAdapter)
    
    // Resolve dependencies
    router, _ := app.Container.Make((*Router)(nil))
    bus, _ := app.Container.Make((*MessageBus)(nil))
}
```

### Auto-Wiring with Reflection

```go
// Component with dependencies injected via struct tags
type UserHandler struct {
    Storage StorageAdapter `inject:""`
    Bus     MessageBus     `inject:""`
    Logger  Logger         `inject:"logger"`
}

// Container auto-wires dependencies
handler := &UserHandler{}
container.AutoWire(handler)
```

### Service Providers

```go
// ServiceProvider interface for organizing bindings
type ServiceProvider interface {
    Register(container Container) error
    Boot(container Container) error
}

// Example: Auth service provider
type AuthServiceProvider struct{}

func (p *AuthServiceProvider) Register(c Container) error {
    c.Singleton((*AuthService)(nil), NewAuthService)
    c.Bind((*TokenGenerator)(nil), NewJWTGenerator)
    return nil
}
```

## CLI Tool (Extensible Command System)

### CLI Architecture

The Toutā CLI is built on Cobra and provides an extensible command system where:
- **Framework** provides core ogam
- **Recipes** can add their own ogam
- **Nemetons** can contribute ogam
- **Projects** can define custom ogam

### Command Interface

```go
// Command interface - all CLI ogam implement this
type Command interface {
    Name() string
    Description() string
    Execute(ctx *CommandContext) error
    Flags() []Flag
}

type CommandContext struct {
    Args      []string
    Flags     map[string]interface{}
    Container Container
    Config    Config
}

type Flag struct {
    Name        string
    Shorthand   string
    Description string
    Default     interface{}
    Required    bool
}
```

### Core Framework Ogam

```bash
# Project management
touta new <project-name>         # Create new project
touta init                       # Initialize in existing directory
touta serve                      # Start development server
touta build                      # Build for production

# Dependency management
touta install                    # Install all dependencies
touta update                     # Update dependencies
touta require <nemeton>          # Add a nemeton

# Nemeton management
touta nemeton:create <name>      # Create new nemeton
touta nemeton:list               # List installed nemetons
touta nemeton:info <name>        # Show nemeton information

# Recipe management
touta recipe:init <recipe>       # Initialize a recipe
touta recipe:list                # List available recipes
touta recipe:extend <recipe>     # Extend a recipe

# Code generation
touta make:component <name>      # Generate a component
touta make:message <name>        # Generate a message
touta make:template <name>       # Generate a template
```

### Nemeton-Provided Ogam

Nemetons can register ogam in `nemeton.yaml`:

```yaml
---
name: auth
---
cli_commands:
  - name: auth:create-user
    description: Create a new user
    handler: ogam.CreateUserCommand
    flags:
      - name: email
        required: true
      - name: admin
        type: bool
        default: false
```

Implementation:

```go
// ogam/create_user.go
nemeton ogam

type CreateUserCommand struct {
    Storage StorageAdapter `inject:""`
}

func (c *CreateUserCommand) Name() string {
    return "auth:create-user"
}

func (c *CreateUserCommand) Execute(ctx *CommandContext) error {
    email := ctx.Flags["email"].(string)
    isAdmin := ctx.Flags["admin"].(bool)
    
    // Create user logic
    return nil
}
```

### Command Registration

```go
// In nemeton bootstrap
func (pkg *AuthPackage) RegisterCommands(cli *CLI) {
    cli.Register(&CreateUserCommand{})
    cli.Register(&ListUsersCommand{})
    cli.Register(&DeleteUserCommand{})
}

// Ogam are automatically discovered and registered
// via DI container and nemeton manifests
```

### Recipe Ogam

Recipes can bundle ogam:

```yaml
# blog recipe.yaml
---
name: blog
---
cli_commands:
  - blog:seed              # Seed sample data
  - blog:publish <post>    # Publish a post
  - blog:import <file>     # Import from another platform
```

### Project-Specific Ogam

Projects can add custom ogam in `touta.yaml`:

```yaml
---
name: my-blog
---
cli_commands:
  - name: deploy
    description: Deploy to production
    script: scripts/deploy.sh
  
  - name: backup
    handler: ogam.BackupCommand
```

### Command Discovery Flow

```
Application Start
       ▼
Load Framework Ogam (core)
       ▼
Load Recipe Ogam (from active recipe)
       ▼
Load Nemeton Ogam (from all nemetons)
       ▼
Load Project Ogam (from touta.yaml)
       ▼
CLI Ready with All Ogam
```

### Example CLI Usage

```bash
# Framework ogam
$ touta serve --port 8080
Starting Toutā development server on :8080...

# Nemeton command
$ touta auth:create-user --email admin@example.com --admin
User created successfully!

# Recipe command
$ touta blog:seed
Seeding blog with sample posts...
✓ Created 10 posts
✓ Created 5 categories
✓ Created 20 comments

# Project command
$ touta deploy --env production
Deploying to production...
✓ Build completed
✓ Uploaded to server
✓ Deployment successful!
```

## Key Go Design Patterns

### 1. Interfaces Everywhere
All core components are interface-driven for maximum pluggability. **Critical principle: Code against interfaces, not concrete types.**

Core interfaces:
- `Router` - HTTP routing (Chi, Fiber, Gin implementations)
- `Container` - Dependency injection
- `ConfigLoader` - Configuration parsing
- `TemplateRenderer` - Template rendering
- `StorageAdapter` - Data persistence
- `MessageBus` - Message routing
- `MessageHandler` - Message processing
- `Command` - CLI ogam

**Example: Everything injected as interface**
```go
type UserService struct {
    storage  StorageAdapter  // Interface, not *FilesystemAdapter
    bus      MessageBus      // Interface, not *DefaultBus
    renderer TemplateRenderer // Interface, not *ChiRenderer
}

// Wrong: Depends on concrete type
type WrongService struct {
    storage *FilesystemAdapter  // ❌ Coupled to implementation
}
```

### 2. Message Objects Over Concrete Classes

All communication happens through message interfaces, never concrete classes directly.

```go
// Message interface - all messages implement this
type Message interface {
    Slug() string
    Type() string
    Metadata() map[string]interface{}
}

// Handlers work with message interfaces
type Handler interface {
    Handle(msg Message) (Message, error)
    Supports(msgType string) bool
}

// Usage: Pass messages, not concrete types
func ProcessRequest(msg Message, handler Handler) {
    if handler.Supports(msg.Type()) {
        response, _ := handler.Handle(msg)
        // Work with response as Message interface
    }
}
```

### 3. Dependency Injection Pattern

All dependencies are injected via the DI container, never instantiated directly.

```go
// Component with injected dependencies
type AuthHandler struct {
    storage StorageAdapter `inject:"storage"`
    bus     MessageBus     `inject:""`
    config  Config         `inject:"config"`
}

// Container resolves and injects
container.AutoWire(&AuthHandler{})

// Or explicit resolution
storage, _ := container.Make((*StorageAdapter)(nil))
```

### 4. Functional Options Pattern

```go
func NewServer(container Container, opts ...ServerOption) *Server {
    s := &Server{
        port: 8080,
        host: "localhost",
        container: container,
    }
    for _, opt := range opts {
        opt(s)
    }
    return s
}

func WithPort(port int) ServerOption {
    return func(s *Server) {
        s.port = port
    }
}

// Usage
server := NewServer(
    container,
    WithPort(3000),
    WithHost("0.0.0.0"),
)
```

### 5. Middleware Chains

```go
type Middleware func(HandlerFunc) HandlerFunc

func Chain(h HandlerFunc, middlewares ...Middleware) HandlerFunc {
    for i := len(middlewares) - 1; i >= 0; i-- {
        h = middlewares[i](h)
    }
    return h
}
```

### 4. Context Passing

```go
type Context interface {
    Request() *http.Request
    Response() http.ResponseWriter
    Param(key string) string
    Get(key string) interface{}
    Set(key string, value interface{})
    Container() Container
    Message() Message
}

// Concrete implementation (internal)
type defaultContext struct {
    req       *http.Request
    res       http.ResponseWriter
    params    map[string]string
    container Container
    message   Message
    data      map[string]interface{}
}

// Handlers work with Context interface
func MyHandler(ctx Context) error {
    // Resolve dependencies from context's container
    storage, _ := ctx.Container().Make((*StorageAdapter)(nil))
    return nil
}
```

### 5. Channels for Message Bus

```go
type MessageBus interface {
    Publish(msg Message) error
    Subscribe(msgType string, handler Handler) error
    Unsubscribe(msgType string, handler Handler) error
}

type defaultMessageBus struct {
    subscribers map[string][]Handler
    messages    chan Message
    mu          sync.RWMutex
}

func (mb *defaultMessageBus) Publish(msg Message) error {
    mb.messages <- msg
    return nil
}

// Process messages concurrently
func (mb *defaultMessageBus) process() {
    for msg := range mb.messages {
        mb.mu.RLock()
        handlers := mb.subscribers[msg.Type()]
        mb.mu.RUnlock()
        
        for _, h := range handlers {
            go h.Handle(msg)  // Concurrent processing
        }
    }
}
```

### 6. Reflection for Metadata and Auto-Wiring

```go
func extractMessageMeta(msg Message) map[string]interface{} {
    val := reflect.ValueOf(msg)
    typ := val.Type()
    
    meta := make(map[string]interface{})
    for i := 0; i < val.NumField(); i++ {
        field := typ.Field(i)
        if tag := field.Tag.Get("route"); tag != "" {
            meta[tag] = val.Field(i).Interface()
        }
    }
    return meta
}
```

### 7. Struct Tags for Configuration

```go
type UserMessage struct {
    BaseMessage
    Email    string `validate:"required,email" route:"user.email"`
    Username string `validate:"required,min=3" route:"user.name"`
}
```

## Development Phases

### Phase 1: Foundation (Months 1-2)
**Goal:** Core infrastructure for message-based routing with DI

- [ ] **Dependency Injection Container**
  - [ ] Container interface definition
  - [ ] Default container implementation
  - [ ] Auto-wiring with reflection
  - [ ] Service provider pattern
- [ ] **Router Interface & Implementation**
  - [ ] Router interface definition
  - [ ] Chi router adapter implementation
  - [ ] HTTP server abstraction
- [ ] **Core Message Bus**
  - [ ] Message interface and BaseMessage
  - [ ] MessageBus interface
  - [ ] Pub/Sub mechanism with channels
  - [ ] Config-based message routing
- [ ] **Configuration System**
  - [ ] ConfigLoader interface
  - [ ] YAML frontmatter loader implementation
  - [ ] Config watching/hot reload
- [ ] **CLI Tool Foundation**
  - [ ] Command interface
  - [ ] Cobra-based CLI framework
  - [ ] Core ogam (new, init, serve)
  - [ ] Command discovery system
- [ ] **Component Registry**
  - [ ] Component registration interface
  - [ ] Nemeton manifest parsing
- [ ] **Basic Template Renderer**
  - [ ] TemplateRenderer interface
  - [ ] html/template wrapper implementation
- [ ] **Development Tools**
  - [ ] Hot reload (air integration)
  - [ ] Development server command

**Deliverable:** Echo server that routes messages with DI container and extensible CLI

### Phase 2: Frontend Engine (Months 2-3)
**Goal:** Working SSR with basic interactivity

- [ ] Custom template dialect (v1)
  - [ ] Parser for `<box:*>` tags
  - [ ] Component system
  - [ ] Template compilation
- [ ] Postback mechanism (HTTP first)
  - [ ] Event binding in templates
  - [ ] Backend handler mapping
  - [ ] Response format for DOM updates
- [ ] SSR implementation
- [ ] Basic reactivity/DOM updates on client
- [ ] Asset pipeline (CSS, JS bundling)

**Deliverable:** Interactive form with server validation

### Phase 3: Nemeton System (Months 3-4)
**Goal:** Reusable, distributable components with centralized dependency management

- [ ] **Nemeton System**
  - [ ] Nemeton discovery and loading
  - [ ] Local nemeton scanning
  - [ ] Nemeton.yaml parsing
  - [ ] Nemeton registry (in-memory)
  - [ ] Local vs external nemeton handling
- [ ] **Dependency Management (Composer-style)**
  - [ ] Dependency resolution algorithm
  - [ ] Go module integration
  - [ ] touta.lock file generation
  - [ ] Centralized dependency installation
  - [ ] Version conflict resolution
- [ ] **Nemeton CLI Ogam**
  - [ ] `touta install` - Install all dependencies
  - [ ] `touta update` - Update dependencies
  - [ ] `touta require <nemeton>` - Add nemeton
  - [ ] `touta nemeton:create <name>` - Create new nemeton
  - [ ] `touta nemeton:list` - List nemetons
  - [ ] `touta nemeton:info <name>` - Nemeton info
- [ ] **CLI Extension System**
  - [ ] Nemeton-provided command registration
  - [ ] Command auto-discovery from nemetons
  - [ ] Command namespacing (nemeton:command)

**Deliverable:** Auth nemeton usable in any project with automatic dependency management

### Phase 4: First Recipe (Month 4-5)
**Goal:** Complete deployable blog application

- [ ] Simple blog recipe
  - [ ] Posts nemeton
  - [ ] Comments nemeton
  - [ ] Media/uploads nemeton
- [ ] Recipe composition/extension
- [ ] Recipe initialization CLI
  - [ ] `touta recipe init blog`
  - [ ] `touta recipe extend blog my-custom-blog`

**Deliverable:** Working blog with posts and comments

### Phase 5: Advanced Features (Month 5+)
**Goal:** Production-ready features

- [ ] WebSocket support for postbacks
- [ ] CSR option (hybrid rendering)
- [ ] Storage adapter interface
  - [ ] Filesystem implementation
  - [ ] In-memory implementation
- [ ] Alternative config loaders (TOML, JSON)
- [ ] Alternative template dialects
- [ ] Nemeton repository system (Git-based)
- [ ] Documentation site (using Toutā itself)

**Deliverable:** Feature-complete framework v1.0

## Target Use Cases

### Initial Target
Small to medium projects:
- Personal blogs
- Small business websites
- Internal tools
- Prototypes and MVPs

### Future Scaling
Architecture designed to support:
- High-traffic applications
- Microservices (nemetons as services)
- Multi-tenant systems
- Enterprise applications

## Inspiration and Differentiation

### Inspired By
- **Laravel:** Developer experience, convention over configuration
- **Next.js:** Hybrid rendering strategy
- **Smalltalk:** Message-passing OOP philosophy
- **Go kit:** Microservices patterns

### Key Differentiators
- **Extreme Decoupling:** Nemetons truly independent
- **Message-First:** All interactions through typed messages
- **Delivery Flexibility:** Code in-repo or external nemetons
- **Database Agnostic:** No ORM lock-in
- **Configuration Pluggability:** Even config format is replaceable
- **Recipe System:** Complete solutions, not just components

## Performance Considerations

### Go Advantages
- Compiled binaries (fast startup, small footprint)
- Goroutines for concurrent message processing
- Native HTTP/2 and WebSocket support
- Minimal runtime overhead

### Optimization Strategies
- Message routing via type maps (O(1) lookup)
- Template compilation and caching
- Lazy nemeton loading
- Asset bundling and minification
- Optional response caching layer

## Security Principles

### Built-in Security
- Input validation on all messages (validator tags)
- CSRF protection for postbacks
- XSS prevention in template rendering
- Secure WebSocket authentication
- Nemeton signature verification (future)

### Stateless Benefits
- No session hijacking risk
- Easier horizontal scaling
- Simpler security model

## Testing Strategy

### Framework Testing
- Unit tests for core components
- Integration tests for message flow
- E2E tests for template rendering
- Benchmark tests for performance

### Nemeton Testing
Each nemeton includes:
- Component unit tests
- Message validation tests
- Template rendering tests

### Recipe Testing
- Deployment tests
- Smoke tests for critical paths
- Performance baseline tests

## Documentation Plan

- **Quick Start Guide:** Get running in 5 minutes
- **Architecture Guide:** This document
- **Message System Guide:** Deep dive into messages
- **Nemeton Development:** Create your first nemeton
- **Recipe Creation:** Build complete applications
- **API Reference:** Generated from code
- **Migration Guides:** Version upgrade paths

## Community and Ecosystem

### Initial Phase
- Open source from day one
- GitHub-based development
- Public roadmap and issues
- Contributing guidelines

### Future Growth
- Nemeton registry/directory
- Recipe showcase
- Community forums
- Video tutorials
- Conference talks

---

**Version:** 0.1.0 (Draft)  
**Last Updated:** 2025-12-20  
**Status:** Planning Phase

## Interface-First Design Principles

### Core Tenet
**Everything is an interface unless there's a compelling reason otherwise.** This is not optional—it's fundamental to Toutā's architecture.

### Rules of Engagement

1. **Code Against Interfaces, Not Implementations**
   ```go
   // ✅ Correct
   type UserService struct {
       storage StorageAdapter    // Interface
       bus     MessageBus        // Interface
   }
   
   // ❌ Wrong
   type UserService struct {
       storage *FilesystemStorage  // Concrete type
       bus     *DefaultBus         // Concrete type
   }
   ```

2. **Dependency Injection Required**
   - Never use `new()` or direct instantiation for dependencies
   - Always inject via constructor or DI container
   - Dependencies must be interfaces

3. **Message Objects for Communication**
   - All inter-component communication via Message interfaces
   - No direct method calls between components (except within a nemeton)
   - Message bus coordinates all interactions

4. **Exceptions** (rare cases where concrete types are acceptable):
   - Simple value objects (no behavior)
   - DTOs for serialization
   - Configuration structs
   - Test fixtures

### Benefits

- **Zero coupling:** Swap implementations without changing code
- **Easy testing:** Mock any dependency
- **Runtime flexibility:** Choose implementations via config
- **Nemeton independence:** True modular architecture

### Example: Complete Component

```go
// messages/user_messages.go
type UserRegistered struct {
    BaseMessage
    Email    string `validate:"required,email"`
    Username string `validate:"required"`
}

// handlers/user_handler.go
type UserHandler struct {
    storage  StorageAdapter    `inject:""`
    bus      MessageBus        `inject:""`
    renderer TemplateRenderer  `inject:""`
}

func (h *UserHandler) Handle(msg Message) (Message, error) {
    // Type assertion (safe because bus routes correctly)
    userMsg := msg.(*UserRegistered)
    
    // Use injected dependencies (all interfaces)
    err := h.storage.Set("user:"+userMsg.Email, userMsg)
    if err != nil {
        return nil, err
    }
    
    // Publish response message
    response := &UserCreatedResponse{
        BaseMessage: BaseMessage{
            MessageSlug: "user.created",
            MessageType: "event",
        },
        UserID: generateID(),
    }
    
    h.bus.Publish(response)
    return response, nil
}

// bootstrap/providers.go
type UserServiceProvider struct{}

func (p *UserServiceProvider) Register(c Container) error {
    // Register handler (will be auto-wired)
    c.Bind((*Handler)(nil), &UserHandler{})
    return nil
}
```

This approach ensures that even if you swap:
- Chi router for Fiber
- Filesystem storage for PostgreSQL
- YAML config for TOML
- Default DI container for wire/dig

...your component code **never changes**.

