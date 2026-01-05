# Toutā

A message-driven Go web framework emphasizing interface-first design, dependency injection, and developer experience.

> **Version 0.3.0**: Now featuring the powerful Scéla message bus with advanced event-driven capabilities!

## Features

✅ **Dependency Injection Container** - Powered by [toutago-nasc-dependency-injector](https://github.com/toutaio/toutago-nasc-dependency-injector)  
✅ **HTTP Router** - Powered by [toutago-cosan-router](https://github.com/toutaio/toutago-cosan-router) with path parameters and middleware  
✅ **Template Engine** - Powered by [toutago-fith-renderer](https://github.com/toutaio/toutago-fith-renderer) with Jinja2-style syntax  
✅ **Data Mapping** - Powered by [toutago-datamapper](https://github.com/toutaio/toutago-datamapper) with pluggable database adapters  
✅ **Database Migrations** - Powered by [toutago-sil-migrator](https://github.com/toutaio/toutago-sil-migrator)  
✅ **Message Bus** - Powered by [toutago-scela-bus](https://github.com/toutaio/toutago-scela-bus) for in-process pub/sub messaging  
✅ **Configuration System** - YAML frontmatter loader with environment variable support  
✅ **CLI Framework** - Cobra-based ogam (commands) for project scaffolding and development  
✅ **Component Registry** - Nemeton (package) manifest parsing and component registration  
🚧 **Ritual Grove** - Application lifecycle system (integrated plugin) - Powered by [toutago-ritual-grove](https://github.com/toutaio/toutago-ritual-grove)  

## Component Architecture

Toutā v2.0 has been refactored to use specialized, production-ready component libraries:

```
┌─────────────────────────────────────────────────┐
│           Toutā Framework (v0.2.0)                │
│                                                 │
│  ┌──────────────────────────────────────────┐  │
│  │    Integration Layer (Adapters)          │  │
│  └───┬──────┬──────┬─────────┬──────┬───────┘  │
│      │      │      │         │      │           │
└──────┼──────┼──────┼─────────┼──────┼───────────┘
       │      │      │         │      │
   ┌───┴──┐ ┌─┴───┐ ┌┴─────┐  ┌┴──┐  ┌┴─────┐
   │ nasc │ │cosan│ │ fith │  │data│  │scela │
   │  DI  │ │HTTP │ │Tmpl  │  │+sil│  │ Bus  │
   └──────┘ └─────┘ └──────┘  └────┘  └──────┘
```

### Why This Architecture?

- **Best-in-class components**: Each library is mature, well-tested (80%+ coverage), and production-ready
- **Flexibility**: Swap implementations easily via interfaces
- **Reduced complexity**: ~500+ lines removed from the main framework
- **Ecosystem**: Components can be used standalone in any Go project
- **SOLID principles**: Clear separation of concerns and responsibilities

### Component Libraries

| Component | Purpose | Repository |
|-----------|---------|------------|
| **nasc** | Dependency injection container | [toutago-nasc-dependency-injector](https://github.com/toutaio/toutago-nasc-dependency-injector) |
| **cosan** | HTTP router with middleware | [toutago-cosan-router](https://github.com/toutaio/toutago-cosan-router) |
| **fith** | Template engine (Jinja2-style) | [toutago-fith-renderer](https://github.com/toutaio/toutago-fith-renderer) |
| **datamapper** | Database abstraction layer | [toutago-datamapper](https://github.com/toutaio/toutago-datamapper) |
| **sil** | Database migration tool | [toutago-sil-migrator](https://github.com/toutaio/toutago-sil-migrator) |
| **scela** | In-process message bus | [toutago-scela-bus](https://github.com/toutaio/toutago-scela-bus) |
| **ritual-grove** 🚧 | Application lifecycle/templates | [toutago-ritual-grove](https://github.com/toutaio/toutago-ritual-grove) |

## Quick Start

### Installation

The Go way - simple and direct:

```bash
# Install directly from source (recommended)
go install github.com/toutaio/toutago/cmd/touta@latest

# Or clone and install locally for development
git clone https://github.com/toutaio/toutago
cd toutago
go install ./cmd/touta
```

**Note:** The `touta` binary includes embedded default rituals (project templates), so it works from any installation location without requiring additional setup.

Make sure `$GOPATH/bin` (usually `~/go/bin`) is in your `PATH`:
```bash
export PATH="$HOME/go/bin:$PATH"
```

### Verify Installation

```bash
# Check version
touta version

# List available rituals
touta ritual list
```

### Create Your First Project

```bash
# Create a new project from a ritual
touta ritual init minimal
cd my-project

# Initialize as a Toutā project
touta init

# Run the application
go run main.go
```

The server will start on `http://localhost:8080`.

### Using Docker (Recommended)

Docker provides a consistent development environment across all platforms.

#### Framework Development with Docker

```bash
# Clone the repository
git clone https://github.com/toutaio/toutago
cd toutago

# Start the development environment
docker-compose up

# The server will start with hot-reload enabled on http://localhost:8080
```

#### Creating Projects with Docker

```bash
# Create a new project
touta new my-app
cd my-app

# Start with Docker
docker-compose up

# The server starts on http://localhost:8080 with hot-reload
```

#### Docker Commands Reference

```bash
# Start services
docker-compose up

# Start in detached mode
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down

# Rebuild after dependency changes
docker-compose up --build
```

## Project Structure

```
my-app/
├── touta.yaml          # Configuration file
├── main.go             # Application entry point
├── Dockerfile          # Docker image configuration
├── docker-compose.yml  # Docker development setup
├── .dockerignore       # Docker build exclusions
├── .air.toml           # Hot-reload configuration
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

### Message Bus (Scéla)

Event-driven communication using the Scéla message bus integrated with Toutā interfaces:

```go
import (
    "context"
    "github.com/toutaio/toutago/pkg/touta"
    "github.com/toutaio/toutago/pkg/touta/integration"
)

// Create message bus
bus := integration.NewScelaBus()
defer bus.Close()

// Subscribe to events with pattern matching
bus.Subscribe("user.*", touta.HandlerFunc(
    func(ctx context.Context, msg touta.Message) error {
        // Handle all user events
        log.Printf("User event: %s - %v", msg.Topic(), msg.Payload())
        return nil
    },
))

// Subscribe to specific events
bus.Subscribe("user.created", touta.HandlerFunc(
    func(ctx context.Context, msg touta.Message) error {
        data := msg.Payload().(map[string]interface{})
        // Send welcome email, create profile, etc.
        return nil
    },
))

// Publish events (async)
ctx := context.Background()
bus.Publish(ctx, "user.created", map[string]interface{}{
    "id":    "123",
    "email": "user@example.com",
    "name":  "John Doe",
})

// Publish events synchronously (wait for handlers)
bus.PublishSync(ctx, "order.completed", orderData)
```

**Features:**
- **Async & sync publishing** - Fire-and-forget or wait for completion
- **Pattern matching** - Subscribe to `user.*`, `*.created`, `app.**`
- **Middleware support** - Logging, validation, retry logic
- **Priority levels** - High-priority messages processed first
- **Dead letter queue** - Handle failed messages gracefully
- **Persistence** - Optional Redis or filesystem backing
- **Metrics hooks** - Monitor message flow and performance

**Advanced Usage:**
```go
// Add global middleware
bus := integration.NewScelaBusWithMiddleware(
    LoggingMiddleware,
    MetricsMiddleware,
)

// Publish with priority
bus.PublishWithPriority(ctx, "alert.critical", alertData, scela.PriorityHigh)

// Handler-specific middleware
handler := ValidationMiddleware(AuthMiddleware(actualHandler))
bus.Subscribe("payment.*", handler)
```

**See:** [Scéla Documentation](https://github.com/toutaio/toutago-scela-bus) | [Message Bus Guide](docs/message-bus.md) | [Examples](examples/with-scela/)

### Dependency Injection (Nasc)

Using the nasc-powered container with integrated components:

```go
import "github.com/toutaio/toutago/pkg/touta/integration"

// Create container
container := integration.NewContainer()

// Bind an interface to an implementation
container.Bind((*Logger)(nil), &ConsoleLogger{})

// Resolve dependencies
logger, _ := container.Make((*Logger)(nil))

// Auto-wire into structs
type MyHandler struct {
    Logger Logger `inject:""`
}
handler := &MyHandler{}
container.AutoWire(handler)
```

### HTTP Router

Using the cosan-powered router with integrated DI:

```go
// Create router with container
router := integration.NewRouter(container)

// Register routes
router.GET("/", func(ctx touta.Context) error {
    return ctx.HTML(200, "<h1>Hello World</h1>")
})

router.POST("/users", func(ctx touta.Context) error {
    return ctx.JSON(201, map[string]string{"status": "created"})
})

// Add middleware
router.Use(func(next touta.HandlerFunc) touta.HandlerFunc {
    return func(ctx touta.Context) error {
        // Middleware logic
        return next(ctx)
    }
})

router.Listen(":8080")
```

### Template Rendering

Using the fith template engine:

```go
import "github.com/toutaio/toutago-fith-renderer"

// Create renderer
renderer, _ := integration.NewTemplateRenderer(&fith.Config{
    TemplateDir: "templates",
    Extensions:  []string{".html"},
})

// Use in routes
router.GET("/", func(ctx touta.Context) error {
    data := map[string]interface{}{
        "Title": "Welcome",
        "User":  user,
    }
    output, _ := renderer.Render("home", data)
    return ctx.HTML(200, string(output))
})
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

## Troubleshooting

### Docker Issues

#### Port Already in Use
```bash
# Find and kill the process using port 8080
lsof -ti:8080 | xargs kill -9

# Or change the port in docker-compose.yml
ports:
  - "3000:8080"  # Use port 3000 instead
```

#### Permission Errors on Linux
```bash
# Add your user to the docker group
sudo usermod -aG docker $USER

# Log out and back in for changes to take effect
```

#### Changes Not Hot-Reloading
```bash
# Ensure volumes are correctly mounted
docker-compose down
docker-compose up --build

# Check that .air.toml is present in the project
```

#### Container Build Fails
```bash
# Clear Docker cache and rebuild
docker-compose down
docker system prune -f
docker-compose up --build
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

Run benchmarks:

```bash
go test ./pkg/touta/integration -bench=. -benchmem
```

See [BENCHMARKS.md](docs/BENCHMARKS.md) for detailed performance metrics.

## Documentation

- **[Router Patterns Guide](docs/ROUTER_PATTERNS.md)** - Advanced routing patterns and best practices
- **[Performance Benchmarks](docs/BENCHMARKS.md)** - Router performance metrics and optimization tips
- **[Architecture](docs/architecture.md)** - Framework architecture and design decisions
- **[Message Bus Guide](docs/message-bus.md)** - Scéla event bus usage and patterns

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
- [x] Router abstraction (Cosan)
- [x] Configuration system
- [x] CLI framework
- [x] Template renderer
- [x] Component registry
- [x] **Test coverage: 85.9%** ✅
- [x] **Hot reload integration** ✅
- [x] Example project
- [x] CI/CD pipeline

**All Phase 1 objectives completed!** 🎉

## Etymology

**Toutā** (pronounced "TOO-tah") comes from Proto-Celtic **\*toutā**, meaning "tribe" or "people." This ancient word is the root of many modern Celtic language words:

- Irish: *tuath* (people, territory, petty kingdom)
- Welsh: *tud* (country, land)
- Gaulish: *teuta* (tribe)

The name reflects the framework's philosophy:

- **Community-driven**: Built by and for developers, emphasizing collaboration
- **Tribal knowledge**: Encoding best practices and patterns into reusable components
- **Territory**: A cohesive ecosystem where all parts work together harmoniously
- **Celtic roots**: Honoring the naming tradition of using Irish/Celtic terms (like our components: Nasc, Cosan, Fíth, Scéla)

Just as a *toutā* was a self-sufficient community with defined roles and shared resources, Toutā framework provides a complete, integrated ecosystem for web development in Go.

## Contributing

Contributions are welcome! Please read our contributing guidelines and code of conduct.

## License

MIT License - see LICENSE file for details.

## Version

v0.1.0 - Phase 1 Foundation Implementation
