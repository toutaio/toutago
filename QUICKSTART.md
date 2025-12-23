# Toutā Framework - Quick Reference

## Installation & Setup

```bash
# Clone and build
git clone https://github.com/toutaio/toutago
cd toutago
go build -o touta cmd/touta/main.go

# Create new project
./touta new my-app
cd my-app
go mod init my-app
go mod tidy
go run main.go
```

## Core Interfaces

### Container (Dependency Injection)
```go
container := di.NewContainer()

// Bind interface to implementation
container.Bind((*StorageAdapter)(nil), &FileStorage{})

// Singleton
container.Singleton((*Logger)(nil), &DefaultLogger{})

// Factory
container.Factory((*DB)(nil), func(c Container) (interface{}, error) {
    return &PostgresDB{}, nil
})

// Resolve
instance, _ := container.Make(reflect.TypeOf((*StorageAdapter)(nil)))

// Auto-wire struct
type Handler struct {
    Storage StorageAdapter `inject:""`
}
handler := &Handler{}
container.AutoWire(handler)
```

### MessageBus (Pub/Sub)
```go
bus := message.NewBus()
bus.Start(context.Background())

// Define a message
type UserCreated struct {
    message.BaseMessage
    UserID string
    Email  string
}

// Subscribe
handler := &MyHandler{}
bus.Subscribe("user.created", handler)

// Publish async
bus.Publish(ctx, &UserCreated{
    BaseMessage: message.BaseMessage{
        MessageSlug: "user.created",
        MessageType: "event",
    },
    UserID: "123",
    Email:  "user@example.com",
})

// Publish sync (wait for handlers)
bus.PublishSync(ctx, msg)
```

### Router (HTTP)
```go
router := router.NewChiRouter(container)

// Routes
router.GET("/", func(ctx touta.Context) error {
    return ctx.HTML(200, "<h1>Home</h1>")
})

router.POST("/users", func(ctx touta.Context) error {
    return ctx.JSON(201, map[string]string{"id": "123"})
})

router.PUT("/users/:id", func(ctx touta.Context) error {
    id := ctx.Param("id")
    return ctx.String(200, "Updated "+id)
})

// Groups
api := router.Group("/api")
api.GET("/status", statusHandler)

// Start server
router.Listen(":8080")
```

### Context (HTTP Request/Response)
```go
func handler(ctx touta.Context) error {
    // Request data
    id := ctx.Param("id")
    name := ctx.Query("name")
    
    // Store/retrieve context values
    ctx.Set("user", user)
    user := ctx.Get("user")
    
    // Responses
    ctx.JSON(200, data)           // JSON response
    ctx.HTML(200, html)            // HTML response
    ctx.String(200, text)          // Plain text
    ctx.Redirect(302, "/login")    // Redirect
    
    // Access container
    container := ctx.Container()
    service, _ := container.Make(...)
    
    return nil
}
```

### Config (Configuration)
```go
// touta.yaml
framework:
  mode: development
  debug: true
  log_level: info
server:
  host: localhost
  port: 8080

// Load config
config, err := config.LoadOrDefault("touta.yaml")

// Access values
port := config.Server.Port
mode := config.Framework.Mode

// Environment variables
server:
  host: ${HOST:localhost}
  port: ${PORT:8080}
```

### Template (Rendering)
```go
renderer := template.NewHTMLRenderer()

// Parse templates
renderer.Parse("templates/*.html")

// Register functions
renderer.RegisterFunction("upper", strings.ToUpper)

// Render
html, err := renderer.Render("home", data)

// Render to HTTP response
renderer.Execute("home", data, w)
```

## CLI Commands

```bash
# Create new project
touta new <name>

# Initialize in existing directory
touta init

# Start development server
touta serve [--port 8080] [--host localhost]

# Show version
touta version
```

## Project Structure

```
my-app/
├── touta.yaml          # Configuration
├── main.go             # Entry point
├── handlers/           # Message handlers
├── templates/          # HTML templates
├── static/             # Static assets
└── config/             # Additional configs
```

## Testing

```go
// DI Container
func TestHandler(t *testing.T) {
    container := di.NewContainer()
    container.Bind((*Storage)(nil), &MockStorage{})
    
    handler := &MyHandler{}
    container.AutoWire(handler)
    
    // Test handler
}

// Message Bus
func TestMessage(t *testing.T) {
    bus := message.NewBus()
    bus.Start(context.Background())
    defer bus.Stop(context.Background())
    
    handler := &TestHandler{}
    bus.Subscribe("test", handler)
    bus.PublishSync(context.Background(), msg)
    
    // Assert handler was called
}

// HTTP Router
func TestRoute(t *testing.T) {
    container := di.NewContainer()
    router := router.NewChiRouter(container)
    
    router.GET("/test", myHandler)
    
    req := httptest.NewRequest("GET", "/test", nil)
    w := httptest.NewRecorder()
    router.Native().(*chi.Mux).ServeHTTP(w, req)
    
    // Assert response
}
```

## Common Patterns

### Message Handler
```go
type UserHandler struct {
    DB      Database      `inject:""`
    Logger  Logger        `inject:""`
}

func (h *UserHandler) Handle(ctx context.Context, msg touta.Message) (touta.Message, error) {
    userMsg := msg.(*UserCreated)
    
    // Business logic
    err := h.DB.Save(userMsg)
    if err != nil {
        h.Logger.Error("Failed to save user", err)
        return nil, err
    }
    
    return nil, nil
}
```

### HTTP Handler with Message
```go
router.POST("/users", func(ctx touta.Context) error {
    // Get dependencies from container
    bus, _ := ctx.Container().Make(reflect.TypeOf((*touta.MessageBus)(nil)))
    
    // Publish message
    bus.(touta.MessageBus).Publish(
        ctx.Request().Context(),
        &UserCreated{UserID: "123"},
    )
    
    return ctx.JSON(201, map[string]string{"status": "created"})
})
```

### Service Provider
```go
type DatabaseProvider struct{}

func (p *DatabaseProvider) Register(c touta.Container) error {
    c.Singleton((*Database)(nil), &PostgresDB{})
    return nil
}

func (p *DatabaseProvider) Boot(c touta.Container) error {
    db, _ := c.Make(reflect.TypeOf((*Database)(nil)))
    db.(*PostgresDB).Connect()
    return nil
}
```

## Tips

1. **Always use interfaces** - Never depend on concrete types
2. **Tag dependencies** - Use `inject:""` for auto-wiring
3. **Message slugs** - Use dot notation: `domain.action` (e.g., `user.created`)
4. **Error handling** - Always check and log errors
5. **Testing** - Mock all dependencies via container

## Resources

- **README.md** - Getting started guide
- **IMPLEMENTATION_SUMMARY.md** - Detailed technical overview
- **PHASE1_COMPLETE.md** - Implementation status
- **pkg/touta/interfaces.go** - All interface definitions

## Version

Toutā v0.1.0 - Phase 1 Foundation
