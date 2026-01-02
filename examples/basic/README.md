# Basic Example

This example demonstrates basic usage of the Toutā framework v2.0 with integrated component libraries:

- **DI Container**: Using `toutago-nasc-dependency-injector`
- **HTTP Router**: Using `toutago-cosan-router`  
- **Middleware**: Logging middleware example

## Running

```bash
cd examples/basic
go run main.go
```

Then visit:
- http://localhost:8080/ - Welcome message
- http://localhost:8080/health - Health check

## Features Demonstrated

1. **Dependency Injection** - Binding and resolving services
2. **HTTP Routing** - GET routes with JSON responses
3. **Middleware** - Request logging
4. **Context** - Using the framework context API
