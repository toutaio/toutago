package template

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"sync"

	"github.com/toutaio/toutago/pkg/touta"
)

// htmlRenderer implements TemplateRenderer using html/template.
type htmlRenderer struct {
	templates *template.Template
	funcs     template.FuncMap
	mu        sync.RWMutex
}

// NewHTMLRenderer creates a new HTML template renderer.
func NewHTMLRenderer() touta.TemplateRenderer {
	return &htmlRenderer{
		funcs: make(template.FuncMap),
	}
}

// Render executes a template with the given data.
func (r *htmlRenderer) Render(name string, data interface{}) ([]byte, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.templates == nil {
		return nil, fmt.Errorf("no templates loaded")
	}

	var buf bytes.Buffer
	if err := r.templates.ExecuteTemplate(&buf, name, data); err != nil {
		return nil, fmt.Errorf("failed to execute template %s: %w", name, err)
	}

	return buf.Bytes(), nil
}

// RegisterFunction adds a custom template function.
func (r *htmlRenderer) RegisterFunction(name string, fn interface{}) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.funcs[name] = fn
}

// Parse loads templates from a pattern.
func (r *htmlRenderer) Parse(pattern string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	tmpl := template.New(filepath.Base(pattern)).Funcs(r.funcs)

	var err error
	r.templates, err = tmpl.ParseGlob(pattern)
	if err != nil {
		return fmt.Errorf("failed to parse templates: %w", err)
	}

	return nil
}

// Execute renders a template to a writer.
func (r *htmlRenderer) Execute(name string, data interface{}, w http.ResponseWriter) error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.templates == nil {
		return fmt.Errorf("no templates loaded")
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err := r.templates.ExecuteTemplate(w, name, data); err != nil {
		return fmt.Errorf("failed to execute template %s: %w", name, err)
	}

	return nil
}
