// Package main demonstrates advanced Cosan router features including hooks,
// error handling, and route introspection.
package main

import (
"fmt"
"log"
"net/http"
"time"

"github.com/toutaio/toutago/pkg/touta"
"github.com/toutaio/toutago/pkg/touta/integration"
)

// APIError represents a custom API error
type APIError struct {
Code    int    `json:"code"`
Message string `json:"message"`
}

func (e *APIError) Error() string {
return e.Message
}

func main() {
container := integration.NewContainer()
router := integration.NewRouter(container)

// ========================================================================
// 1. Request/Response Hooks
// ========================================================================

// BeforeRequest hook - runs before every request
router.BeforeRequest(func(req *http.Request) error {
log.Printf("[REQUEST] %s %s from %s", req.Method, req.URL.Path, req.RemoteAddr)

// Example: Check authentication
if req.Header.Get("X-API-Key") == "" && req.URL.Path != "/public" {
return &APIError{
Code:    401,
Message: "Missing API key",
}
}
return nil
})

// AfterResponse hook - runs after every request
router.AfterResponse(func(req *http.Request, statusCode int) {
log.Printf("[RESPONSE] %s %s -> %d", req.Method, req.URL.Path, statusCode)

// Example: Record metrics
if statusCode >= 500 {
log.Printf("[ALERT] Server error on %s %s", req.Method, req.URL.Path)
}
})

// ========================================================================
// 2. Custom Error Handler
// ========================================================================

router.SetErrorHandler(func(ctx touta.Context, err error) {
// Check if it's a custom API error
if apiErr, ok := err.(*APIError); ok {
ctx.JSON(apiErr.Code, map[string]interface{}{
"error": apiErr.Message,
"code":  apiErr.Code,
})
return
}

// Default error response
log.Printf("[ERROR] %v", err)
ctx.JSON(500, map[string]string{
"error": "Internal server error",
})
})

// ========================================================================
// 3. Route Handlers
// ========================================================================

// Public endpoint - no auth required
router.GET("/public", func(ctx touta.Context) error {
return ctx.JSON(200, map[string]string{
"message": "This is a public endpoint",
})
})

// Protected endpoint
router.GET("/users", func(ctx touta.Context) error {
return ctx.JSON(200, []map[string]interface{}{
{"id": 1, "name": "Alice"},
{"id": 2, "name": "Bob"},
})
})

// Endpoint that returns an error
router.GET("/error", func(ctx touta.Context) error {
return &APIError{
Code:    400,
Message: "This is a custom API error",
}
})

// Endpoint with parameter
router.GET("/users/:id", func(ctx touta.Context) error {
id := ctx.Param("id")
if id == "999" {
return &APIError{
Code:    404,
Message: "User not found",
}
}
return ctx.JSON(200, map[string]interface{}{
"id":   id,
"name": "User " + id,
})
})

// ========================================================================
// 4. Route Introspection
// ========================================================================

// List all routes at startup
routes := router.GetRoutes()
fmt.Printf("\n=== Registered Routes ===\n")
for _, route := range routes {
fmt.Printf("%s\t%s\n", route.Method, route.Pattern)
}
fmt.Printf("Total: %d routes\n\n", len(routes))

// ========================================================================
// 5. Middleware for specific routes
// ========================================================================

api := router.Group("/api")

// Add timing middleware to API group
api.Use(func(next touta.HTTPHandlerFunc) touta.HTTPHandlerFunc {
return func(ctx touta.Context) error {
start := time.Now()
err := next(ctx)
duration := time.Since(start)
log.Printf("[TIMING] %s took %v", ctx.Request().URL.Path, duration)
return err
}
})

api.GET("/stats", func(ctx touta.Context) error {
return ctx.JSON(200, map[string]interface{}{
"uptime":       "24h",
"requests":     12345,
"errors":       42,
"total_routes": len(routes),
})
})

// ========================================================================
// 6. Start Server
// ========================================================================

fmt.Println("Server starting on :8080")
fmt.Println("\nTry these requests:")
fmt.Println("  curl http://localhost:8080/public")
fmt.Println("  curl -H 'X-API-Key: secret' http://localhost:8080/users")
fmt.Println("  curl -H 'X-API-Key: secret' http://localhost:8080/users/123")
fmt.Println("  curl -H 'X-API-Key: secret' http://localhost:8080/error")
fmt.Println("  curl -H 'X-API-Key: secret' http://localhost:8080/api/stats")
fmt.Println()

log.Fatal(router.Listen(":8080"))
}
