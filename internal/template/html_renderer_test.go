package template

import (
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestHTMLRenderer_Parse(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Create test template
	templatePath := filepath.Join(tmpDir, "test.html")
	content := `<h1>{{.Title}}</h1>`
	if err := os.WriteFile(templatePath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test template: %v", err)
	}

	renderer := NewHTMLRenderer()
	err := renderer.Parse(filepath.Join(tmpDir, "*.html"))

	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
}

func TestHTMLRenderer_Render(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Create test template
	templatePath := filepath.Join(tmpDir, "test.html")
	content := `<h1>{{.Title}}</h1>`
	os.WriteFile(templatePath, []byte(content), 0644)

	renderer := NewHTMLRenderer()
	renderer.Parse(filepath.Join(tmpDir, "*.html"))

	data := map[string]string{"Title": "Hello World"}
	result, err := renderer.Render("test.html", data)

	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	expected := "<h1>Hello World</h1>"
	if string(result) != expected {
		t.Errorf("Expected %s, got %s", expected, string(result))
	}
}

func TestHTMLRenderer_RenderWithoutTemplates(t *testing.T) {
	renderer := NewHTMLRenderer()

	_, err := renderer.Render("nonexistent", nil)
	if err == nil {
		t.Error("Should fail when no templates loaded")
	}
}

func TestHTMLRenderer_RegisterFunction(t *testing.T) {
	tmpDir := t.TempDir()
	
	templatePath := filepath.Join(tmpDir, "test.html")
	content := `{{upper .Name}}`
	os.WriteFile(templatePath, []byte(content), 0644)

	renderer := NewHTMLRenderer()
	renderer.RegisterFunction("upper", func(s string) string {
		return "UPPERCASE"
	})
	renderer.Parse(filepath.Join(tmpDir, "*.html"))

	data := map[string]string{"Name": "test"}
	result, err := renderer.Render("test.html", data)

	if err != nil {
		t.Fatalf("Render with function failed: %v", err)
	}

	if string(result) != "UPPERCASE" {
		t.Errorf("Custom function not working, got: %s", string(result))
	}
}

func TestHTMLRenderer_Execute(t *testing.T) {
	tmpDir := t.TempDir()
	
	templatePath := filepath.Join(tmpDir, "test.html")
	content := `<h1>{{.Title}}</h1>`
	os.WriteFile(templatePath, []byte(content), 0644)

	renderer := NewHTMLRenderer()
	renderer.Parse(filepath.Join(tmpDir, "*.html"))

	w := httptest.NewRecorder()
	data := map[string]string{"Title": "Test"}
	
	err := renderer.Execute("test.html", data, w)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if w.Header().Get("Content-Type") != "text/html; charset=utf-8" {
		t.Error("Content-Type should be text/html")
	}

	if w.Body.String() != "<h1>Test</h1>" {
		t.Errorf("Unexpected output: %s", w.Body.String())
	}
}

func TestHTMLRenderer_ExecuteWithoutTemplates(t *testing.T) {
	renderer := NewHTMLRenderer()
	w := httptest.NewRecorder()

	err := renderer.Execute("nonexistent", nil, w)
	if err == nil {
		t.Error("Should fail when no templates loaded")
	}
}
