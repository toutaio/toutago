package integration_test

import (
	"testing"

	"github.com/toutaio/toutago/pkg/touta"
	"github.com/toutaio/toutago/pkg/touta/integration"
)

// Test interfaces and implementations
type Logger interface {
	Log(msg string)
}

type ConsoleLogger struct {
	messages []string
}

func (l *ConsoleLogger) Log(msg string) {
	l.messages = append(l.messages, msg)
}

type Database interface {
	Connect() error
}

type MySQL struct {
	connected bool
}

func (m *MySQL) Connect() error {
	m.connected = true
	return nil
}

type Service struct {
	Logger Logger `inject:""`
}

// TestNewContainer verifies that NewContainer creates a valid container.
func TestNewContainer(t *testing.T) {
	container := integration.NewContainer()
	
	if container == nil {
		t.Fatal("NewContainer returned nil")
	}
	
	// Verify it implements touta.Container interface
	var _ touta.Container = container
}

// TestContainerBind verifies dependency binding works.
func TestContainerBind(t *testing.T) {
	container := integration.NewContainer()
	
	err := container.Bind((*Logger)(nil), &ConsoleLogger{})
	if err != nil {
		t.Fatalf("Bind failed: %v", err)
	}
	
	instance, err := container.Make((*Logger)(nil))
	if err != nil {
		t.Fatalf("Make failed: %v", err)
	}
	
	if instance == nil {
		t.Fatal("Make returned nil instance")
	}
	
	if _, ok := instance.(*ConsoleLogger); !ok {
		t.Fatal("Instance is not a ConsoleLogger")
	}
}

// TestContainerSingleton verifies singleton binding works.
func TestContainerSingleton(t *testing.T) {
	container := integration.NewContainer()
	
	err := container.Singleton((*Database)(nil), &MySQL{})
	if err != nil {
		t.Fatalf("Singleton failed: %v", err)
	}
	
	instance1, err := container.Make((*Database)(nil))
	if err != nil {
		t.Fatalf("First Make failed: %v", err)
	}
	
	instance2, err := container.Make((*Database)(nil))
	if err != nil {
		t.Fatalf("Second Make failed: %v", err)
	}
	
	// For singletons, both should be the same instance
	if instance1 != instance2 {
		t.Fatal("Singleton instances are not the same")
	}
}

// TestContainerAutoWire verifies AutoWire functionality.
func TestContainerAutoWire(t *testing.T) {
	container := integration.NewContainer()
	
	err := container.Bind((*Logger)(nil), &ConsoleLogger{})
	if err != nil {
		t.Fatalf("Bind failed: %v", err)
	}
	
	service := &Service{}
	err = container.AutoWire(service)
	if err != nil {
		t.Fatalf("AutoWire failed: %v", err)
	}
	
	if service.Logger == nil {
		t.Fatal("Logger was not injected")
	}
}
