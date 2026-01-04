// Package integration provides adapters for integrating external components with Toutā.
//
// The Scéla adapter provides a simplified interface to the Scéla message bus
// (toutago-scela-bus) for use in Toutā applications.
package integration

import (
	"context"
	"time"

	"github.com/toutaio/toutago-scela-bus/pkg/scela"
	"github.com/toutaio/toutago/pkg/touta"
)

// ScelaBus wraps scela.Bus to provide a more integrated Toutā experience.
//
// This adapter exposes Scéla's full feature set while implementing Toutā's
// Bus interface. It can be used directly or accessed through the Container.
//
// Example:
//
//	bus := integration.NewScelaBus(
//	    scela.WithWorkers(20),
//	    scela.WithMaxRetries(3),
//	)
//	defer bus.Close()
//
//	// Use with Toutā interfaces
//	bus.Subscribe("user.*", touta.HandlerFunc(func(ctx context.Context, msg touta.Message) error {
//	    log.Printf("User event: %s", msg.Topic())
//	    return nil
//	}))
type ScelaBus struct {
	bus scela.Bus
}

// NewScelaBus creates a new Scéla message bus with the given options.
//
// Common options include:
//   - scela.WithWorkers(n) - Set number of worker goroutines
//   - scela.WithMaxRetries(n) - Set max retries for failed messages
//   - scela.WithDeadLetterHandler(h) - Set handler for failed messages
//
// Example:
//
//	bus := NewScelaBus(
//	    scela.WithWorkers(20),
//	    scela.WithMaxRetries(5),
//	)
func NewScelaBus(opts ...scela.Option) *ScelaBus {
	return &ScelaBus{
		bus: scela.New(opts...),
	}
}

// NewScelaWithDefaults creates a Scéla bus with sensible defaults for web apps.
//
// Default configuration:
//   - 20 worker goroutines
//   - 3 max retries
//   - 1000 buffer size
//
// This is a convenience constructor for typical use cases.
func NewScelaWithDefaults() *ScelaBus {
	return NewScelaBus(
		scela.WithWorkers(20),
		scela.WithMaxRetries(3),
	)
}

// NewScelaBusWithMiddleware creates a new Scéla bus with global middleware applied.
//
// The middleware will be applied to all handlers registered on the bus.
//
// Example:
//
//	bus := NewScelaBusWithMiddleware(
//	    LoggingMiddleware,
//	    ValidationMiddleware,
//	)
func NewScelaBusWithMiddleware(middleware ...touta.Middleware) *ScelaBus {
	bus := NewScelaBus()
	bus.Use(middleware...)
	return bus
}

// Publish publishes a message asynchronously to all matching subscribers.
//
// The message is queued and processed by worker goroutines. This method
// returns immediately without waiting for handlers to complete.
//
// Example:
//
//	err := bus.Publish(ctx, "user.registered", UserData{
//	    ID:    "123",
//	    Email: "user@example.com",
//	})
func (s *ScelaBus) Publish(ctx context.Context, topic string, payload interface{}) error {
	return s.bus.Publish(ctx, topic, payload)
}

// PublishSync publishes a message synchronously and waits for all handlers.
//
// This method blocks until all matching handlers have processed the message
// or an error occurs. Use this when you need immediate feedback or ordering
// guarantees.
//
// Example:
//
//	if err := bus.PublishSync(ctx, "order.created", orderData); err != nil {
//	    log.Printf("Order processing failed: %v", err)
//	}
func (s *ScelaBus) PublishSync(ctx context.Context, topic string, payload interface{}) error {
	return s.bus.PublishSync(ctx, topic, payload)
}

// PublishWithPriority publishes a message with a specific priority level.
//
// Higher priority messages are processed before lower priority ones.
// Priority levels: PriorityLow, PriorityNormal, PriorityHigh, PriorityUrgent
//
// Example:
//
//	err := bus.PublishWithPriority(ctx, "alert.critical", alertData, scela.PriorityUrgent)
func (s *ScelaBus) PublishWithPriority(ctx context.Context, topic string, payload interface{}, priority scela.Priority) error {
	return s.bus.PublishWithPriority(ctx, topic, payload, priority)
}

// Subscribe registers a handler for messages matching the topic pattern.
//
// Patterns support wildcards:
//   - "*" matches a single segment (e.g., "user.*" matches "user.created")
//   - "**" matches multiple segments (e.g., "user.**" matches "user.order.created")
//
// The returned Subscription can be used to unsubscribe later.
//
// Example:
//
//	sub, err := bus.Subscribe("user.*", touta.HandlerFunc(func(ctx context.Context, msg touta.Message) error {
//	    log.Printf("User event: %s", msg.Topic())
//	    return nil
//	}))
//	defer sub.Unsubscribe()
func (s *ScelaBus) Subscribe(pattern string, handler touta.Handler) (touta.Subscription, error) {
	// Wrap the touta.Handler as scela.Handler
	scelaHandler := scela.HandlerFunc(func(ctx context.Context, msg scela.Message) error {
		return handler.Handle(ctx, &messageAdapter{msg})
	})
	return s.bus.Subscribe(pattern, scelaHandler)
}

// SubscribeScela subscribes using a native Scéla handler directly.
//
// Use this when you need direct access to Scéla's Message interface
// without the adapter overhead.
//
// Example:
//
//	sub, err := bus.SubscribeScela("metrics.**", scela.HandlerFunc(func(ctx context.Context, msg scela.Message) error {
//	    // Direct access to scela.Message
//	    return nil
//	}))
func (s *ScelaBus) SubscribeScela(pattern string, handler scela.Handler) (scela.Subscription, error) {
	return s.bus.Subscribe(pattern, handler)
}

// Use adds middleware to the bus for cross-cutting concerns.
//
// Middleware is applied to all handlers in the order registered.
//
// Example:
//
//	bus.Use(
//	    LoggingMiddleware(),
//	    RetryMiddleware(3),
//	    MetricsMiddleware(),
//	)
func (s *ScelaBus) Use(middleware ...touta.Middleware) {
	// Convert touta.Middleware to scela.Middleware
	for _, mw := range middleware {
		scelaMW := func(next scela.Handler) scela.Handler {
			return scela.HandlerFunc(func(ctx context.Context, msg scela.Message) error {
				// Wrap scela.Handler as touta.Handler
				toutaNext := touta.HandlerFunc(func(ctx context.Context, m touta.Message) error {
					// Unwrap back to scela.Message
					if adapter, ok := m.(*messageAdapter); ok {
						return next.Handle(ctx, adapter.msg)
					}
					// Fallback - shouldn't happen
					return next.Handle(ctx, msg)
				})
				
				// Apply touta middleware
				wrapped := mw(toutaNext)
				
				// Execute with adapted message
				return wrapped.Handle(ctx, &messageAdapter{msg})
			})
		}
		s.bus.Use(scelaMW)
	}
}

// UseScela adds native Scéla middleware to the bus.
//
// Use this when you have middleware that needs direct access to
// Scéla's Message interface without adaptation overhead.
//
// Example:
//
//	bus.UseScela(
//	    scela.LoggingMiddleware(),
//	    scela.RetryMiddleware(3, 1*time.Second),
//	)
func (s *ScelaBus) UseScela(middleware ...scela.Middleware) {
	s.bus.Use(middleware...)
}

// Close gracefully shuts down the bus.
//
// This waits for all in-flight messages to complete processing
// before returning. Always call Close() when shutting down to
// ensure messages are not lost.
//
// Example:
//
//	defer bus.Close()
func (s *ScelaBus) Close() error {
	return s.bus.Close()
}

// Unwrap returns the underlying scela.Bus instance.
//
// Use this when you need direct access to Scéla-specific features
// not exposed through the adapter.
func (s *ScelaBus) Unwrap() scela.Bus {
	return s.bus
}

// messageAdapter adapts scela.Message to touta.Message interface.
type messageAdapter struct {
	msg scela.Message
}

func (m *messageAdapter) Topic() string                       { return m.msg.Topic() }
func (m *messageAdapter) Payload() interface{}                { return m.msg.Payload() }
func (m *messageAdapter) Metadata() map[string]interface{}    { return m.msg.Metadata() }
func (m *messageAdapter) ID() string                          { return m.msg.ID() }
func (m *messageAdapter) Timestamp() time.Time                { return m.msg.Timestamp() }

// Unwrap returns the underlying scela.Message.
func (m *messageAdapter) Unwrap() scela.Message {
	return m.msg
}

