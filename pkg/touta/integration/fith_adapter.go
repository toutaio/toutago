package integration

import (
	"net/http"

	"github.com/toutaio/toutago-fith-renderer"
	"github.com/toutaio/toutago-fith-renderer/runtime"
	"github.com/toutaio/toutago/pkg/touta"
)

// FithRendererAdapter adapts fith.Engine to implement touta.TemplateRenderer interface.
type FithRendererAdapter struct {
	engine *fith.Engine
}

// NewTemplateRenderer creates a new template renderer using fith as the underlying implementation.
func NewTemplateRenderer(config *fith.Config) (touta.TemplateRenderer, error) {
	engine, err := fith.New(config)
	if err != nil {
		return nil, err
	}
	
	return &FithRendererAdapter{
		engine: engine,
	}, nil
}

// NewTemplateRendererWithDefaults creates a template renderer with default fith configuration.
func NewTemplateRendererWithDefaults() (touta.TemplateRenderer, error) {
	engine, err := fith.NewWithDefaults()
	if err != nil {
		return nil, err
	}
	
	return &FithRendererAdapter{
		engine: engine,
	}, nil
}

// Render executes a template with the given data.
// Note: Fith returns strings, so we convert to []byte.
func (a *FithRendererAdapter) Render(name string, data interface{}) ([]byte, error) {
	result, err := a.engine.Render(name, data)
	if err != nil {
		return nil, err
	}
	return []byte(result), nil
}

// RegisterFunction adds a custom template function.
func (a *FithRendererAdapter) RegisterFunction(name string, fn interface{}) {
	// Fith expects runtime.Function type
	if runtimeFn, ok := fn.(runtime.Function); ok {
		a.engine.RegisterFunction(name, runtimeFn)
	}
	// Note: For non-runtime.Function types, users should use Native() for direct access
}

// Parse loads templates from a pattern (e.g., "templates/*.html").
// Note: Fith uses configuration-based template loading.
// This method is a no-op as templates are loaded via the Config at initialization.
func (a *FithRendererAdapter) Parse(pattern string) error {
	// Fith loads templates based on Config, not dynamically via patterns
	// Users should configure TemplateDir in fith.Config
	return nil
}

// Execute renders a template to a writer.
func (a *FithRendererAdapter) Execute(name string, data interface{}, w http.ResponseWriter) error {
	output, err := a.engine.Render(name, data)
	if err != nil {
		return err
	}
	
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, err = w.Write([]byte(output))
	return err
}

// Native returns the underlying fith engine for advanced usage.
func (a *FithRendererAdapter) Native() *fith.Engine {
	return a.engine
}
