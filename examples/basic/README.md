# Basic Example

This example demonstrates basic usage of the Toutā framework v2.0 with integrated component libraries:

- **DI Container**: Using `toutago-nasc-dependency-injector`
- **HTTP Router**: Using `toutago-cosan-router`
- **Message Bus**: Using `toutago-scela-bus`
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

1. **Dependency Injection** - Binding and resolving services with Nasc
2. **HTTP Routing** - GET routes with JSON responses using Cosán
3. **Event-Driven Messaging** - Publishing and subscribing to application events with Scéla
4. **Middleware** - Request logging
5. **Context** - Using the framework context API

## Message Bus Events

The example demonstrates Scéla message bus by:
- Publishing an `app.started` event on startup
- Subscribing to all `app.*` events with a logging handler
- Publishing `app.request` events for each HTTP request

This shows how you can add event-driven behavior to your application without coupling components directly.
