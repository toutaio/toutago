# Toutā Router - Advanced Patterns Guide

This guide covers advanced routing patterns using the Cosan-based router in Toutā.

## Table of Contents

- [Basic Routing](#basic-routing)
- [Route Parameters](#route-parameters)
- [Route Groups](#route-groups)
- [Middleware](#middleware)
- [Request/Response Hooks](#requestresponse-hooks)
- [Error Handling](#error-handling)
- [Route Introspection](#route-introspection)
- [Ecosystem Integration](#ecosystem-integration)
- [Performance](#performance)

## Basic Routing

### HTTP Methods

```go
router := integration.NewRouter(container)

router.GET("/users", listUsers)
router.POST("/users", createUser)
router.PUT("/users/:id", updateUser)
router.DELETE("/users/:id", deleteUser)
router.PATCH("/users/:id", patchUser)
router.OPTIONS("/users", optionsUsers)
router.HEAD("/users", headUsers)
```

### Handler Signature

Handlers follow the `touta.HTTPHandlerFunc` signature:

```go
type HTTPHandlerFunc func(Context) error

func listUsers(ctx touta.Context) error {
    users := []User{...}
    return ctx.JSON(200, users)
}
```

## Route Parameters

### Single Parameter

```go
router.GET("/users/:id", func(ctx touta.Context) error {
    id := ctx.Param("id")
    return ctx.JSON(200, map[string]string{"id": id})
})
```

### Multiple Parameters

```go
router.GET("/users/:userId/posts/:postId", func(ctx touta.Context) error {
    params := ctx.Params()
    return ctx.JSON(200, params)
})
```

### Query Parameters

```go
router.GET("/search", func(ctx touta.Context) error {
    q := ctx.Query("q")           // Single value
    tags := ctx.QueryAll("tag")   // Multiple values
    
    return ctx.JSON(200, map[string]interface{}{
        "query": q,
        "tags":  tags,
    })
})
```

## Route Groups

### Basic Grouping

```go
api := router.Group("/api")
api.GET("/users", listUsers)       // /api/users
api.POST("/users", createUser)     // /api/users
```

### Nested Groups

```go
api := router.Group("/api")
v1 := api.Group("/v1")
v1.GET("/users", listUsersV1)      // /api/v1/users

v2 := api.Group("/v2")
v2.GET("/users", listUsersV2)      // /api/v2/users
```

### Group Middleware

```go
api := router.Group("/api")

// Middleware applies only to /api/* routes
api.Use(authMiddleware)
api.Use(rateLimitMiddleware)

api.GET("/users", listUsers)
```

## Middleware

### Simple Middleware

```go
func loggingMiddleware(next touta.HTTPHandlerFunc) touta.HTTPHandlerFunc {
    return func(ctx touta.Context) error {
        log.Printf("%s %s", ctx.Request().Method, ctx.Request().URL.Path)
        return next(ctx)
    }
}

router.Use(loggingMiddleware)
```

### Middleware with State

```go
func timingMiddleware(next touta.HTTPHandlerFunc) touta.HTTPHandlerFunc {
    return func(ctx touta.Context) error {
        start := time.Now()
        err := next(ctx)
        duration := time.Since(start)
        
        ctx.Header().Set("X-Response-Time", duration.String())
        return err
    }
}
```

### Authentication Middleware

```go
func authMiddleware(next touta.HTTPHandlerFunc) touta.HTTPHandlerFunc {
    return func(ctx touta.Context) error {
        token := ctx.Request().Header.Get("Authorization")
        
        if token == "" {
            return ctx.JSON(401, map[string]string{
                "error": "Unauthorized",
            })
        }
        
        // Validate token and set user in context
        ctx.Set("user", userFromToken(token))
        return next(ctx)
    }
}
```

## Request/Response Hooks

### BeforeRequest Hook

Runs before every request, can abort the request:

```go
router.BeforeRequest(func(req *http.Request) error {
    // Check rate limit
    if rateLimitExceeded(req.RemoteAddr) {
        return errors.New("rate limit exceeded")
    }
    
    // Log request
    log.Printf("[REQUEST] %s %s", req.Method, req.URL.Path)
    return nil
})
```

### AfterResponse Hook

Runs after every response, cannot abort:

```go
router.AfterResponse(func(req *http.Request, statusCode int) {
    // Record metrics
    metrics.RecordResponse(req.URL.Path, statusCode)
    
    // Alert on errors
    if statusCode >= 500 {
        alert.ServerError(req.URL.Path)
    }
})
```

### Multiple Hooks

Hooks execute in registration order:

```go
router.BeforeRequest(hook1)  // Runs first
router.BeforeRequest(hook2)  // Runs second
router.BeforeRequest(hook3)  // Runs third
```

## Error Handling

### Custom Error Handler

```go
type APIError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
}

func (e *APIError) Error() string {
    return e.Message
}

router.SetErrorHandler(func(ctx touta.Context, err error) {
    // Handle custom errors
    if apiErr, ok := err.(*APIError); ok {
        ctx.JSON(apiErr.Code, map[string]interface{}{
            "error": apiErr.Message,
            "code":  apiErr.Code,
        })
        return
    }
    
    // Default error response
    ctx.JSON(500, map[string]string{
        "error": "Internal server error",
    })
})
```

### Using Custom Errors

```go
router.GET("/users/:id", func(ctx touta.Context) error {
    id := ctx.Param("id")
    
    user, err := db.FindUser(id)
    if err != nil {
        return &APIError{
            Code:    404,
            Message: "User not found",
        }
    }
    
    return ctx.JSON(200, user)
})
```

## Route Introspection

### List All Routes

```go
routes := router.GetRoutes()
for _, route := range routes {
    fmt.Printf("%s %s\n", route.Method, route.Pattern)
}
```

### Find Route by Name

```go
route := router.FindRoute("users.show")
if route != nil {
    fmt.Printf("Found route: %s %s\n", route.Method, route.Pattern)
}
```

### Generate API Documentation

```go
routes := router.GetRoutes()
for _, route := range routes {
    fmt.Printf("## %s %s\n", route.Method, route.Pattern)
    if route.Description != "" {
        fmt.Printf("%s\n", route.Description)
    }
    if len(route.Tags) > 0 {
        fmt.Printf("Tags: %v\n", route.Tags)
    }
    fmt.Println()
}
```

## Ecosystem Integration

### Template Rendering (Fith)

```go
// Register renderer in container
container.Singleton((*touta.TemplateRenderer)(nil), renderer)

router.GET("/page", func(ctx touta.Context) error {
    return ctx.Render("home", map[string]string{
        "title": "Home Page",
    })
})
```

### Event Publishing (Scéla)

```go
// Register bus in container
container.Singleton((*touta.Bus)(nil), bus)

router.POST("/users", func(ctx touta.Context) error {
    var user User
    if err := ctx.Bind(&user); err != nil {
        return err
    }
    
    // Save user
    db.Save(&user)
    
    // Publish event
    ctx.Publish("user.created", user)
    
    return ctx.JSON(201, user)
})
```

### Request Binding (Datamapper)

```go
type CreateUserRequest struct {
    Name  string `json:"name" validate:"required"`
    Email string `json:"email" validate:"required,email"`
}

router.POST("/users", func(ctx touta.Context) error {
    var req CreateUserRequest
    if err := ctx.Bind(&req); err != nil {
        return err
    }
    
    // Use validated data
    user := createUser(req)
    return ctx.JSON(201, user)
})
```

## Performance

### Benchmarks

Results on reference hardware (16 CPU):

```
BenchmarkRouterSimpleRoute        979,614 ops    1,096 ns/op
BenchmarkRouterParametricRoute    869,697 ops    1,300 ns/op
BenchmarkRouterMultiParam         753,536 ops    1,366 ns/op
BenchmarkRouterManyRoutes       1,484,053 ops      798 ns/op
BenchmarkRouterWithMiddleware   1,384,324 ops      873 ns/op
BenchmarkRouterWithHooks        1,582,950 ops      764 ns/op
BenchmarkContextParamAccess     3,500,922 ops      341 ns/op
BenchmarkContextJSONResponse      618,445 ops    1,730 ns/op
BenchmarkRouterGrouping         1,561,191 ops      771 ns/op
```

### Performance Tips

1. **Minimize Middleware**: Each middleware adds overhead
2. **Use Hooks Wisely**: BeforeRequest/AfterResponse run on every request
3. **Cache Route Lookups**: GetRoutes() is expensive, cache the result
4. **Pool Buffers**: Reuse buffers for JSON encoding
5. **Keep Handlers Thin**: Move business logic to services

### Memory Usage

- Simple route: ~1,554 bytes, 18 allocations
- Parametric route: ~1,842 bytes, 19 allocations
- JSON response: ~1,538 bytes, 24 allocations

## Best Practices

### 1. Structure Routes Logically

```go
// Group by resource
users := router.Group("/users")
users.GET("", listUsers)
users.POST("", createUser)
users.GET("/:id", getUser)
users.PUT("/:id", updateUser)
users.DELETE("/:id", deleteUser)
```

### 2. Use Middleware for Cross-Cutting Concerns

```go
router.Use(
    recoveryMiddleware,   // Panic recovery
    loggingMiddleware,    // Request logging
    corsMiddleware,       // CORS headers
)

api := router.Group("/api")
api.Use(
    authMiddleware,       // Authentication
    rateLimitMiddleware,  // Rate limiting
)
```

### 3. Validate Input Early

```go
router.POST("/users", func(ctx touta.Context) error {
    var req CreateUserRequest
    if err := ctx.Bind(&req); err != nil {
        return &APIError{Code: 400, Message: "Invalid request"}
    }
    
    // Business logic here
})
```

### 4. Use Custom Error Types

```go
type ValidationError struct {
    Field   string `json:"field"`
    Message string `json:"message"`
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("%s: %s", e.Field, e.Message)
}
```

### 5. Document Routes

```go
// Use route introspection for documentation
routes := router.GetRoutes()
generateSwaggerDoc(routes)
```

## Complete Example

See `examples/advanced-routing/main.go` for a complete working example demonstrating:
- Request/response hooks
- Custom error handling
- Route introspection
- Middleware
- Authentication
- Metrics collection

## See Also

- [Cosan Router Documentation](../../toutago-cosan-router/README.md)
- [Ecosystem Integration Guide](../README.md)
- [Examples](../examples/)
