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
// Message Bus Interfaces
// ============================================================================

// Message represents a message that flows through the system.
// All messages must have a unique slug, type, and optional metadata.
type Message interface {
	// Slug returns the unique identifier for this message (e.g., "user.registered")
	Slug() string

	// Type returns the message category (e.g., "event", "command", "query")
	Type() string

	// Metadata returns additional message metadata
	Metadata() map[string]interface{}
}

// MessageHandler processes incoming messages and optionally returns a response.
type MessageHandler interface {
	// Handle processes a message and returns an optional response message
	Handle(ctx context.Context, msg Message) (Message, error)
}

// MessageBus coordinates message publishing and subscription.
// It supports both synchronous and asynchronous message dispatch.
type MessageBus interface {
	// Publish sends a message asynchronously to all subscribers
	Publish(ctx context.Context, msg Message) error

	// PublishSync sends a message synchronously and waits for handlers to complete
	PublishSync(ctx context.Context, msg Message) error

	// Subscribe registers a handler for messages of a specific type or slug
	Subscribe(pattern string, handler MessageHandler) error

	// Unsubscribe removes a handler for a specific pattern
	Unsubscribe(pattern string, handler MessageHandler) error

	// Start begins processing messages (for async bus implementations)
	Start(ctx context.Context) error

	// Stop gracefully shuts down the message bus
	Stop(ctx context.Context) error
}

// ============================================================================
// Router Interfaces
// ============================================================================

// HandlerFunc is the signature for HTTP request handlers.
type HandlerFunc func(Context) error

// MiddlewareFunc wraps a HandlerFunc to provide cross-cutting concerns.
type MiddlewareFunc func(HandlerFunc) HandlerFunc

// Router provides HTTP routing abstraction.
// The default implementation uses Chi, but other routers can be swapped in.
type Router interface {
	// GET registers a handler for GET requests
	GET(path string, handler HandlerFunc)

	// POST registers a handler for POST requests
	POST(path string, handler HandlerFunc)

	// PUT registers a handler for PUT requests
	PUT(path string, handler HandlerFunc)

	// DELETE registers a handler for DELETE requests
	DELETE(path string, handler HandlerFunc)

	// PATCH registers a handler for PATCH requests
	PATCH(path string, handler HandlerFunc)

	// Group creates a route group with a prefix
	Group(prefix string) Router

	// Use adds middleware to the router
	Use(middleware ...MiddlewareFunc)

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

	// Query retrieves a query string parameter
	Query(key string) string

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

	// Redirect redirects to another URL
	Redirect(status int, url string) error

	// Status sets the response status code
	Status(status int) Context
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
	Mode        string `yaml:"mode"`         // development, production
	Debug       bool   `yaml:"debug"`        // enable debug logging
	HotReload   bool   `yaml:"hot_reload"`   // enable hot reload in dev
	LogLevel    string `yaml:"log_level"`    // trace, debug, info, warn, error
	Timezone    string `yaml:"timezone"`     // default timezone
}

// RouterConfig contains HTTP router settings.
type RouterConfig struct {
	BasePath     string            `yaml:"base_path"`     // base URL path
	Middleware   []string          `yaml:"middleware"`    // global middleware
	CORS         CORSConfig        `yaml:"cors"`          // CORS settings
	RateLimit    RateLimitConfig   `yaml:"rate_limit"`    // rate limiting
	Static       []StaticConfig    `yaml:"static"`        // static file serving
}

// ServerConfig contains HTTP server settings.
type ServerConfig struct {
	Host            string `yaml:"host"`              // bind host
	Port            int    `yaml:"port"`              // bind port
	ReadTimeout     int    `yaml:"read_timeout"`      // seconds
	WriteTimeout    int    `yaml:"write_timeout"`     // seconds
	IdleTimeout     int    `yaml:"idle_timeout"`      // seconds
	MaxHeaderBytes  int    `yaml:"max_header_bytes"`  // bytes
	TLS             TLSConfig `yaml:"tls"`            // TLS settings
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
	Path   string `yaml:"path"`   // URL path
	Dir    string `yaml:"dir"`    // filesystem directory
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
	Name        string                 `yaml:"name"`
	Version     string                 `yaml:"version"`
	Type        string                 `yaml:"type"` // package, recipe, component
	Path        string                 `yaml:"path"`
	Handlers    []string               `yaml:"handlers"`
	Templates   []string               `yaml:"templates"`
	Routes      []string               `yaml:"routes"`
	Assets      []string               `yaml:"assets"`
	Metadata    map[string]interface{} `yaml:"metadata"`
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
