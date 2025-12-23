package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

// NewCommand creates a new project.
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "new [project-name]",
		Short: "Create a new Toutā project",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			projectName := args[0]
			return createProject(projectName)
		},
	}
	return cmd
}

// InitCommand initializes Toutā in an existing directory.
func InitCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize Toutā in the current directory",
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}
			return initProject(cwd)
		},
	}
	return cmd
}

// ServeCommand starts the development server.
func ServeCommand() *cobra.Command {
	var port int
	var host string

	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the development server",
		RunE: func(cmd *cobra.Command, args []string) error {
			return serve(host, port)
		},
	}

	cmd.Flags().IntVarP(&port, "port", "p", 8080, "Port to listen on")
	cmd.Flags().StringVarP(&host, "host", "H", "localhost", "Host to bind to")

	return cmd
}

// VersionCommand shows version information.
func VersionCommand(version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Toutā v%s\n", version)
		},
	}
	return cmd
}

// createProject scaffolds a new project.
func createProject(name string) error {
	if err := os.MkdirAll(name, 0755); err != nil {
		return fmt.Errorf("failed to create project directory: %w", err)
	}

	if err := initProject(name); err != nil {
		return err
	}

	fmt.Printf("✓ Created new Toutā project: %s\n", name)
	fmt.Printf("\nNext steps:\n")
	fmt.Printf("  cd %s\n", name)
	fmt.Printf("  touta serve\n")

	return nil
}

// initProject initializes Toutā in a directory.
func initProject(dir string) error {
	dirs := []string{
		filepath.Join(dir, "handlers"),
		filepath.Join(dir, "templates"),
		filepath.Join(dir, "static"),
		filepath.Join(dir, "config"),
	}

	for _, d := range dirs {
		if err := os.MkdirAll(d, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", d, err)
		}
	}

	// Initialize Go module if go.mod doesn't exist
	goModPath := filepath.Join(dir, "go.mod")
	if _, err := os.Stat(goModPath); os.IsNotExist(err) {
		projectName := filepath.Base(dir)
		cmd := exec.Command("go", "mod", "init", projectName)
		cmd.Dir = dir
		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to initialize go module: %w\nOutput: %s", err, output)
		}

		// Add chi router dependency (the only external dependency needed for basic projects)
		cmd = exec.Command("go", "get", "github.com/go-chi/chi/v5@latest")
		cmd.Dir = dir
		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to add chi dependency: %w\nOutput: %s", err, output)
		}
	}

	configPath := filepath.Join(dir, "touta.yaml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		config := `framework:
  mode: development
  debug: true
  hot_reload: true
  log_level: info

server:
  host: localhost
  port: 8080

router:
  base_path: /
`
		if err := os.WriteFile(configPath, []byte(config), 0644); err != nil {
			return fmt.Errorf("failed to create touta.yaml: %w", err)
		}
	}

	handlerPath := filepath.Join(dir, "handlers", "hello.go")
	if _, err := os.Stat(handlerPath); os.IsNotExist(err) {
		handler := `package handlers

import (
	"net/http"
)

func Hello(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte("<h1>Hello from Toutā!</h1>"))
}
`
		if err := os.WriteFile(handlerPath, []byte(handler), 0644); err != nil {
			return fmt.Errorf("failed to create example handler: %w", err)
		}
	}

	mainPath := filepath.Join(dir, "main.go")
	if _, err := os.Stat(mainPath); os.IsNotExist(err) {
		main := `package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	
	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Routes
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("<h1>Welcome to Toutā!</h1>"))
	})

	addr := "localhost:8080"
	fmt.Printf("🚀 Server starting on http://%s\n", addr)
	
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
`
		if err := os.WriteFile(mainPath, []byte(main), 0644); err != nil {
			return fmt.Errorf("failed to create main.go: %w", err)
		}
	}

	fmt.Printf("✓ Initialized Toutā project in %s\n", dir)
	return nil
}

// serve starts the development server.
func serve(host string, port int) error {
	// Find project root
	projectRoot, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	// Check if main.go exists
	mainPath := filepath.Join(projectRoot, "main.go")
	if _, err := os.Stat(mainPath); os.IsNotExist(err) {
		fmt.Printf("⚠️  No main.go found in current directory\n")
		fmt.Printf("   Make sure you're in a Toutā project directory\n")
		return nil
	}

	// Check if touta.yaml exists for hot reload config
	configPath := filepath.Join(projectRoot, "touta.yaml")
	hotReloadEnabled := true
	if _, err := os.Stat(configPath); err == nil {
		// Could parse config to check hot_reload setting
		// For now, default to enabled in development
	}

	fmt.Printf("🚀 Starting Toutā development server\n")
	fmt.Printf("   Project: %s\n", filepath.Base(projectRoot))
	fmt.Printf("   Host: %s\n", host)
	fmt.Printf("   Port: %d\n", port)
	fmt.Printf("\n")

	if hotReloadEnabled {
		// Use hot reload
		hr := NewHotReload(projectRoot)
		return hr.Start()
	}

	// Fallback: run without hot reload
	fmt.Printf("🔧 Running without hot reload\n")
	fmt.Printf("   Run 'go run main.go' to start the server\n\n")

	cmd := exec.Command("go", "run", "main.go")
	cmd.Dir = projectRoot
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
