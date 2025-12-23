package registry

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/toutaio/toutago/pkg/touta"
)

func TestComponentRegistry_Register(t *testing.T) {
	registry := NewComponentRegistry()

	component := &touta.Component{
		Name:    "test-component",
		Version: "1.0.0",
		Type:    "package",
	}

	err := registry.Register(component)
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	if !registry.Has("test-component") {
		t.Error("Component should be registered")
	}
}

func TestComponentRegistry_RegisterWithoutName(t *testing.T) {
	registry := NewComponentRegistry()

	component := &touta.Component{
		Version: "1.0.0",
	}

	err := registry.Register(component)
	if err == nil {
		t.Error("Should fail when component has no name")
	}
}

func TestComponentRegistry_Get(t *testing.T) {
	registry := NewComponentRegistry()

	component := &touta.Component{
		Name:    "test-component",
		Version: "1.0.0",
	}

	registry.Register(component)

	retrieved, err := registry.Get("test-component")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if retrieved.Name != "test-component" {
		t.Errorf("Expected test-component, got %s", retrieved.Name)
	}
}

func TestComponentRegistry_GetNonexistent(t *testing.T) {
	registry := NewComponentRegistry()

	_, err := registry.Get("nonexistent")
	if err == nil {
		t.Error("Should fail for nonexistent component")
	}
}

func TestComponentRegistry_List(t *testing.T) {
	registry := NewComponentRegistry()

	component1 := &touta.Component{Name: "comp1", Version: "1.0.0"}
	component2 := &touta.Component{Name: "comp2", Version: "2.0.0"}

	registry.Register(component1)
	registry.Register(component2)

	components, err := registry.List()
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(components) != 2 {
		t.Errorf("Expected 2 components, got %d", len(components))
	}
}

func TestComponentRegistry_LoadFromManifest(t *testing.T) {
	tmpDir := t.TempDir()
	manifestPath := filepath.Join(tmpDir, "package.yaml")

	manifestContent := `name: test-package
version: 1.0.0
type: package
components:
  - name: component1
    handlers:
      - handler1
    templates:
      - template1
  - name: component2
    handlers:
      - handler2
`

	if err := os.WriteFile(manifestPath, []byte(manifestContent), 0644); err != nil {
		t.Fatalf("Failed to create manifest: %v", err)
	}

	registry := NewComponentRegistry()
	err := registry.LoadFromManifest(manifestPath)

	if err != nil {
		t.Fatalf("LoadFromManifest failed: %v", err)
	}

	if !registry.Has("component1") {
		t.Error("component1 should be registered")
	}

	if !registry.Has("component2") {
		t.Error("component2 should be registered")
	}

	// Verify component details
	comp1, _ := registry.Get("component1")
	if len(comp1.Handlers) != 1 || comp1.Handlers[0] != "handler1" {
		t.Error("Component handlers not loaded correctly")
	}
}

func TestComponentRegistry_LoadFromInvalidManifest(t *testing.T) {
	tmpDir := t.TempDir()
	manifestPath := filepath.Join(tmpDir, "invalid.yaml")

	invalidContent := `invalid: yaml: content:`
	os.WriteFile(manifestPath, []byte(invalidContent), 0644)

	registry := NewComponentRegistry()
	err := registry.LoadFromManifest(manifestPath)

	if err == nil {
		t.Error("Should fail for invalid manifest")
	}
}

func TestComponentRegistry_LoadFromNonexistentFile(t *testing.T) {
	registry := NewComponentRegistry()
	err := registry.LoadFromManifest("/nonexistent/path.yaml")

	if err == nil {
		t.Error("Should fail for nonexistent file")
	}
}

func TestComponentRegistry_Has(t *testing.T) {
	registry := NewComponentRegistry()

	component := &touta.Component{Name: "test", Version: "1.0.0"}
	registry.Register(component)

	if !registry.Has("test") {
		t.Error("Should return true for registered component")
	}

	if registry.Has("nonexistent") {
		t.Error("Should return false for nonexistent component")
	}
}
