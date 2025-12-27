package templates

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

//go:embed project
var projectTemplates embed.FS

// TemplateLoader handles loading and writing project templates
type TemplateLoader struct {
	fs embed.FS
}

// NewProjectTemplateLoader creates a loader for project templates
func NewProjectTemplateLoader() *TemplateLoader {
	return &TemplateLoader{
		fs: projectTemplates,
	}
}

// WriteTemplate writes a template file to the destination path
func (tl *TemplateLoader) WriteTemplate(templatePath, destPath string) error {
	content, err := tl.fs.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("failed to read template %s: %w", templatePath, err)
	}

	// Ensure parent directory exists
	dir := filepath.Dir(destPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// Write the file
	if err := os.WriteFile(destPath, content, 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", destPath, err)
	}

	return nil
}

// ReadTemplate reads a template file and returns its content
func (tl *TemplateLoader) ReadTemplate(templatePath string) ([]byte, error) {
	content, err := tl.fs.ReadFile(templatePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read template %s: %w", templatePath, err)
	}
	return content, nil
}

// ListTemplates returns all available template paths
func (tl *TemplateLoader) ListTemplates() ([]string, error) {
	var templates []string
	err := fs.WalkDir(tl.fs, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && filepath.Ext(path) == ".tmpl" {
			templates = append(templates, path)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list templates: %w", err)
	}
	return templates, nil
}

// TemplateExists checks if a template file exists
func (tl *TemplateLoader) TemplateExists(templatePath string) bool {
	_, err := tl.fs.ReadFile(templatePath)
	return err == nil
}
