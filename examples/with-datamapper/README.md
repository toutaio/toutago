# DataMapper with Scéla Example

This example demonstrates using the Scéla message bus for event-driven data operations.

## Features Demonstrated

- **Event-driven architecture**: Publish events before and after data operations
- **Multiple subscribers**: Different handlers react to the same events
- **Audit logging**: Track all data changes
- **Email notifications**: Send emails on user creation
- **Cache invalidation**: Invalidate cache on updates
- **Analytics tracking**: Track user signups
- **Search indexing**: Update search index on data changes

## Running the Example

```bash
cd examples/with-datamapper
go run main.go
```

## Output

You'll see events being published and handled for:
- User creation (with welcome email, analytics, etc.)
- User updates (with cache invalidation)
- User deletion (with audit logs)

## Event-Driven Pattern

This pattern is especially useful for:
- Decoupling business logic
- Adding new features without modifying existing code
- Background job processing
- Real-time notifications
- Data synchronization across services

## Integration with DataMapper

In a real application, you would:
1. Use the actual DataMapper component to interact with the database
2. Publish events in repository methods
3. Subscribe to events for side effects (emails, cache, etc.)
4. Use middleware for cross-cutting concerns (logging, validation)

See the [Message Bus documentation](../../docs/message-bus.md) for more details.
