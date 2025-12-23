package message

import (
	"context"
	"fmt"
	"sync"

	"github.com/toutaio/toutago/pkg/touta"
)

// BaseMessage provides a default implementation of the Message interface.
type BaseMessage struct {
	MessageSlug string                 `yaml:"slug" json:"slug"`
	MessageType string                 `yaml:"type" json:"type"`
	Meta        map[string]interface{} `yaml:"metadata,omitempty" json:"metadata,omitempty"`
}

// Slug returns the message slug.
func (m *BaseMessage) Slug() string {
	return m.MessageSlug
}

// Type returns the message type.
func (m *BaseMessage) Type() string {
	return m.MessageType
}

// Metadata returns the message metadata.
func (m *BaseMessage) Metadata() map[string]interface{} {
	if m.Meta == nil {
		m.Meta = make(map[string]interface{})
	}
	return m.Meta
}

// bus implements the MessageBus interface using channels.
type bus struct {
	subscribers map[string][]touta.MessageHandler
	messages    chan messageEnvelope
	wg          sync.WaitGroup
	mu          sync.RWMutex
	ctx         context.Context
	cancel      context.CancelFunc
	started     bool
}

// messageEnvelope wraps a message with its context.
type messageEnvelope struct {
	ctx  context.Context
	msg  touta.Message
	sync bool
	done chan error
}

// NewBus creates a new message bus.
func NewBus() touta.MessageBus {
	ctx, cancel := context.WithCancel(context.Background())
	return &bus{
		subscribers: make(map[string][]touta.MessageHandler),
		messages:    make(chan messageEnvelope, 100),
		ctx:         ctx,
		cancel:      cancel,
	}
}

// Publish sends a message asynchronously to all subscribers.
func (b *bus) Publish(ctx context.Context, msg touta.Message) error {
	if !b.started {
		return fmt.Errorf("message bus not started")
	}

	envelope := messageEnvelope{
		ctx:  ctx,
		msg:  msg,
		sync: false,
	}

	select {
	case b.messages <- envelope:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// PublishSync sends a message synchronously and waits for handlers to complete.
func (b *bus) PublishSync(ctx context.Context, msg touta.Message) error {
	if !b.started {
		return fmt.Errorf("message bus not started")
	}

	done := make(chan error, 1)
	envelope := messageEnvelope{
		ctx:  ctx,
		msg:  msg,
		sync: true,
		done: done,
	}

	select {
	case b.messages <- envelope:
		// Wait for processing to complete
		select {
		case err := <-done:
			return err
		case <-ctx.Done():
			return ctx.Err()
		}
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Subscribe registers a handler for messages matching a pattern.
func (b *bus) Subscribe(pattern string, handler touta.MessageHandler) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.subscribers[pattern] == nil {
		b.subscribers[pattern] = make([]touta.MessageHandler, 0)
	}
	b.subscribers[pattern] = append(b.subscribers[pattern], handler)
	return nil
}

// Unsubscribe removes a handler for a specific pattern.
func (b *bus) Unsubscribe(pattern string, handler touta.MessageHandler) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	handlers := b.subscribers[pattern]
	for i, h := range handlers {
		if fmt.Sprintf("%p", h) == fmt.Sprintf("%p", handler) {
			b.subscribers[pattern] = append(handlers[:i], handlers[i+1:]...)
			break
		}
	}
	return nil
}

// Start begins processing messages.
func (b *bus) Start(ctx context.Context) error {
	if b.started {
		return fmt.Errorf("message bus already started")
	}

	b.started = true
	b.wg.Add(1)
	go b.process()
	return nil
}

// Stop gracefully shuts down the message bus.
func (b *bus) Stop(ctx context.Context) error {
	if !b.started {
		return nil
	}

	b.cancel()
	close(b.messages)

	// Wait for processing to complete with timeout
	done := make(chan struct{})
	go func() {
		b.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// process is the main message processing loop.
func (b *bus) process() {
	defer b.wg.Done()

	for envelope := range b.messages {
		handlers := b.getHandlers(envelope.msg)

		if envelope.sync {
			// Synchronous processing
			var errs []error
			for _, handler := range handlers {
				if _, err := handler.Handle(envelope.ctx, envelope.msg); err != nil {
					errs = append(errs, err)
				}
			}

			if len(errs) > 0 {
				envelope.done <- errs[0] // Return first error
			} else {
				envelope.done <- nil
			}
		} else {
			// Asynchronous processing
			for _, handler := range handlers {
				h := handler // Capture for goroutine
				b.wg.Add(1)
				go func() {
					defer b.wg.Done()
					h.Handle(envelope.ctx, envelope.msg)
				}()
			}
		}
	}
}

// getHandlers returns all handlers matching the message.
func (b *bus) getHandlers(msg touta.Message) []touta.MessageHandler {
	b.mu.RLock()
	defer b.mu.RUnlock()

	var handlers []touta.MessageHandler

	// Match by exact slug
	if slugHandlers, ok := b.subscribers[msg.Slug()]; ok {
		handlers = append(handlers, slugHandlers...)
	}

	// Match by type
	if typeHandlers, ok := b.subscribers[msg.Type()]; ok {
		handlers = append(handlers, typeHandlers...)
	}

	// Match by wildcard
	if wildcardHandlers, ok := b.subscribers["*"]; ok {
		handlers = append(handlers, wildcardHandlers...)
	}

	return handlers
}

// HandlerFunc is a function adapter for MessageHandler.
type HandlerFunc func(context.Context, touta.Message) (touta.Message, error)

// Handle implements MessageHandler.
func (f HandlerFunc) Handle(ctx context.Context, msg touta.Message) (touta.Message, error) {
	return f(ctx, msg)
}
