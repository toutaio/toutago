# Message Bus Guide

Toutā uses the [Scéla message bus](https://github.com/toutaio/toutago-scela-bus) for event-driven communication within your application. This guide explains how to use it effectively.

## Table of Contents

- [Quick Start](#quick-start)
- [Core Concepts](#core-concepts)
- [Publishing Messages](#publishing-messages)
- [Subscribing to Messages](#subscribing-to-messages)
- [Pattern Matching](#pattern-matching)
- [Middleware](#middleware)
- [Priority Messages](#priority-messages)
- [Persistence](#persistence)
- [Best Practices](#best-practices)
- [Troubleshooting](#troubleshooting)

## Quick Start

```go
import (
    "context"
    "github.com/toutaio/toutago/pkg/touta"
    "github.com/toutaio/toutago/pkg/touta/integration"
)

func main() {
    // Create the message bus
    bus := integration.NewScelaBus()
    defer bus.Close()
    
    // Subscribe to events
    bus.Subscribe("user.created", touta.HandlerFunc(
        func(ctx context.Context, msg touta.Message) error {
            log.Printf("New user: %v", msg.Payload())
            return nil
        },
    ))
    
    // Publish an event
    bus.Publish(context.Background(), "user.created", map[string]interface{}{
        "id":    "123",
        "email": "user@example.com",
    })
}
```

## Core Concepts

### Messages

A message has:
- **Topic**: A hierarchical identifier (e.g., `user.created`, `order.payment.failed`)
- **Payload**: Any Go value (struct, map, primitive, etc.)
- **Metadata**: Additional key-value data
- **ID**: Unique identifier
- **Timestamp**: When the message was created

### Handlers

Handlers process messages. They implement the `touta.Handler` interface:

```go
type Handler interface {
    Handle(ctx context.Context, msg Message) error
}
```

Use `touta.HandlerFunc` for convenience:

```go
bus.Subscribe("topic", touta.HandlerFunc(
    func(ctx context.Context, msg touta.Message) error {
        // Process message
        return nil
    },
))
```

### Subscriptions

When you subscribe, you get a `Subscription` that can be used to unsubscribe:

```go
sub, err := bus.Subscribe("user.*", handler)
if err != nil {
    log.Fatal(err)
}

// Later...
sub.Unsubscribe()
```

## Publishing Messages

### Async Publishing (Default)

Messages are queued and processed asynchronously. The `Publish` method returns immediately:

```go
err := bus.Publish(ctx, "user.registered", userData)
```

**Use when:**
- You don't need to wait for handlers
- High throughput is important
- Eventual consistency is acceptable

### Sync Publishing

Wait for all handlers to complete before returning:

```go
err := bus.PublishSync(ctx, "order.created", orderData)
if err != nil {
    // One or more handlers failed
}
```

**Use when:**
- You need immediate feedback
- Order of operations matters
- Strong consistency is required

### Publishing with Priority

Process high-priority messages before normal ones:

```go
bus.PublishWithPriority(ctx, "alert.critical", alertData, scela.PriorityHigh)
```

Priority levels:
- `PriorityLow` - Background tasks, analytics
- `PriorityNormal` - Regular events (default)
- `PriorityHigh` - Important notifications
- `PriorityUrgent` - Critical alerts, system errors

## Subscribing to Messages

### Simple Subscription

```go
handler := touta.HandlerFunc(func(ctx context.Context, msg touta.Message) error {
    log.Printf("Received: %s", msg.Topic())
    return nil
})

sub, err := bus.Subscribe("user.created", handler)
```

### Multiple Handlers

Multiple handlers can subscribe to the same topic. They all receive the message:

```go
// Send welcome email
bus.Subscribe("user.created", sendWelcomeEmail)

// Create user profile
bus.Subscribe("user.created", createUserProfile)

// Send analytics event
bus.Subscribe("user.created", trackUserRegistration)
```

### Unsubscribing

```go
sub, _ := bus.Subscribe("topic", handler)

// Later...
sub.Unsubscribe()
```

## Pattern Matching

Subscribe to multiple topics using wildcards:

### Single-level Wildcard (`*`)

Matches exactly one segment:

```go
// Matches: user.created, user.updated, user.deleted
// Doesn't match: user.order.created
bus.Subscribe("user.*", handler)
```

### Multi-level Wildcard (`**`)

Matches zero or more segments:

```go
// Matches: app.error, app.database.error, app.api.auth.error
bus.Subscribe("app.**", handler)
```

### Combined Patterns

```go
// All creation events
bus.Subscribe("*.created", handler)

// All user-related events at any depth
bus.Subscribe("user.**", handler)

// Specific patterns
bus.Subscribe("order.payment.*", handler)  // order.payment.success, order.payment.failed
```

## Middleware

Middleware wraps handlers to add cross-cutting functionality.

### Creating Middleware

```go
func LoggingMiddleware(next touta.Handler) touta.Handler {
    return touta.HandlerFunc(func(ctx context.Context, msg touta.Message) error {
        start := time.Now()
        log.Printf("Processing %s", msg.Topic())
        
        err := next.Handle(ctx, msg)
        
        log.Printf("Completed %s in %v", msg.Topic(), time.Since(start))
        return err
    })
}
```

### Global Middleware

Applied to all handlers:

```go
bus := integration.NewScelaBusWithMiddleware(
    LoggingMiddleware,
    MetricsMiddleware,
    RecoveryMiddleware,
)
```

### Handler-Specific Middleware

Applied to specific handlers only:

```go
// Validation only for payment events
handler := ValidationMiddleware(paymentHandler)
bus.Subscribe("payment.*", handler)
```

### Common Middleware Patterns

#### Validation

```go
func ValidationMiddleware(next touta.Handler) touta.Handler {
    return touta.HandlerFunc(func(ctx context.Context, msg touta.Message) error {
        payload, ok := msg.Payload().(map[string]interface{})
        if !ok {
            return fmt.Errorf("invalid payload type")
        }
        
        // Validate required fields
        if _, ok := payload["id"]; !ok {
            return fmt.Errorf("missing required field: id")
        }
        
        return next.Handle(ctx, msg)
    })
}
```

#### Retry Logic

```go
func RetryMiddleware(maxRetries int) touta.Middleware {
    return func(next touta.Handler) touta.Handler {
        return touta.HandlerFunc(func(ctx context.Context, msg touta.Message) error {
            var err error
            for i := 0; i < maxRetries; i++ {
                err = next.Handle(ctx, msg)
                if err == nil {
                    return nil
                }
                time.Sleep(time.Second * time.Duration(i+1))
            }
            return err
        })
    }
}
```

#### Authentication

```go
func AuthMiddleware(next touta.Handler) touta.Handler {
    return touta.HandlerFunc(func(ctx context.Context, msg touta.Message) error {
        // Check authentication
        user := ctx.Value("user")
        if user == nil {
            return fmt.Errorf("unauthenticated")
        }
        
        return next.Handle(ctx, msg)
    })
}
```

## Priority Messages

Process important messages before others:

```go
// Subscribe with default priority
bus.Subscribe("notification.email", emailHandler)

// High-priority notifications processed first
bus.Subscribe("notification.sms", smsHandler)

// Publish with priority
bus.PublishWithPriority(ctx, "alert.critical", data, scela.PriorityHigh)
```

**Priority Queue Behavior:**
- Higher priority messages are dequeued first
- Within the same priority, FIFO order is maintained
- Useful for ensuring critical messages are processed quickly

## Persistence

For production deployments, persist messages to survive restarts.

### Redis Persistence

```go
import "github.com/toutaio/toutago-scela-bus/pkg/scela/persistence/redis"

store := redis.NewStore("localhost:6379", redis.Options{
    Password: "",
    DB:       0,
})

bus := integration.NewScelaBus(
    scela.WithPersistence(store),
)
```

### Filesystem Persistence

```go
import "github.com/toutaio/toutago-scela-bus/pkg/scela/persistence/fs"

store := fs.NewStore("./data/messages")

bus := integration.NewScelaBus(
    scela.WithPersistence(store),
)
```

## Best Practices

### Topic Naming

Use hierarchical, descriptive names:

```go
// Good
"user.created"
"order.payment.success"
"email.sent"

// Avoid
"uc"
"event123"
"stuff"
```

### Payload Design

Use structured data:

```go
// Good - typed struct
type UserCreatedEvent struct {
    ID        string
    Email     string
    Timestamp time.Time
}

// Good - map with clear keys
payload := map[string]interface{}{
    "id":    "123",
    "email": "user@example.com",
}

// Avoid - ambiguous data
payload := "123,user@example.com"
```

### Error Handling

Always handle errors from Publish:

```go
if err := bus.Publish(ctx, "event", data); err != nil {
    log.Printf("Failed to publish: %v", err)
    // Consider fallback behavior
}
```

Return errors from handlers to trigger retries:

```go
handler := touta.HandlerFunc(func(ctx context.Context, msg touta.Message) error {
    if err := processMessage(msg); err != nil {
        return err  // Will be retried
    }
    return nil
})
```

### Context Usage

Always respect context cancellation:

```go
handler := touta.HandlerFunc(func(ctx context.Context, msg touta.Message) error {
    select {
    case <-ctx.Done():
        return ctx.Err()
    default:
        // Process message
    }
    return nil
})
```

### Graceful Shutdown

Always close the bus on shutdown:

```go
bus := integration.NewScelaBus()
defer bus.Close()  // Waits for in-flight messages
```

### Avoid Blocking

Don't perform long-running operations in handlers:

```go
// Bad - blocks the worker
handler := touta.HandlerFunc(func(ctx context.Context, msg touta.Message) error {
    time.Sleep(10 * time.Second)  // Blocks worker
    return nil
})

// Good - spawn goroutine for long operations
handler := touta.HandlerFunc(func(ctx context.Context, msg touta.Message) error {
    go func() {
        time.Sleep(10 * time.Second)
        // Long operation
    }()
    return nil
})
```

## Troubleshooting

### Messages Not Being Processed

**Check:**
1. Is the bus properly initialized?
2. Are handlers registered before publishing?
3. Is the topic pattern correct?
4. Are there any errors in handler execution?

**Debug:**
```go
// Add logging middleware
bus.Use(LoggingMiddleware)
```

### Slow Processing

**Solutions:**
1. Increase worker count:
   ```go
   bus := integration.NewScelaBus(scela.WithWorkers(50))
   ```

2. Use async publishing:
   ```go
   bus.Publish(ctx, topic, data)  // Don't use PublishSync
   ```

3. Profile your handlers for bottlenecks

### Memory Issues

**Solutions:**
1. Limit buffer size:
   ```go
   bus := integration.NewScelaBus(scela.WithBufferSize(1000))
   ```

2. Use persistence to offload to disk/Redis
3. Process messages faster or increase workers

### Lost Messages

**Solutions:**
1. Enable persistence:
   ```go
   bus := integration.NewScelaBus(scela.WithPersistence(store))
   ```

2. Ensure proper shutdown:
   ```go
   defer bus.Close()  // Wait for pending messages
   ```

3. Check handler errors and enable retries

## Advanced Topics

### Dead Letter Queue

Handle messages that fail after max retries:

```go
dlqHandler := touta.HandlerFunc(func(ctx context.Context, msg touta.Message) error {
    log.Printf("Message failed permanently: %v", msg)
    // Send alert, log to database, etc.
    return nil
})

bus := integration.NewScelaBus(
    scela.WithDeadLetterHandler(dlqHandler),
    scela.WithMaxRetries(3),
)
```

### Metrics and Monitoring

Track message flow:

```go
func MetricsMiddleware(next touta.Handler) touta.Handler {
    return touta.HandlerFunc(func(ctx context.Context, msg touta.Message) error {
        start := time.Now()
        err := next.Handle(ctx, msg)
        duration := time.Since(start)
        
        // Record metrics
        recordMessageMetric(msg.Topic(), duration, err)
        
        return err
    })
}

bus := integration.NewScelaBusWithMiddleware(MetricsMiddleware)
```

### Integration with DI Container

```go
import "github.com/toutaio/toutago/pkg/touta/integration"

container := integration.NewContainer()
bus := integration.NewScelaBus()

// Register bus in container
container.Singleton((*touta.Bus)(nil), bus)

// Inject into handlers
type UserHandler struct {
    Bus touta.Bus `inject:""`
}
```

## Further Reading

- [Scéla Bus Repository](https://github.com/toutaio/toutago-scela-bus)
- [Basic Example](../examples/with-scela/)
- [Architecture Documentation](./architecture.md)
- [Toutā Main Documentation](../README.md)
