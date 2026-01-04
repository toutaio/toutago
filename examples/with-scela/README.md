# Scéla Message Bus Example

This example demonstrates advanced features of the Scéla message bus integrated with Toutā.

## Features Demonstrated

### 1. Pattern Matching
Subscribe to multiple topics using wildcard patterns:
```go
bus.Subscribe("user.*", handler)  // Matches user.created, user.deleted, etc.
```

### 2. Middleware
Chain middleware for cross-cutting concerns:
```go
// Global middleware applied to all messages
bus := integration.NewScelaBusWithMiddleware(LoggingMiddleware)

// Handler-specific middleware
handler := ValidationMiddleware(actualHandler)
bus.Subscribe("user.created", handler)
```

### 3. Priority Levels
Process high-priority messages first:
```go
bus.Subscribe("notification.urgent", handler, scela.WithPriority(scela.PriorityHigh))
bus.Publish(ctx, "notification.urgent", data, scela.WithPriority(scela.PriorityHigh))
```

### 4. Async vs Sync Publishing
- **Async**: Non-blocking, handlers run in background
- **Sync**: Blocks until all handlers complete

```go
// Async - returns immediately
bus.Publish(ctx, "user.created", data)

// Sync - waits for handlers
bus.PublishSync(ctx, "user.updated", data)
```

## Running the Example

```bash
cd examples/with-scela
go run main.go
```

## Testing

In another terminal:

```bash
# Create a user (async event)
curl -X POST http://localhost:8080/users

# Delete a user (async event)
curl -X DELETE http://localhost:8080/users

# Send urgent notification (high priority)
curl -X POST http://localhost:8080/notify

# Trigger synchronous event
curl -X POST http://localhost:8080/sync-event
```

## Output

You'll see output showing:
- Middleware logging for each message
- Pattern matching (user.* catches all user events)
- Specific handlers for user.created and user.deleted
- Validation middleware in action
- High-priority message processing
- Synchronous vs asynchronous behavior

## Key Concepts

### Middleware Chain
Middleware wraps handlers to add functionality:
1. Request logging
2. Validation
3. Error handling
4. Metrics collection
5. Authentication/authorization

### Pattern Matching
Use wildcards to subscribe to multiple topics:
- `user.*` - All user events
- `*.created` - All creation events
- `app.*.error` - All app errors

### Priority Processing
Messages with higher priority are processed first:
- `PriorityHigh` - Urgent notifications, critical errors
- `PriorityNormal` - Regular events (default)
- `PriorityLow` - Background tasks, analytics

### Sync vs Async
- **Async** (default): Fast response, eventual consistency
- **Sync**: Guaranteed processing, slower response

## Advanced Usage

For more advanced features including:
- Persistence (Redis, filesystem)
- Dead letter queues
- Retry policies
- Monitoring and metrics

See the `examples/scela-advanced/` directory.

## Documentation

- [Toutā Documentation](../../README.md)
- [Scéla Bus Documentation](https://github.com/toutaio/toutago-scela-bus)
- [Message Bus Guide](../../docs/message-bus.md)
