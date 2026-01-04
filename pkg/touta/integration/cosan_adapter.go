package integration

import (
	"net/http"

	"github.com/toutaio/toutago-cosan-router/pkg/cosan"
	"github.com/toutaio/toutago/pkg/touta"
)

// CosanRouterAdapter adapts cosan.Router to implement touta.Router interface.
type CosanRouterAdapter struct {
	router    cosan.Router
	container touta.Container
}

// NewRouter creates a new HTTP router using cosan as the underlying implementation.
func NewRouter(container touta.Container, options ...cosan.Option) touta.Router {
	return &CosanRouterAdapter{
		router:    cosan.New(options...),
		container: container,
	}
}

// GET registers a handler for GET requests.
func (a *CosanRouterAdapter) GET(path string, handler touta.HTTPHandlerFunc) {
	a.router.GET(path, a.adaptHandler(handler))
}

// POST registers a handler for POST requests.
func (a *CosanRouterAdapter) POST(path string, handler touta.HTTPHandlerFunc) {
	a.router.POST(path, a.adaptHandler(handler))
}

// PUT registers a handler for PUT requests.
func (a *CosanRouterAdapter) PUT(path string, handler touta.HTTPHandlerFunc) {
	a.router.PUT(path, a.adaptHandler(handler))
}

// DELETE registers a handler for DELETE requests.
func (a *CosanRouterAdapter) DELETE(path string, handler touta.HTTPHandlerFunc) {
	a.router.DELETE(path, a.adaptHandler(handler))
}

// PATCH registers a handler for PATCH requests.
func (a *CosanRouterAdapter) PATCH(path string, handler touta.HTTPHandlerFunc) {
	a.router.PATCH(path, a.adaptHandler(handler))
}

// OPTIONS registers a handler for OPTIONS requests.
func (a *CosanRouterAdapter) OPTIONS(path string, handler touta.HTTPHandlerFunc) {
	a.router.OPTIONS(path, a.adaptHandler(handler))
}

// HEAD registers a handler for HEAD requests.
func (a *CosanRouterAdapter) HEAD(path string, handler touta.HTTPHandlerFunc) {
	a.router.HEAD(path, a.adaptHandler(handler))
}

// Group creates a route group with a prefix.
func (a *CosanRouterAdapter) Group(prefix string) touta.Router {
	return &CosanRouterAdapter{
		router:    a.router.Group(prefix),
		container: a.container,
	}
}

// Use adds middleware to the router.
func (a *CosanRouterAdapter) Use(middleware ...touta.MiddlewareFunc) {
	for _, mw := range middleware {
		cosanMw := a.adaptMiddleware(mw)
		a.router.Use(cosanMw)
	}
}

// Listen starts the HTTP server on the given address.
func (a *CosanRouterAdapter) Listen(addr string) error {
	return a.router.Listen(addr)
}

// Native returns the underlying cosan router for advanced usage.
func (a *CosanRouterAdapter) Native() interface{} {
	return a.router
}

// adaptHandler converts a touta.HTTPHandlerFunc to a cosan.HandlerFunc.
func (a *CosanRouterAdapter) adaptHandler(handler touta.HTTPHandlerFunc) cosan.HandlerFunc {
	return func(ctx cosan.Context) error {
		toutaCtx := &CosanContextAdapter{
			ctx:       ctx,
			container: a.container,
		}
		return handler(toutaCtx)
	}
}

// adaptMiddleware converts a touta.MiddlewareFunc to a cosan.Middleware.
func (a *CosanRouterAdapter) adaptMiddleware(mw touta.MiddlewareFunc) cosan.Middleware {
	return cosan.MiddlewareFunc(func(next cosan.HandlerFunc) cosan.HandlerFunc {
		return func(ctx cosan.Context) error {
			toutaCtx := &CosanContextAdapter{
				ctx:       ctx,
				container: a.container,
			}
			
			nextHandler := func(tc touta.Context) error {
				return next(ctx)
			}
			
			wrappedHandler := mw(nextHandler)
			return wrappedHandler(toutaCtx)
		}
	})
}

// CosanContextAdapter adapts cosan.Context to implement touta.Context interface.
type CosanContextAdapter struct {
	ctx       cosan.Context
	container touta.Container
}

// Request returns the HTTP request.
func (a *CosanContextAdapter) Request() *http.Request {
	return a.ctx.Request()
}

// Response returns the HTTP response writer.
func (a *CosanContextAdapter) Response() http.ResponseWriter {
	return a.ctx.Response()
}

// Param retrieves a URL parameter by name.
func (a *CosanContextAdapter) Param(key string) string {
	return a.ctx.Param(key)
}

// Query retrieves a query string parameter.
func (a *CosanContextAdapter) Query(key string) string {
	return a.ctx.Query(key)
}

// Get retrieves a value from the context.
func (a *CosanContextAdapter) Get(key string) interface{} {
	return a.ctx.Get(key)
}

// Set stores a value in the context.
func (a *CosanContextAdapter) Set(key string, value interface{}) {
	a.ctx.Set(key, value)
}

// Container returns the DI container.
func (a *CosanContextAdapter) Container() touta.Container {
	return a.container
}

// JSON sends a JSON response.
func (a *CosanContextAdapter) JSON(status int, data interface{}) error {
	return a.ctx.JSON(status, data)
}

// String sends a plain text response.
func (a *CosanContextAdapter) String(status int, text string) error {
	return a.ctx.String(status, text)
}

// HTML sends an HTML response.
func (a *CosanContextAdapter) HTML(status int, html string) error {
	return a.ctx.HTML(status, html)
}

// Redirect redirects to another URL.
func (a *CosanContextAdapter) Redirect(status int, url string) error {
	http.Redirect(a.Response(), a.Request(), url, status)
	return nil
}

// Status sets the response status code.
func (a *CosanContextAdapter) Status(status int) touta.Context {
	a.ctx.Status(status)
	return a
}

// Params returns all URL parameters as a map.
func (a *CosanContextAdapter) Params() map[string]string {
	return a.ctx.Params()
}

// QueryAll retrieves all values for a query string parameter.
func (a *CosanContextAdapter) QueryAll(key string) []string {
	return a.ctx.QueryAll(key)
}

// Bind parses the request body into the provided struct.
func (a *CosanContextAdapter) Bind(v interface{}) error {
	return a.ctx.Bind(v)
}

// BodyBytes returns the raw request body as bytes.
func (a *CosanContextAdapter) BodyBytes() ([]byte, error) {
	return a.ctx.BodyBytes()
}

// Header returns the response header map.
func (a *CosanContextAdapter) Header() http.Header {
	return a.ctx.Header()
}

// Write writes the response body bytes.
func (a *CosanContextAdapter) Write(b []byte) (int, error) {
	return a.ctx.Write(b)
}

// Render renders a template with the given data.
// This method provides integration with Fith renderer if available in the container.
func (a *CosanContextAdapter) Render(template string, data interface{}) error {
	// Try to get renderer from container
	if a.container != nil {
		if a.container.Has((*touta.TemplateRenderer)(nil)) {
			rendererInterface, err := a.container.Make((*touta.TemplateRenderer)(nil))
			if err == nil && rendererInterface != nil {
				if renderer, ok := rendererInterface.(touta.TemplateRenderer); ok && renderer != nil {
					rendered, err := renderer.Render(template, data)
					if err != nil {
						return err
					}
					return a.HTML(200, string(rendered))
				}
			}
		}
	}
	// No renderer available
	return a.HTML(500, "Template renderer not configured")
}

// Publish publishes a message to the event bus.
// This method provides integration with Scéla bus if available in the container.
func (a *CosanContextAdapter) Publish(topic string, payload interface{}) error {
	// Try to get bus from container
	if a.container != nil && a.container.Has((*touta.Bus)(nil)) {
		busInterface, err := a.container.Make((*touta.Bus)(nil))
		if err == nil {
			if bus, ok := busInterface.(touta.Bus); ok {
				// Check if it's a pointer and not nil
				if bus != nil {
					return bus.Publish(a.Request().Context(), topic, payload)
				}
			}
		}
	}
	// No bus available - silently ignore
	return nil
}
