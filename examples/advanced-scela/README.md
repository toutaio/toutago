# Advanced Scéla Example

This example demonstrates advanced Scéla features for building production-ready event-driven systems.

## Features Demonstrated

### 1. Event-Driven Order Processing
A complete e-commerce order workflow using events:
- Order creation
- Inventory reservation
- Payment processing
- Shipping creation
- Email notifications

### 2. Middleware Stack
Multiple middleware layers working together:
- **Logging Middleware**: Track all message processing with timing
- **Validation Middleware**: Ensure message payload is valid
- **Retry Middleware**: Automatic retry with exponential backoff

### 3. Error Handling
Comprehensive error handling:
- Automatic retries for transient failures
- Dead Letter Queue (DLQ) for permanently failed messages
- Failure notifications to customers

### 4. Message Persistence
Messages are persisted to ensure delivery:
- In-memory store for development
- Can be swapped for Redis, PostgreSQL, etc. in production

### 5. Priority Processing
High-priority orders are processed first:
- Normal priority for regular orders
- High priority for VIP/urgent orders

### 6. Pattern Matching
Analytics subscriber uses wildcard patterns:
- Subscribe to `order.*` to track all order events
- Enables cross-cutting concerns

## Running the Example

```bash
cd examples/advanced-scela
go run main.go
```

## Output

You'll see a complete order processing workflow with:
- Order creation and validation
- Inventory reservation
- Payment processing (with retries on failure)
- Shipping creation
- Email confirmations
- Analytics tracking
- Middleware logging

## Real-World Usage

This pattern is used in production for:
- **E-commerce**: Order processing, inventory management
- **Financial systems**: Transaction processing, fraud detection
- **SaaS applications**: User onboarding, billing, notifications
- **IoT platforms**: Device events, telemetry processing

## Architecture

```
Order Created Event
    ↓
Inventory Reserved Event
    ↓
Payment Processing (with retries)
    ↓
Payment Completed Event
    ↓
Shipment Created + Email Sent
    ↓
Order Complete

(Failed messages → Dead Letter Queue → Ops Alert)
```

## Middleware Chain

```
Message → Logging → Validation → Retry → Handler
```

Each middleware can:
- Modify the message
- Stop the chain (validation failure)
- Retry on errors
- Add context/metadata

## Configuration Options

In production, you would configure:

### Persistence
```go
scela.WithPersistence(scela.NewRedisStore(redis.Client))
scela.WithPersistence(scela.NewPostgresStore(db))
```

### Performance
```go
scela.WithBufferSize(10000)
scela.WithWorkerCount(20)
```

### Retry Strategy
```go
middleware.RetryWithMaxAttempts(5)
middleware.RetryWithBackoff(500*time.Millisecond, 2.0)
middleware.RetryWithMaxBackoff(30*time.Second)
```

### Dead Letter Queue
```go
scela.WithDeadLetterQueue("failed-orders")
scela.WithDLQMaxRetries(3)
```

## Best Practices

1. **Use middleware for cross-cutting concerns**: logging, validation, metrics
2. **Implement idempotency**: handlers should be safe to retry
3. **Monitor the DLQ**: failed messages need attention
4. **Use persistence in production**: ensure message delivery
5. **Set appropriate timeouts**: prevent resource exhaustion
6. **Implement circuit breakers**: for failing external services
7. **Use pattern matching wisely**: for analytics and monitoring

## See Also

- [Basic Scéla Example](../with-scela/README.md)
- [Message Bus Documentation](../../docs/message-bus.md)
- [Scéla Repository](https://github.com/toutaio/toutago-scela-bus)
