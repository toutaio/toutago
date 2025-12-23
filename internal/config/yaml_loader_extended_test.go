package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/toutaio/toutago/pkg/touta"
)

func TestYAMLLoader_Watch(t *testing.T) {
	loader := NewYAMLLoader()

	called := false
	callback := func(cfg *touta.Config) {
		called = true
	}

	err := loader.Watch(callback)
	if err != nil {
		t.Fatalf("Watch failed: %v", err)
	}

	// Watch doesn't trigger callback immediately, just registers it
	if called {
		t.Error("Callback should not be called on Watch()")
	}
}

func TestYAMLLoader_LoadWithFrontmatter(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	content := `---
title: My App
version: 1.0.0
---
framework:
  mode: production
  debug: false

server:
  port: 3000
`

	os.WriteFile(configPath, []byte(content), 0644)

	loader := NewYAMLLoader()
	config, err := loader.Load(configPath)

	if err != nil {
		t.Fatalf("Load with frontmatter failed: %v", err)
	}

	if config.Framework.Mode != "production" {
		t.Error("Should parse YAML body after frontmatter")
	}

	if config.Server.Port != 3000 {
		t.Errorf("Expected port 3000, got %d", config.Server.Port)
	}
}

func TestYAMLLoader_EnvSubstitution(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	// Set environment variable
	os.Setenv("TEST_HOST", "testhost")
	defer os.Unsetenv("TEST_HOST")

	content := `framework:
  mode: production

server:
  host: ${TEST_HOST}
  port: 8080
`

	os.WriteFile(configPath, []byte(content), 0644)

	loader := NewYAMLLoader()
	config, err := loader.Load(configPath)

	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if config.Server.Host != "testhost" {
		t.Errorf("Expected 'testhost', got '%s'", config.Server.Host)
	}
}

func TestYAMLLoader_LoadNonexistentFile(t *testing.T) {
	loader := NewYAMLLoader()
	_, err := loader.Load("/nonexistent/path.yaml")

	if err == nil {
		t.Error("Should fail for nonexistent file")
	}
}

func TestYAMLLoader_LoadInvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "invalid.yaml")

	invalidContent := `this is: not: valid: yaml:`
	os.WriteFile(configPath, []byte(invalidContent), 0644)

	loader := NewYAMLLoader()
	_, err := loader.Load(configPath)

	if err == nil {
		t.Error("Should fail for invalid YAML")
	}
}

func TestFindConfig(t *testing.T) {
	tmpDir := t.TempDir()

	// Create config in temp dir
	configPath := filepath.Join(tmpDir, "touta.yaml")
	os.WriteFile(configPath, []byte("framework:\n  mode: test\n"), 0644)

	// Change to temp dir
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tmpDir)

	found, err := FindConfig()
	if err != nil {
		t.Fatalf("FindConfig failed: %v", err)
	}

	if filepath.Base(found) != "touta.yaml" {
		t.Errorf("Expected touta.yaml, got %s", filepath.Base(found))
	}
}

func TestFindConfigNotFound(t *testing.T) {
	tmpDir := t.TempDir()

	// Change to empty temp dir
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tmpDir)

	_, err := FindConfig()
	if err == nil {
		t.Error("Should fail when no config file found")
	}
}

func TestMergeConfig(t *testing.T) {
	dst := LoadDefaults()
	src := &touta.Config{
		Framework: touta.FrameworkConfig{
			Mode: "production",
		},
		Server: touta.ServerConfig{
			Port: 9000,
		},
	}

	mergeConfig(dst, src)

	if dst.Framework.Mode != "production" {
		t.Error("Should merge framework config")
	}

	if dst.Server.Port != 9000 {
		t.Error("Should merge server config")
	}
}

func TestLoadDefaults_AllFields(t *testing.T) {
	config := LoadDefaults()

	// Check all default values are set
	if config.Framework.Mode == "" {
		t.Error("Framework mode should have default")
	}

	if config.Server.Port == 0 {
		t.Error("Server port should have default")
	}

	if config.Router.BasePath == "" {
		t.Error("Router base path should have default")
	}

	if config.Packages == nil {
		t.Error("Packages map should be initialized")
	}

	if config.App == nil {
		t.Error("App map should be initialized")
	}
}

func TestYAMLLoader_ValidatePortRange(t *testing.T) {
	loader := NewYAMLLoader()

	tests := []struct {
		name    string
		port    int
		wantErr bool
	}{
		{"valid port", 8080, false},
		{"minimum port", 1, false},
		{"maximum port", 65535, false},
		{"invalid negative", -1, true},
		{"invalid too high", 70000, true},
		{"zero port", 0, false}, // 0 is valid (random port)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &touta.Config{
				Server: touta.ServerConfig{Port: tt.port},
			}
			err := loader.Validate(config)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
