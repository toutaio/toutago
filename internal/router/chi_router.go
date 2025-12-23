package router

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/toutaio/toutago/pkg/touta"
)

// chiRouter implements Router using the Chi router.
type chiRouter struct {
	mux       *chi.Mux
	container touta.Container
}

// NewChiRouter creates a new Chi-based router.
func NewChiRouter(container touta.Container) touta.Router {
	return &chiRouter{
		mux:       chi.NewRouter(),
		container: container,
	}
}

// GET registers a handler for GET requests.
func (r *chiRouter) GET(path string, handler touta.HandlerFunc) {
	r.mux.Get(path, r.adapt(handler))
}

// POST registers a handler for POST requests.
func (r *chiRouter) POST(path string, handler touta.HandlerFunc) {
	r.mux.Post(path, r.adapt(handler))
}

// PUT registers a handler for PUT requests.
func (r *chiRouter) PUT(path string, handler touta.HandlerFunc) {
	r.mux.Put(path, r.adapt(handler))
}

// DELETE registers a handler for DELETE requests.
func (r *chiRouter) DELETE(path string, handler touta.HandlerFunc) {
	r.mux.Delete(path, r.adapt(handler))
}

// PATCH registers a handler for PATCH requests.
func (r *chiRouter) PATCH(path string, handler touta.HandlerFunc) {
	r.mux.Patch(path, r.adapt(handler))
}

// Group creates a route group with a prefix.
func (r *chiRouter) Group(prefix string) touta.Router {
	subRouter := &chiRouter{
		mux:       chi.NewRouter(),
		container: r.container,
	}
	r.mux.Mount(prefix, subRouter.mux)
	return subRouter
}

// Use adds middleware to the router.
func (r *chiRouter) Use(middleware ...touta.MiddlewareFunc) {
	for _, mw := range middleware {
		r.mux.Use(r.adaptMiddleware(mw))
	}
}

// Listen starts the HTTP server on the given address.
func (r *chiRouter) Listen(addr string) error {
	server := &http.Server{
		Addr:         addr,
		Handler:      r.mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	return server.ListenAndServe()
}

// Native returns the underlying Chi router.
func (r *chiRouter) Native() interface{} {
	return r.mux
}

// adapt converts a touta.HandlerFunc to http.HandlerFunc.
func (r *chiRouter) adapt(handler touta.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := NewContext(w, req, r.container)
		if err := handler(ctx); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// adaptMiddleware converts touta.MiddlewareFunc to Chi middleware.
func (r *chiRouter) adaptMiddleware(mw touta.MiddlewareFunc) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			ctx := NewContext(w, req, r.container)

			// Wrap next handler
			wrappedHandler := func(c touta.Context) error {
				next.ServeHTTP(w, req)
				return nil
			}

			// Call middleware
			handler := mw(wrappedHandler)
			if err := handler(ctx); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		})
	}
}

// defaultContext implements the Context interface.
type defaultContext struct {
	req       *http.Request
	res       http.ResponseWriter
	container touta.Container
	data      map[string]interface{}
}

// NewContext creates a new request context.
func NewContext(w http.ResponseWriter, req *http.Request, container touta.Container) touta.Context {
	return &defaultContext{
		req:       req,
		res:       w,
		container: container,
		data:      make(map[string]interface{}),
	}
}

// Request returns the HTTP request.
func (c *defaultContext) Request() *http.Request {
	return c.req
}

// Response returns the HTTP response writer.
func (c *defaultContext) Response() http.ResponseWriter {
	return c.res
}

// Param retrieves a URL parameter by name.
func (c *defaultContext) Param(key string) string {
	return chi.URLParam(c.req, key)
}

// Query retrieves a query string parameter.
func (c *defaultContext) Query(key string) string {
	return c.req.URL.Query().Get(key)
}

// Get retrieves a value from the context.
func (c *defaultContext) Get(key string) interface{} {
	// First check our data map
	if val, ok := c.data[key]; ok {
		return val
	}
	// Then check request context
	return c.req.Context().Value(key)
}

// Set stores a value in the context.
func (c *defaultContext) Set(key string, value interface{}) {
	c.data[key] = value
	// Also store in request context
	ctx := context.WithValue(c.req.Context(), key, value)
	c.req = c.req.WithContext(ctx)
}

// Container returns the DI container.
func (c *defaultContext) Container() touta.Container {
	return c.container
}

// JSON sends a JSON response.
func (c *defaultContext) JSON(status int, data interface{}) error {
	c.res.Header().Set("Content-Type", "application/json")
	c.res.WriteHeader(status)

	// Simple JSON encoding (could use encoding/json for real implementation)
	fmt.Fprintf(c.res, "%v", data)
	return nil
}

// String sends a plain text response.
func (c *defaultContext) String(status int, text string) error {
	c.res.Header().Set("Content-Type", "text/plain")
	c.res.WriteHeader(status)
	_, err := c.res.Write([]byte(text))
	return err
}

// HTML sends an HTML response.
func (c *defaultContext) HTML(status int, html string) error {
	c.res.Header().Set("Content-Type", "text/html; charset=utf-8")
	c.res.WriteHeader(status)
	_, err := c.res.Write([]byte(html))
	return err
}

// Redirect redirects to another URL.
func (c *defaultContext) Redirect(status int, url string) error {
	http.Redirect(c.res, c.req, url, status)
	return nil
}

// Status sets the response status code.
func (c *defaultContext) Status(status int) touta.Context {
	c.res.WriteHeader(status)
	return c
}
