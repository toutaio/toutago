package integration

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/toutaio/toutago-scela-bus/pkg/scela"
)

func TestScelaBus_PublishSubscribe(t *testing.T) {
	bus := NewScelaBus()
	defer bus.Close()

	var received int32
	handler := scela.HandlerFunc(func(ctx context.Context, msg scela.Message) error {
		atomic.AddInt32(&received, 1)
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

	handler := scela.HandlerFunc(func(ctx context.Context, msg scela.Message) error {
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

func TestScelaBus_Middleware(t *testing.T) {
	bus := NewScelaBus()
	defer bus.Close()

	var middlewareCalled bool
	middleware := func(next scela.Handler) scela.Handler {
		return scela.HandlerFunc(func(ctx context.Context, msg scela.Message) error {
			middlewareCalled = true
			return next.Handle(ctx, msg)
		})
	}

	bus.Use(middleware)

	var handlerCalled bool
	handler := scela.HandlerFunc(func(ctx context.Context, msg scela.Message) error {
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
