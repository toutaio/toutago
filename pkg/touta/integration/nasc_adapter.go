package integration

import (
	"github.com/toutaio/toutago-nasc-dependency-injector"
	"github.com/toutaio/toutago/pkg/touta"
)

// NascContainerAdapter adapts nasc.Nasc to implement touta.Container interface.
type NascContainerAdapter struct {
	container *nasc.Nasc
}

// NewContainer creates a new DI container using nasc as the underlying implementation.
func NewContainer(options ...nasc.Option) touta.Container {
	return &NascContainerAdapter{
		container: nasc.New(options...),
	}
}

// Bind registers an interface to a concrete implementation.
func (a *NascContainerAdapter) Bind(abstract interface{}, concrete interface{}) error {
	return a.container.Bind(abstract, concrete)
}

// Singleton registers an interface to a singleton instance.
func (a *NascContainerAdapter) Singleton(abstract interface{}, concrete interface{}) error {
	return a.container.Singleton(abstract, concrete)
}

// Factory registers a factory function for creating instances.
func (a *NascContainerAdapter) Factory(abstract interface{}, factory func(touta.Container) (interface{}, error)) error {
	return a.container.Factory(abstract, func(c *nasc.Nasc) (interface{}, error) {
		return factory(a)
	})
}

// Make resolves and returns an instance of the given interface.
func (a *NascContainerAdapter) Make(abstract interface{}) (interface{}, error) {
	result, err := a.container.MakeSafe(abstract)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// MakeWith resolves an instance with additional parameters.
// Note: nasc doesn't have MakeWith, so we use MakeNamed if params contains "name".
func (a *NascContainerAdapter) MakeWith(abstract interface{}, params map[string]interface{}) (interface{}, error) {
	if name, ok := params["name"].(string); ok {
		result, err := a.container.MakeNamedSafe(abstract, name)
		if err != nil {
			return nil, err
		}
		return result, nil
	}
	return a.Make(abstract)
}

// Has checks if a binding exists for the given interface.
// Note: nasc doesn't have Has, so we check by trying to make it.
func (a *NascContainerAdapter) Has(abstract interface{}) bool {
	_, err := a.container.MakeSafe(abstract)
	return err == nil
}

// AutoWire injects dependencies into a struct using reflection.
func (a *NascContainerAdapter) AutoWire(target interface{}) error {
	return a.container.AutoWire(target)
}

// Tagged returns all instances registered with the given tag.
func (a *NascContainerAdapter) Tagged(tag string) ([]interface{}, error) {
	instances := a.container.MakeWithTag(tag)
	return instances, nil
}

// Native returns the underlying nasc container for advanced usage.
func (a *NascContainerAdapter) Native() *nasc.Nasc {
	return a.container
}
