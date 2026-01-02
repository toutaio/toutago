// Package integration provides adapters for integrating external components with Toutā.
//
// The scela adapter wraps the Scéla message bus (toutago-scela-bus) for use in Toutā applications.
package integration

import (
	"context"

	"github.com/toutaio/toutago-scela-bus/pkg/scela"
)

// ScelaBus wraps scela.Bus for framework integration.
type ScelaBus struct {
	bus scela.Bus
}

// NewScelaBus creates a new Scéla message bus with the given options.
func NewScelaBus(opts ...scela.Option) *ScelaBus {
	return &ScelaBus{
		bus: scela.New(opts...),
	}
}

// Publish publishes a message asynchronously.
func (s *ScelaBus) Publish(ctx context.Context, topic string, payload interface{}) error {
	return s.bus.Publish(ctx, topic, payload)
}

// PublishSync publishes a message synchronously.
func (s *ScelaBus) PublishSync(ctx context.Context, topic string, payload interface{}) error {
	return s.bus.PublishSync(ctx, topic, payload)
}

// Subscribe subscribes a handler to a topic pattern.
func (s *ScelaBus) Subscribe(pattern string, handler scela.Handler) (scela.Subscription, error) {
	return s.bus.Subscribe(pattern, handler)
}

// Use adds middleware to the bus.
func (s *ScelaBus) Use(middleware ...scela.Middleware) {
	s.bus.Use(middleware...)
}

// Close gracefully shuts down the bus.
func (s *ScelaBus) Close() error {
	return s.bus.Close()
}

// Unwrap returns the underlying scela.Bus.
func (s *ScelaBus) Unwrap() scela.Bus {
	return s.bus
}
