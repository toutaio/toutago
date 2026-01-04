package integration

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/toutaio/toutago-scela-bus/pkg/scela"
	"github.com/toutaio/toutago/pkg/touta"
)

func TestScelaBus_PublishSubscribe(t *testing.T) {
	bus := NewScelaBus()
	defer bus.Close()

	var received int32
	handler := touta.HandlerFunc(func(ctx context.Context, msg touta.Message) error {
		atomic.AddInt32(&received, 1)
		if msg.Topic() != "test.topic" {
			t.Errorf("Expected topic 'test.topic', got '%s'", msg.Topic())
		}
		return nil
	})

	_, err := bus.Subscribe("test.topic", handler)
	if err != nil {
		t.Fatalf("Subscribe() error = %v", err)
	}

	ctx := context.Background()
	if err := bus.PublishSync(ctx, "test.topic", "test"); err != nil {
		t.Fatalf("PublishSync() error = %v", err)
	}

	if got := atomic.LoadInt32(&received); got != 1 {
		t.Errorf("Expected 1 message received, got %d", got)
	}
}

func TestScelaBus_Async(t *testing.T) {
	bus := NewScelaBus()
	defer bus.Close()

	var received int32
	done := make(chan struct{})

	handler := touta.HandlerFunc(func(ctx context.Context, msg touta.Message) error {
		atomic.AddInt32(&received, 1)
		close(done)
		return nil
	})

	_, err := bus.Subscribe("test.async", handler)
	if err != nil {
		t.Fatalf("Subscribe() error = %v", err)
	}

	ctx := context.Background()
	if err := bus.Publish(ctx, "test.async", "payload"); err != nil {
		t.Fatalf("Publish() error = %v", err)
	}

	select {
	case <-done:
		// Success
	case <-time.After(2 * time.Second):
		t.Fatal("Timeout waiting for async message")
	}

	if got := atomic.LoadInt32(&received); got != 1 {
		t.Errorf("Expected 1 message received, got %d", got)
	}
}

func TestScelaBus_ToutaMiddleware(t *testing.T) {
	bus := NewScelaBus()
	defer bus.Close()

	var middlewareCalled bool
	middleware := func(next touta.Handler) touta.Handler {
		return touta.HandlerFunc(func(ctx context.Context, msg touta.Message) error {
			middlewareCalled = true
			return next.Handle(ctx, msg)
		})
	}

	bus.Use(middleware)

	var handlerCalled bool
	handler := touta.HandlerFunc(func(ctx context.Context, msg touta.Message) error {
		handlerCalled = true
		return nil
	})

	_, err := bus.Subscribe("test.middleware", handler)
	if err != nil {
		t.Fatalf("Subscribe() error = %v", err)
	}

	ctx := context.Background()
	if err := bus.PublishSync(ctx, "test.middleware", nil); err != nil {
		t.Fatalf("PublishSync() error = %v", err)
	}

	if !middlewareCalled {
		t.Error("Middleware was not called")
	}

	if !handlerCalled {
		t.Error("Handler was not called")
	}
}

func TestScelaBus_ScelaMiddleware(t *testing.T) {
	bus := NewScelaBus()
	defer bus.Close()

	var middlewareCalled bool
	middleware := func(next scela.Handler) scela.Handler {
		return scela.HandlerFunc(func(ctx context.Context, msg scela.Message) error {
			middlewareCalled = true
			return next.Handle(ctx, msg)
		})
	}

	bus.UseScela(middleware)

	var handlerCalled bool
	handler := scela.HandlerFunc(func(ctx context.Context, msg scela.Message) error {
		handlerCalled = true
		return nil
	})

	_, err := bus.SubscribeScela("test.scela.middleware", handler)
	if err != nil {
		t.Fatalf("SubscribeScela() error = %v", err)
	}

	ctx := context.Background()
	if err := bus.PublishSync(ctx, "test.scela.middleware", nil); err != nil {
		t.Fatalf("PublishSync() error = %v", err)
	}

	if !middlewareCalled {
		t.Error("Middleware was not called")
	}

	if !handlerCalled {
		t.Error("Handler was not called")
	}
}

func TestScelaBus_PatternMatching(t *testing.T) {
	bus := NewScelaBus()
	defer bus.Close()

	var received int32
	handler := touta.HandlerFunc(func(ctx context.Context, msg touta.Message) error {
		atomic.AddInt32(&received, 1)
		return nil
	})

	// Subscribe to wildcard pattern
	_, err := bus.Subscribe("user.*", handler)
	if err != nil {
		t.Fatalf("Subscribe() error = %v", err)
	}

	ctx := context.Background()
	
	// Publish to different topics that match
	bus.PublishSync(ctx, "user.created", nil)
	bus.PublishSync(ctx, "user.updated", nil)
	bus.PublishSync(ctx, "user.deleted", nil)

	if got := atomic.LoadInt32(&received); got != 3 {
		t.Errorf("Expected 3 messages received, got %d", got)
	}
}

func TestScelaBus_WithDefaults(t *testing.T) {
	bus := NewScelaWithDefaults()
	defer bus.Close()

	// Verify it works
	var received bool
	handler := touta.HandlerFunc(func(ctx context.Context, msg touta.Message) error {
		received = true
		return nil
	})

	_, err := bus.Subscribe("test", handler)
	if err != nil {
		t.Fatalf("Subscribe() error = %v", err)
	}

	ctx := context.Background()
	if err := bus.PublishSync(ctx, "test", nil); err != nil {
		t.Fatalf("PublishSync() error = %v", err)
	}

	if !received {
		t.Error("Handler was not called")
	}
}

func TestScelaBus_PublishWithPriority(t *testing.T) {
	bus := NewScelaBus()
	defer bus.Close()

	var received int32
	handler := touta.HandlerFunc(func(ctx context.Context, msg touta.Message) error {
		atomic.AddInt32(&received, 1)
		return nil
	})

	_, err := bus.Subscribe("urgent", handler)
	if err != nil {
		t.Fatalf("Subscribe() error = %v", err)
	}

	ctx := context.Background()
	if err := bus.PublishWithPriority(ctx, "urgent", "critical data", scela.PriorityUrgent); err != nil {
		t.Fatalf("PublishWithPriority() error = %v", err)
	}

	// Give it time to process
	time.Sleep(100 * time.Millisecond)

	if got := atomic.LoadInt32(&received); got != 1 {
		t.Errorf("Expected 1 message received, got %d", got)
	}
}

func TestScelaBus_MessageAdapter(t *testing.T) {
	bus := NewScelaBus()
	defer bus.Close()

	handler := touta.HandlerFunc(func(ctx context.Context, msg touta.Message) error {
		// Verify message interface methods
		if msg.Topic() == "" {
			t.Error("Topic() returned empty string")
		}
		if msg.Payload() == nil {
			t.Error("Payload() returned nil")
		}
		if msg.Metadata() == nil {
			t.Error("Metadata() returned nil")
		}
		if msg.ID() == "" {
			t.Error("ID() returned empty string")
		}
		if msg.Timestamp().IsZero() {
			t.Error("Timestamp() returned zero time")
		}
		return nil
	})

	_, err := bus.Subscribe("test.adapter", handler)
	if err != nil {
		t.Fatalf("Subscribe() error = %v", err)
	}

	ctx := context.Background()
	if err := bus.PublishSync(ctx, "test.adapter", "test payload"); err != nil {
		t.Fatalf("PublishSync() error = %v", err)
	}
}

func TestScelaBus_Unwrap(t *testing.T) {
	bus := NewScelaBus()
	defer bus.Close()

	underlying := bus.Unwrap()
	if underlying == nil {
		t.Error("Unwrap() returned nil")
	}
}

func TestScelaBus_WithOptions(t *testing.T) {
	bus := NewScelaBus(
		scela.WithWorkers(5),
		scela.WithMaxRetries(2),
	)
	defer bus.Close()

	// Just verify it doesn't panic
	ctx := context.Background()
	bus.Publish(ctx, "test", nil)
}

