package di

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/toutaio/toutago/pkg/touta"
)

// binding stores information about how to resolve a dependency.
type binding struct {
	concrete interface{}
	factory  func(touta.Container) (interface{}, error)
	shared   bool // singleton flag
	tags     []string
}

// container implements the Container interface.
type container struct {
	bindings   map[string]*binding
	singletons map[string]interface{}
	mu         sync.RWMutex
}

// NewContainer creates a new dependency injection container.
func NewContainer() touta.Container {
	return &container{
		bindings:   make(map[string]*binding),
		singletons: make(map[string]interface{}),
	}
}

// Bind registers an interface to a concrete implementation.
func (c *container) Bind(abstract interface{}, concrete interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	key := c.getKey(abstract)
	c.bindings[key] = &binding{
		concrete: concrete,
		shared:   false,
	}
	return nil
}

// Singleton registers an interface to a singleton instance.
func (c *container) Singleton(abstract interface{}, concrete interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	key := c.getKey(abstract)
	c.bindings[key] = &binding{
		concrete: concrete,
		shared:   true,
	}
	return nil
}

// Factory registers a factory function for creating instances.
func (c *container) Factory(abstract interface{}, factory func(touta.Container) (interface{}, error)) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	key := c.getKey(abstract)
	c.bindings[key] = &binding{
		factory: factory,
		shared:  false,
	}
	return nil
}

// Make resolves and returns an instance of the given interface.
func (c *container) Make(abstract interface{}) (interface{}, error) {
	return c.MakeWith(abstract, nil)
}

// MakeWith resolves an instance with additional parameters.
func (c *container) MakeWith(abstract interface{}, params map[string]interface{}) (interface{}, error) {
	c.mu.RLock()
	key := c.getKey(abstract)
	b, exists := c.bindings[key]
	c.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("no binding found for %s", key)
	}

	// Check if singleton already instantiated
	if b.shared {
		c.mu.RLock()
		if instance, ok := c.singletons[key]; ok {
			c.mu.RUnlock()
			return instance, nil
		}
		c.mu.RUnlock()
	}

	var instance interface{}
	var err error

	// Resolve using factory or direct instantiation
	if b.factory != nil {
		instance, err = b.factory(c)
	} else {
		instance, err = c.build(b.concrete, params)
	}

	if err != nil {
		return nil, err
	}

	// Store singleton
	if b.shared {
		c.mu.Lock()
		c.singletons[key] = instance
		c.mu.Unlock()
	}

	return instance, nil
}

// Has checks if a binding exists for the given interface.
func (c *container) Has(abstract interface{}) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	key := c.getKey(abstract)
	_, exists := c.bindings[key]
	return exists
}

// AutoWire injects dependencies into a struct using reflection.
func (c *container) AutoWire(target interface{}) error {
	val := reflect.ValueOf(target)
	if val.Kind() != reflect.Ptr {
		return fmt.Errorf("target must be a pointer")
	}

	elem := val.Elem()
	if elem.Kind() != reflect.Struct {
		return fmt.Errorf("target must be a pointer to a struct")
	}

	typ := elem.Type()
	for i := 0; i < elem.NumField(); i++ {
		field := elem.Field(i)
		fieldType := typ.Field(i)

		// Check for inject tag
		tag := fieldType.Tag.Get("inject")
		if tag == "" && !fieldType.Anonymous {
			continue
		}

		// Skip if already set
		if !field.IsZero() {
			continue
		}

		// Resolve dependency
		var abstractType interface{}
		if field.Kind() == reflect.Interface {
			abstractType = reflect.New(field.Type()).Elem().Interface()
		} else if field.Kind() == reflect.Ptr {
			abstractType = reflect.New(field.Type().Elem()).Interface()
		} else {
			continue // Can only inject interfaces or pointers
		}

		instance, err := c.Make(reflect.TypeOf(abstractType))
		if err != nil {
			if tag == "optional" {
				continue // Skip optional dependencies
			}
			return fmt.Errorf("failed to resolve %s: %w", fieldType.Name, err)
		}

		// Set the field
		if field.CanSet() {
			instanceVal := reflect.ValueOf(instance)
			if instanceVal.Type().AssignableTo(field.Type()) {
				field.Set(instanceVal)
			}
		}
	}

	return nil
}

// Tagged returns all instances registered with the given tag.
func (c *container) Tagged(tag string) ([]interface{}, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var instances []interface{}
	for key, b := range c.bindings {
		for _, t := range b.tags {
			if t == tag {
				// Check singleton cache first
				if b.shared {
					if instance, ok := c.singletons[key]; ok {
						instances = append(instances, instance)
						continue
					}
				}

				// Build instance
				var instance interface{}
				var err error
				if b.factory != nil {
					instance, err = b.factory(c)
				} else {
					instance, err = c.build(b.concrete, nil)
				}
				if err != nil {
					return nil, err
				}
				instances = append(instances, instance)
			}
		}
	}

	return instances, nil
}

// build creates a new instance using reflection.
func (c *container) build(concrete interface{}, params map[string]interface{}) (interface{}, error) {
	val := reflect.ValueOf(concrete)
	typ := reflect.TypeOf(concrete)

	// If it's already an instance, return it
	if typ.Kind() != reflect.Func {
		// If it's a pointer to a struct, try to auto-wire
		if typ.Kind() == reflect.Ptr && typ.Elem().Kind() == reflect.Struct {
			if err := c.AutoWire(concrete); err != nil {
				return nil, err
			}
		}
		return concrete, nil
	}

	// If it's a constructor function
	if typ.NumOut() != 1 && typ.NumOut() != 2 {
		return nil, fmt.Errorf("constructor must return 1 or 2 values")
	}

	// Build constructor arguments
	args := make([]reflect.Value, typ.NumIn())
	for i := 0; i < typ.NumIn(); i++ {
		argType := typ.In(i)

		// Check params first
		if params != nil {
			for key, value := range params {
				if key == argType.String() {
					args[i] = reflect.ValueOf(value)
					continue
				}
			}
		}

		// Resolve from container
		instance, err := c.Make(argType)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve constructor arg %d: %w", i, err)
		}
		args[i] = reflect.ValueOf(instance)
	}

	// Call constructor
	results := val.Call(args)

	// Check for error return
	if len(results) == 2 {
		if !results[1].IsNil() {
			return nil, results[1].Interface().(error)
		}
	}

	instance := results[0].Interface()

	// Auto-wire if it's a struct pointer
	if reflect.TypeOf(instance).Kind() == reflect.Ptr {
		if err := c.AutoWire(instance); err != nil {
			return nil, err
		}
	}

	return instance, nil
}

// getKey returns a unique key for an interface or type.
func (c *container) getKey(abstract interface{}) string {
	if t, ok := abstract.(reflect.Type); ok {
		return t.String()
	}
	return reflect.TypeOf(abstract).String()
}

// BindTagged registers a binding with tags.
func (c *container) BindTagged(abstract interface{}, concrete interface{}, tags []string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	key := c.getKey(abstract)
	c.bindings[key] = &binding{
		concrete: concrete,
		shared:   false,
		tags:     tags,
	}
	return nil
}
