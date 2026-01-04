// Package touta provides the core interfaces and contracts for the Toutā framework.
//
// Toutā is a Go web framework that emphasizes:
//   - Interface-first design for pluggability
//   - Message-passing architecture inspired by Smalltalk
//   - Dependency injection for testability
//   - Developer experience with hot reload and CLI tools
package touta

import (
	"context"
	"net/http"
	"time"
)

// ============================================================================
// Dependency Injection Interfaces
// ============================================================================

// Container manages dependency injection and service resolution.
// It supports binding interfaces to concrete implementations, singletons,
// factories, and auto-wiring via reflection.
type Container interface {
	// Bind registers an interface to a concrete implementation
	Bind(abstract interface{}, concrete interface{}) error

	// Singleton registers an interface to a singleton instance
	Singleton(abstract interface{}, concrete interface{}) error

	// Factory registers a factory function for creating instances
	Factory(abstract interface{}, factory func(Container) (interface{}, error)) error

	// Make resolves and returns an instance of the given interface
	Make(abstract interface{}) (interface{}, error)

	// MakeWith resolves an instance with additional parameters
	MakeWith(abstract interface{}, params map[string]interface{}) (interface{}, error)

	// Has checks if a binding exists for the given interface
	Has(abstract interface{}) bool

	// AutoWire injects dependencies into a struct using reflection
	AutoWire(target interface{}) error

	// Tagged returns all instances registered with the given tag
	Tagged(tag string) ([]interface{}, error)
}

// ServiceProvider registers services into the container during bootstrap.
type ServiceProvider interface {
	// Register binds services into the container
	Register(container Container) error

	// Boot is called after all providers are registered
	Boot(container Container) error
}

// ============================================================================
// Message Bus Interfaces (Scéla Integration)
// ============================================================================

// Message represents a message that flows through the Scéla message bus.
//
// Messages in Toutā use a topic-based routing system with support for
// wildcard patterns. Topics typically follow a hierarchical structure
// like "user.registered", "order.created", etc.
//
// Example:
//
//	type UserRegisteredEvent struct {
//	    UserID string
//	    Email  string
//	}
//
//	// Publish the message
//	bus.Publish(ctx, "user.registered", UserRegisteredEvent{...})
type Message interface {
	// Topic returns the message topic (e.g., "user.registered", "order.created")
	Topic() string

	// Payload returns the message payload data
	Payload() interface{}

	// Metadata returns additional message metadata
	Metadata() map[string]interface{}

	// ID returns the unique message identifier
	ID() string

	// Timestamp returns when the message was created
	Timestamp() time.Time
}

// Handler processes messages from the message bus.
//
// Handlers are registered with topic patterns and executed when matching
// messages are published. Patterns support wildcards:
//   - "*" matches any single segment (e.g., "user.*" matches "user.created", "user.updated")
//   - "**" matches any number of segments (e.g., "user.**" matches "user.created.admin")
//
// Example:
//
//	handler := touta.HandlerFunc(func(ctx context.Context, msg touta.Message) error {
//	    event := msg.Payload().(UserRegisteredEvent)
//	    log.Printf("User registered: %s", event.Email)
//	    return nil
//	})
//	bus.Subscribe("user.registered", handler)
type Handler interface {
	// Handle processes a message and returns an error if processing fails
	Handle(ctx context.Context, msg Message) error
}

// HandlerFunc is a function adapter that implements the Handler interface.
//
// This allows using plain functions as message handlers without
// creating a separate type.
//
// Example:
//
//	bus.Subscribe("user.*", touta.HandlerFunc(func(ctx context.Context, msg touta.Message) error {
//	    log.Printf("User event: %s", msg.Topic())
//	    return nil
//	}))
type HandlerFunc func(ctx context.Context, msg Message) error

// Handle implements the Handler interface.
func (f HandlerFunc) Handle(ctx context.Context, msg Message) error {
	return f(ctx, msg)
}

// Bus is the message bus interface for pub/sub messaging in Toutā.
//
// The Bus provides topic-based message routing with support for:
//   - Asynchronous and synchronous message dispatch
//   - Pattern-based subscriptions with wildcards
//   - Middleware for cross-cutting concerns
//   - Priority-based message handling
//   - Graceful shutdown
//
// Example usage:
//
//	// Create a bus with options
//	bus := integration.NewScelaBus(
//	    scela.WithWorkers(20),
//	    scela.WithMaxRetries(3),
//	)
//
//	// Subscribe to messages
//	bus.Subscribe("user.*", touta.HandlerFunc(func(ctx context.Context, msg touta.Message) error {
//	    log.Printf("User event: %s", msg.Topic())
//	    return nil
//	}))
//
//	// Publish messages
//	bus.Publish(ctx, "user.registered", userData)
type Bus interface {
	// Publish publishes a message asynchronously to matching subscribers.
	// Returns immediately without waiting for handlers to complete.
	Publish(ctx context.Context, topic string, payload interface{}) error

	// PublishSync publishes a message synchronously, waiting for all handlers to complete.
	// Returns an error if any handler fails.
	PublishSync(ctx context.Context, topic string, payload interface{}) error

	// Subscribe registers a handler for messages matching the topic pattern.
	// Patterns support wildcards: "*" for single segment, "**" for multiple segments.
	// Returns a Subscription that can be used to unsubscribe later.
	Subscribe(pattern string, handler Handler) (Subscription, error)

	// Use adds middleware to the bus for cross-cutting concerns.
	// Middleware is applied to all handlers in the order they are registered.
	Use(middleware ...Middleware)

	// Close gracefully shuts down the bus, waiting for in-flight messages to complete.
	Close() error
}

// Subscription represents a subscription to messages on the bus.
//
// Subscriptions can be used to unsubscribe from a topic pattern when
// the handler is no longer needed.
type Subscription interface {
	// Topic returns the subscription pattern
	Topic() string

	// Unsubscribe removes this subscription from the bus
	Unsubscribe() error
}

// Middleware wraps message handlers to provide cross-cutting concerns.
//
// Middleware can be used for logging, metrics, retry logic, circuit breaking,
// rate limiting, and other concerns that apply across multiple handlers.
//
// Example:
//
//	loggingMiddleware := func(next touta.Handler) touta.Handler {
//	    return touta.HandlerFunc(func(ctx context.Context, msg touta.Message) error {
//	        log.Printf("Processing: %s", msg.Topic())
//	        err := next.Handle(ctx, msg)
//	        if err != nil {
//	            log.Printf("Error: %v", err)
//	        }
//	        return err
//	    })
//	}
//	bus.Use(loggingMiddleware)
type Middleware func(Handler) Handler

// ============================================================================
// Router Interfaces
// ============================================================================

// HTTPHandlerFunc is the signature for HTTP request handlers.
type HTTPHandlerFunc func(Context) error

// MiddlewareFunc wraps an HTTPHandlerFunc to provide cross-cutting concerns.
type MiddlewareFunc func(HTTPHandlerFunc) HTTPHandlerFunc

// RequestHook is a function that runs before request processing.
type RequestHook func(req *http.Request) error

// ResponseHook is a function that runs after response is written.
type ResponseHook func(req *http.Request, statusCode int)

// ErrorHandler is a custom error handling function for the router.
type ErrorHandler func(ctx Context, err error)

// RouteInfo contains metadata about a registered route.
type RouteInfo struct {
	Method      string
	Pattern     string
	Name        string
	Description string
	Tags        []string
	Deprecated  bool
	Version     string
}

// Router provides HTTP routing abstraction.
// The default implementation uses Cosan, but other routers can be swapped in.
type Router interface {
	// GET registers a handler for GET requests
	GET(path string, handler HTTPHandlerFunc)

	// POST registers a handler for POST requests
	POST(path string, handler HTTPHandlerFunc)

	// PUT registers a handler for PUT requests
	PUT(path string, handler HTTPHandlerFunc)

	// DELETE registers a handler for DELETE requests
	DELETE(path string, handler HTTPHandlerFunc)

	// PATCH registers a handler for PATCH requests
	PATCH(path string, handler HTTPHandlerFunc)

	// OPTIONS registers a handler for OPTIONS requests
	OPTIONS(path string, handler HTTPHandlerFunc)

	// HEAD registers a handler for HEAD requests
	HEAD(path string, handler HTTPHandlerFunc)

	// Group creates a route group with a prefix
	Group(prefix string) Router

	// Use adds middleware to the router
	Use(middleware ...MiddlewareFunc)

	// BeforeRequest registers a hook to run before each request
	BeforeRequest(hook RequestHook)

	// AfterResponse registers a hook to run after each response
	AfterResponse(hook ResponseHook)

	// SetErrorHandler sets a custom error handler for the router
	SetErrorHandler(handler ErrorHandler)

	// GetRoutes returns all registered routes with metadata
	GetRoutes() []RouteInfo

	// FindRoute finds a route by name from its metadata
	FindRoute(name string) *RouteInfo

	// Listen starts the HTTP server on the given address
	Listen(addr string) error

	// Native returns the underlying router implementation
	Native() interface{}
}

// Context provides access to the HTTP request/response and framework services.
type Context interface {
	// Request returns the HTTP request
	Request() *http.Request

	// Response returns the HTTP response writer
	Response() http.ResponseWriter

	// Param retrieves a URL parameter by name
	Param(key string) string

	// Params returns all URL parameters as a map
	Params() map[string]string

	// Query retrieves a query string parameter
	Query(key string) string

	// QueryAll retrieves all values for a query string parameter
	QueryAll(key string) []string

	// Bind parses the request body into the provided struct
	Bind(v interface{}) error

	// BodyBytes returns the raw request body as bytes
	BodyBytes() ([]byte, error)

	// Get retrieves a value from the context
	Get(key string) interface{}

	// Set stores a value in the context
	Set(key string, value interface{})

	// Container returns the DI container
	Container() Container

	// JSON sends a JSON response
	JSON(status int, data interface{}) error

	// String sends a plain text response
	String(status int, text string) error

	// HTML sends an HTML response
	HTML(status int, html string) error

	// Render renders a template with the given data (uses Fith if available)
	Render(template string, data interface{}) error

	// Redirect redirects to another URL
	Redirect(status int, url string) error

	// Status sets the response status code
	Status(status int) Context

	// Header returns the response header map
	Header() http.Header

	// Write writes the response body bytes
	Write([]byte) (int, error)

	// Publish publishes a message to the event bus (uses Scéla if available)
	Publish(topic string, payload interface{}) error
}

// ============================================================================
// Configuration Interfaces
// ============================================================================

// ConfigLoader loads and manages configuration from various sources.
type ConfigLoader interface {
	// Load parses configuration from a file or source
	Load(source string) (*Config, error)

	// Watch monitors configuration for changes and triggers reload
	Watch(callback func(*Config)) error

	// Validate checks if the configuration is valid
	Validate(config *Config) error
}

// Config represents the framework configuration.
type Config struct {
	// Framework core settings
	Framework FrameworkConfig `yaml:"framework"`

	// Router settings
	Router RouterConfig `yaml:"router"`

	// Server settings
	Server ServerConfig `yaml:"server"`

	// Packages and components
	Packages map[string]interface{} `yaml:"packages"`

	// Custom application config
	App map[string]interface{} `yaml:"app"`
}

// FrameworkConfig contains core framework settings.
type FrameworkConfig struct {
	Mode      string `yaml:"mode"`       // development, production
	Debug     bool   `yaml:"debug"`      // enable debug logging
	HotReload bool   `yaml:"hot_reload"` // enable hot reload in dev
	LogLevel  string `yaml:"log_level"`  // trace, debug, info, warn, error
	Timezone  string `yaml:"timezone"`   // default timezone
}

// RouterConfig contains HTTP router settings.
type RouterConfig struct {
	BasePath   string          `yaml:"base_path"`  // base URL path
	Middleware []string        `yaml:"middleware"` // global middleware
	CORS       CORSConfig      `yaml:"cors"`       // CORS settings
	RateLimit  RateLimitConfig `yaml:"rate_limit"` // rate limiting
	Static     []StaticConfig  `yaml:"static"`     // static file serving
}

// ServerConfig contains HTTP server settings.
type ServerConfig struct {
	Host           string    `yaml:"host"`             // bind host
	Port           int       `yaml:"port"`             // bind port
	ReadTimeout    int       `yaml:"read_timeout"`     // seconds
	WriteTimeout   int       `yaml:"write_timeout"`    // seconds
	IdleTimeout    int       `yaml:"idle_timeout"`     // seconds
	MaxHeaderBytes int       `yaml:"max_header_bytes"` // bytes
	TLS            TLSConfig `yaml:"tls"`              // TLS settings
}

// CORSConfig contains CORS settings.
type CORSConfig struct {
	Enabled          bool     `yaml:"enabled"`
	AllowedOrigins   []string `yaml:"allowed_origins"`
	AllowedMethods   []string `yaml:"allowed_methods"`
	AllowedHeaders   []string `yaml:"allowed_headers"`
	ExposeHeaders    []string `yaml:"expose_headers"`
	AllowCredentials bool     `yaml:"allow_credentials"`
	MaxAge           int      `yaml:"max_age"`
}

// RateLimitConfig contains rate limiting settings.
type RateLimitConfig struct {
	Enabled  bool `yaml:"enabled"`
	Requests int  `yaml:"requests"` // requests per window
	Window   int  `yaml:"window"`   // window in seconds
}

// StaticConfig defines static file serving.
type StaticConfig struct {
	Path   string `yaml:"path"`    // URL path
	Dir    string `yaml:"dir"`     // filesystem directory
	MaxAge int    `yaml:"max_age"` // cache max age in seconds
}

// TLSConfig contains TLS/SSL settings.
type TLSConfig struct {
	Enabled  bool   `yaml:"enabled"`
	CertFile string `yaml:"cert_file"`
	KeyFile  string `yaml:"key_file"`
}

// ============================================================================
// Template Renderer Interface
// ============================================================================

// TemplateRenderer handles template parsing and rendering.
type TemplateRenderer interface {
	// Render executes a template with the given data
	Render(name string, data interface{}) ([]byte, error)

	// RegisterFunction adds a custom template function
	RegisterFunction(name string, fn interface{})

	// Parse loads templates from a pattern (e.g., "templates/*.html")
	Parse(pattern string) error

	// Execute renders a template to a writer
	Execute(name string, data interface{}, w http.ResponseWriter) error
}

// ============================================================================
// Component Registry Interface
// ============================================================================

// Component represents a registered package component.
type Component struct {
	Name      string                 `yaml:"name"`
	Version   string                 `yaml:"version"`
	Type      string                 `yaml:"type"` // package, recipe, component
	Path      string                 `yaml:"path"`
	Handlers  []string               `yaml:"handlers"`
	Templates []string               `yaml:"templates"`
	Routes    []string               `yaml:"routes"`
	Assets    []string               `yaml:"assets"`
	Metadata  map[string]interface{} `yaml:"metadata"`
}

// ComponentRegistry manages package and component registration.
type ComponentRegistry interface {
	// Register adds a component to the registry
	Register(component *Component) error

	// Get retrieves a component by name
	Get(name string) (*Component, error)

	// List returns all registered components
	List() ([]*Component, error)

	// LoadFromManifest parses a package.yaml manifest and registers components
	LoadFromManifest(path string) error

	// Has checks if a component is registered
	Has(name string) bool
}

// ============================================================================
// CLI Command Interface
// ============================================================================

// Command represents a CLI command that can be registered.
type Command interface {
	// Name returns the command name
	Name() string

	// Description returns the command description
	Description() string

	// Execute runs the command with the given context
	Execute(ctx CommandContext) error

	// Flags returns command-specific flags
	Flags() []Flag
}

// CommandContext provides access to command execution context.
type CommandContext interface {
	// Args returns command arguments
	Args() []string

	// Flag retrieves a flag value
	Flag(name string) interface{}

	// Container returns the DI container
	Container() Container

	// Config returns the framework configuration
	Config() *Config
}

// Flag represents a command-line flag.
type Flag struct {
	Name        string
	Short       string
	Description string
	Default     interface{}
	Required    bool
}
