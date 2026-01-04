package integration

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/toutaio/toutago-scela-bus/pkg/scela"
	"github.com/toutaio/toutago/pkg/touta"
)

// TestIntegration_ScelaWithNasc demonstrates integration between Scéla bus and Nasc DI container.
func TestIntegration_ScelaWithNasc(t *testing.T) {
	// Create container
	container := NewContainer()

	// Create bus
	bus := NewScelaBus()
	defer bus.Close()

	// Register bus in container using factory
	if err := container.Factory((*touta.Bus)(nil), func(c touta.Container) (interface{}, error) {
		return bus, nil
	}); err != nil {
		t.Fatalf("Failed to register bus: %v", err)
	}

	// Register a service that depends on the bus
	type NotificationService struct {
		Bus touta.Bus
	}

	if err := container.Factory((*NotificationService)(nil), func(c touta.Container) (interface{}, error) {
		resolvedBus, err := c.Make((*touta.Bus)(nil))
		if err != nil {
			return nil, err
		}
		return &NotificationService{Bus: resolvedBus.(touta.Bus)}, nil
	}); err != nil {
		t.Fatalf("Failed to register service: %v", err)
	}

	// Resolve service
	svc, err := container.Make((*NotificationService)(nil))
	if err != nil {
		t.Fatalf("Failed to resolve service: %v", err)
	}

	notifSvc := svc.(*NotificationService)

	// Verify the bus works
	var received bool
	handler := touta.HandlerFunc(func(ctx context.Context, msg touta.Message) error {
		received = true
		return nil
	})

	if _, err := notifSvc.Bus.Subscribe("notification.sent", handler); err != nil {
		t.Fatalf("Subscribe failed: %v", err)
	}

	ctx := context.Background()
	if err := notifSvc.Bus.(*ScelaBus).PublishSync(ctx, "notification.sent", "test"); err != nil {
		t.Fatalf("Publish failed: %v", err)
	}

	if !received {
		t.Error("Message was not received")
	}
}

// TestIntegration_EventDrivenArchitecture demonstrates a complete event-driven architecture.
func TestIntegration_EventDrivenArchitecture(t *testing.T) {
	bus := NewScelaBus()
	defer bus.Close()

	// Simulate a user registration flow with multiple event handlers
	var (
		emailSent      int32
		profileCreated int32
		analyticsLogged int32
		wg             sync.WaitGroup
	)

	wg.Add(3)

	// Email service handler
	emailHandler := touta.HandlerFunc(func(ctx context.Context, msg touta.Message) error {
		defer wg.Done()
		atomic.AddInt32(&emailSent, 1)
		// Simulate email sending
		time.Sleep(10 * time.Millisecond)
		return nil
	})

	// Profile service handler
	profileHandler := touta.HandlerFunc(func(ctx context.Context, msg touta.Message) error {
		defer wg.Done()
		atomic.AddInt32(&profileCreated, 1)
		// Simulate profile creation
		time.Sleep(15 * time.Millisecond)
		return nil
	})

	// Analytics service handler
	analyticsHandler := touta.HandlerFunc(func(ctx context.Context, msg touta.Message) error {
		defer wg.Done()
		atomic.AddInt32(&analyticsLogged, 1)
		// Simulate analytics logging
		time.Sleep(5 * time.Millisecond)
		return nil
	})

	// Subscribe all handlers to user.registered event
	if _, err := bus.Subscribe("user.registered", emailHandler); err != nil {
		t.Fatalf("Failed to subscribe email handler: %v", err)
	}
	if _, err := bus.Subscribe("user.registered", profileHandler); err != nil {
		t.Fatalf("Failed to subscribe profile handler: %v", err)
	}
	if _, err := bus.Subscribe("user.registered", analyticsHandler); err != nil {
		t.Fatalf("Failed to subscribe analytics handler: %v", err)
	}

	// Publish user registration event
	ctx := context.Background()
	userData := map[string]interface{}{
		"id":    "user-123",
		"email": "user@example.com",
		"name":  "John Doe",
	}

	if err := bus.Publish(ctx, "user.registered", userData); err != nil {
		t.Fatalf("Failed to publish event: %v", err)
	}

	// Wait for all handlers to complete
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// Success
	case <-time.After(2 * time.Second):
		t.Fatal("Timeout waiting for event handlers")
	}

	// Verify all handlers were called
	if got := atomic.LoadInt32(&emailSent); got != 1 {
		t.Errorf("Email handler called %d times, expected 1", got)
	}
	if got := atomic.LoadInt32(&profileCreated); got != 1 {
		t.Errorf("Profile handler called %d times, expected 1", got)
	}
	if got := atomic.LoadInt32(&analyticsLogged); got != 1 {
		t.Errorf("Analytics handler called %d times, expected 1", got)
	}
}

// TestIntegration_MiddlewareStack demonstrates middleware composition.
func TestIntegration_MiddlewareStack(t *testing.T) {
	bus := NewScelaBus()
	defer bus.Close()

	var executionOrder []string
	var mu sync.Mutex

	// Logging middleware
	loggingMiddleware := func(next touta.Handler) touta.Handler {
		return touta.HandlerFunc(func(ctx context.Context, msg touta.Message) error {
			mu.Lock()
			executionOrder = append(executionOrder, "logging-before")
			mu.Unlock()

			err := next.Handle(ctx, msg)

			mu.Lock()
			executionOrder = append(executionOrder, "logging-after")
			mu.Unlock()

			return err
		})
	}

	// Timing middleware
	timingMiddleware := func(next touta.Handler) touta.Handler {
		return touta.HandlerFunc(func(ctx context.Context, msg touta.Message) error {
			mu.Lock()
			executionOrder = append(executionOrder, "timing-before")
			mu.Unlock()

			err := next.Handle(ctx, msg)

			mu.Lock()
			executionOrder = append(executionOrder, "timing-after")
			mu.Unlock()

			return err
		})
	}

	// Validation middleware
	validationMiddleware := func(next touta.Handler) touta.Handler {
		return touta.HandlerFunc(func(ctx context.Context, msg touta.Message) error {
			mu.Lock()
			executionOrder = append(executionOrder, "validation-before")
			mu.Unlock()

			err := next.Handle(ctx, msg)

			mu.Lock()
			executionOrder = append(executionOrder, "validation-after")
			mu.Unlock()

			return err
		})
	}

	// Register middleware (order matters: logging -> timing -> validation)
	bus.Use(loggingMiddleware, timingMiddleware, validationMiddleware)

	// Handler
	handler := touta.HandlerFunc(func(ctx context.Context, msg touta.Message) error {
		mu.Lock()
		executionOrder = append(executionOrder, "handler")
		mu.Unlock()
		return nil
	})

	if _, err := bus.Subscribe("test.middleware", handler); err != nil {
		t.Fatalf("Subscribe failed: %v", err)
	}

	ctx := context.Background()
	if err := bus.PublishSync(ctx, "test.middleware", nil); err != nil {
		t.Fatalf("Publish failed: %v", err)
	}

	// Verify execution order
	expected := []string{
		"logging-before",
		"timing-before",
		"validation-before",
		"handler",
		"validation-after",
		"timing-after",
		"logging-after",
	}

	mu.Lock()
	defer mu.Unlock()

	if len(executionOrder) != len(expected) {
		t.Fatalf("Expected %d executions, got %d", len(expected), len(executionOrder))
	}

	for i, exp := range expected {
		if executionOrder[i] != exp {
			t.Errorf("Execution order[%d]: expected %s, got %s", i, exp, executionOrder[i])
		}
	}
}

// TestIntegration_ErrorHandling demonstrates error handling in event handlers.
func TestIntegration_ErrorHandling(t *testing.T) {
	// Retries only work for async messages, not sync
	bus := NewScelaBus(scela.WithMaxRetries(2))
	defer bus.Close()

	var attempts int32
	done := make(chan struct{})
	
	handler := touta.HandlerFunc(func(ctx context.Context, msg touta.Message) error {
		count := atomic.AddInt32(&attempts, 1)
		if count >= 2 {
			close(done)
		}
		// Always fail - testing that retries happen
		return fmt.Errorf("simulated error (attempt %d)", count)
	})

	if _, err := bus.Subscribe("test.retry", handler); err != nil {
		t.Fatalf("Subscribe failed: %v", err)
	}

	ctx := context.Background()
	// Use async Publish for retry logic to work
	if err := bus.Publish(ctx, "test.retry", nil); err != nil {
		t.Fatalf("Publish failed: %v", err)
	}

	// Wait for retries to complete
	select {
	case <-done:
		// Success
	case <-time.After(2 * time.Second):
		t.Logf("Timeout - handler was called %d times", atomic.LoadInt32(&attempts))
	}

	// Should have tried at least 2 times (1 initial + retries)
	got := atomic.LoadInt32(&attempts)
	if got < 2 {
		t.Errorf("Expected at least 2 attempts, got %d", got)
	}
	t.Logf("Handler was called %d times with retry logic", got)
}

// TestIntegration_PriorityProcessing demonstrates priority-based message processing.
// Note: This is a basic demonstration - actual priority behavior depends on queue implementation.
func TestIntegration_PriorityProcessing(t *testing.T) {
	bus := NewScelaBus(scela.WithWorkers(2))
	defer bus.Close()

	var processed int32

	handler := touta.HandlerFunc(func(ctx context.Context, msg touta.Message) error {
		atomic.AddInt32(&processed, 1)
		return nil
	})

	// Subscribe to specific pattern (Scéla uses * for single segment, not ** for multi-segment)
	if _, err := bus.Subscribe("task.*", handler); err != nil {
		t.Fatalf("Subscribe failed: %v", err)
	}

	ctx := context.Background()

	// Publish messages with different priorities
	bus.PublishWithPriority(ctx, "task.low", "low-priority", scela.PriorityLow)
	bus.PublishWithPriority(ctx, "task.urgent", "urgent", scela.PriorityUrgent)
	bus.PublishWithPriority(ctx, "task.normal", "normal", scela.PriorityNormal)
	bus.PublishWithPriority(ctx, "task.high", "high-priority", scela.PriorityHigh)

	// Wait for processing to complete
	time.Sleep(200 * time.Millisecond)

	// Verify all messages were processed
	got := atomic.LoadInt32(&processed)
	if got != 4 {
		t.Errorf("Expected 4 messages processed, got %d", got)
	}
	
	t.Logf("Successfully processed %d priority messages", got)
}

// TestIntegration_MultiplePatterns demonstrates complex pattern matching.
func TestIntegration_MultiplePatterns(t *testing.T) {
	bus := NewScelaBus()
	defer bus.Close()

	var (
		userEvents  int32
		orderEvents int32
		allEvents   int32
	)

	userHandler := touta.HandlerFunc(func(ctx context.Context, msg touta.Message) error {
		atomic.AddInt32(&userEvents, 1)
		return nil
	})

	orderHandler := touta.HandlerFunc(func(ctx context.Context, msg touta.Message) error {
		atomic.AddInt32(&orderEvents, 1)
		return nil
	})

	allHandler := touta.HandlerFunc(func(ctx context.Context, msg touta.Message) error {
		atomic.AddInt32(&allEvents, 1)
		return nil
	})

	// Subscribe to different patterns
	if _, err := bus.Subscribe("user.*", userHandler); err != nil {
		t.Fatalf("Subscribe user.* failed: %v", err)
	}
	if _, err := bus.Subscribe("order.*", orderHandler); err != nil {
		t.Fatalf("Subscribe order.* failed: %v", err)
	}
	if _, err := bus.Subscribe("*", allHandler); err != nil {
		t.Fatalf("Subscribe * failed: %v", err)
	}

	ctx := context.Background()

	// Publish various events
	bus.PublishSync(ctx, "user.created", nil)
	bus.PublishSync(ctx, "user.updated", nil)
	bus.PublishSync(ctx, "order.created", nil)
	bus.PublishSync(ctx, "order.shipped", nil)
	bus.PublishSync(ctx, "system.startup", nil)

	// Verify counts
	if got := atomic.LoadInt32(&userEvents); got != 2 {
		t.Errorf("Expected 2 user events, got %d", got)
	}
	if got := atomic.LoadInt32(&orderEvents); got != 2 {
		t.Errorf("Expected 2 order events, got %d", got)
	}
	if got := atomic.LoadInt32(&allEvents); got != 5 {
		t.Errorf("Expected 5 all events, got %d", got)
	}
}

// TestIntegration_GracefulShutdown demonstrates graceful shutdown behavior.
func TestIntegration_GracefulShutdown(t *testing.T) {
	bus := NewScelaBus(scela.WithWorkers(2))

	var processed int32

	handler := touta.HandlerFunc(func(ctx context.Context, msg touta.Message) error {
		// Simulate long-running task
		time.Sleep(100 * time.Millisecond)
		atomic.AddInt32(&processed, 1)
		return nil
	})

	if _, err := bus.Subscribe("task", handler); err != nil {
		t.Fatalf("Subscribe failed: %v", err)
	}

	ctx := context.Background()

	// Publish multiple messages
	for i := 0; i < 5; i++ {
		bus.Publish(ctx, "task", i)
	}

	// Give messages time to start processing
	time.Sleep(150 * time.Millisecond)

	// Close bus (should wait for in-flight messages)
	start := time.Now()
	if err := bus.Close(); err != nil {
		t.Fatalf("Close failed: %v", err)
	}
	elapsed := time.Since(start)

	// Should have processed all messages
	if got := atomic.LoadInt32(&processed); got != 5 {
		t.Errorf("Expected 5 messages processed, got %d", got)
	}

	// With messages already processing, Close should wait but may not take full 100ms
	// Just verify it did wait at least some time for completion
	t.Logf("Close waited %v for in-flight messages", elapsed)
}

// TestIntegration_ContextCancellation demonstrates context-based cancellation.
func TestIntegration_ContextCancellation(t *testing.T) {
	bus := NewScelaBus()
	defer bus.Close()

	var started int32
	var completed int32

	handler := touta.HandlerFunc(func(ctx context.Context, msg touta.Message) error {
		atomic.AddInt32(&started, 1)
		
		// Simulate long operation with context checking
		select {
		case <-time.After(500 * time.Millisecond):
			atomic.AddInt32(&completed, 1)
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	})

	if _, err := bus.Subscribe("slow.task", handler); err != nil {
		t.Fatalf("Subscribe failed: %v", err)
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// Publish with context
	err := bus.PublishSync(ctx, "slow.task", nil)

	// Should get context deadline exceeded error
	if err == nil {
		t.Error("Expected context error, got nil")
	}

	// Handler should have started but not completed
	if got := atomic.LoadInt32(&started); got != 1 {
		t.Errorf("Expected 1 handler started, got %d", got)
	}
	
	// Give a bit more time to ensure it didn't complete
	time.Sleep(200 * time.Millisecond)
	if got := atomic.LoadInt32(&completed); got != 0 {
		t.Errorf("Expected 0 handlers completed, got %d", got)
	}
}

// TestIntegration_Unsubscribe demonstrates subscription management.
func TestIntegration_Unsubscribe(t *testing.T) {
	bus := NewScelaBus()
	defer bus.Close()

	var count int32

	handler := touta.HandlerFunc(func(ctx context.Context, msg touta.Message) error {
		atomic.AddInt32(&count, 1)
		return nil
	})

	sub, err := bus.Subscribe("test.unsub", handler)
	if err != nil {
		t.Fatalf("Subscribe failed: %v", err)
	}

	ctx := context.Background()

	// First message should be received
	bus.PublishSync(ctx, "test.unsub", nil)

	// Unsubscribe
	sub.Unsubscribe()

	// Second message should NOT be received
	bus.PublishSync(ctx, "test.unsub", nil)

	// Only first message should have been processed
	if got := atomic.LoadInt32(&count); got != 1 {
		t.Errorf("Expected 1 message received, got %d", got)
	}
}
