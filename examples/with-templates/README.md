# Templates with Scéla Example

This example demonstrates using Scéla message bus for event-driven template rendering and analytics.

## Features Demonstrated

- **Page view tracking**: Automatically track all page views
- **Render event publishing**: Publish events when templates are rendered
- **Performance monitoring**: Detect and alert on slow renders
- **Caching opportunities**: Identify templates that should be cached
- **Real-time analytics**: Track page views in real-time

## Running the Example

```bash
cd examples/with-templates
go run main.go
```

Then visit:
- `http://localhost:8080/` - Home page with template
- `http://localhost:8080/stats` - View analytics

## Event-Driven Patterns

This example shows how to use events for:

### 1. Analytics
Every page view publishes a `page.view` event that can be consumed by:
- Analytics services
- Real-time dashboards
- A/B testing frameworks

### 2. Caching
Every template render publishes a `template.render` event that can trigger:
- Cache warming
- Cache invalidation
- Cache hit/miss tracking

### 3. Monitoring
Slow renders publish `template.render.slow` events for:
- Performance alerts
- Resource optimization
- Bottleneck identification

## Integration with Fíth

In a real application, you would:
1. Use the actual Fíth template renderer
2. Publish events for all render operations
3. Subscribe to events for analytics and monitoring
4. Use middleware for cross-cutting concerns

## Benefits

This event-driven approach provides:
- **Decoupling**: Analytics code separate from rendering logic
- **Extensibility**: Easy to add new event handlers
- **Observability**: Full visibility into template performance
- **Flexibility**: Different handlers can react differently to the same events

See the [Message Bus documentation](../../docs/message-bus.md) for more details.
