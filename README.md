# Toutā

A message-driven Go web framework emphasizing interface-first design, dependency injection, and developer experience.

## Features

✅ **Dependency Injection Container** - Interface-based DI with auto-wiring and singleton support  
✅ **Message Bus** - Pub/sub system for message-based communication  
✅ **Router Abstraction** - HTTP router interface with Chi as default implementation  
✅ **Configuration System** - YAML frontmatter loader with environment variable support  
✅ **CLI Framework** - Cobra-based ogam (commands) for project scaffolding and development  
✅ **Template Renderer** - HTML template wrapper with custom function support  
✅ **Component Registry** - Nemeton (package) manifest parsing and component registration  

## Quick Start

### Installation

```bash
# Clone the repository
git clone https://github.com/toutaio/toutago
cd toutago

# Build the CLI
go build -o touta cmd/touta/main.go

# Make it available globally (optional)
sudo mv touta /usr/local/bin/
# Or to the user bin directory
mv touta ~/go/bin/
```

### Create Your First Project

```bash
# Create a new project
touta new my-app
cd my-app

# Run the application
go run main.go
```

The server will start on `http://localhost:8080`.

## Project Structure

```
my-app/
├── touta.yaml          # Configuration file
├── main.go             # Application entry point
├── handlers/           # Message handlers
│   └── hello.go
├── templates/          # HTML templates
├── static/             # Static assets
└── config/             # Additional config files
```

## Configuration

`touta.yaml` example:

```yaml
framework:
  mode: development
  debug: true
  hot_reload: true
  log_level: info

server:
  host: localhost
  port: 8080

router:
  base_path: /
```

## Core Concepts

### Message Bus

Messages flow through the system via a pub/sub bus:

```go
// Define a message
type UserRegistered struct {
    message.BaseMessage
    Email    string
    Username string
}

// Create a handler
type UserHandler struct{}

func (h *UserHandler) Handle(ctx context.Context, msg touta.Message) (touta.Message, error) {
    // Process the message
    return nil, nil
}

// Subscribe and publish
bus.Subscribe("user.registered", &UserHandler{})
bus.Publish(ctx, &UserRegistered{
    Email: "user@example.com",
    Username: "john",
})
```

### Dependency Injection

All components use interface-based dependency injection:

```go
// Bind an interface to an implementation
container.Bind((*StorageAdapter)(nil), &FileStorage{})

// Resolve dependencies
storage, _ := container.Make(reflect.TypeOf((*StorageAdapter)(nil)))

// Auto-wire into structs
type MyHandler struct {
    Storage StorageAdapter `inject:""`
}
handler := &MyHandler{}
container.AutoWire(handler)
```

### HTTP Router

Clean router abstraction with Chi underneath:

```go
router := router.NewChiRouter(container)

router.GET("/", func(ctx touta.Context) error {
    return ctx.HTML(200, "<h1>Hello World</h1>")
})

router.POST("/users", func(ctx touta.Context) error {
    return ctx.JSON(201, map[string]string{"status": "created"})
})

router.Listen(":8080")
```

## CLI Commands (Commands)

```bash
# Create a new project
touta new <project-name>

# Initialize in existing directory
touta init

# Start development server
touta serve [--port 8080] [--host localhost]

# Show version
touta version
```

## Testing

Run all tests:

```bash
go test ./...
```

Run specific package tests:

```bash
go test ./internal/di -v
go test ./internal/message -v
go test ./internal/router -v
```

## Architecture

Toutā follows these principles:

1. **Interface-First Design** - All dependencies are interfaces, never concrete types
2. **Message-Passing Architecture** - Components communicate via messages
3. **Dependency Injection** - No global state or direct instantiation
4. **Pluggability** - Swap implementations without code changes

### Core Interfaces

- `Container` - Dependency injection and service resolution
- `MessageBus` - Message publishing and subscription
- `Router` - HTTP routing abstraction
- `ConfigLoader` - Configuration loading and watching
- `TemplateRenderer` - Template parsing and rendering
- `ComponentRegistry` - Nemeton (package) and component management

## Development Status

**Phase 1: Foundation** ✅ **COMPLETE**

- [x] Core interfaces defined
- [x] DI container implementation  
- [x] Message bus implementation
- [x] Router abstraction (Chi)
- [x] Configuration system
- [x] CLI framework
- [x] Template renderer
- [x] Component registry
- [x] **Test coverage: 85.9%** ✅
- [x] **Hot reload integration** ✅
- [x] Example project
- [x] CI/CD pipeline

**All Phase 1 objectives completed!** 🎉

## Contributing

Contributions are welcome! Please read our contributing guidelines and code of conduct.

## License

MIT License - see LICENSE file for details.

## Version

v0.1.0 - Phase 1 Foundation Implementation
