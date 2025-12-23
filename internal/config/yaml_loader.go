package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/adrg/frontmatter"
	"github.com/toutaio/toutago/pkg/touta"
	"gopkg.in/yaml.v3"
)

// yamlLoader implements ConfigLoader using YAML with frontmatter.
type yamlLoader struct {
	watchers []func(*touta.Config)
}

// NewYAMLLoader creates a new YAML configuration loader.
func NewYAMLLoader() touta.ConfigLoader {
	return &yamlLoader{
		watchers: make([]func(*touta.Config), 0),
	}
}

// Load parses configuration from a file.
func (l *yamlLoader) Load(source string) (*touta.Config, error) {
	data, err := os.ReadFile(source)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	config := &touta.Config{}

	// Try to parse with frontmatter first
	var meta map[string]interface{}
	rest, err := frontmatter.Parse(strings.NewReader(string(data)), &meta)
	if err == nil && len(meta) > 0 {
		// Has frontmatter, parse the rest as YAML
		if err := yaml.Unmarshal(rest, config); err != nil {
			return nil, fmt.Errorf("failed to parse YAML body: %w", err)
		}
	} else {
		// No frontmatter, parse entire file as YAML
		if err := yaml.Unmarshal(data, config); err != nil {
			return nil, fmt.Errorf("failed to parse YAML: %w", err)
		}
	}

	// Apply environment variable substitution
	l.substituteEnv(config)

	return config, nil
}

// Watch monitors configuration for changes.
func (l *yamlLoader) Watch(callback func(*touta.Config)) error {
	l.watchers = append(l.watchers, callback)
	// TODO: Implement file watching with fsnotify in Phase 1 completion
	return nil
}

// Validate checks if the configuration is valid.
func (l *yamlLoader) Validate(config *touta.Config) error {
	if config == nil {
		return fmt.Errorf("config is nil")
	}

	// Validate framework settings
	if config.Framework.Mode != "" {
		if config.Framework.Mode != "development" && config.Framework.Mode != "production" {
			return fmt.Errorf("invalid framework mode: %s", config.Framework.Mode)
		}
	}

	// Validate server settings
	if config.Server.Port < 0 || config.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", config.Server.Port)
	}

	return nil
}

// substituteEnv replaces ${VAR} patterns with environment variables.
func (l *yamlLoader) substituteEnv(config *touta.Config) {
	// Substitute in framework config
	config.Framework.Mode = l.expandEnv(config.Framework.Mode)
	config.Framework.LogLevel = l.expandEnv(config.Framework.LogLevel)
	config.Framework.Timezone = l.expandEnv(config.Framework.Timezone)

	// Substitute in server config
	config.Server.Host = l.expandEnv(config.Server.Host)

	// Substitute in TLS config
	config.Server.TLS.CertFile = l.expandEnv(config.Server.TLS.CertFile)
	config.Server.TLS.KeyFile = l.expandEnv(config.Server.TLS.KeyFile)
}

// expandEnv expands environment variables in a string.
func (l *yamlLoader) expandEnv(s string) string {
	return os.ExpandEnv(s)
}

// LoadDefaults returns a configuration with sensible defaults.
func LoadDefaults() *touta.Config {
	return &touta.Config{
		Framework: touta.FrameworkConfig{
			Mode:      "development",
			Debug:     true,
			HotReload: true,
			LogLevel:  "info",
			Timezone:  "UTC",
		},
		Router: touta.RouterConfig{
			BasePath:   "/",
			Middleware: []string{},
			CORS: touta.CORSConfig{
				Enabled:          false,
				AllowedOrigins:   []string{"*"},
				AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
				AllowedHeaders:   []string{"Content-Type", "Authorization"},
				AllowCredentials: false,
				MaxAge:           3600,
			},
			RateLimit: touta.RateLimitConfig{
				Enabled:  false,
				Requests: 100,
				Window:   60,
			},
			Static: []touta.StaticConfig{},
		},
		Server: touta.ServerConfig{
			Host:           "localhost",
			Port:           8080,
			ReadTimeout:    15,
			WriteTimeout:   15,
			IdleTimeout:    60,
			MaxHeaderBytes: 1 << 20, // 1MB
			TLS: touta.TLSConfig{
				Enabled: false,
			},
		},
		Packages: make(map[string]interface{}),
		App:      make(map[string]interface{}),
	}
}

// LoadOrDefault loads configuration from a file or returns defaults.
func LoadOrDefault(path string) (*touta.Config, error) {
	loader := NewYAMLLoader()

	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return LoadDefaults(), nil
	}

	// Load from file
	config, err := loader.Load(path)
	if err != nil {
		return nil, err
	}

	// Merge with defaults
	defaults := LoadDefaults()
	mergeConfig(defaults, config)

	return config, nil
}

// mergeConfig merges src into dst, preferring src values.
func mergeConfig(dst, src *touta.Config) {
	if src.Framework.Mode != "" {
		dst.Framework = src.Framework
	}
	if src.Server.Port != 0 {
		dst.Server = src.Server
	}
	if len(src.Router.Middleware) > 0 {
		dst.Router = src.Router
	}
	if len(src.Packages) > 0 {
		dst.Packages = src.Packages
	}
	if len(src.App) > 0 {
		dst.App = src.App
	}
}

// FindConfig searches for a config file in the current directory and parents.
func FindConfig() (string, error) {
	names := []string{"touta.yaml", "touta.yml", ".touta.yaml", ".touta.yml"}

	// Start from current directory
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Search up to root
	for {
		for _, name := range names {
			path := filepath.Join(dir, name)
			if _, err := os.Stat(path); err == nil {
				return path, nil
			}
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break // Reached root
		}
		dir = parent
	}

	return "", fmt.Errorf("no config file found")
}
