package registry

import (
	"fmt"
	"os"
	"sync"

	"github.com/toutaio/toutago/pkg/touta"
	"gopkg.in/yaml.v3"
)

// componentRegistry implements ComponentRegistry.
type componentRegistry struct {
	components map[string]*touta.Component
	mu         sync.RWMutex
}

// NewComponentRegistry creates a new component registry.
func NewComponentRegistry() touta.ComponentRegistry {
	return &componentRegistry{
		components: make(map[string]*touta.Component),
	}
}

// Register adds a component to the registry.
func (r *componentRegistry) Register(component *touta.Component) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if component.Name == "" {
		return fmt.Errorf("component name is required")
	}

	r.components[component.Name] = component
	return nil
}

// Get retrieves a component by name.
func (r *componentRegistry) Get(name string) (*touta.Component, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	component, exists := r.components[name]
	if !exists {
		return nil, fmt.Errorf("component %s not found", name)
	}

	return component, nil
}

// List returns all registered components.
func (r *componentRegistry) List() ([]*touta.Component, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	components := make([]*touta.Component, 0, len(r.components))
	for _, c := range r.components {
		components = append(components, c)
	}

	return components, nil
}

// LoadFromManifest parses a package.yaml manifest and registers components.
func (r *componentRegistry) LoadFromManifest(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read manifest: %w", err)
	}

	var manifest struct {
		Name       string                 `yaml:"name"`
		Version    string                 `yaml:"version"`
		Type       string                 `yaml:"type"`
		Components []*touta.Component     `yaml:"components"`
		Metadata   map[string]interface{} `yaml:"metadata"`
	}

	if err := yaml.Unmarshal(data, &manifest); err != nil {
		return fmt.Errorf("failed to parse manifest: %w", err)
	}

	// Register each component
	for _, component := range manifest.Components {
		if component.Name == "" {
			component.Name = manifest.Name
		}
		if component.Version == "" {
			component.Version = manifest.Version
		}
		if component.Type == "" {
			component.Type = manifest.Type
		}
		if component.Path == "" {
			component.Path = path
		}

		if err := r.Register(component); err != nil {
			return fmt.Errorf("failed to register component %s: %w", component.Name, err)
		}
	}

	return nil
}

// Has checks if a component is registered.
func (r *componentRegistry) Has(name string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.components[name]
	return exists
}
