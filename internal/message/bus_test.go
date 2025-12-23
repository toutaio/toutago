package message

import (
	"context"
	"testing"
	"time"

	"github.com/toutaio/toutago/pkg/touta"
)

type testMessage struct {
	BaseMessage
	payload string
}

type testHandler struct {
	received bool
	msg      touta.Message
}

func (h *testHandler) Handle(ctx context.Context, msg touta.Message) (touta.Message, error) {
	h.received = true
	h.msg = msg
	return nil, nil
}

func TestBus_PublishAndSubscribe(t *testing.T) {
	bus := NewBus()
	if err := bus.Start(context.Background()); err != nil {
		t.Fatalf("Failed to start bus: %v", err)
	}
	defer bus.Stop(context.Background())

	handler := &testHandler{}
	bus.Subscribe("test.message", handler)

	msg := &testMessage{
		BaseMessage: BaseMessage{
			MessageSlug: "test.message",
			MessageType: "event",
		},
		payload: "hello",
	}

	if err := bus.Publish(context.Background(), msg); err != nil {
		t.Fatalf("Publish failed: %v", err)
	}

	time.Sleep(50 * time.Millisecond)

	if !handler.received {
		t.Fatal("Handler should have received message")
	}
}

func TestBus_PublishSync(t *testing.T) {
	bus := NewBus()
	if err := bus.Start(context.Background()); err != nil {
		t.Fatalf("Failed to start bus: %v", err)
	}
	defer bus.Stop(context.Background())

	handler := &testHandler{}
	bus.Subscribe("test.sync", handler)

	msg := &testMessage{
		BaseMessage: BaseMessage{
			MessageSlug: "test.sync",
			MessageType: "command",
		},
	}

	if err := bus.PublishSync(context.Background(), msg); err != nil {
		t.Fatalf("PublishSync failed: %v", err)
	}

	if !handler.received {
		t.Fatal("Handler should have received message synchronously")
	}
}

func TestBus_Unsubscribe(t *testing.T) {
	bus := NewBus()
	if err := bus.Start(context.Background()); err != nil {
		t.Fatalf("Failed to start bus: %v", err)
	}
	defer bus.Stop(context.Background())

	handler := &testHandler{}
	bus.Subscribe("test.unsub", handler)
	bus.Unsubscribe("test.unsub", handler)

	msg := &testMessage{
		BaseMessage: BaseMessage{
			MessageSlug: "test.unsub",
			MessageType: "event",
		},
	}

	bus.Publish(context.Background(), msg)
	time.Sleep(50 * time.Millisecond)

	if handler.received {
		t.Fatal("Handler should not receive after unsubscribe")
	}
}

func TestBus_MultipleHandlers(t *testing.T) {
	bus := NewBus()
	if err := bus.Start(context.Background()); err != nil {
		t.Fatalf("Failed to start bus: %v", err)
	}
	defer bus.Stop(context.Background())

	handler1 := &testHandler{}
	handler2 := &testHandler{}

	bus.Subscribe("test.multi", handler1)
	bus.Subscribe("test.multi", handler2)

	msg := &testMessage{
		BaseMessage: BaseMessage{
			MessageSlug: "test.multi",
			MessageType: "event",
		},
	}

	bus.Publish(context.Background(), msg)
	time.Sleep(50 * time.Millisecond)

	if !handler1.received || !handler2.received {
		t.Fatal("Both handlers should receive message")
	}
}
