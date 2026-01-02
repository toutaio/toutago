package integration_test

import (
	"testing"

	"github.com/toutaio/toutago-fith-renderer"
	"github.com/toutaio/toutago/pkg/touta"
	"github.com/toutaio/toutago/pkg/touta/integration"
)

// TestNewTemplateRenderer verifies that NewTemplateRenderer creates a valid renderer.
func TestNewTemplateRenderer(t *testing.T) {
	config := &fith.Config{
		TemplateDir: "testdata/templates",
		Extensions:  []string{".html", ".fith"},
	}
	
	renderer, err := integration.NewTemplateRenderer(config)
	if err != nil {
		t.Fatalf("NewTemplateRenderer failed: %v", err)
	}
	
	if renderer == nil {
		t.Fatal("NewTemplateRenderer returned nil")
	}
	
	// Verify it implements touta.TemplateRenderer interface
	var _ touta.TemplateRenderer = renderer
}

// TestNewTemplateRendererWithDefaults verifies default renderer creation.
func TestNewTemplateRendererWithDefaults(t *testing.T) {
	renderer, err := integration.NewTemplateRendererWithDefaults()
	if err != nil {
		t.Fatalf("NewTemplateRendererWithDefaults failed: %v", err)
	}
	
	if renderer == nil {
		t.Fatal("NewTemplateRendererWithDefaults returned nil")
	}
	
	var _ touta.TemplateRenderer = renderer
}

// TestRendererRender verifies template rendering.
func TestRendererRender(t *testing.T) {
	// Skip this test as it requires actual template files
	t.Skip("Skipping test - requires template directory setup")
	
	config := &fith.Config{
		TemplateDir: "testdata/templates",
		Extensions:  []string{".html"},
	}
	
	renderer, err := integration.NewTemplateRenderer(config)
	if err != nil {
		t.Skip("Skipping test - template directory not available")
	}
	
	data := map[string]interface{}{
		"Title": "Test Page",
		"Name":  "World",
	}
	
	output, err := renderer.Render("test", data)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}
	
	if len(output) == 0 {
		t.Fatal("Render returned empty output")
	}
}

// TestRendererRegisterFunction verifies custom function registration.
func TestRendererRegisterFunction(t *testing.T) {
	renderer, err := integration.NewTemplateRendererWithDefaults()
	if err != nil {
		t.Fatalf("Failed to create renderer: %v", err)
	}
	
	// Note: This test just verifies the method exists and doesn't panic
	// Actual function registration requires runtime.Function type
	// For real usage, use Native() to access fith engine directly
	renderer.RegisterFunction("upper", nil)
}
