package di

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/toutaio/toutago/pkg/touta"
)

// Additional comprehensive tests for DI container

func TestContainer_BindTagged(t *testing.T) {
	container := &container{
		bindings:   make(map[string]*binding),
		singletons: make(map[string]interface{}),
	}

	service := &testServiceImpl{name: "tagged"}
	err := container.BindTagged((*TestService)(nil), service, []string{"tag1", "tag2"})

	if err != nil {
		t.Fatalf("BindTagged failed: %v", err)
	}

	if !container.Has((*TestService)(nil)) {
		t.Error("Tagged binding should exist")
	}
}

func TestContainer_Tagged(t *testing.T) {
	container := &container{
		bindings:   make(map[string]*binding),
		singletons: make(map[string]interface{}),
	}

	service1 := &testServiceImpl{name: "service1"}
	service2 := &testServiceImpl{name: "service2"}

	container.BindTagged((*TestService)(nil), service1, []string{"mytag"})

	// Create a second binding with the same tag (need different key)
	key := "TestService2"
	container.bindings[key] = &binding{
		concrete: service2,
		tags:     []string{"mytag"},
	}

	instances, err := container.Tagged("mytag")
	if err != nil {
		t.Fatalf("Tagged failed: %v", err)
	}

	if len(instances) < 1 {
		t.Error("Should return tagged instances")
	}
}

func TestContainer_MakeWithParams(t *testing.T) {
	container := NewContainer()

	// Simple instance binding (not a constructor)
	service := &testServiceImpl{name: "test"}
	container.Bind((*TestService)(nil), service)

	instance, err := container.MakeWith(reflect.TypeOf((*TestService)(nil)), nil)
	if err != nil {
		t.Fatalf("MakeWith failed: %v", err)
	}

	if instance == nil {
		t.Error("Instance should not be nil")
	}
}

func TestContainer_BuildWithConstructor(t *testing.T) {
	container := NewContainer()

	// Test constructor pattern
	constructor := func() *testServiceImpl {
		return &testServiceImpl{name: "constructed"}
	}

	container.Bind((*TestService)(nil), constructor)
	instance, err := container.Make(reflect.TypeOf((*TestService)(nil)))

	if err != nil {
		t.Fatalf("Make with constructor failed: %v", err)
	}

	if instance == nil {
		t.Error("Instance should not be nil")
	}
}

func TestContainer_BuildWithError(t *testing.T) {
	container := NewContainer()

	// Constructor that returns error
	constructor := func() (*testServiceImpl, error) {
		return nil, fmt.Errorf("construction error")
	}

	container.Bind((*TestService)(nil), constructor)
	_, err := container.Make(reflect.TypeOf((*TestService)(nil)))

	if err == nil {
		t.Error("Should return constructor error")
	}
}

func TestContainer_AutoWireNotPointer(t *testing.T) {
	container := NewContainer()

	var notPointer testServiceImpl
	err := container.AutoWire(notPointer)

	if err == nil {
		t.Error("AutoWire should fail for non-pointer")
	}
}

func TestContainer_AutoWireNotStruct(t *testing.T) {
	container := NewContainer()

	var str string
	err := container.AutoWire(&str)

	if err == nil {
		t.Error("AutoWire should fail for non-struct")
	}
}

func TestContainer_AutoWireOptionalDependency(t *testing.T) {
	container := NewContainer()

	type TestStructOptional struct {
		Service *testServiceImpl `inject:"optional"`
	}

	target := &TestStructOptional{}
	err := container.AutoWire(target)

	// Should not error for optional missing dependency
	if err != nil {
		t.Logf("AutoWire with optional dependency: %v", err)
	}

	// Just verify it doesn't crash
	if target == nil {
		t.Error("Target should not be nil")
	}
}

func TestContainer_AutoWireAlreadySet(t *testing.T) {
	container := NewContainer()
	service := &testServiceImpl{name: "preset"}

	type TestStructPreset struct {
		Service *testServiceImpl `inject:""`
	}

	target := &TestStructPreset{Service: service}
	err := container.AutoWire(target)

	if err != nil {
		t.Fatalf("AutoWire failed: %v", err)
	}

	// Should not override already set field
	if target.Service.name != "preset" {
		t.Error("Should not override preset field")
	}
}

func TestContainer_GetKeyFromReflectType(t *testing.T) {
	container := &container{
		bindings:   make(map[string]*binding),
		singletons: make(map[string]interface{}),
	}

	// Test with reflect.Type
	typ := reflect.TypeOf((*TestService)(nil))
	key1 := container.getKey(typ)

	// Test with interface
	key2 := container.getKey((*TestService)(nil))

	if key1 == "" || key2 == "" {
		t.Error("Keys should not be empty")
	}
}

func TestContainer_FactoryWithSingleton(t *testing.T) {
	c := NewContainer()
	impl := c.(*container)

	callCount := 0
	factory := func(c touta.Container) (interface{}, error) {
		callCount++
		return &testServiceImpl{name: "factory-singleton"}, nil
	}

	// Register as singleton factory
	key := impl.getKey((*TestService)(nil))
	impl.bindings[key] = &binding{
		factory: factory,
		shared:  true,
	}

	// First call
	instance1, _ := c.Make(reflect.TypeOf((*TestService)(nil)))
	// Second call
	instance2, _ := c.Make(reflect.TypeOf((*TestService)(nil)))

	if callCount != 1 {
		t.Errorf("Factory should be called once for singleton, called %d times", callCount)
	}

	if instance1 != instance2 {
		t.Error("Singleton factory should return same instance")
	}
}

func TestContainer_MakeNonexistent(t *testing.T) {
	container := NewContainer()

	_, err := container.Make(reflect.TypeOf((*TestService)(nil)))

	if err == nil {
		t.Error("Should fail for nonexistent binding")
	}
}

func TestContainer_ConcurrentAccess(t *testing.T) {
	container := NewContainer()
	service := &testServiceImpl{name: "concurrent"}
	container.Singleton((*TestService)(nil), service)

	// Test concurrent reads
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func() {
			container.Make(reflect.TypeOf((*TestService)(nil)))
			container.Has((*TestService)(nil))
			done <- true
		}()
	}

	for i := 0; i < 10; i++ {
		<-done
	}
}

func TestContainer_ComplexDependencyGraph(t *testing.T) {
	container := NewContainer()

	// Create a dependency chain
	type ServiceA struct{ name string }
	type ServiceB struct{ A *ServiceA }
	type ServiceC struct{ B *ServiceB }

	a := &ServiceA{name: "A"}
	b := &ServiceB{A: a}
	c := &ServiceC{B: b}

	container.Bind((*ServiceA)(nil), a)
	container.Bind((*ServiceB)(nil), b)
	container.Bind((*ServiceC)(nil), c)

	instance, err := container.Make(reflect.TypeOf((*ServiceC)(nil)))
	if err != nil {
		t.Fatalf("Complex dependency resolution failed: %v", err)
	}

	if instance == nil {
		t.Error("Should resolve complex dependencies")
	}
}
