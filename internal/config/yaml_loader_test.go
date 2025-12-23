package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/toutaio/toutago/pkg/touta"
)

func TestLoadDefaults(t *testing.T) {
	config := LoadDefaults()

	if config == nil {
		t.Fatal("LoadDefaults should return a config")
	}

	if config.Framework.Mode != "development" {
		t.Errorf("Expected development mode, got %s", config.Framework.Mode)
	}

	if config.Server.Port != 8080 {
		t.Errorf("Expected port 8080, got %d", config.Server.Port)
	}
}

func TestYAMLLoader_Load(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test.yaml")

	yamlContent := `framework:
  mode: production
  debug: false

server:
  host: 0.0.0.0
  port: 3000
`

	if err := os.WriteFile(configPath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	loader := NewYAMLLoader()
	config, err := loader.Load(configPath)

	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if config.Framework.Mode != "production" {
		t.Errorf("Expected production mode, got %s", config.Framework.Mode)
	}

	if config.Server.Port != 3000 {
		t.Errorf("Expected port 3000, got %d", config.Server.Port)
	}
}

func TestYAMLLoader_Validate(t *testing.T) {
	loader := NewYAMLLoader()

	tests := []struct {
		name    string
		config  *touta.Config
		wantErr bool
	}{
		{
			name:    "nil config",
			config:  nil,
			wantErr: true,
		},
		{
			name:    "valid config",
			config:  LoadDefaults(),
			wantErr: false,
		},
		{
			name: "invalid mode",
			config: &touta.Config{
				Framework: touta.FrameworkConfig{Mode: "invalid"},
			},
			wantErr: true,
		},
		{
			name: "invalid port",
			config: &touta.Config{
				Server: touta.ServerConfig{Port: -1},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := loader.Validate(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLoadOrDefault(t *testing.T) {
	tmpDir := t.TempDir()

	t.Run("nonexistent file returns defaults", func(t *testing.T) {
		config, err := LoadOrDefault(filepath.Join(tmpDir, "nonexistent.yaml"))
		if err != nil {
			t.Fatalf("LoadOrDefault failed: %v", err)
		}
		if config.Framework.Mode != "development" {
			t.Error("Should return defaults for nonexistent file")
		}
	})

	t.Run("existing file loads config", func(t *testing.T) {
		configPath := filepath.Join(tmpDir, "config.yaml")
		yamlContent := `framework:
  mode: production
`
		os.WriteFile(configPath, []byte(yamlContent), 0644)

		config, err := LoadOrDefault(configPath)
		if err != nil {
			t.Fatalf("LoadOrDefault failed: %v", err)
		}
		if config.Framework.Mode != "production" {
			t.Error("Should load config from file")
		}
	})
}
